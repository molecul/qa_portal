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
	ctx.HTML(http.StatusOK, "index.html", gin.H{"title": "Main website",
		"user": current_user})
}

func LoginHandler(ctx *gin.Context) {
	auth.Login(ctx, "/")
}

func LogoutHandler(ctx *gin.Context) {
	auth.Logout(ctx, "/")
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
