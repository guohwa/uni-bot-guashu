package handlers

import (
	"bot/config"
	"bot/forms"
	"bot/log"
	"bot/models"
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var tvHandler = &tv{}

type tv struct {
}

func (handler *tv) Handle(router *gin.Engine) {
	router.POST("/tv/:token", func(ctx *gin.Context) {
		if !handler.authenticate(ctx) {
			return
		}

		token := strings.TrimLeft(ctx.Param("id"), "/")
		filter := bson.M{
			"token": token,
		}

		var customer models.Customer
		if err := models.CustomerCollection.FindOne(context.Background(), filter).Decode(&customer); err != nil {
			if err == mongo.ErrNoDocuments {
				log.Error(err)
				ctx.String(http.StatusForbidden, "forbidden")
			} else {
				log.Error(err)
				ctx.String(http.StatusInternalServerError, "internal error")
			}
			return
		}

		var form forms.Command
		err := ctx.ShouldBindJSON(&form)
		if err != nil {
			log.Error(err)
			ctx.String(http.StatusInternalServerError, "internal error")
			return
		}

	})
}

func (handler *tv) authenticate(ctx *gin.Context) bool {
	remoteAddr := ctx.Request.RemoteAddr
	if ip := ctx.Request.Header.Get("X-Real-IP"); ip != "" {
		remoteAddr = ip
	} else if ip := ctx.Request.Header.Get("X-Forwarded-For"); ip != "" {
		remoteAddr = ip
	} else {
		ip, _, err := net.SplitHostPort(remoteAddr)
		if err != nil {
			log.Error(err)
			ctx.String(http.StatusInternalServerError, "internal error")
			return false
		}
		remoteAddr = ip
	}

	if !strings.Contains(config.App.WhiteList, remoteAddr) {
		log.Errorf("illegal access, remote address: %s", remoteAddr)
		ctx.String(http.StatusInternalServerError, "internal error")
		return false
	}

	return true
}
