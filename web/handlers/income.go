package handlers

import (
	"context"
	"net/http"
	"sort"
	"strings"

	"bot/log"
	"bot/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/uncle-gua/gobinance/futures"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var incomeHandler = &income{}

type income struct {
	base
}

func (handler *income) Handle(router *gin.Engine) {
	router.GET("/income/*id", func(ctx *gin.Context) {
		user, ok := getUser(ctx)
		if !ok {
			ctx.Redirect(http.StatusFound, "/account/login")
			ctx.Abort()
			return
		}

		filter := bson.M{
			"userId": user.ID,
			"status": "Enable",
		}
		cursor, err := models.CustomerCollection.Find(
			context.TODO(),
			filter, options.Find(),
		)
		if err != nil {
			handler.Error(ctx, err)
			return
		}

		var items []models.Customer
		if err = cursor.All(context.TODO(), &items); err != nil {
			handler.Error(ctx, err)
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
			handler.HTML(ctx, "income/index.html", Context{
				"items":    items,
				"customer": customer,
				"incomes":  nil,
			})
			return
		}

		client := futures.NewClient(customer.ApiKey, customer.ApiSecret)
		incomes, err := client.NewGetIncomeHistoryService().Do(context.Background())
		if err != nil {
			handler.Error(ctx, err)
			return
		}
		sort.SliceStable(incomes, func(i, j int) bool {
			return incomes[i].Time > incomes[j].Time
		})

		handler.HTML(ctx, "income/index.html", Context{
			"items":    items,
			"customer": customer,
			"incomes":  incomes,
		})
	})

}
