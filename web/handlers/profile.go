package handlers

import (
	"github.com/gin-gonic/gin"
)

var profileHandler = &profile{}

type profile struct {
	base
}

func (handler *profile) Handle(router *gin.Engine) {
	router.GET("/profile", func(ctx *gin.Context) {
		handler.HTML(ctx, "profile/index.html", nil)
	})
}
