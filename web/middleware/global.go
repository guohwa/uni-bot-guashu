package middleware

import (
	"bot/config"
	"time"

	"github.com/gin-gonic/gin"
)

func Global() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("title", config.Config.App.Title)
		ctx.Set("now", time.Now())
	}
}
