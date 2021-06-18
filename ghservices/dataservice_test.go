package ghservices

import (
	// "context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"
	// secretmanager "cloud.google.com/go/secretmanager/apiv1"
	// secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

// var ctx context.Context
var ghWebhookSecret []byte
var ghToken []byte

// Get credentials locally. Used only in this test file.
func init() {
	data, err := ioutil.ReadFile("../../secret.json")
	if err != nil {
		fmt.Print(err)
	}
	type GhCredentials struct {
		Webhook string `json:"gh-cards-qa-webook"` // Currently unused in tests. 6/18/21
		Token   string `json:"gh-studygroup-bot-tok"`
	}
	var obj GhCredentials
	err = json.Unmarshal(data, &obj)
	if err != nil {
		fmt.Println("error:", err)
	}
	ghWebhookSecret = []byte(obj.Webhook)
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

func TestGetWeekdayInLoc(t *testing.T) {
	testCases := []struct {
		dt     string
		region string
		want   string
	}{
		{"2021-06-08T01:37:41Z", "America/New_York", "Monday"},
		{"2021-06-08T18:32:54Z", "America/New_York", "Tuesday"},
	}

	for _, tc := range testCases {
		got := getWeekdayInLoc(tc.dt, tc.region)
		if got != tc.want {
			t.Errorf("%q/%q ➡ %q; want: %q", tc.dt, tc.region, got, tc.want)
		}
	}
}

func TestCreateCard(t *testing.T) {
	issue := &IssueShort{
		Title:  "Test local from go server-title",
		Body:   "Test local from go server-body",
		Labels: []string{"daily accomplishment"},
	}
	issueReturn := CreateCard(ghToken, issue)
	t.Log(">> Created issue title:", issueReturn.Title)
	t.Log(">> Created issue body:", issueReturn.Body)
	t.Log(">> Created issue label:", (*issueReturn.Labels)[0].Name)

	testCases := []struct {
		desc string
		got  string
		want string
	}{
		{"Title", issueReturn.Title, issue.Title},
		{"Body", issueReturn.Body, issue.Body},
		{"Label", (*issueReturn.Labels)[0].Name, issue.Labels[0]},
	}
	for _, tc := range testCases {
		if tc.got != tc.want {
			t.Errorf("%q incorrect; got: %q, want: %q", tc.desc, tc.got, tc.want)
		}
	}
}

func TestUpdateCard(t *testing.T) {
	issueInput := &IssueShort{
		Title:  "Test local from go server-title",
		Body:   "Test local from go server-body",
		Labels: []string{"daily accomplishment"},
	}
	issueReturn := CreateCard(ghToken, issueInput)
	issueReturn = UpdateCard(ghToken, issueReturn)

	t.Log(">> Updated issue title:", issueReturn.Title)
	t.Log(">> Updated issue body:", issueReturn.Body)
	t.Log(">> Updated issue label:", (*issueReturn.Labels)[0].Name)
	t.Log(">> Updated issue milestone:", issueReturn.Milestone.Title)

	testCases := []struct {
		desc string
		got  string
		want string
	}{
		{"Title", issueReturn.Title, issueInput.Title},
		{"Body", issueReturn.Body, issueInput.Body},
		{"Label", (*issueReturn.Labels)[0].Name, strings.ToLower(fmt.Sprint(time.Now().Weekday()))},
		{"Milestone", issueReturn.Milestone.Title, "Daily Accomplishment"},
	}
	for _, tc := range testCases {
		if tc.got != tc.want {
			t.Errorf("%q incorrect; got: %q, want: %q", tc.desc, tc.got, tc.want)
		}
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
