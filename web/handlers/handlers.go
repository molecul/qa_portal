package handlers

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/molecul/qa_portal/checker"
	"github.com/molecul/qa_portal/model"
	"github.com/molecul/qa_portal/web/middleware/auth"
	"github.com/molecul/qa_portal/web/views"
	"github.com/shurcooL/github_flavored_markdown"
)

func MainPageHandler(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "pages/index", gin.H{"user": auth.GetUser(ctx)})
}

func ChallengesWebHandler(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "pages/challenges", gin.H{"user": auth.GetUser(ctx)})
}

func ProfileHandler(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "pages/profile", gin.H{"user": auth.GetUser(ctx)})
}

func ScoreboardHandler(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "pages/scoreboard", gin.H{"user": auth.GetUser(ctx)})
}

func TestHandler(ctx *gin.Context) {
	testId, err := strconv.ParseInt(ctx.Param("test"), 10, 0)
	if err != nil {
		views.RenderError(ctx, err)
		return
	}
	test, err := model.GetTestById(testId)
	if err != nil {
		views.RenderError(ctx, err)
		return
	}
	user := auth.GetUser(ctx)
	if test == nil || test.UserId != user.Id {
		ctx.JSON(http.StatusNotFound, 404)
	}
	input, _ := ioutil.ReadFile(test.GetInputFileName())
	output, _ := ioutil.ReadFile(test.GetOutputFileName())

	ctx.HTML(http.StatusOK, "pages/test",
		gin.H{"user": user,
			"test":        test,
			"test_input":  string(input),
			"test_output": string(output)})
}

func SolveHandlerGet(ctx *gin.Context) {
	challengeId, err := strconv.ParseInt(ctx.Param("challenge"), 10, 0)
	if err != nil {
		views.RenderError(ctx, err)
		return
	}
	challenge, err := model.GetChallengeById(challengeId)
	if err != nil {
		views.RenderError(ctx, err)
		return
	}
	if challenge == nil {
		ctx.Redirect(http.StatusFound, "/challenges")
		return
	}
	description := template.HTML(github_flavored_markdown.Markdown([]byte(challenge.Description)))
	ctx.HTML(http.StatusOK, "pages/solve", gin.H{
		"user":        auth.GetUser(ctx),
		"description": description})
}

func SolveHandlerPost(ctx *gin.Context) {
	code := ctx.PostForm("code")
	if code == "" {
		views.RenderRedirect(ctx, ctx.Request.URL.String())
		return
	}
	challengeId, err := strconv.ParseInt(ctx.Param("challenge"), 10, 0)
	if err != nil {
		views.RenderError(ctx, err)
		return
	}
	challenge, err := model.GetChallengeById(challengeId)
	if err != nil {
		views.RenderError(ctx, err)
		return
	}
	user := auth.GetUser(ctx)
	test := &model.Test{
		ChallengeId: challenge.Id,
		UserId:      user.Id,
	}

	if err := model.CreateTest(test, []byte(code)); err != nil {
		views.RenderError(ctx, err)
		return
	}
	checker.Get().CollectorUpdater <- true

	ctx.Redirect(http.StatusFound, fmt.Sprintf("/test/%d", test.Id))
}

func LoginHandler(ctx *gin.Context) {
	auth.Login(ctx, "/")
}

func LogoutHandler(ctx *gin.Context) {
	auth.Logout(ctx, "/")
}
