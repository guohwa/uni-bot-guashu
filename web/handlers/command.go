package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"bot/models"
	"bot/web/handlers/response"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var commandHandler = &command{}

type command struct {
}

func (handler *command) Handle(router *gin.Engine) {
	router.GET("/command", func(ctx *gin.Context) {
		user, ok := getUser(ctx)
		if !ok {
			ctx.Redirect(http.StatusFound, "/account/login")
			ctx.Abort()
			return
		}

		resp := response.New(ctx)
		filter := bson.M{
			"userId": user.ID,
		}

		count, err := models.CommandCollection.CountDocuments(
			context.TODO(),
			filter,
			options.Count().SetMaxTime(2*time.Second))
		if err != nil {
			resp.Error(err)
			return
		}

		page, err := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
		if err != nil {
			resp.Error(err)
			return
		}
		limit, err := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)
		if err != nil {
			resp.Error(err)
			return
		}
		cursor, err := models.CommandCollection.Find(
			context.TODO(),
			filter, options.Find().SetSkip((page-1)*limit).SetLimit(limit),
		)
		if err != nil {
			resp.Error(err)
			return
		}

		var items []models.Command
		if err = cursor.All(context.TODO(), &items); err != nil {
			resp.Error(err)
			return
		}

		resp.HTML("command/index.html", response.Context{
			"page":  page,
			"count": count,
			"limit": limit,
			"items": items,
		})
	})
}
