package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

var ctx context.Context
var ghToken string

func init() {
	// Create the client.
	ctx = context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatalf("failed to setup client: %v", err)
	}
	defer client.Close()

	// Build the request.
	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: "projects/r002-cloud/secrets/8xg3vE8Ie_E/versions/latest",
	}

	// Call the API.
	result, err := client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		log.Fatalf("failed to access secret version: %v", err)
	}
	ghToken = string(result.Payload.Data)
}

type Issue struct {
	// Title     string   `json:"title"`
	// Body      string   `json:"body"`
	// Labels    []string `json:"labels"`
	Milestone int `json:"milestone"`
}

func main() {
	url := "https://api.github.com/repos/r002/codenewbie/issues/78"

	// Create a Bearer string by appending string access token
	bearer := "token " + ghToken

	issue := &Issue{
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

	header, err := json.MarshalIndent(resp.Header, "", "  ")
	if err != nil {
		fmt.Println("Error while reading the response header map:", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error while reading the response body bytes:", err)
	}
	fmt.Println(string(header))
	fmt.Println(string(body))
}
