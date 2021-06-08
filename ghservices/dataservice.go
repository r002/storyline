package ghservices

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	firebase "firebase.google.com/go"
)

type Payload struct {
	Id      int       `json:"id" firestore:"id"`
	Action  string    `json:"action" firestore:"action"`
	Kind    string    `json:"kind" firestore:"kind"`
	Dt      time.Time `json:"dt" firestore:"dt"`
	Issue   Issue     `json:"issue" firestore:"issue"`
	Comment *Comment  `json:"comment,omitempty" firestore:"comment,omitempty"`
}

type Comment struct {
	Created string `json:"created_at" firestore:"created_at"`
	Updated string `json:"updated_at" firestore:"updated_at"`
	Id      int    `json:"id" firestore:"id"`
	Body    string `json:"body" firestore:"body"`
	User    User   `json:"user" firestore:"user"`
}

type Issue struct {
	Number    int        `json:"number" firestore:"number"`
	Title     string     `json:"title" firestore:"title"`
	Id        int        `json:"id" firestore:"id"`
	Body      string     `json:"body" firestore:"body"`
	Created   string     `json:"created_at" firestore:"created_at"`
	Updated   string     `json:"updated_at" firestore:"updated_at"`
	Comments  int        `json:"comments" firestore:"comments"`
	User      User       `json:"user" firestore:"user"`
	Labels    *[]Label   `json:"labels" firestore:"labels"`
	Milestone *Milestone `json:"milestone" firestore:"milestone"`
}

type Milestone struct {
	Title string `json:"title" firestore:"title"`
}

type Label struct {
	Id   int    `json:"id" firestore:"id"`
	Name string `json:"name" firestore:"name"`
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
	payload.Dt = time.Now()

	return payload
}

type IssueShort struct {
	Title     string   `json:"title"`
	Body      string   `json:"body"`
	Labels    []string `json:"labels"`
	Milestone int      `json:"milestone"`
}

type UpdateIssue struct {
	Labels    []string `json:"labels"`
	Milestone int      `json:"milestone"`
}

// This function updates the card with the "Daily Accomplishment" milestone
// and also labels the card with the day it was created. Eg. "Monday"
//
func UpdateCard(ghToken []byte, issueNo int, createdAt string) Issue {
	url := "https://api.github.com/repos/studydash/cards-qa/issues/" + fmt.Sprint(issueNo)
	bearer := "token " + string(ghToken)
	tm, _ := time.Parse(time.RFC3339, createdAt)
	loc, _ := time.LoadLocation("America/New_York") // Hack: Assumes all users are ET. Fix later. 6/8/21
	issue := &UpdateIssue{
		Labels:    []string{strings.ToLower(fmt.Sprint(tm.In(loc).Weekday()))},
		Milestone: 1,
	}
	postBody, _ := json.Marshal(issue)
	responseBody := bytes.NewBuffer(postBody)

	// Create a new request using http
	req, _ := http.NewRequest("POST", url, responseBody)

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)
	req.Header.Add("accept", "application/vnd.github.v3+json")

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error while reading the response body bytes:", err)
	}

	// Print debug payload return
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	s, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(s))

	var issueReturn Issue
	json.Unmarshal(body, &issueReturn)
	return issueReturn
}

func CreateCard(ghToken []byte, issue *IssueShort) Issue {
	url := "https://api.github.com/repos/studydash/cards-qa/issues"
	bearer := "token " + string(ghToken)
	postBody, _ := json.Marshal(issue)
	responseBody := bytes.NewBuffer(postBody)

	// Create a new request using http
	req, _ := http.NewRequest("POST", url, responseBody)

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)
	req.Header.Add("accept", "application/vnd.github.v3+json")

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	// header, err := json.MarshalIndent(resp.Header, "", "  ")
	// if err != nil {
	// 	fmt.Println("Error while reading the response header map:", err)
	// }
	// fmt.Println(string(header))

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error while reading the response body bytes:", err)
	}
	// var result map[string]interface{}
	// json.Unmarshal(body, &result)
	// s, _ := json.MarshalIndent(result, "", "  ")
	// fmt.Println(string(s))

	var issueReturn Issue
	json.Unmarshal(body, &issueReturn)
	return issueReturn
}

func WriteToGitHub(ghToken []byte, issueNo int) {
	url := "https://api.github.com/repos/r002/codenewbie/issues/" + fmt.Sprint(issueNo)

	// Create a Bearer string by appending string access token
	bearer := "token " + string(ghToken)

	issue := &IssueShort{
		// Title:     "Test from go server - title7",
		// Body:      "Test from go server - body7",
		// Labels:    []string{"invalid", "duplicate"},
		Milestone: 1,
	}
	postBody, _ := json.Marshal(issue)
	responseBody := bytes.NewBuffer(postBody)

	// Create a new request using http
	req, _ := http.NewRequest("POST", url, responseBody)

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)
	req.Header.Add("accept", "application/vnd.github.v3+json")

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	// header, err := json.MarshalIndent(resp.Header, "", "  ")
	// if err != nil {
	// 	fmt.Println("Error while reading the response header map:", err)
	// }
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Println("Error while reading the response body bytes:", err)
	// }
	// fmt.Println(string(header))
	// fmt.Println(string(body))
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

	_, err = client.Collection("ghUpdates").Doc("latestUpdate").Set(ctx, payload)
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
