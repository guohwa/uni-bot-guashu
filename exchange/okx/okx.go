package okx

import (
	"bot/exchange/common"
	"bot/log"
	"bot/models"
	"context"
	"errors"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrApi   = errors.New("API Error")
	ErrEmpty = errors.New("position empty")
	ErrHold  = errors.New("position hold")
	ErrRisk  = errors.New("position risk")
)

type Action func(customer models.Customer, command models.Command) error

func New(customer models.Customer) common.Exchange {
	return &Okx{
		customer,
	}
}

type Okx struct {
	customer models.Customer
}

var actions = map[string]Action{
	"OPEN":  open,
	"CLOSE": close,
	"INCR":  incr,
	"DECR":  decr,
}

func (ex *Okx) Execute(command models.Command) {
	if action, ok := actions[command.Action]; ok {
		err := action(ex.customer, command)

		filter := bson.M{"_id": command.ID}

		var update bson.M
		if err != nil {
			if err == ErrEmpty || err == ErrHold || err == ErrRisk {
				update = bson.M{"$set": bson.M{
					"status": "FAILED",
					"reason": err.Error(),
				}}
			} else {
				log.Error(err)
				update = bson.M{"$set": bson.M{
					"status": "FAILED",
					"reason": "Api Error",
				}}
			}
		} else {
			update = bson.M{"$set": bson.M{
				"status": "SUCCESS",
			}}
		}

		if err := models.CommandCollection.FindOneAndUpdate(
			context.TODO(),
			filter,
			update,
			options.FindOneAndUpdate(),
		).Err(); err != nil {
			log.Error(err)
		}
	} else {
		log.Errorf("unsupported action: %s", command.Action)
	}
}

func (ex *Okx) FormatSize(symbol string, size float64) string {
	return strconv.FormatFloat(size, 'f', 2, 64)
}

func open(customer models.Customer, command models.Command) error {
	log.Printf("action: %s, exchange: %s, symbol: %s", command.Action, command.Exchange, command.Symbol)
	return nil
}

func close(customer models.Customer, command models.Command) error {
	log.Printf("action: %s, exchange: %s, symbol: %s", command.Action, command.Exchange, command.Symbol)
	return nil
}

func incr(customer models.Customer, command models.Command) error {
	log.Printf("action: %s, exchange: %s, symbol: %s", command.Action, command.Exchange, command.Symbol)
	return nil
}

func decr(customer models.Customer, command models.Command) error {
	log.Printf("action: %s, exchange: %s, symbol: %s", command.Action, command.Exchange, command.Symbol)
	return nil
}
