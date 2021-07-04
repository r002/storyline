package fbservices

import (
	"fmt"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/r002/storyline-api/ghservices"
)

func TestReadCollection(t *testing.T) {
	collection := "testing"
	arr := ReadCollection(collection)
	if len(arr) == 0 {
		t.Logf("Firestore collection %q unexpectedly empty. Test failed!", collection)
		t.Fail()
	}
	fmt.Printf("%q size: %d\n", collection, len(arr))
}

func TestCreateDoc(t *testing.T) {
	payload := map[string]interface{}{
		"title":        "Test - title",
		"desc":         "Test - desc",
		"randomNumber": 555,
		"dt":           firestore.ServerTimestamp,
	}
	err := CreateDoc("testing", "testDoc", payload)
	if err != nil {
		t.Log("Firestore doc creation test failed!", err)
		t.Fail()
	}
}

func TestReadDoc(t *testing.T) {
	doc := ReadDoc("testing", "testDoc")
	if doc == nil {
		t.Log("Firestore doc read failed!")
		t.Fail()
	}
	fmt.Printf("Document data: %#v\n", doc)
}

func TestDeleteDoc(t *testing.T) {
	err := DeleteDoc("testing", "testDoc")
	if err != nil {
		t.Log("Firestore doc deletion test failed!", err)
		t.Fail()
	}
}

// This test sends a payload to Firestore
func TestSendPayload(t *testing.T) {
	mockUser := ghservices.User{
		Login: "testUser",
		Id:    789,
	}
	mockIssue := ghservices.Issue{
		Number:    234,
		Title:     "TestPayload Title",
		Id:        345,
		Body:      "TestPayload Body",
		Created:   "TestPayload Created",
		Updated:   "TestPayload Updated",
		Comments:  0,
		User:      mockUser,
		Labels:    nil,
		Milestone: nil,
	}
	mockPayload := ghservices.Payload{
		Id:      123,
		Action:  "Test Creation",
		Kind:    "issue",
		Dt:      time.Now(),
		Issue:   mockIssue,
		Comment: nil,
	}
	err := SendPayload("ghUpdatesQa", "latestUpdate", mockPayload)
	if err != nil {
		t.Fatalf("Failed sending payload: %v", err)
	}
}

func TestIncrementMemberStreak(t *testing.T) {
	mockUser := ghservices.User{
		Login: "r002",
		Id:    45280066,
	}
	mockIssue := ghservices.Issue{
		Number:    167,
		Title:     "Submit GitHub `Issues` Feedback to azenMatt",
		Id:        936163216,
		Body:      "Dummy Body",
		Created:   "2021-07-03T03:55:43Z",
		Updated:   "2021-07-03T03:55:45Z",
		Comments:  0,
		User:      mockUser,
		Labels:    nil,
		Milestone: nil,
	}
	IncrementMemberStreak(mockIssue)
}

// // Use channels to test this later.
// func TestListenToDoc(t *testing.T) {
// 	ListenToDoc("r002-cloud", "ghUpdatesQa", "123")
// }
