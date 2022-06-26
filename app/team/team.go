package team

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"project/services/database"
)

type Team struct {
	Id         string    `json:"id"`
	Name       string    `json:"name"`
	HomePoints []float64 `json:"homePoints"`
	AwayPoints []float64 `json:"awayPoints"`
}

func CreateFixtureMatrix() [][]Team {
	var teams []Team
	teams, err := GetTeams()

	if err != nil {
		panic(err)
	}

	fixtureMatrix := make([][]Team, len(teams))
	for i := range fixtureMatrix {
		fixtureMatrix[i] = make([]Team, len(teams))
	}
	for i := 0; i < len(teams); i++ {
		firstTeam := teams[0]
		teams = teams[1:]
		teams = append(teams, firstTeam)
		fixtureMatrix[i] = teams
	}

	return fixtureMatrix
}

func CreateCombinations() [][]Team {
	var teams []Team
	teams, err := GetTeams()
	if err != nil {
		panic(err)
	}

	combinations := make([][]Team, 0)
	for i := range combinations {
		combinations[i] = make([]Team, 2)
	}

	for _, team1 := range teams {
		for _, team2 := range teams {
			if team1.Id != team2.Id {
				tmp := make([]Team, 0)
				tmp = append(tmp, team1, team2)
				combinations = append(combinations, tmp)
			}
		}
	}

	k := len(teams) * (len(teams) - 1)

	fixture := make([][]Team, 0)
	for i := range fixture {
		fixture[i] = make([]Team, len(teams))
	}

	for k > 0 {
		a := rand.Intn(k)
		firstGroup := combinations[a]
		combinations = append(combinations[:a], combinations[a+1:]...)
		k -= 1
		var b int
		b = rand.Intn(k)
		secondGroup := combinations[b]
		for secondGroup[0].Id == firstGroup[0].Id || secondGroup[0].Id == firstGroup[1].Id || secondGroup[1].Id == firstGroup[0].Id || secondGroup[1].Id == firstGroup[1].Id {
			b = rand.Intn(k)
			secondGroup = combinations[b]
		}
		combinations = append(combinations[:b], combinations[b+1:]...)
		k -= 1

		tempGroups := make([]Team, 0)
		tempGroups = append(tempGroups, firstGroup[0], firstGroup[1], secondGroup[0], secondGroup[1])

		fixture = append(fixture, tempGroups)
	}
	return fixture
}

func GetTeams() ([]Team, error) {
	rows, err := database.Connection.Query("select * from insiderDb.team")
	if err != nil {
		panic(err)
	}
	fmt.Println(rows)

	var teams []Team
	for rows.Next() {
		var newTeam Team
		var homePoints string
		var awayPoints string

		if err := rows.Scan(&newTeam.Id, &newTeam.Name, &homePoints, &awayPoints); err != nil {
			return nil, fmt.Errorf("error occured. %v", err)
		}
		json.Unmarshal([]byte(homePoints), &newTeam.HomePoints)
		json.Unmarshal([]byte(awayPoints), &newTeam.AwayPoints)

		teams = append(teams, newTeam)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occured %v", err)
	}

	return teams, nil
}
