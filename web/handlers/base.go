package handlers

import (
	"context"
	"net/http"
	"regexp"
	"strings"

	"bot/log"
	"bot/models"

	"github.com/flosch/pongo2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var matcher = regexp.MustCompile("^[a-zA-Z0-9_]+$")

type Context map[string]interface{}

type AjaxResponse struct {
	Code int         `json:"code"`
	Msg  interface{} `json:"msg"`
	URL  string      `json:"url"`
	Wait int         `json:"wait"`
}

type base struct {
}

func (handler *base) getUser(ctx *gin.Context) *models.User {
	u, exists := ctx.Get("user")

	if exists {
		if user, ok := u.(models.User); ok {
			return &user
		}
	}

	return nil
}

func (handler *base) getCustomer(ctx *gin.Context) (*models.Customer, []models.Customer) {
	var items []models.Customer

	user := handler.getUser(ctx)
	if user == nil {
		return nil, items
	}

	customerFilter := bson.M{
		"userId": user.ID,
		"status": "Enable",
	}
	customerCursor, err := models.CustomerCollection.Find(
		context.TODO(),
		customerFilter, options.Find(),
	)
	if err != nil {
		handler.Error(ctx, err)
		return nil, items
	}

	if err = customerCursor.All(context.TODO(), &items); err != nil {
		handler.Error(ctx, err)
		return nil, items
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

	if sId != "" {
		for _, item := range items {
			if item.ID.Hex() == sId {
				return &item, items
			}
		}
	}
	if len(items) > 0 {
		return &items[0], items
	}

	return nil, items
}

func (handler *base) with(ctx *gin.Context, data Context) pongo2.Context {
	context := pongo2.Context{}
	for k, v := range ctx.Keys {
		if matcher.MatchString(k) {
			context[k] = v
		}
	}

	for k, v := range data {
		if matcher.MatchString(k) {
			context[k] = v
		}
	}

	return context
}

func (handler *base) HTML(ctx *gin.Context, name string, data Context) {
	ctx.HTML(http.StatusOK, name, handler.with(ctx, data))
}

func (handler *base) String(ctx *gin.Context, data string) {
	ctx.String(http.StatusOK, data)
}

func (handler *base) OK(ctx *gin.Context) {
	ctx.String(http.StatusOK, "OK")
}

func (handler *base) Success(ctx *gin.Context, msg interface{}, url string) {
	ctx.JSON(http.StatusOK, AjaxResponse{
		Msg:  msg,
		URL:  url,
		Wait: 2,
	})
}

func (handler *base) Error(ctx *gin.Context, msg interface{}) {
	if err, ok := msg.(error); ok {
		log.Error(err)
		ctx.JSON(http.StatusOK, AjaxResponse{
			Code: 1,
			Msg:  err.Error(),
			Wait: 2,
		})
		return
	}

	ctx.JSON(http.StatusOK, AjaxResponse{
		Code: 1,
		Msg:  msg,
		Wait: 2,
	})
}
