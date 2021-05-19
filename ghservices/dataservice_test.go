package ghservices

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestGetCards(t *testing.T) {
	cards := GetCards()
	b, err := json.MarshalIndent(cards, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)
}
