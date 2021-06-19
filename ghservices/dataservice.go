package ghservices

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/r002/storyline-api/config"
)

var GH_REPO_ENDPOINT string
var FIRESTORE_ENDPOINT string

func init() {
	APP_ENV := config.GetEnvVars().Env
	GH_REPO_ENDPOINT = config.GetEnvVars().GhRepoEndpoint
	FIRESTORE_ENDPOINT = config.GetEnvVars().FirestoreEndpoint

	log.Println(">> APP_ENV:", APP_ENV)
	log.Println(">> GH_REPO_ENDPOINT:", GH_REPO_ENDPOINT)
	log.Println(">> FIRESTORE_ENDPOINT:", FIRESTORE_ENDPOINT)
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
	Title  string   `json:"title"`
	Body   string   `json:"body"`
	Labels []string `json:"labels"`
}

type UpdateIssue struct {
	Labels    []string `json:"labels"`
	Milestone int      `json:"milestone"`
}

func getWeekdayInLoc(dt string, region string) string {
	tm, _ := time.Parse(time.RFC3339, dt) // Eg. "2021-06-08T01:37:41Z"
	loc, _ := time.LoadLocation(region)   // Eg. "America/New_York"
	return fmt.Sprint(tm.In(loc).Weekday())
}

// This function updates the card with the "Daily Accomplishment" milestone
// and also labels the card with the day it was created. Eg. "Monday"
func UpdateCard(ghToken []byte, issue Issue) Issue {
	url := GH_REPO_ENDPOINT + "/issues/" + fmt.Sprint(issue.Number)
	bearer := "token " + string(ghToken)
	weekday := getWeekdayInLoc(issue.Created, "America/New_York") // HACK: Assumes all users are ET. TODO: Fix later. 6/8/21
	updateIssue := &UpdateIssue{
		Labels:    []string{strings.ToLower(weekday)},
		Milestone: 1,
	}
	postBody, _ := json.Marshal(updateIssue)
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

	// // Print debug payload return
	// var result map[string]interface{}
	// json.Unmarshal(body, &result)
	// s, _ := json.MarshalIndent(result, "", "  ")
	// fmt.Println(string(s))

	var issueReturn Issue
	json.Unmarshal(body, &issueReturn)
	return issueReturn
}

func CreateCard(ghToken []byte, issue *IssueShort) Issue {
	url := GH_REPO_ENDPOINT + "/issues"
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

func GetCards() []Card {
	uri := GH_REPO_ENDPOINT + "/issues?since=2021-05-03&milestone=1&sort=created&direction=desc&per_page=100"
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
