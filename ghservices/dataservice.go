package ghservices

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Card struct {
	Title   string `json:"title"`
	Created string `json:"created_at"`
	Updated string `json:"updated_at"`
	Number  int    `json:"number"`
	User    User
}

type User struct {
	Login string `json:"login"`
	Id    int    `json:"id"`
}

func GetCards() []Card {
	uri := "https://api.github.com/repos/r002/codenewbie/issues?since=2021-05-03&labels=daily%20accomplishment&sort=created&direction=desc&per_page=100"
	resp, err := http.Get(uri)
	if err != nil {
		log.Fatalln(err)
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	// Convert the body to type string
	// sb := string(body)
	// log.Print(sb)

	var cards []Card
	json.Unmarshal(body, &cards)

	fmt.Println(">> len(cards):", len(cards))
	// fmt.Println(">> cards0 Title:", cards[0].Title)
	// fmt.Println(">> cards0 UserHandle:", cards[0].User.Login)
	// fmt.Println(">> cards1 title", cards[1].Title)
	// fmt.Println(">> cards2 title", cards[2].Title)

	return cards
}
