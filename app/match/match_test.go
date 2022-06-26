package match

import (
	"github.com/google/uuid"
	"project/app/enums/match"
	"project/app/team"
	"testing"
	"time"
)

func TestReturnNil_WhenMatchStatusIsIdle(t *testing.T) {
	testMatch := Match{
		Id: uuid.New().String(),
		HomeTeam: team.Team{
			Id:         uuid.New().String(),
			Name:       "Test Team 1",
			HomePoints: nil,
			AwayPoints: nil,
		},
		AwayTeam: team.Team{
			Id:         uuid.New().String(),
			Name:       "Test Team 2",
			HomePoints: nil,
			AwayPoints: nil,
		},
		Status:    match.Finished,
		Minute:    0,
		HomeScore: 0,
		AwayScore: 0,
		Date:      time.Now(),
	}
	result := PlayOneMatch(&testMatch)

	var expectedResult *Match

	if result != expectedResult {
		t.Errorf("Test Failed. Because expected value is %v, but got %v", expectedResult, result)
	}
	t.Logf("Test Passed. Because expected value is %v, got %v", expectedResult, result)
}

func Test(t *testing.T) {
	testMatch := Match{
		Id: uuid.New().String(),
		HomeTeam: team.Team{
			Id:         uuid.New().String(),
			Name:       "Test Team 1",
			HomePoints: nil,
			AwayPoints: nil,
		},
		AwayTeam: team.Team{
			Id:         uuid.New().String(),
			Name:       "Test Team 2",
			HomePoints: nil,
			AwayPoints: nil,
		},
		Status:    match.Idle,
		Minute:    91,
		HomeScore: 2,
		AwayScore: 1,
		Date:      time.Now(),
	}
	//
	//db, mock, err := sqlmock.New()
	//if err != nil {
	//	t.Fatalf("an error'%s' was not expected when opening a stub database connection", err)
	//}
	//defer db.Close()
	//mock.ExpectBegin()
	//mock.ExpectExec("insert into score").WithArgs(testMatch.Id, testMatch.HomeScore, testMatch.AwayScore).WillReturnResult(sqlmock.NewResult(1, 1))
	//mock.ExpectCommit()

	PlayOneMatch(&testMatch)

}
