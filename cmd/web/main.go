package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/r002/storyline-api/ghservices"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

var ctx context.Context
var ghSecret []byte
var ghToken []byte

func init() {
	// Create the client.
	ctx = context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatalf("failed to setup client: %v", err)
	}
	defer client.Close()

	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: "projects/r002-cloud/secrets/ghSecret/versions/latest",
	}
	result, err := client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		log.Fatalf("failed to access secret ghSecret version: %v", err)
	}
	ghSecret = result.Payload.Data

	accessRequest = &secretmanagerpb.AccessSecretVersionRequest{
		Name: "projects/r002-cloud/secrets/8xg3vE8Ie_E/versions/latest",
	}
	result, err = client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		log.Fatalf("failed to access secret token version: %v", err)
	}
	ghToken = result.Payload.Data
}

func main() {
	http.HandleFunc("/", indexHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "4567"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe("127.0.0.1:"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/payload" {
		buf := new(strings.Builder)
		_, err := io.Copy(buf, r.Body)
		if err != nil {
			panic(err)
		}

		headerSig := r.Header.Get("X-Hub-Signature-256")
		key := ghSecret
		sig := hmac.New(sha256.New, key)
		sig.Write([]byte(buf.String()))
		verificationSig := "sha256=" + hex.EncodeToString(sig.Sum(nil))
		verified := subtle.ConstantTimeCompare([]byte(headerSig), []byte(verificationSig))

		// fmt.Println(`>> r.Header.Get("X-Hub-Signature-256"):`, headerSig)
		// fmt.Println(`>> Verification signature:`, verificationSig)
		// fmt.Println(`>> Signature comparison: `, verified)

		if verified == 0 {
			fmt.Println("Error: Signatures don't match!")
			http.Error(w, "Error: Signatures don't match!", 500)
			return
		}

		// Unmarshal the json payload
		var result map[string]interface{}
		json.Unmarshal([]byte(buf.String()), &result)

		// Uncomment this section to print the payload to stdout. Eventually, properly log this in 'verbose' mode. 6/2/21
		// entirePayload, err := json.MarshalIndent(result, "", "  ")
		// if err != nil {
		// 	fmt.Println("error:", err)
		// 	http.Error(w, "Error parsing the GitHub Webhook JSON", 500)
		// 	return
		// }
		// os.Stdout.Write(entirePayload)

		payload := ghservices.TransformIssue(buf.String())

		// Milestone all newly "opened" cards that are labeled "daily accomplishment" as "Daily Accomplishment"
		if payload.Action == "opened" {
			for _, label := range *payload.Issue.Labels {
				// fmt.Println(">> Label:", label.Name)
				if label.Name == "daily accomplishment" {
					ghservices.WriteToGitHub(ghToken, payload.Issue.Number)
					w.Write([]byte(">> New issue opened & milestoned as 'Daily Accomplishment'."))
					os.Stdout.Write([]byte(">> New issue opened & milestoned as 'Daily Accomplishment'.\n"))
					return
				}
			}
		}

		// Only act on "Daily Accomplishment" milestone cards
		if _, ok := result["issue"].(map[string]interface{})["milestone"].(map[string]interface{}); ok {
			if result["issue"].(map[string]interface{})["milestone"].(map[string]interface{})["title"].(string) == "Daily Accomplishment" {
				// payload := ghservices.TransformIssue(buf.String())
				ghservices.WriteToFirestore(payload, ctx)
				w.Write([]byte(">> Payload received & processed."))
				os.Stdout.Write([]byte(">> Payload received & processed.\n"))
			} else {
				w.Write([]byte(">> Payload received & ignored. Milestone isn't 'Daily Accomplishment'."))
				os.Stdout.Write([]byte(">> Payload received & ignored. Milestone isn't 'Daily Accomplishment'.\n"))
			}
		} else {
			w.Write([]byte(">> Payload received & ignored. Milestone isn't 'Daily Accomplishment'."))
			os.Stdout.Write([]byte(">> Payload received & ignored. Milestone isn't 'Daily Accomplishment'.\n"))
		}
		return
	}

	if r.URL.Path == "/cards" {
		cards := ghservices.GetCards()
		b, err := json.MarshalIndent(cards, "", "  ")
		if err != nil {
			fmt.Println("error:", err)
			http.Error(w, "Error parsing the GitHub API JSON", 500)
			return
		}
		w.Write(b)
		return
	}

	// If no paths match, return "Not Found"
	http.NotFound(w, r)
}
