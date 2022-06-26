package helpers

import (
	"math"
	"project/app/fixture"
	"project/app/team"
)

type Scores struct {
	Team    team.Team
	For     float64
	Against float64
	Games   float64
	Attack  float64
	Defence float64
}

type Poisson struct {
	Goal  int
	Team1 float64
	Team2 float64
}

func expectedGoals(attackDefenceMap map[string]Scores, totalAverageFor float64, team1 team.Team, team2 team.Team) float64 {
	t1Attack := attackDefenceMap[team1.Name].Attack
	t2Defence := attackDefenceMap[team2.Name].Defence

	return t1Attack * t2Defence * totalAverageFor
}

func factorial(n int) float64 {
	m := 1
	for i := 1; i <= n; i++ {
		m *= i
	}
	return float64(m)
}

func poissonProbability(x int, mean float64) float64 {
	return (math.Pow(mean, float64(x))) / factorial(x) * math.Exp(-mean)
}

func PredictMatch(t1 team.Team, t2 team.Team) (float64, float64, float64) {
	totalGames := 8.0
	depth := 10
	teams, _ := team.GetTeams()
	fix, _ := fixture.GetFixture()
	weekIndex, _ := fixture.FindTheWeekToPlay()

	scoreMap := make(map[string]*Scores)

	for _, predictTeam := range teams {
		scoreMap[predictTeam.Name] = &Scores{
			Team:    predictTeam,
			For:     0,
			Against: 0,
			Games:   0,
			Attack:  0,
			Defence: 0,
		}
	}
	if weekIndex == 0 {
		return 0.33, 0.33, 0.33
	}
	for _, week := range fix.Weeks[:weekIndex] {
		for _, match := range week.Matches {
			scoreMap[match.HomeTeam.Name].For += float64(match.HomeScore)
			scoreMap[match.HomeTeam.Name].Against += float64(match.AwayScore)
			scoreMap[match.HomeTeam.Name].Games += 1

			scoreMap[match.AwayTeam.Name].For += float64(match.AwayScore)
			scoreMap[match.AwayTeam.Name].Against += float64(match.HomeScore)
			scoreMap[match.AwayTeam.Name].Games += 1
		}
	}

	totalFor := 0.0
	totalAgainst := 0.0
	for _, element := range scoreMap {
		totalFor += element.For
		totalAgainst += element.Against
	}

	avgScores := make(map[string]Scores)
	avgTotalFor := totalFor / (totalGames * float64(len(scoreMap)))
	avgTotalAgainst := totalAgainst / (totalGames * float64(len(scoreMap)))
	for key, element := range scoreMap {
		avgScores[key] = Scores{
			For:     element.For / totalGames,
			Against: element.Against / totalGames,
			Attack:  (element.For / totalGames) / avgTotalFor,
			Defence: (element.Against / totalGames) / avgTotalAgainst,
		}
	}

	var poissonArray []Poisson

	t1ExpectedGoals := expectedGoals(avgScores, avgTotalFor, t1, t2)
	t2ExpectedGoals := expectedGoals(avgScores, avgTotalFor, t2, t1)
	for i := 0; i < depth; i++ {
		poissonArray = append(poissonArray, Poisson{
			Goal:  i,
			Team1: poissonProbability(i, t1ExpectedGoals),
			Team2: poissonProbability(i, t2ExpectedGoals),
		})
	}

	w1 := 0.0
	w2 := 0.0

	for i := 1; i < depth; i++ {
		sumTeam1 := 0.0
		sumTeam2 := 0.0
		for j := 0; j < i; j++ {
			sumTeam1 += poissonArray[j].Team1
			sumTeam2 += poissonArray[j].Team2
		}
		probabilityOfWinTeam1 := poissonArray[i].Team1 * sumTeam2
		probabilityOfWinTeam2 := poissonArray[i].Team2 * sumTeam1

		w1 += probabilityOfWinTeam1
		w2 += probabilityOfWinTeam2
	}
	draw := 1 - w1 - w2

	scoreMap[t1.Name].Games += 1
	scoreMap[t1.Name].For += t1ExpectedGoals
	scoreMap[t1.Name].Against += t2ExpectedGoals

	scoreMap[t2.Name].Games += 1
	scoreMap[t2.Name].For += t2ExpectedGoals
	scoreMap[t2.Name].Against += t1ExpectedGoals

	return w1, w2, draw
}
