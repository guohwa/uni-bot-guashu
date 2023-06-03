package handlers

import (
	"bot/web/handlers/response"

	"github.com/gin-gonic/gin"
)

var homeHandler = &home{}

type home struct {
}

func (handler *home) Handle(router *gin.Engine) {
	router.GET("/", func(ctx *gin.Context) {
		ctx.Request.URL.Path = "/home"
		router.HandleContext(ctx)
	})

	router.GET("/home", func(ctx *gin.Context) {
		resp := response.New(ctx)
		resp.HTML("home/index.html", response.Context{})
	})
}
