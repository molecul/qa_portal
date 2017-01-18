package webHandlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/zalando/gin-oauth2/google"
)

func UserInfoHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"Hello": "from private", "user": ctx.MustGet("user").(google.User)})
}

func MainPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Main website",
	})
}