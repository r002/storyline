package models

import (
	"fmt"
	"time"

	"github.com/r002/storyline-api/ghservices"
)

type Member struct {
	Fullname      string
	Handle        string
	StartDate     string
	Uid           int
	Repo          string
	Active        bool
	StreakCurrent int
	StreakMax     int // TODO
	RecordCount   int
	Record        map[string]interface{}
	Updated       time.Time `firestore:"Updated,serverTimestamp"`
}

func (m *Member) BuildMember() {
	m.buildRecord()
	m.RecordCount = len(m.Record)
	m.calcStreakCurrent()
}

// Read all of the member's cards by GitHub REST API and build their record
func (m *Member) buildRecord() {
	loc, _ := time.LoadLocation("America/New_York")
	cards := ghservices.GetCards(m.Handle)
	record := make(map[string]interface{})
	for _, card := range cards {
		t, _ := time.Parse(time.RFC3339, card.Created)
		k := t.In(loc).Format("2006-01-02") // Eg. Output: "2021-05-03"
		record[k] = card.Number
	}
	m.Record = record
}

func (m *Member) calcStreakCurrent() {
	streak := 0
	dateCursor := time.Now()
	key := dateCursor.Format("2006-01-02")
	if _, ok := m.Record[key]; ok {
		streak = 1
	}

	for range m.Record {
		dateCursor = dateCursor.Add(-24 * time.Hour)
		key = dateCursor.Format("2006-01-02")
		fmt.Println(">> key:", key)
		if _, ok := m.Record[key]; ok {
			streak++
		} else {
			fmt.Printf(">> Streak broken on: %q; Streak: %v", key, streak)
			break
		}
	}
	m.StreakCurrent = streak
}
