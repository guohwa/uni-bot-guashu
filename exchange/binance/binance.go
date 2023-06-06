package binance

import (
	"bot/exchange/common"
	"bot/log"
	"bot/models"
	"bot/utils"
	"context"
	"errors"

	"github.com/uncle-gua/gobinance/futures"
)

func New(customer models.Customer, command models.Command) common.Exchange {
	return &Binance{
		customer,
		command,
	}
}

type Binance struct {
	customer models.Customer
	command  models.Command
}

type Action func(customer models.Customer, command models.Command) error

var actions = map[string]Action{
	"OPEN":  open,
	"CLOSE": close,
	"INCR":  incr,
	"DECR":  decr,
}

func (ex *Binance) Execute() {
	if action, ok := actions[ex.command.Action]; ok {
		if err := action(ex.customer, ex.command); err != nil {
			log.Error(err)
		}
	} else {
		log.Errorf("unsupported action: %s", ex.command.Action)
	}
}

func open(customer models.Customer, command models.Command) error {
	client := futures.NewClient(customer.ApiKey, customer.ApiSecret)
	account, err := client.NewGetAccountService().Do(context.Background())
	if err != nil {
		return err
	}

	amount1 := amount(account.Positions, command.Symbol, futures.PositionSideType(command.Side))
	if !utils.IsZero(amount1) {
		return errors.New("position hold")
	}

	opposite := func(side string) futures.PositionSideType {
		if side == "LONG" {
			return futures.PositionSideTypeShort
		}

		return futures.PositionSideTypeLong
	}(command.Side)
	amount2 := amount(account.Positions, command.Symbol, opposite)
	if !utils.IsZero(amount2) {
		_, err := client.NewCreateOrderService().
			Symbol(command.Symbol).
			Type(futures.OrderTypeMarket).
			Side("SELL").
			PositionSide(opposite).
			Quantity(utils.Abs(amount2)).
			Do(context.Background())
		if err != nil {
			return err
		}
	}

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

func amount(positions []*futures.AccountPosition, symbol string, side futures.PositionSideType) string {
	for _, p := range positions {
		if p.Symbol == symbol && p.PositionSide == side {
			return p.PositionAmt
		}
	}

	return "0"
}
