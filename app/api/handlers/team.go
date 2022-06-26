package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"project/app/team"
)

type TeamResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
type AllTeamResponse struct {
	Teams []TeamResponse `json:"teams"`
}

// get team ıds and names from db
func GetTeams(context *gin.Context) {
	allTeamsResponse := AllTeamResponse{
		Teams: nil,
	}

	teamsDb, err := team.GetTeams()
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "An error occured."})
		return
	}

	for _, team := range teamsDb {
		teamResponse := TeamResponse{
			Id:   team.Id,
			Name: team.Name,
		}
		allTeamsResponse.Teams = append(allTeamsResponse.Teams, teamResponse)

	}
	context.IndentedJSON(http.StatusOK, allTeamsResponse)
}
