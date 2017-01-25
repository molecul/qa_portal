package views

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RenderError(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusInternalServerError, gin.H{"Internal error": err.Error()})
	panic(err)
}
