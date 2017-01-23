package webHandlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/molecul/qa_portal/checker"
	"github.com/molecul/qa_portal/model"
	"github.com/molecul/qa_portal/web/middleware"
	"github.com/zalando/gin-oauth2/google"
)

func doError(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusInternalServerError, gin.H{"Internal error": err.Error()})
	panic(err)
}

func UserLoginHandler(ctx *gin.Context) {
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
		if usr == nil {
			doError(ctx, nil)
		}
		// ================================
	} else {
		if err = usr.FillFromGoogle(&gu).Update(); err != nil {
			doError(ctx, err)
			return
		}
	}
	middleware.UserSessionSet(ctx, usr.ID)
	ctx.JSON(http.StatusOK, gin.H{"Hello": "from private", "user": gu, "internal_user": usr})
}

func UserLogoutHandler(ctx *gin.Context) {
	middleware.UserLogout(ctx)
	ctx.Redirect(http.StatusMovedPermanently, "/")
}

func MainPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Main website",
	})
}

func DockerHealthCheckHandler(c *gin.Context) {
	temp_id, _ := strconv.Atoi(time.Now().Format("20060102150405"))
	task := checker.Get().NewTask(&model.Challenge{
		ID: 	    int64(temp_id),
		Image:      "python:2.7",
		TargetPath: "/tmp/task.py",
		Cmd:        "echo \"Inside $CHECKER_FILE:\"; cat $CHECKER_FILE",
		InternalName: "test",
	}, &model.Test{
		ID:          int64(temp_id),
		ChallengeID: int64(temp_id),
		InputFile:   "healthcheck_"+strconv.FormatInt(int64(temp_id), 10),
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	if err := task.Do(ctx); err != nil {
		c.JSON(http.StatusBadGateway, err)
		logrus.Print(err)
	} else {

		c.JSON(http.StatusOK,
			gin.H{
				"exitcode": task.Result.ExitCode,
				"stdout": task.Result.Stdout.String(),
				"stderr": task.Result.Stderr.String()})
	}
	cancel()
}