package handlers

import (
	"context"
	"net/http"

	"bot/forms"
	"bot/models"
	"bot/utils"
	"bot/web/handlers/response"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var passwordHandler = &password{}

type password struct {
}

func (handler *password) Handle(router *gin.Engine) {
	router.GET("/password", func(ctx *gin.Context) {
		resp := response.New(ctx)
		resp.HTML("password/index.html", response.Context{})
	})

	router.POST("/password", func(ctx *gin.Context) {
		user, ok := getUser(ctx)
		if !ok {
			ctx.Redirect(http.StatusFound, "/account/login")
			ctx.Abort()
			return
		}

		resp := response.New(ctx)
		form := forms.Password{}
		if user.Role == "Demo" {
			resp.Error("Demo user can not change passowrd")
			return
		}

		if err := ctx.ShouldBind(&form); err != nil {
			resp.Error(err)
			return
		}

		if user.Password != utils.Encrypt(form.Password) {
			resp.Error("Invalid password")
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
			resp.Error(err)
			return
		}

		resp.Success("password update successful", "")
	})
}
