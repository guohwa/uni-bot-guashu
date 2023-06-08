package handlers

import (
	"bot/config"
	"bot/exchange"
	"bot/forms"
	"bot/log"
	"bot/models"
	"context"
	"errors"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var tvHandler = &tv{}

type tv struct {
}

func (handler *tv) Handle(router *gin.Engine) {
	router.POST("/tv/:token", func(ctx *gin.Context) {
		if !handler.authenticate(ctx) {
			return
		}

		token := strings.TrimLeft(ctx.Param("token"), "/")
		filter := bson.M{
			"token": token,
		}

		var customer models.Customer
		if err := models.CustomerCollection.FindOne(context.Background(), filter).Decode(&customer); err != nil {
			if err == mongo.ErrNoDocuments {
				log.Error(err)
				ctx.String(http.StatusForbidden, "forbidden")
			} else {
				log.Error(err)
				ctx.String(http.StatusInternalServerError, "internal error")
			}
			return
		}

		var form forms.Command
		err := ctx.ShouldBindJSON(&form)
		if err != nil {
			log.Error(err)
			ctx.String(http.StatusInternalServerError, "internal error")
			return
		}

		status := "NEW"
		reason := ""
		exname, symbol, err := handler.split(form.Symbol)
		if err != nil {
			status = "FAILED"
			reason = err.Error()
		}

		if exname != customer.Exchange {
			status = "FAILED"
			reason = "exchange mismatch"
		}

		command := models.Command{
			ID:         primitive.NewObjectID(),
			CustomerID: customer.ID,
			Exchange:   exname,
			Action:     form.Action,
			Symbol:     symbol,
			Side:       form.Side,
			Size:       form.Size,
			Quantity:   form.Size * customer.Scale,
			Comment:    form.Comment,
			Status:     status,
			Reason:     reason,
			Time:       time.Now().UTC().UnixMilli(),
		}
		if _, err := models.CommandCollection.InsertOne(
			context.TODO(),
			command,
			options.InsertOne(),
		); err != nil {
			log.Error(err)
			ctx.String(http.StatusInternalServerError, "internal error")
			return
		}

		if status == "NEW" {
			go exchange.New(exname, customer, command).Execute()
		}

		ctx.String(http.StatusOK, "ok")
	})
}

func (handler *tv) split(s string) (exchange, symbol string, err error) {
	sections := strings.Split(s, ":")
	if len(sections) < 2 {
		err = errors.New("invalid symbol")
		return
	}
	exchange = sections[0]
	symbol = strings.TrimSuffix(sections[1], ".P")
	return
}

func (handler *tv) authenticate(ctx *gin.Context) bool {
	addr, port, err := func(req *http.Request) (string, string, error) {
		addr, port, err := net.SplitHostPort(req.RemoteAddr)
		if err != nil {
			return addr, port, err
		}
		if addr := req.Header.Get("X-Real-IP"); addr != "" {
			return addr, port, nil
		}
		if addr := req.Header.Get("X-Forwarded-For"); addr != "" {
			return addr, port, nil
		}
		return addr, port, err
	}(ctx.Request)

	if err != nil {
		log.Error(err)
		ctx.String(http.StatusInternalServerError, "internal error")
		return false
	}

	if !strings.Contains(config.App.WhiteList, addr) {
		log.Errorf("unauthorized, remote address: %s, port: %s", addr, port)
		ctx.String(http.StatusUnauthorized, "unauthorized")
		return false
	}

	return true
}
