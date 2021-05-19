package main

// [START import]
import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/r002/storyline-api/ghservices"
)

// [END import]
// [START main_func]

func main() {
	http.HandleFunc("/", indexHandler)

	// [START setting_port]
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe("127.0.0.1:"+port, nil); err != nil {
		log.Fatal(err)
	}
	// [END setting_port]
}

// [END main_func]

// [START indexHandler]

// indexHandler responds to requests with our greeting.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	// fmt.Fprint(w, "Hello, GCP! 5/18/21")
	cards := ghservices.GetCards()
	b, err := json.MarshalIndent(cards, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
		http.Error(w, "Error pasring the GitHub API JSON", 123)
		return
	}
	w.Write(b)
}

// [END indexHandler]
