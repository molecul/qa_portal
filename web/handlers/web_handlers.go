package webHandlers

import (
	"context"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/molecul/qa_portal/checker"
	"github.com/molecul/qa_portal/model"
	"github.com/zalando/gin-oauth2/google"
	"github.com/Sirupsen/logrus"
)

func doError(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusInternalServerError, gin.H{"Internal error": err.Error()})
	panic(err)
}

func UserInfoHandler(ctx *gin.Context) {
	gu := ctx.MustGet("user").(google.User)
	usr, err := model.GetUserByEmail(gu.Email)
	if err != nil {
		doError(ctx, err)
		return
	}
	if usr == nil {
		if err = model.CreateUser(new(model.User).FillFromGoogle(&gu)); err != nil {
			doError(ctx, err)
			return
		}
		// TODO ONLY FOR DEBUG. NEED REMOVE
		usr, err = model.GetUserByEmail(gu.Email)
		if err != nil {
			doError(ctx, err)
			return
		}
		// ================================
	} else {
		if err = usr.FillFromGoogle(&gu).Update(); err != nil {
			doError(ctx, err)
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"Hello": "from private", "user": gu, "internal_user": usr})
}

func MainPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Main website",
	})
}

func DockerHealthCheckHandler(c *gin.Context) {
	task := checker.Get().NewTask(&model.Challenge{
		ID:         rand.Int63(),
		Image:      "python:2.7",
		TargetPath: "/tmp/task.py",
		Cmd:        "echo \"Inside $CHECKER_FILE:\"; cat $CHECKER_FILE",
		InternalName: "test",
	}, &model.Test{
		ID:          rand.Int63(),
		ChallengeID: rand.Int63(),
		InputFile:   "healthcheck",
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	if err := task.Do(ctx); err != nil {
		c.JSON(http.StatusBadGateway, err)
		logrus.Print(err)
	} else {
		c.JSON(http.StatusOK, task.Result)
	}
	cancel()
}