package server

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func SetMode(mode string) {
	gin.SetMode(mode)
}

func Redirect(target string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, target)
	}
}
