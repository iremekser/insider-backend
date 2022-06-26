package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"project/app/enums/match"
	"project/app/fixture"
	"project/app/team"
	"project/helpers"
	"project/services/cache"
	"project/services/database"
)

type StartOneWeekMatchResponse struct {
	MatchId   string `json:"matchId"`
	HomeScore int    `json:"homeScore"`
	AwayScore int    `json:"awayScore"`
}

type Prediction struct {
	HomeTeam string  `json:"homeTeam"`
	AwayTeam string  `json:"awayTeam"`
	HomeWin  float64 `json:"homeWin"`
	Draw     float64 `json:"draw"`
	AwayWin  float64 `json:"awayWin"`
}

type StartOneWeekResponse struct {
	WeekIndex   int                         `json:"weekIndex"`
	Matches     []StartOneWeekMatchResponse `json:"matches"`
	Predictions []Prediction                `json:"predictions"`
}

type NextWeekOneMatchResponse struct {
	MatchId  string `json:"matchId"`
	HomeTeam string `json:"homeTeam"`
	AwayTeam string `json:"awayTeam"`
}

type NextWeekResponse struct {
	WeekIndex int                        `json:"weekIndex"`
	Matches   []NextWeekOneMatchResponse `json:"matches"`
}

type StartAllWeeksResponse struct {
	Weeks []StartOneWeekResponse `json:"weeks"`
}

type ScoreResponse struct {
	TeamId   string `json:"teamId"`
	TeamName string `json:"teamName"`
	Win      int    `json:"win"`
	Lose     int    `json:"lose"`
	Draw     int    `json:"draw"`
	Point    int    `json:"point"`
}

type AllScoresResponse struct {
	TeamScores  []*ScoreResponse `json:"teamScores"`
	Predictions []Prediction     `json:"predictions"`
}

func GetFixture(context *gin.Context) {
	fix, err := fixture.GetFixture()
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "An error occured."})
		return
	}
	context.IndentedJSON(http.StatusOK, fix)
}

func StartOneWeek(context *gin.Context) {
	fix, _ := fixture.GetFixture()
	weekIndexToPlay, _ := fixture.FindTheWeekToPlay()
	if weekIndexToPlay >= len(fix.Weeks) {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "There is no match to play. You can reset the fixture."})
		return
	}
	weekToPlay := &fix.Weeks[weekIndexToPlay]
	fixture.PlayOneWeek(weekToPlay)

	data, _ := json.Marshal(fix)
	err := cache.Set("fixture", string(data))
	if err != nil {
		panic(err)
	}

	response := StartOneWeekResponse{
		WeekIndex:   weekIndexToPlay,
		Matches:     nil,
		Predictions: nil,
	}
	for _, playedMatch := range weekToPlay.Matches {
		matchResponse := StartOneWeekMatchResponse{
			MatchId:   playedMatch.Id,
			HomeScore: playedMatch.HomeScore,
			AwayScore: playedMatch.AwayScore,
		}
		response.Matches = append(response.Matches, matchResponse)
	}
	for _, nextMatch := range weekToPlay.Matches {
		homeWin, awayWin, draw := helpers.PredictMatch(nextMatch.HomeTeam, nextMatch.AwayTeam)
		nextMatchPred := Prediction{
			HomeTeam: nextMatch.HomeTeam.Name,
			AwayTeam: nextMatch.AwayTeam.Name,
			HomeWin:  homeWin,
			Draw:     draw,
			AwayWin:  awayWin,
		}
		response.Predictions = append(response.Predictions, nextMatchPred)
	}

	context.IndentedJSON(http.StatusOK, response)
}

func StartAllWeek(context *gin.Context) {
	fix, _ := fixture.GetFixture()

	weekIndexToPlay, _ := fixture.FindTheWeekToPlay()
	if weekIndexToPlay == len(fix.Weeks) {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "All matches played"})
		return
	}

	response := StartAllWeeksResponse{
		Weeks: nil,
	}
	for i := weekIndexToPlay; i < len(fix.Weeks); i++ {
		fixture.PlayOneWeek(&fix.Weeks[i])

		data, _ := json.Marshal(fix)
		err := cache.Set("fixture", string(data))
		if err != nil {
			panic(err)
		}

		responseOneWeek := StartOneWeekResponse{
			WeekIndex: func() int {
				weekIndex, _ := fixture.FindTheWeekToPlay()
				return weekIndex - 1
			}(),
			Matches: nil,
		}
		for _, playedMatch := range fix.Weeks[i].Matches {
			matchResponse := StartOneWeekMatchResponse{
				MatchId:   playedMatch.Id,
				HomeScore: playedMatch.HomeScore,
				AwayScore: playedMatch.AwayScore,
			}
			responseOneWeek.Matches = append(responseOneWeek.Matches, matchResponse)
		}
		response.Weeks = append(response.Weeks, responseOneWeek)
	}
	context.IndentedJSON(http.StatusOK, response)
}

func ClearData(context *gin.Context) {
	_, cacheErr := cache.ClearKeys()
	_, dbErr := database.Connection.Exec("truncate score;")

	if cacheErr != nil || dbErr != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "An error occured."})
		return
	}
	context.IndentedJSON(http.StatusOK, gin.H{"message": "all clean"})
}

func GetScores(context *gin.Context) {
	response := AllScoresResponse{
		TeamScores:  nil,
		Predictions: nil,
	}

	teams, _ := team.GetTeams()

	teamMap := make(map[string]*ScoreResponse)

	for _, team := range teams {
		teamScoreResponse := &ScoreResponse{
			TeamId:   team.Id,
			TeamName: team.Name,
			Lose:     0,
			Win:      0,
			Draw:     0,
			Point:    0,
		}
		teamMap[team.Id] = teamScoreResponse
		response.TeamScores = append(response.TeamScores, teamScoreResponse)
	}

	fix, _ := fixture.GetFixture()

	for i1 := 0; i1 < len(fix.Weeks); i1++ {
		for _, _match := range fix.Weeks[i1].Matches {
			if _match.Status == match.Finished {
				if _match.HomeScore > _match.AwayScore {
					teamMap[_match.HomeTeam.Id].Win += 1
					teamMap[_match.HomeTeam.Id].Point += 3

					teamMap[_match.AwayTeam.Id].Lose += 1
				} else if _match.HomeScore < _match.AwayScore {
					teamMap[_match.AwayTeam.Id].Win += 1
					teamMap[_match.AwayTeam.Id].Point += 3

					teamMap[_match.HomeTeam.Id].Lose += 1
				} else {
					teamMap[_match.AwayTeam.Id].Draw += 1
					teamMap[_match.AwayTeam.Id].Point += 1

					teamMap[_match.HomeTeam.Id].Draw += 1
					teamMap[_match.HomeTeam.Id].Point += 1
				}
			}
		}
	}

	weekIndex, _ := fixture.FindTheWeekToPlay()

	if weekIndex != len(fix.Weeks) {
		for _, nextMatch := range fix.Weeks[weekIndex].Matches {
			homeWin, awayWin, draw := helpers.PredictMatch(nextMatch.HomeTeam, nextMatch.AwayTeam)
			nextMatchPred := Prediction{
				HomeTeam: nextMatch.HomeTeam.Name,
				AwayTeam: nextMatch.AwayTeam.Name,
				HomeWin:  homeWin,
				Draw:     draw,
				AwayWin:  awayWin,
			}
			response.Predictions = append(response.Predictions, nextMatchPred)
		}
	}

	context.IndentedJSON(http.StatusOK, response)
}

func GetCurrentWeek(context *gin.Context) {
	fix, _ := fixture.GetFixture()
	weekIndex, _ := fixture.FindTheWeekToPlay()

	response := NextWeekResponse{
		WeekIndex: weekIndex,
		Matches:   nil,
	}

	if weekIndex == len(fix.Weeks) {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "All matches played."})
		return
	}
	for _, nextMatch := range fix.Weeks[weekIndex].Matches {
		matchResponse := NextWeekOneMatchResponse{
			MatchId:  nextMatch.Id,
			HomeTeam: nextMatch.HomeTeam.Name,
			AwayTeam: nextMatch.AwayTeam.Name,
		}
		response.Matches = append(response.Matches, matchResponse)
	}

	context.IndentedJSON(http.StatusOK, response)
}
