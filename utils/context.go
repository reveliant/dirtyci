package utils

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type Context struct {
	*gin.Context
}

func NewContext(ginCtx *gin.Context) *Context {
	return &Context{ginCtx}
}

func (ctx *Context) RequireHeader(name, value string) int {
	var content = ctx.GetHeader(name)

	switch content {
	// Header is required to the parent handler but not present
	case "": return http.StatusBadRequest

	// Is a push event
	case value: return 0

	// Not a push event, not handled
	default: return http.StatusNotImplemented
	}
}

func (ctx *Context) BindPushEvent(data interface{}) (int, error) {
	var err = ctx.Bind(data)
	if err != nil {
		return http.StatusInternalServerError, err
	} else {
		return 0, nil
	}
}
