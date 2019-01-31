package main

import (
	"github.com/gin-gonic/gin"
	"github.com/reveliant/dirtyci/utils"
)

func Handler(ginCtx *gin.Context) *string {
	var code int
	var err error
	var data = new(PushEvent)
	var ctx = utils.NewContext(ginCtx)

	code = ctx.RequireHeader("X-Gitlab-Event", "Push Hook")
	if code != 0 {
		ctx.AbortWithStatus(code)
		return nil
	}

	code, err = ctx.BindPushEvent(data)
	if code != 0 {
		ctx.AbortWithError(code, err)
		return nil
	}

	return &data.Project.WebUrl
}
