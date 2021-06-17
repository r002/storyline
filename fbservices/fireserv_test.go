package fbservices

import (
	"fmt"
	"testing"

	"cloud.google.com/go/firestore"
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

// // Use channels to test this later.
// func TestListenToDoc(t *testing.T) {
// 	ListenToDoc("r002-cloud", "ghUpdatesQa", "123")
// }
