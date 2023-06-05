package handlers

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"bot/log"
	"bot/models"
	"bot/web/handlers/response"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var commandHandler = &command{}

type command struct {
}

func (handler *command) Handle(router *gin.Engine) {
	router.GET("/command/*id", func(ctx *gin.Context) {
		user, ok := getUser(ctx)
		if !ok {
			ctx.Redirect(http.StatusFound, "/account/login")
			ctx.Abort()
			return
		}

		resp := response.New(ctx)

		customerFilter := bson.M{
			"userId": user.ID,
			"status": "Enable",
		}
		customerCursor, err := models.CustomerCollection.Find(
			context.TODO(),
			customerFilter, options.Find(),
		)
		if err != nil {
			resp.Error(err)
			return
		}

		var items []models.Customer
		if err = customerCursor.All(context.TODO(), &items); err != nil {
			resp.Error(err)
			return
		}

		sId := strings.TrimLeft(ctx.Param("id"), "/")
		session := sessions.Default(ctx)
		if sId == "" {
			cId := session.Get("customer-id")
			if cId != nil {
				if id, ok := cId.(string); ok {
					sId = id
				}
			}
		} else {
			session.Set("customer-id", sId)
			if err := session.Save(); err != nil {
				log.Error(err)
			}
		}

		var customer models.Customer
		if sId != "" {
			for _, item := range items {
				if item.ID.Hex() == sId {
					customer = item
				}
			}
		} else {
			if len(items) > 0 {
				customer = items[0]
			}
		}

		filter := bson.M{
			"customerId": customer.ID,
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

		var commands []models.Command
		if err = cursor.All(context.TODO(), &commands); err != nil {
			resp.Error(err)
			return
		}

		resp.HTML("command/index.html", response.Context{
			"items":    items,
			"page":     page,
			"count":    count,
			"limit":    limit,
			"commands": commands,
			"customer": customer,
		})
	})
}
