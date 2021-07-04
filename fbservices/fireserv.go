// Houses all code that talks with Firestore

package fbservices

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/r002/storyline-api/ghservices"
	"github.com/r002/storyline-api/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var client *firestore.Client
var ctx context.Context

func getClient() *firestore.Client {
	if client == nil {
		log.Println(">> Creating a new client!")
		// projectID := "r002-cloud"
		projectID := "studydash-qa"
		ctx = context.Background()
		c, err := firestore.NewClient(ctx, projectID)
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}
		// Close client when done with
		// defer client.Close()
		client = c
	}
	return client
}

func ReadCollection(collection string) []*firestore.DocumentSnapshot {
	client = getClient()
	// defer client.Close()
	iter := client.Collection(collection).Documents(ctx)
	all, err := iter.GetAll()
	if err != nil {
		log.Fatalf("Error fetching collection: %v", err)
	}
	return all
}

func SendPayload(collection string, doc string, payload ghservices.Payload) error {
	client = getClient()
	// defer client.Close()
	_, err := client.Collection(collection).Doc(doc).Set(ctx, payload)

	if err != nil {
		log.Fatalf("Failed sending payload to %s/%s: %v", collection, doc, err)
	}
	return err
}

func CreateDoc(collection string, doc string, payload map[string]interface{}) error {
	client = getClient()
	// defer client.Close()
	_, err := client.Collection(collection).Doc(doc).Set(ctx, payload, firestore.MergeAll)
	// fmt.Printf(">> rs: %v", rs) // Output: &{2021-06-17 12:13:20.465879 +0000 UTC}

	if err != nil {
		log.Fatalf("Failed adding to %s/%s: %v", collection, doc, err)
	}
	return err
}

func IncrementMemberStreak(issue ghservices.Issue) {
	// Read the member from Firestore to see if the member has yet contributed a card for today.
	handle := issue.User.Login
	m := GetMember(handle)
	loc, _ := time.LoadLocation("America/New_York")
	t, _ := time.Parse(time.RFC3339, issue.Created)
	date := t.In(loc).Format("2006-01-02")

	if _, ok := m.Record[date]; !ok {
		m.Record[date] = issue.Number
		m.RecordCount = len(m.Record)
		m.CalcStreakCurrent()
		m.CalcMaxStreak()
		m.CalcDaysJoined()

		AddMember("studyMembers", m.Handle, m)
		log.Println(">> Member updated", handle)
	} else {
		log.Panicln(">> Member not updated; card for date already exists", handle)
	}
}

func AddMember(collection string, doc string, member models.Member) error {
	client = getClient()
	// defer client.Close()
	_, err := client.Collection(collection).Doc(doc).Set(ctx, member)
	// fmt.Printf(">> rs: %v", rs) // Output: &{2021-06-17 12:13:20.465879 +0000 UTC}

	if err != nil {
		log.Fatalf("Failed adding to %s/%s: %v", collection, doc, err)
	}
	return err
}

func GetMember(userHandle string) models.Member {
	client = getClient()
	// defer client.Close()
	dsnap, _ := client.Collection("studyMembers").Doc(userHandle).Get(ctx)
	var m models.Member
	dsnap.DataTo(&m)
	return m
}

func ReadDoc(collection string, doc string) map[string]interface{} {
	client = getClient()
	// defer client.Close()
	dsnap, err := client.Collection(collection).Doc(doc).Get(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	return dsnap.Data()
}

func DeleteDoc(collection string, doc string) error {
	client = getClient()
	// defer client.Close()
	_, err := client.Collection(collection).Doc(doc).Delete(ctx)
	return err
}

func ListenToDoc(projectId string, collection string, doc string) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		fmt.Printf("firestore.NewClient: %v", err)
	}
	defer client.Close()

	it := client.Collection(collection).Doc(doc).Snapshots(ctx)
	for {
		snap, err := it.Next()
		// DeadlineExceeded will be returned when ctx is cancelled.
		if status.Code(err) == codes.DeadlineExceeded {
			fmt.Println("Timeout exceeded")
			return
		}
		if err != nil {
			fmt.Printf("Snapshots.Next: %v", err)
		}
		if !snap.Exists() {
			fmt.Printf("Document no longer exists\n")
		}
		fmt.Printf("Received document snapshot: %v\n", snap.Data())
	}
}
