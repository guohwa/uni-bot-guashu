package handlers

import (
	"context"
	"net/http"
	"strings"

	"bot/exchange"
	orderform "bot/forms/order"
	"bot/log"
	"bot/models"
	"bot/web/handlers/response"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/uncle-gua/gobinance/futures"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var positionHandler = &position{}

type position struct {
}

func (handler *position) Handle(router *gin.Engine) {
	router.GET("/position/*id", func(ctx *gin.Context) {
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
				"account":  new(account),
				"orders":   make([]*futures.Order, 0),
			})
			return
		}

		client := futures.NewClient(customer.ApiKey, customer.ApiSecret)
		account, err := client.NewGetAccountService().Do(context.Background())
		if err != nil {
			resp.Error(err)
			return
		}
		orders, err := client.NewListOpenOrdersService().Do(context.Background())
		if err != nil {
			resp.Error(err)
			return
		}

		resp.HTML("position/index.html", response.Context{
			"items":    items,
			"customer": customer,
			"account":  account,
			"orders":   orders,
		})
	})

	router.POST("/position/close", func(ctx *gin.Context) {
		user, ok := getUser(ctx)
		if !ok {
			ctx.Redirect(http.StatusFound, "/account/login")
			ctx.Abort()
			return
		}

		resp := response.New(ctx)
		form := orderform.Close{}
		if user.Role == "Demo" {
			resp.Error("Demo user can not close position")
			return
		}

		if err := ctx.ShouldBind(&form); err != nil {
			resp.Error(err)
			return
		}

		id, err := primitive.ObjectIDFromHex(form.CustomerID)
		if err != nil {
			resp.Error(err)
			return
		}

		customer := models.Customer{}
		if err := models.CustomerCollection.FindOne(context.TODO(), bson.M{
			"_id":    id,
			"userId": user.ID,
		}).Decode(&customer); err != nil {
			resp.Error(err)
			return
		}

		side := futures.SideTypeBuy
		if form.PositionSide == "LONG" {
			side = futures.SideTypeSell
		}
		positionSide := futures.PositionSideType(form.PositionSide)
		quantity := func(s string) string {
			if s[0:1] == "-" {
				return s[1:]
			}
			return s
		}(form.PositionAmt)

		client := futures.NewClient(customer.ApiKey, customer.ApiSecret)
		_, err = client.NewCreateOrderService().
			Symbol(form.Symbol).
			Type(futures.OrderTypeMarket).
			Side(side).
			PositionSide(positionSide).
			Quantity(quantity).
			Do(context.Background())
		if err != nil {
			resp.Error(err)
			return
		}

		resp.Success("Position close successful", "")
	})

	router.POST("/position/cancel", func(ctx *gin.Context) {
		user, ok := getUser(ctx)
		if !ok {
			ctx.Redirect(http.StatusFound, "/account/login")
			ctx.Abort()
			return
		}

		resp := response.New(ctx)
		form := orderform.Cancel{}

		if err := ctx.ShouldBind(&form); err != nil {
			resp.Error(err)
			return
		}

		id, err := primitive.ObjectIDFromHex(form.CustomerID)
		if err != nil {
			resp.Error(err)
			return
		}

		customer := models.Customer{}
		if err := models.CustomerCollection.FindOne(context.TODO(), bson.M{
			"_id":    id,
			"userId": user.ID,
		}).Decode(&customer); err != nil {
			resp.Error(err)
			return
		}

		client := futures.NewClient(customer.ApiKey, customer.ApiSecret)
		_, err = client.NewCancelOrderService().Symbol(form.Symbol).OrderID(form.OrderID).Do(context.Background())
		if err != nil {
			resp.Error(err)
			return
		}

		resp.Success("Order cancel successful", "")
	})

	router.POST("/position/create", func(ctx *gin.Context) {
		user, ok := getUser(ctx)
		if !ok {
			ctx.Redirect(http.StatusFound, "/account/login")
			ctx.Abort()
			return
		}

		resp := response.New(ctx)
		form := orderform.Create{}

		if err := ctx.ShouldBind(&form); err != nil {
			resp.Error(err)
			return
		}

		id, err := primitive.ObjectIDFromHex(form.CustomerID)
		if err != nil {
			resp.Error(err)
			return
		}

		customer := models.Customer{}
		if err := models.CustomerCollection.FindOne(context.TODO(), bson.M{
			"_id":    id,
			"userId": user.ID,
		}).Decode(&customer); err != nil {
			resp.Error(err)
			return
		}

		exchange := exchange.New(customer)
		if exchange == nil {
			resp.Error("exchange mismatch")
			return
		}
		size := exchange.FormatSize(form.Symbol, form.Size)

		var side futures.SideType
		var positionSide futures.PositionSideType
		switch strings.ToUpper(form.Side) {
		case "LONG":
			positionSide = futures.PositionSideTypeLong
			switch strings.ToUpper(form.Action) {
			case "OPEN":
				side = futures.SideTypeBuy
			case "CLOSE":
				side = futures.SideTypeSell
			}
		case "SHORT":
			positionSide = futures.PositionSideTypeShort
			switch strings.ToUpper(form.Action) {
			case "OPEN":
				side = futures.SideTypeSell
			case "CLOSE":
				side = futures.SideTypeBuy
			}
		}

		client := futures.NewClient(customer.ApiKey, customer.ApiSecret)
		if _, err = client.NewCreateOrderService().
			Symbol(form.Symbol).
			Type(futures.OrderTypeMarket).
			Side(side).
			PositionSide(positionSide).
			Quantity(size).
			Do(context.Background()); err != nil {
			resp.Error(err)
			return
		}

		resp.Success("Order create successful", "")
	})
}
