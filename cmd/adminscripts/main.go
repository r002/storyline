// This admin script is to be used only for 'initial-loading' new study members.
// Once members are in the system, there is separate logic that updates their
// metrics in the daily job.  Also: This current 'initial-load' logic is limited
// to only getting the first 100 cards per each member.
//
// Last modified: Sat - June 26, 2021

package main

import (
	"github.com/r002/storyline-api/fbservices"
	"github.com/r002/storyline-api/models"
)

var m models.Member

func main() {
	// m = models.Member{
	// 	Fullname:      "Robert Lin",
	// 	Handle:        "r002",
	// 	StartDate:     "2021-05-03T04:00:00Z",
	// 	Uid:           45280066,
	// 	Repo:          "https://github.com/studydash/cards",
	// 	Active:        true,
	// 	StreakCurrent: -1,
	// 	StreakMax:     -1,
	// 	RecordCount:   -1,
	// 	Record:        nil,
	// }
	// m.BuildMember()
	// fbservices.AddMember("testing", m.Handle, m)

	// m = models.Member{
	// 	Fullname:      "Anita Beauchamp",
	// 	Handle:        "anitabe404",
	// 	StartDate:     "2021-05-04T04:00:00Z",
	// 	Uid:           9167395,
	// 	Repo:          "https://github.com/studydash/cards",
	// 	Active:        true,
	// 	StreakCurrent: -1,
	// 	StreakMax:     -1,
	// 	RecordCount:   -1,
	// 	Record:        nil,
	// }
	// m.BuildMember()
	// fbservices.AddMember("testing", m.Handle, m)

	m = models.Member{
		Fullname:      "Shaza Huang",
		Handle:        "shazahuang",
		StartDate:     "2021-06-18T04:00:00Z",
		Uid:           85973779,
		Repo:          "https://github.com/studydash/cards",
		Active:        true,
		StreakCurrent: -1,
		StreakMax:     -1,
		RecordCount:   -1,
		Record:        nil,
	}
	m.BuildMember()
	fbservices.AddMember("testing", m.Handle, m)
}
