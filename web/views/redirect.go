package views

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RenderRedirect(ctx *gin.Context, path string) {
	ctx.Header("P3P", "CP='INT NAV UNI'")
	ctx.Header("Pragma", "no-cache")
	ctx.Header("Cache-Control", "no-cache")
	ctx.HTML(http.StatusOK, "pages/redirect", gin.H{"Path": path})
}
