package handlers

import (
	"context"
	"net/http"
	"sort"
	"strings"

	"bot/log"
	"bot/models"
	"bot/web/handlers/response"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/uncle-gua/gobinance/futures"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var tradeHandler = &trade{}

type trade struct {
}

func (handler *trade) Handle(router *gin.Engine) {
	router.GET("/trade/*id", func(ctx *gin.Context) {
		user, ok := getUser(ctx)
		if !ok {
			ctx.Redirect(http.StatusFound, "/account/login")
			ctx.Abort()
			return
		}

		resp := response.New(ctx)

		filter := bson.M{
			"userId": user.ID,
			"status": "Enable",
		}
		cursor, err := models.CustomerCollection.Find(
			context.TODO(),
			filter, options.Find(),
		)
		if err != nil {
			resp.Error(err)
			return
		}

		var items []models.Customer
		if err = cursor.All(context.TODO(), &items); err != nil {
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

		if customer.ApiKey == "" || customer.ApiSecret == "" {
			resp.HTML("income/index.html", response.Context{
				"items":    items,
				"customer": customer,
				"trades":   make([]*futures.AccountTrade, 0),
			})
			return
		}

		client := futures.NewClient(customer.ApiKey, customer.ApiSecret)
		trades, err := client.NewListAccountTradeService().Do(context.Background())
		if err != nil {
			resp.Error(err)
			return
		}
		sort.SliceStable(trades, func(i, j int) bool {
			return trades[i].Time > trades[j].Time
		})

		resp.HTML("trade/index.html", response.Context{
			"items":    items,
			"customer": customer,
			"trades":   trades,
		})
	})

}
