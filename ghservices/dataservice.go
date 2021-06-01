package ghservices

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
)

type Payload struct {
	Id      int      `json:"id" firestore:"id"`
	Action  string   `json:"action" firestore:"action"`
	Kind    string   `json:"kind" firestore:"kind"`
	Issue   Issue    `json:"issue" firestore:"issue"`
	Comment *Comment `json:"comment,omitempty" firestore:"comment,omitempty"`
}

type Comment struct {
	Created string `json:"created_at" firestore:"created_at"`
	Updated string `json:"updated_at" firestore:"updated_at"`
	Id      int    `json:"id" firestore:"id"`
	Body    string `json:"body" firestore:"body"`
	User    User   `json:"user" firestore:"user"`
}

type Issue struct {
	Number  int    `json:"number" firestore:"number"`
	Title   string `json:"title" firestore:"title"`
	Id      int    `json:"id" firestore:"id"`
	Body    string `json:"body" firestore:"body"`
	Created string `json:"created_at" firestore:"created_at"`
	Updated string `json:"updated_at" firestore:"updated_at"`
	User    User   `json:"user" firestore:"user"`
}

type User struct {
	Login string `json:"login" firestore:"login"`
	Id    int    `json:"id" firestore:"id"`
}

type Card struct {
	Title   string `json:"title"`
	Created string `json:"created_at"`
	Updated string `json:"updated_at"`
	Number  int    `json:"number"`
	Id      int    `json:"id"`
	User    User
}

func TransformIssue(buf string) Payload {
	var result map[string]interface{}
	json.Unmarshal([]byte(buf), &result)

	var payload Payload
	json.Unmarshal([]byte(buf), &payload)

	if _, ok := result["comment"]; ok {
		payload.Kind = "issue_comment"
		payload.Id = payload.Comment.Id
	} else {
		payload.Kind = "issue"
		payload.Id = payload.Issue.Id
	}

	return payload
}

func WriteToFirestore(payload Payload, ctx context.Context) {
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatalln(err)
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	_, err = client.Collection("ghUpdates").Doc(fmt.Sprint(payload.Id)).Set(ctx, payload)
	if err != nil {
		log.Fatalf("Failed adding ghUpdate: %v", err)
	}
}

func GetCards() []Card {
	uri := "https://api.github.com/repos/r002/codenewbie/issues?since=2021-05-03&labels=daily%20accomplishment&sort=created&direction=desc&per_page=100"
	resp, err := http.Get(uri)
	if err != nil {
		log.Fatalln(err)
	}
	// Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	// Convert the body to type string
	// sb := string(body)
	// log.Print(sb)

	var cards []Card
	json.Unmarshal(body, &cards)

	fmt.Println(">> len(cards):", len(cards))
	// fmt.Println(">> cards0 Title:", cards[0].Title)
	// fmt.Println(">> cards0 UserHandle:", cards[0].User.Login)
	// fmt.Println(">> cards1 title", cards[1].Title)
	// fmt.Println(">> cards2 title", cards[2].Title)

	return cards
}
