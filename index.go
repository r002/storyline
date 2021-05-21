package main

import (
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

	"github.com/r002/storyline-api/ghservices"
)

type Config struct {
	Secret string `json:"secret"`
}

func LoadConfiguration(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}

var config Config

func main() {
	// Read json config file
	config = LoadConfiguration("../secret.json")
	// log.Printf(">> secret: %s", config.Secret)

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
		key := []byte(config.Secret)
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
		b, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			fmt.Println("error:", err)
			http.Error(w, "Error parsing the GitHub Webhook JSON", 500)
			return
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
