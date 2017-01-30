package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/molecul/qa_portal/checker"
	"github.com/molecul/qa_portal/model"
	"github.com/molecul/qa_portal/web/middleware/auth"
)

func UsersHandler(ctx *gin.Context) {
	order := ctx.Param("order")
	users, _ := model.Users(0, 1000, order)
	ctx.JSON(http.StatusOK, users)
}

func ChallengesHandler(ctx *gin.Context) {
	order := ctx.Param("order")
	challenges, _ := model.Challenges(0, 1000, order)
	ctx.JSON(http.StatusOK, challenges)
}

func DockerHealthCheckHandler(ctx *gin.Context) {
	status := "OK"
	if err := checker.Get().PingDocker(); err != nil {
		status = err.Error()
	}
	ctx.JSON(http.StatusOK, gin.H{"status": status})
}

type HistoryApiElement struct {
	Id          int64
	Challenge   string
	ChallengeId int64
	Duration    string
	State       int8
}

func ProfileHistoryHandler(ctx *gin.Context) {
	// TODO Fix this ugly place
	tests, err := model.Tests(0, 20, false, auth.GetUser(ctx).Id)
	if err != nil {
		panic(err)
	}
	challengeCache := make(map[int64]*model.Challenge)
	result := make([]HistoryApiElement, len(tests))
	for i, test := range tests {
		var challenge *model.Challenge
		challenge, ex := challengeCache[test.ChallengeId]
		if !ex {
			challenge, err = model.GetChallengeById(test.ChallengeId)
			if err != nil {
				panic(err)
			}
			challengeCache[test.ChallengeId] = challenge
		}
		result[i] = HistoryApiElement{
			Id:          test.Id,
			Challenge:   challenge.Name,
			ChallengeId: challenge.Id,
			Duration:    test.Duration.String(),
		}
		if test.Checked != nil {
			if test.IsSucess {
				result[i].State = 1
			} else {
				result[i].State = 2
			}
		}
	}
	ctx.JSON(http.StatusOK, result)
}
