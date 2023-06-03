package handlers

import (
	"context"

	"bot/forms"
	"bot/models"
	"bot/utils"
	"bot/web/handlers/response"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var accountHandler = &account{}

type account struct {
}

func (handler *account) Handle(router *gin.Engine) {
	router.GET("/account/login", func(ctx *gin.Context) {
		response.New(ctx).HTML("account/login.html", nil)
	})

	router.POST("/account/login", func(ctx *gin.Context) {
		resp := response.New(ctx)

		form := forms.Login{}
		if err := ctx.ShouldBind(&form); err != nil {
			resp.Error(err)
			return
		}

		if !captchaHandler.Verify(ctx, form.Verify) {
			resp.Error("Invalid verify code")
			return
		}

		user := models.User{}
		if err := models.UserCollection.FindOne(context.Background(), bson.M{"username": form.Username}).Decode(&user); err != nil {
			if err == mongo.ErrNoDocuments {
				resp.Error("User Does not exists.")
			} else {
				resp.Error(err)
			}
			return
		}

		if user.Password != utils.Encrypt(form.Password) {
			resp.Error("Invalid password")
			return
		}

		if user.Status == "Disable" {
			resp.Error("User disabled")
			return
		}

		session := sessions.Default(ctx)
		session.Set("user-id", user.ID.Hex())

		if err := session.Save(); err == nil {
			resp.Success("Login Successful", "/")
		} else {
			resp.Error("Login Failed")
		}
	})

	router.GET("/account/logout", func(ctx *gin.Context) {
		resp := response.New(ctx)

		session := sessions.Default(ctx)
		session.Clear()

		if err := session.Save(); err == nil {
			resp.Success("Logout successful", "/")
		} else {
			resp.Error("Logout failed")
		}
	})

}
