package match

import (
	"fmt"
	"math/rand"
	"project/app/enums/match"
	"project/app/team"
	"project/services/database"
	"time"
)

type Match struct {
	Id        string
	HomeTeam  team.Team
	AwayTeam  team.Team
	Status    match.MatchStatus
	Minute    int
	HomeScore int
	AwayScore int
	Date      time.Time
}

func PlayOneMatch(oneMatch *Match) *Match {
	if oneMatch.Status != match.Idle {
		return nil
	}

	oneMatch.Status = match.Playing
	var i = 0
	for oneMatch.Minute <= 90 && i < len(oneMatch.HomeTeam.HomePoints) {
		possibleHomeTeamScore := rand.Float64()
		possibleAwayTeamScore := rand.Float64()
		if possibleHomeTeamScore < oneMatch.HomeTeam.HomePoints[i] {
			oneMatch.HomeScore += 1
			//fmt.Printf("%s gol attı!!!\n", oneMatch.HomeTeam.Name)
		}
		if possibleAwayTeamScore < oneMatch.AwayTeam.AwayPoints[i] {
			oneMatch.AwayScore += 1
			//fmt.Printf("%s gol attı!!!\n", oneMatch.AwayTeam.Name)
		}
		//time.Sleep(10 * time.Millisecond)
		oneMatch.Minute += 5
		i += 1
	}
	oneMatch.Status = match.Finished
	fmt.Printf("Skor : %d - %d \n", oneMatch.HomeScore, oneMatch.AwayScore)

	_, err := database.Connection.Exec("insert into score (matchId, homeScore, awayScore) values (?, ?, ?)", oneMatch.Id, oneMatch.HomeScore, oneMatch.AwayScore)
	if err != nil {
		panic(err)
	}

	return oneMatch
}
