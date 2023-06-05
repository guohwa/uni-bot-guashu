package handlers

import (
	"bot/forms"
	"bot/log"
	"bot/models"
	"context"
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
