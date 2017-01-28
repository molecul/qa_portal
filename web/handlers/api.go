package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/molecul/qa_portal/checker"
	"github.com/molecul/qa_portal/model"
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

func DockerHealthCheckHandler(c *gin.Context) {
	status := "OK"
	if err := checker.Get().PingDocker(); err != nil {
		status = err.Error()
	}
	c.JSON(http.StatusOK, gin.H{"status": status})
}
