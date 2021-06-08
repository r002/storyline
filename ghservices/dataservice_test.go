package ghservices

import (
	// "context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	// secretmanager "cloud.google.com/go/secretmanager/apiv1"
	// secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

// var ctx context.Context
var ghWebhookSecret []byte
var ghToken []byte

// Get credentials locally
func init() {
	data, err := ioutil.ReadFile("../../secret.json")
	if err != nil {
		fmt.Print(err)
	}
	type GhCredentials struct {
		WebhookSecret string `json:"gh-codenewbie-webook"`
		Token         string `json:"gh-studygroup-bot-pat"`
	}
	var obj GhCredentials
	err = json.Unmarshal(data, &obj)
	if err != nil {
		fmt.Println("error:", err)
	}
	ghWebhookSecret = []byte(obj.WebhookSecret)
	ghToken = []byte(obj.Token)
}

//// Get credentials from GCP Secret Manager
// func init() {
// 	ctx = context.Background()
// 	client, err := secretmanager.NewClient(ctx)
// 	if err != nil {
// 		log.Fatalf("failed to setup client: %v", err)
// 	}
// 	defer client.Close()
// 	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
// 		Name: "projects/r002-cloud/secrets/ghSecret/versions/latest",
// 	}
// 	result, err := client.AccessSecretVersion(ctx, accessRequest)
// 	if err != nil {
// 		log.Fatalf("failed to access secret ghSecret version: %v", err)
// 	}
// 	ghSecret = result.Payload.Data
// 	accessRequest = &secretmanagerpb.AccessSecretVersionRequest{
// 		Name: "projects/r002-cloud/secrets/8xg3vE8Ie_E/versions/latest",
// 	}
// 	result, err = client.AccessSecretVersion(ctx, accessRequest)
// 	if err != nil {
// 		log.Fatalf("failed to access secret token version: %v", err)
// 	}
// 	ghToken = result.Payload.Data
// }

func TestCreateCard(t *testing.T) {
	issue := &IssueShort{
		Title:     "Test from go server-title",
		Body:      "Test from go server-body",
		Labels:    []string{"question"},
		Milestone: 1,
	}
	issueReturn := CreateCard(ghToken, issue)
	fmt.Println(">> Created issue title:", issueReturn.Title)
	fmt.Println(">> Created issue body:", issueReturn.Body)
	fmt.Println(">> Created issue label:", (*issueReturn.Labels)[0].Name)
	fmt.Println(">> Created issue milestone:", issueReturn.Milestone.Title)

	if issueReturn.Title != issue.Title {
		t.Errorf("Title incorrect; got: %s, want: %s", issueReturn.Title, issue.Title)
	}
	if issueReturn.Body != issue.Body {
		t.Errorf("Body incorrect; got: %s, want: %s", issueReturn.Body, issue.Body)
	}
	if (*issueReturn.Labels)[0].Name != issue.Labels[0] {
		t.Errorf("Label incorrect; got: %s, want: %s", (*issueReturn.Labels)[0].Name, issue.Labels[0])
	}
	if issueReturn.Milestone.Title != "Daily Accomplishment" {
		t.Errorf("Milestone incorrect; got: %s, want: %s", issueReturn.Milestone.Title, "Daily Accomplishment")
	}
}

func TestUpdateCard(t *testing.T) {
	expectedTitle := "Test from go server-title"
	expectedBody := "Test from go server-body"
	expectedLabel := "monday"
	expectedMilestone := "Daily Accomplishment"

	issueReturn := UpdateCard(ghToken, 16, "2021-06-08T01:37:41Z")

	fmt.Println(">> Updated issue title:", issueReturn.Title)
	fmt.Println(">> Updated issue body:", issueReturn.Body)
	fmt.Println(">> Updated issue label:", (*issueReturn.Labels)[0].Name)
	fmt.Println(">> Updated issue milestone:", issueReturn.Milestone.Title)

	if issueReturn.Title != expectedTitle {
		t.Errorf("Title incorrect; got: %s, want: %s", issueReturn.Title, expectedTitle)
	}
	if issueReturn.Body != expectedBody {
		t.Errorf("Body incorrect; got: %s, want: %s", issueReturn.Body, expectedBody)
	}
	if (*issueReturn.Labels)[0].Name != expectedLabel {
		t.Errorf("Label incorrect; got: %s, want: %s", (*issueReturn.Labels)[0].Name, expectedLabel)
	}
	if issueReturn.Milestone.Title != expectedMilestone {
		t.Errorf("Milestone incorrect; got: %s, want: %s", issueReturn.Milestone.Title, expectedMilestone)
	}
}

func TestGetCards(t *testing.T) {
	cards := GetCards()
	b, err := json.MarshalIndent(cards, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)
}
