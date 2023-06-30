package handlers

import (
	"context"
	"strconv"
	"time"

	forms "bot/forms/user"
	"bot/models"
	"bot/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userHandler = &user{}

type user struct {
	base
}

func (handler *user) Handle(router *gin.Engine) {
	router.GET("/user", func(ctx *gin.Context) {
		filter := bson.M{}

		count, err := models.UserCollection.CountDocuments(
			context.TODO(),
			filter,
			options.Count().SetMaxTime(2*time.Second))
		if err != nil {
			handler.Error(ctx, err)
			return
		}

		page, err := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
		if err != nil {
			handler.Error(ctx, err)
			return
		}
		limit, err := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)
		if err != nil {
			handler.Error(ctx, err)
			return
		}
		cursor, err := models.UserCollection.Find(
			context.TODO(),
			filter, options.Find().SetSkip((page-1)*limit).SetLimit(limit),
		)
		if err != nil {
			handler.Error(ctx, err)
			return
		}

		var items []models.User
		if err = cursor.All(context.TODO(), &items); err != nil {
			handler.Error(ctx, err)
			return
		}

		handler.HTML(ctx, "user/index.html", Context{
			"page":  page,
			"count": count,
			"limit": limit,
			"items": items,
		})
	})

	router.GET("/user/add", func(ctx *gin.Context) {
		handler.HTML(ctx, "user/add.html", Context{})
	})

	router.POST("/user/save", func(ctx *gin.Context) {
		form := forms.Save{}

		if err := ctx.ShouldBind(&form); err != nil {
			handler.Error(ctx, err)
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
			handler.Error(ctx, err)
			return
		}

		handler.Success(ctx, "User save successful", "/user")
	})

	router.GET("/user/edit/:id", func(ctx *gin.Context) {
		sId := ctx.Param("id")
		uId, err := primitive.ObjectIDFromHex(sId)
		if err != nil {
			handler.Error(ctx, err)
			return
		}

		user := models.User{}
		if err := models.UserCollection.FindOne(context.TODO(), bson.M{
			"_id": uId,
		}).Decode(&user); err != nil {
			handler.Error(ctx, err)
			return
		}

		handler.HTML(ctx, "user/edit.html", Context{
			"item": user,
		})
	})

	router.POST("/user/update/:id", func(ctx *gin.Context) {
		form := forms.Update{}
		sId := ctx.Param("id")
		uId, err := primitive.ObjectIDFromHex(sId)
		if err != nil {
			handler.Error(ctx, err)
			return
		}

		if err := ctx.ShouldBind(&form); err != nil {
			handler.Error(ctx, err)
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
			handler.Error(ctx, err)
			return
		}

		handler.Success(ctx, "User update successful", "/user")
	})
}
