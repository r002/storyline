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

	"cloud.google.com/go/firestore"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	firebase "firebase.google.com/go"
	"github.com/r002/storyline-api/ghservices"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

var ctx context.Context

func getGitHubSecret() []byte {
	// Create the client.
	ctx = context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatalf("failed to setup client: %v", err)
	}
	defer client.Close()

	// Build the request.
	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: "projects/r002-cloud/secrets/ghSecret/versions/latest",
	}

	// Call the API.
	result, err := client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		log.Fatalf("failed to access secret version: %v", err)
	}
	return result.Payload.Data
}

var ghSecret []byte

func main() {
	ghSecret = getGitHubSecret()
	// log.Printf(">> ghSecret: %s", ghSecret)

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

		fmt.Println(`>> r.Header.Get("X-Hub-Signature-256"):`, headerSig)
		fmt.Println(`>> Verification signature:`, verificationSig)
		fmt.Println(`>> Signature comparison: `, verified)

		if verified == 0 {
			fmt.Println("Error: Signatures don't match!")
			http.Error(w, "Error: Signatures don't match!", 500)
			return
		}

		var result map[string]interface{}
		json.Unmarshal([]byte(buf.String()), &result)
		// b, err := json.MarshalIndent(result, "", "  ")
		// if err != nil {
		// 	fmt.Println("error:", err)
		// 	http.Error(w, "Error parsing the GitHub Webhook JSON", 500)
		// 	return
		// }

		s := result["comment"].(map[string]interface{})["body"].(string)
		b := []byte(">> Received payload - body: " + s)

		app, err := firebase.NewApp(ctx, nil)
		if err != nil {
			log.Fatalln(err)
		}
		client, err := app.Firestore(ctx)
		if err != nil {
			log.Fatalln(err)
		}
		defer client.Close()

		_, _, err = client.Collection("ghUpdates").Add(ctx, map[string]interface{}{
			"body": s,
			"dt":   firestore.ServerTimestamp,
		})
		if err != nil {
			log.Fatalf("Failed adding ghUpdate: %v", err)
		}

		os.Stdout.Write(b)
		w.Write(b)

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
