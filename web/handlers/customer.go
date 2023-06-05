package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	forms "bot/forms/customer"
	"bot/models"
	"bot/web/handlers/response"

	"github.com/gin-gonic/gin"
	"github.com/jaevor/go-nanoid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var customerHandler = &customer{}

type customer struct {
}

func (handler *customer) Handle(router *gin.Engine) {
	router.GET("/customer", func(ctx *gin.Context) {
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

		count, err := models.CustomerCollection.CountDocuments(
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
		cursor, err := models.CustomerCollection.Find(
			context.TODO(),
			filter, options.Find().SetSkip((page-1)*limit).SetLimit(limit),
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

		resp.HTML("customer/index.html", response.Context{
			"page":  page,
			"count": count,
			"limit": limit,
			"items": items,
		})
	})

	router.GET("/customer/add", func(ctx *gin.Context) {
		resp := response.New(ctx)
		resp.HTML("customer/add.html", response.Context{})
	})

	router.POST("/customer/save", func(ctx *gin.Context) {
		user, ok := getUser(ctx)
		if !ok {
			ctx.Redirect(http.StatusFound, "/account/login")
			ctx.Abort()
			return
		}

		resp := response.New(ctx)
		form := forms.Save{}
		if user.Role == "Demo" {
			resp.Error("Demo user can not add customer")
			return
		}

		if err := ctx.ShouldBind(&form); err != nil {
			resp.Error(err)
			return
		}

		token, err := nanoid.CustomASCII("0123456789abcdefghijklmnopqrstuvwxyz", 10)
		if err != nil {
			resp.Error(err)
			return
		}

		saved := bson.M{
			"userId":    user.ID,
			"name":      form.Name,
			"token":     token(),
			"apiKey":    form.ApiKey,
			"apiSecret": form.ApiSecret,
			"capital":   form.Capital,
			"scale":     form.Scale,
			"level1":    form.Level1,
			"level2":    form.Level2,
			"status":    form.Status,
		}
		if _, err := models.CustomerCollection.InsertOne(
			context.TODO(),
			saved,
			options.InsertOne(),
		); err != nil {
			resp.Error(err)
			return
		}

		resp.Success("Customer save successful", "/customer")
	})

	router.GET("/customer/edit/:id", func(ctx *gin.Context) {
		resp := response.New(ctx)

		sId := ctx.Param("id")
		uId, err := primitive.ObjectIDFromHex(sId)
		if err != nil {
			resp.Error(err)
			return
		}

		customer := models.Customer{}
		if err := models.CustomerCollection.FindOne(context.TODO(), bson.M{
			"_id": uId,
		}).Decode(&customer); err != nil {
			resp.Error(err)
			return
		}

		resp.HTML("customer/edit.html", response.Context{
			"item": customer,
		})
	})

	router.POST("/customer/update/:id", func(ctx *gin.Context) {
		user, ok := getUser(ctx)
		if !ok {
			ctx.Redirect(http.StatusFound, "/account/login")
			ctx.Abort()
			return
		}

		resp := response.New(ctx)
		form := forms.Update{}
		if user.Role == "Demo" {
			resp.Error("Demo user can not edit customer")
			return
		}

		sId := ctx.Param("id")
		cId, err := primitive.ObjectIDFromHex(sId)
		if err != nil {
			resp.Error(err)
			return
		}

		if err := ctx.ShouldBind(&form); err != nil {
			resp.Error(err)
			return
		}

		filter := bson.M{
			"_id":    cId,
			"userId": user.ID,
		}

		update := bson.M{"$set": bson.M{
			"name":      form.Name,
			"apiKey":    form.ApiKey,
			"apiSecret": form.ApiSecret,
			"capital":   form.Capital,
			"scale":     form.Scale,
			"level1":    form.Level1,
			"level2":    form.Level2,
			"status":    form.Status,
		}}
		err = models.CustomerCollection.FindOneAndUpdate(
			context.TODO(),
			filter,
			update,
			options.FindOneAndUpdate(),
		).Err()
		if err != nil {
			resp.Error(err)
			return
		}

		resp.Success("Customer update successful", "/customer")
	})
}
