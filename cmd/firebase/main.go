package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
	// "google.golang.org/api/option"
)

type Payload struct {
	Action  string `json:"action"`
	Issue   Issue
	Comment *Comment `json:"comment,omitempty" firestore:"comment,omitempty"`
}

type Comment struct {
	Body string `json:"body,omitempty" firestore:"body,omitempty"`
	User User   `json:"user,omitempty" firestore:"user,omitempty"`
}

type Issue struct {
	Title string `json:"title"`
	User  User
}

type User struct {
	Login string `json:"login,omitempty" firestore:"login,omitempty"`
}

func main() {
	log.Printf(">> Hello from cmd/firebase/main.go")

	buf := `
{
	"action": "opened",
	"issue": {
		"title": "test111",
		"user": {
			"login": "robert"
		}
	}
}
`
	var payload Payload
	json.Unmarshal([]byte(buf), &payload)
	entirePayload, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	os.Stdout.Write(entirePayload)

	// Use a service account
	ctx := context.Background()
	// sa := option.WithCredentialsFile("../service-account.json")
	// app, err := firebase.NewApp(ctx, nil, sa)
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	iter := client.Collection("authorized").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		fmt.Println(doc.Data())
	}

	_, err = client.Collection("ghUpdates").Doc("123").Set(ctx, map[string]interface{}{
		"first": "Ada",
		"last":  "Lovelace",
		"born":  1815,
		"dt":    firestore.ServerTimestamp,
	}, firestore.MergeAll)

	if err != nil {
		log.Fatalf("Failed adding alovelace: %v", err)
	}

	// _, _, err = client.Collection("ghUpdates").Add(ctx, map[string]interface{}{
	// 	"first":  "Alan",
	// 	"middle": "Mathison",
	// 	"last":   "Turing",
	// 	"born":   1912,
	// })
	// if err != nil {
	// 	log.Fatalf("Failed adding aturing: %v", err)
	// }
}
