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
		Title:     "Test from go server-title",
		Body:      "Test from go server-body",
		Labels:    []string{"daily accomplishment"},
		Milestone: 1,
	}
	issueReturn := CreateCard(ghToken, issue)
	fmt.Println(">> Created issue title:", issueReturn.Title)
	fmt.Println(">> Created issue body:", issueReturn.Body)
	fmt.Println(">> Created issue label:", (*issueReturn.Labels)[0].Name)
	fmt.Println(">> Created issue milestone:", issueReturn.Milestone.Title)

	testCases := []struct {
		desc string
		got  string
		want string
	}{
		{"Title", issueReturn.Title, issue.Title},
		{"Body", issueReturn.Body, issue.Body},
		{"Label", (*issueReturn.Labels)[0].Name, issue.Labels[0]},
		{"Milestone", issueReturn.Milestone.Title, "Daily Accomplishment"},
	}
	for _, tc := range testCases {
		if tc.got != tc.want {
			t.Errorf("%q incorrect; got: %q, want: %q", tc.desc, tc.got, tc.want)
		}
	}
}

func TestUpdateCard(t *testing.T) {
	// issueReturn := UpdateCard(ghToken, 16, "2021-06-08T01:37:41Z") // Monday
	issueReturn := UpdateCard(ghToken, 18, "2021-06-08T18:32:54Z") // Tuesday

	fmt.Println(">> Updated issue title:", issueReturn.Title)
	fmt.Println(">> Updated issue body:", issueReturn.Body)
	fmt.Println(">> Updated issue label:", (*issueReturn.Labels)[0].Name)
	fmt.Println(">> Updated issue milestone:", issueReturn.Milestone.Title)

	testCases := []struct {
		desc string
		got  string
		want string
	}{
		{"Title", issueReturn.Title, "Test from go server-title"},
		{"Body", issueReturn.Body, "Test from go server-body"},
		// {"Label", (*issueReturn.Labels)[0].Name, "monday"},
		{"Label", (*issueReturn.Labels)[0].Name, "tuesday"},
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
