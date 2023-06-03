package handlers

import (
	"context"
	"strconv"
	"time"

	forms "bot/forms/user"
	"bot/models"
	"bot/utils"
	"bot/web/handlers/response"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userHandler = &user{}

type user struct {
}

func (handler *user) Handle(router *gin.Engine) {
	router.GET("/user", func(ctx *gin.Context) {
		resp := response.New(ctx)
		filter := bson.M{}

		count, err := models.UserCollection.CountDocuments(
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
		cursor, err := models.UserCollection.Find(
			context.TODO(),
			filter, options.Find().SetSkip((page-1)*limit).SetLimit(limit),
		)
		if err != nil {
			resp.Error(err)
			return
		}

		var items []models.User
		if err = cursor.All(context.TODO(), &items); err != nil {
			resp.Error(err)
			return
		}

		resp.HTML("user/index.html", response.Context{
			"page":  page,
			"count": count,
			"limit": limit,
			"items": items,
		})
	})

	router.GET("/user/add", func(ctx *gin.Context) {
		resp := response.New(ctx)
		resp.HTML("user/add.html", response.Context{})
	})

	router.POST("/user/save", func(ctx *gin.Context) {
		resp := response.New(ctx)
		form := forms.Save{}

		if err := ctx.ShouldBind(&form); err != nil {
			resp.Error(err)
			return
		}

		saved := bson.M{
			"username": form.Username,
			"password": utils.Encrypt(form.Password),
			"role":     form.Role,
			"status":   form.Status,
		}
		if _, err := models.UserCollection.InsertOne(
			context.TODO(),
			saved,
			options.InsertOne(),
		); err != nil {
			resp.Error(err)
			return
		}

		resp.Success("User save successful", "/user")
	})

	router.GET("/user/edit/:id", func(ctx *gin.Context) {
		resp := response.New(ctx)

		sId := ctx.Param("id")
		uId, err := primitive.ObjectIDFromHex(sId)
		if err != nil {
			resp.Error(err)
			return
		}

		user := models.User{}
		if err := models.UserCollection.FindOne(context.TODO(), bson.M{
			"_id": uId,
		}).Decode(&user); err != nil {
			resp.Error(err)
			return
		}

		resp.HTML("user/edit.html", response.Context{
			"item": user,
		})
	})

	router.POST("/user/update/:id", func(ctx *gin.Context) {
		resp := response.New(ctx)
		form := forms.Update{}
		sId := ctx.Param("id")
		uId, err := primitive.ObjectIDFromHex(sId)
		if err != nil {
			resp.Error(err)
			return
		}

		if err := ctx.ShouldBind(&form); err != nil {
			resp.Error(err)
			return
		}

		filter := bson.M{
			"_id": uId,
		}
		var update bson.M
		if form.Password == "" {
			update = bson.M{"$set": bson.M{
				"role":   form.Role,
				"status": form.Status,
			}}
		} else {
			update = bson.M{"$set": bson.M{
				"password": utils.Encrypt(form.Password),
				"role":     form.Role,
				"status":   form.Status,
			}}
		}
		err = models.UserCollection.FindOneAndUpdate(
			context.TODO(),
			filter,
			update,
			options.FindOneAndUpdate(),
		).Err()
		if err != nil {
			resp.Error(err)
			return
		}

		resp.Success("User update successful", "/user")
	})
}
