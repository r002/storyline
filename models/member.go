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
	m.calcMaxStreak()
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

func (m *Member) calcMaxStreak() {
	maxStreak, streak := 0, 0
	dateCursor := time.Now()

	startDate, _ := time.Parse(time.RFC3339, m.StartDate)
	startSeconds := startDate.Unix()

	for dateCursor.Unix() >= startSeconds {
		key := dateCursor.Format("2006-01-02")
		if _, ok := m.Record[key]; ok {
			streak++
			// fmt.Printf(">> key: %v; streak: %v\n", key, streak)
		} else {
			// fmt.Printf(">> Streak broken on: %q; Streak: %v\n", key, streak)
			if streak > maxStreak {
				maxStreak = streak
				// fmt.Printf("\t>> New MaxStreak: %v\n", maxStreak)
			}
			streak = 0
		}
		dateCursor = dateCursor.Add(-24 * time.Hour)
	}
	if streak > maxStreak { // This only happens if member has missed zero days
		maxStreak = streak
		fmt.Printf(">> Study member has never missed a day! MaxStreak: %v\n", maxStreak)
	}
	m.StreakMax = maxStreak
}

func (m *Member) calcStreakCurrent() {
	streak := 0
	dateCursor := time.Now()

	startDate, _ := time.Parse(time.RFC3339, m.StartDate)
	startSeconds := startDate.Unix()

	for dateCursor.Unix() >= startSeconds {
		key := dateCursor.Format("2006-01-02")
		if _, ok := m.Record[key]; ok {
			streak++
			fmt.Printf(">> key: %v; streakCurrent: %v\n", key, streak)
		} else {
			// Do not break streakCurrent if member hasn't yet submitted a card today
			if key != time.Now().Format("2006-01-02") {
				fmt.Printf(">> streakCurrent broken on: %q; streakCurrent: %v\n", key, streak)
				break
			}
		}
		dateCursor = dateCursor.Add(-24 * time.Hour)
	}
	m.StreakCurrent = streak
}
