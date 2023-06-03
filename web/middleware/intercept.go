package middleware

import (
	"context"
	"net/http"
	"strings"

	"bot/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Intercept(auth map[string]map[string]bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := func() string {
			paths := strings.Split(ctx.Request.URL.Path, "/")
			if len(paths) > 1 && paths[1] != "" {
				return paths[1]
			}

			return "home"
		}()
		ctx.Set("path", path)

		session := sessions.Default(ctx)
		uId := session.Get("user-id")
		role := "Public"
		if uId != nil {
			if sId, ok := uId.(string); ok {
				id, err := primitive.ObjectIDFromHex(sId)
				if err != nil {
					panic(err)
				}

				user := models.User{}
				if err := models.UserCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user); err != nil {
					panic(err)
				}

				session.Set("user-id", user.ID.Hex())
				session.Save()

				ctx.Set("user", user)
				role = user.Role
			}
		}

		if allows, ok := auth[role]; ok {
			if allow, ok := allows[path]; !ok || !allow {
				ctx.Redirect(http.StatusFound, "/account/login")
				ctx.Abort()
			}
		}
	}
}
