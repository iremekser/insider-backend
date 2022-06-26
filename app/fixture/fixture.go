package fixture

import (
	"encoding/json"
	"github.com/google/uuid"
	MatchStatus "project/app/enums/match"
	"project/app/match"
	"project/app/team"
	"project/services/cache"
	"strconv"
)

type Week struct {
	Index   int            `json:"index"`
	Matches []*match.Match `json:"matches"`
}

type Fixture struct {
	Weeks []Week `json:"weeks"`
}

// get fixture from cache
//if fixture doesnt exist, create fixture
func GetFixture() (*Fixture, error) {
	var fix Fixture
	fixtureExist, err := cache.KeyExist("fixture")
	if err != nil {
		return nil, err
	}

	if fixtureExist == 0 {
		fix = CreateFixture()

	} else {
		fixtureString, _ := cache.Get("fixture")
		err := json.Unmarshal([]byte(fixtureString), &fix)
		if err != nil {
			return nil, err
		}
	}
	return &fix, nil
}

// create matches and fixture
// create fixture key in cache
func CreateFixture() Fixture {
	fixtureCombinations := team.CreateCombinations()
	fixture := Fixture{Weeks: nil}
	for i := 0; i < len(fixtureCombinations); i++ {
		newWeek := Week{
			Index:   i,
			Matches: nil,
		}

		for j := 0; j < len(fixtureCombinations[0]); j += 2 {
			newMatch := match.Match{
				Id:        uuid.New().String(),
				HomeTeam:  fixtureCombinations[i][j],
				AwayTeam:  fixtureCombinations[i][j+1],
				Status:    MatchStatus.Idle,
				Minute:    0,
				HomeScore: 0,
				AwayScore: 0,
			}
			newWeek.Matches = append(newWeek.Matches, &newMatch)
		}
		fixture.Weeks = append(fixture.Weeks, newWeek)
	}
	data, _ := json.Marshal(fixture)
	err := cache.Set("fixture", string(data))
	if err != nil {
		panic(err)
	}
	return fixture
}

// get week index to play from cache
func FindTheWeekToPlay() (int, error) {
	weekIndex, err := cache.Get("weekIndex")
	if err != nil {
		return 0, nil
	}
	return strconv.Atoi(weekIndex)
}
func FindWhichMatchToPlay() (int, error) {
	matchIndex, err := cache.Get("matchIndex")
	if err != nil {
		return 0, nil
	}
	return strconv.Atoi(matchIndex)
}

func PlayOneWeek(oneWeek *Week) {
	var i, err = FindWhichMatchToPlay()
	var j int
	if err != nil {
		panic(err)
	}

	for j = i; j < len(oneWeek.Matches); j++ {
		match.PlayOneMatch(oneWeek.Matches[j])
		cache.Set("matchIndex", strconv.Itoa(j+1))
	}

	if j == len(oneWeek.Matches) {
		cache.Set("matchIndex", "0")
	}
	cache.Set("weekIndex", strconv.Itoa(oneWeek.Index+1))
}
