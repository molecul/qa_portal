package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/molecul/qa_portal/checker"
	"github.com/molecul/qa_portal/model"
	"github.com/molecul/qa_portal/web/middleware/auth"
)

func MainPageHandler(ctx *gin.Context) {
	current_user := auth.GetUser(ctx)
	ctx.HTML(http.StatusOK, "pages/index", gin.H{"user": current_user})
}

func ChallengesWebHandler(ctx *gin.Context) {
	current_user := auth.GetUser(ctx)
	ctx.HTML(http.StatusOK, "pages/challenges", gin.H{"user": current_user})
}

func ProfileHandler(ctx *gin.Context) {
	current_user := auth.GetUser(ctx)
	ctx.HTML(http.StatusOK, "pages/profile", gin.H{"user": current_user})
}

func ScoreboardHandler(ctx *gin.Context) {
	current_user := auth.GetUser(ctx)
	ctx.HTML(http.StatusOK, "pages/scoreboard", gin.H{"user": current_user})
}

func LoginHandler(ctx *gin.Context) {
	auth.Login(ctx, "/")
}

func LogoutHandler(ctx *gin.Context) {
	auth.Logout(ctx, "/")
}

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

func DockerHealthCheckHandler(c *gin.Context) {
	temp_id := time.Now().Unix()
	task := checker.Get().NewTask(&model.Challenge{
		ID:           temp_id,
		Image:        "python:2.7",
		TargetPath:   "/tmp/task.py",
		Cmd:          "echo \"Inside $CHECKER_FILE:\"; cat $CHECKER_FILE",
		InternalName: "test",
	}, &model.Test{
		ID:          temp_id,
		ChallengeID: temp_id,
		InputFile:   fmt.Sprintf("healthcheck_%v", temp_id),
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	if err := task.Do(ctx); err != nil {
		c.JSON(http.StatusBadGateway, err.Error())
		logrus.Print(err)
	} else {
		c.JSON(http.StatusOK, gin.H{"exitcode": task.Result.ExitCode,
			"stdout": task.Result.Stdout.String(),
			"stderr": task.Result.Stderr.String()})
	}
	cancel()
}
