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
	"github.com/r002/storyline-api/config"
	"github.com/r002/storyline-api/fbservices"
	"github.com/r002/storyline-api/ghservices"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

var ctx context.Context
var ghWebhook []byte
var ghToken []byte
var FIRESTORE_ENDPOINT string

func init() {
	APP_ENV := config.GetEnvVars().Env
	keyGhWebhook := config.GetEnvVars().KeyGhWebhook
	keyGhToken := config.GetEnvVars().KeyGhToken
	FIRESTORE_ENDPOINT = config.GetEnvVars().FirestoreEndpoint

	log.Println(">> Setting up server. Env:", APP_ENV)
	log.Println(">> keyGhWebhook:", keyGhWebhook)
	log.Println(">> keyGhToken:", keyGhToken)
	log.Println(">> FIRESTORE_ENDPOINT:", FIRESTORE_ENDPOINT)

	ctx = context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatalf("failed to setup client: %v", err)
	}
	defer client.Close()

	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/r002-cloud/secrets/%s/versions/latest", keyGhWebhook),
	}
	result, err := client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		log.Fatalf("failed to access secret %q version: %v", keyGhWebhook, err)
	}
	ghWebhook = result.Payload.Data

	accessRequest = &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/r002-cloud/secrets/%s/versions/latest", keyGhToken),
	}
	result, err = client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		log.Fatalf("failed to access secret %q version: %v", keyGhToken, err)
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
	switch path := r.URL.Path; path {
	case "/cards":
		handleCards(w)
	case "/payload":
		handlePayload(w, r)
	case "/info":
		handleInfo(w)
	default:
		http.NotFound(w, r) // If no paths match, return "Not Found"
	}
}

func handleInfo(w http.ResponseWriter) {
	s := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
			<style>
				body {
					margin: 40px;
				}

				.wrapper {
					display: grid;
					grid-template-columns: 100px 300px;
					grid-gap: 10px;
					background-color: #fff;
					color: #444;
				}

				.box {
					background-color: #444;
					color: #fff;
					border-radius: 5px;
					padding: 5px;
					font-size: 12px;
				}
			</style>
		<body>
			<div class="wrapper">
				<div class="box">Env</div>								<div class="box">%[1]s</div>
				<div class="box">KeyGhWebhook</div>				<div class="box">%[2]s</div>
				<div class="box">KeyGhToken</div>					<div class="box">%[3]s</div>
				<div class="box">GhRepoEndpoint</div> 		<div class="box">%[4]s</div>
				<div class="box">GCP Project</div>				<div class="box">%[5]s</div>
				<div class="box">FirestoreEndpoint</div>	<div class="box">%[6]s</div>
				<div class="box">Version</div>						<div class="box">%[7]s</div>
				<div class="box">Last built</div>					<div class="box">%[8]s</div>
				<div class="box">Notes</div>							<div class="box">%[9]s</div>
			</div>
		</body>
		</html>`,
		config.GetEnvVars().Env,
		config.GetEnvVars().KeyGhWebhook,
		config.GetEnvVars().KeyGhToken,
		config.GetEnvVars().GhRepoEndpoint,
		config.GetEnvVars().GcpProject,
		config.GetEnvVars().FirestoreEndpoint,
		"0.0.17",
		"Tuesday - July 20, 2021",
		"🧪 Testing GitHub Action on commit-to-fork & PR to 'main'",
	)
	w.Write([]byte(s))
}

// Get all cards
func handleCards(w http.ResponseWriter) {
	cards := ghservices.GetCards("r002")
	b, err := json.MarshalIndent(cards, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
		http.Error(w, "Error parsing the GitHub API JSON", 500)
		return
	}
	w.Write(b)
}

func handlePayload(w http.ResponseWriter, r *http.Request) {
	buf := new(strings.Builder)
	_, err := io.Copy(buf, r.Body)
	if err != nil {
		panic(err)
	}

	headerSig := r.Header.Get("X-Hub-Signature-256")
	key := ghWebhook
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

	// Set all newly "opened" cards that have a Label="daily accomplishment" to Milestone="Daily Accomplishment"
	// Add the day-of-the-week label (eg. `monday`, `tuesday`, etc)
	// Add label, YYYY-MM (eg. `2021-06`, etc)
	if payload.Action == "opened" {
		for _, label := range *payload.Issue.Labels {
			if label.Name == "daily accomplishment" {
				ghservices.UpdateCard(ghToken, payload.Issue)
				fbservices.IncrementMemberStreak(payload.Issue)
				w.Write([]byte(">> New issue opened & milestoned as 'Daily Accomplishment'."))
				os.Stdout.Write([]byte(">> New issue opened & milestoned as 'Daily Accomplishment'.\n"))
				return
			}
		}
	}

	// Only act on "Daily Accomplishment" milestone cards
	if _, ok := result["issue"].(map[string]interface{})["milestone"].(map[string]interface{}); ok {
		if result["issue"].(map[string]interface{})["milestone"].(map[string]interface{})["title"].(string) == "Daily Accomplishment" {
			fbservices.SendPayload(FIRESTORE_ENDPOINT, "latestUpdate", payload)
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
}
