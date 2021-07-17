package models

import (
	"fmt"
	"math"
	"time"

	"github.com/r002/storyline-api/ghservices"
)

type Streak struct {
	StartDate string
	EndDate   string
	Days      int
}

type Card struct {
	Date   string
	Number int
}

type Member struct {
	Fullname      string
	Handle        string
	StartDate     string
	Uid           int
	Repo          string
	Active        bool
	StreakCurrent Streak
	StreakMax     Streak
	RecordCount   int
	Record        map[string]int
	DaysJoined    int
	Updated       time.Time
	LastCard      Card
}

func (m *Member) BuildMember() {
	m.buildRecord()
	m.CalcStreakCurrent()
	m.CalcMaxStreakAndLastCard()
	m.CalcDaysJoined()
	m.Updated = time.Now()
}

// Read all of the member's cards by GitHub REST API and build their record
func (m *Member) buildRecord() {
	loc, _ := time.LoadLocation("America/New_York")
	cards := ghservices.GetCards(m.Handle)
	record := make(map[string]int)
	for _, card := range cards {
		t, _ := time.Parse(time.RFC3339, card.Created)
		k := t.In(loc).Format("2006-01-02") // Eg. Output: "2021-05-03"
		record[k] = card.Number
	}
	m.Record = record
	m.RecordCount = len(m.Record)

	// Set lastCard to be the first card the studyMember ever submitted.
	// The system assumes that the member's start date *always* has a card submission.
	// Put another way: Until you properly submit a card, you are not in the system!
	// A study member's startDate will *always* be the day of their first card submission.
}

func (m *Member) CalcDaysJoined() {
	startDate, _ := time.Parse(time.RFC3339, m.StartDate)
	m.DaysJoined = int(math.Ceil(time.Since(startDate).Hours() / 24))
}

func (m *Member) CalcMaxStreakAndLastCard() {
	loc, _ := time.LoadLocation("America/New_York")
	dateCursor := time.Now().In(loc)
	streak := Streak{}

	startDate, _ := time.Parse(time.RFC3339, m.StartDate)
	startSeconds := startDate.Unix()

	// By default, a member's LastCard is always their first card
	m.LastCard = Card{
		Date:   startDate.Format("2006-01-02"),
		Number: 0, // dummy place holder for cardNo
	}

	for dateCursor.Unix() >= startSeconds {
		key := dateCursor.Format("2006-01-02")
		if cardNo, ok := m.Record[key]; ok {

			// Set LastCard if it's more recent
			if key > m.LastCard.Date {
				m.LastCard = Card{
					Date:   key,
					Number: cardNo,
				}
			}

			// If this is the beginning of a new streak, track the date
			if streak.Days == 0 {
				streak.EndDate = key
			}
			streak.Days++
			// fmt.Printf(">> key: %v; streak: %v\n", key, streak)
		} else {
			// fmt.Printf(">> Streak broken on: %q; Streak: %v\n", key, streak)
			streak.StartDate = dateCursor.Add(24 * time.Hour).Format("2006-01-02")
			if streak.Days > m.StreakMax.Days {
				m.StreakMax = streak
				// fmt.Printf("\t>> New MaxStreak: %v\n", streak)
			}
			// Reset the streak
			streak = Streak{}
		}
		streak.StartDate = key
		dateCursor = dateCursor.Add(-24 * time.Hour)
	}
	if streak.Days > m.StreakMax.Days { // This only happens if member has missed zero days
		m.StreakMax = streak
		fmt.Printf(">> Study member has never missed a day! MaxStreak: %v\n", m.StreakMax)
	}
}

func (m *Member) CalcStreakCurrent() {
	loc, _ := time.LoadLocation("America/New_York")
	dateCursor := time.Now().In(loc)
	streak := Streak{}

	startDate, _ := time.Parse(time.RFC3339, m.StartDate)
	startSeconds := startDate.Unix()

	i := 0
	for dateCursor.Unix() >= startSeconds {
		key := dateCursor.Format("2006-01-02")
		if _, ok := m.Record[key]; ok {
			// If this is the beginning of a new streak, track the date
			if streak.Days == 0 {
				streak.EndDate = key
			}
			streak.Days++
			// fmt.Printf(">> kv: %v | %v; streakCurrent: %v\n", key, cardNo, streak)
		} else {
			// Do not break streakCurrent if member hasn't yet submitted a card today
			if i != 0 {
				fmt.Printf(">> %v | streakCurrent broken on: %q; streakCurrent: %v\n", m.Handle, key, streak)
				break
			}
		}
		streak.StartDate = key
		dateCursor = dateCursor.Add(-24 * time.Hour)
		i++
	}
	// If there is no current streak, set streak.EndDate = streak.StartDate
	// streak.EndDate and streak.StartDate should both never be empty
	if streak.Days == 0 {
		streak.EndDate = streak.StartDate
	}
	m.StreakCurrent = streak
}
