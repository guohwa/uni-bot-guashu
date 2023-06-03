package handlers

import (
	"bot/web/handlers/response"

	"github.com/gin-gonic/gin"
)

var profileHandler = &profile{}

type profile struct {
}

func (handler *profile) Handle(router *gin.Engine) {
	router.GET("/profile", func(ctx *gin.Context) {
		resp := response.New(ctx)
		resp.HTML("profile/index.html", response.Context{})
	})
}
