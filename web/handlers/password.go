package handlers

import (
	"context"
	"net/http"

	"bot/forms"
	"bot/models"
	"bot/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var passwordHandler = &password{}

type password struct {
	base
}

func (handler *password) Handle(router *gin.Engine) {
	router.GET("/password", func(ctx *gin.Context) {
		handler.HTML(ctx, "password/index.html", Context{})
	})

	router.POST("/password", func(ctx *gin.Context) {
		user, ok := getUser(ctx)
		if !ok {
			ctx.Redirect(http.StatusFound, "/account/login")
			ctx.Abort()
			return
		}

		form := forms.Password{}
		if user.Role == "Demo" {
			handler.Error(ctx, "Demo user can not change passowrd")
			return
		}

		if err := ctx.ShouldBind(&form); err != nil {
			handler.Error(ctx, err)
			return
		}

		if user.Password != utils.Encrypt(form.Password) {
			handler.Error(ctx, "Invalid password")
			return
		}

		filter := bson.M{"_id": user.ID}
		update := bson.M{"$set": bson.M{
			"password": utils.Encrypt(form.NewPassword),
		}}
		if err := models.UserCollection.FindOneAndUpdate(
			context.TODO(),
			filter,
			update,
			options.FindOneAndUpdate(),
		).Err(); err != nil {
			handler.Error(ctx, err)
			return
		}

		handler.Success(ctx, "password update successful", "")
	})
}
