package okx

import (
	"bot/exchange/common"
	"bot/log"
	"bot/models"
)

func New(customer models.Customer, command models.Command) common.Exchange {
	return &Okx{
		customer,
		command,
	}
}

type Okx struct {
	customer models.Customer
	command  models.Command
}

func (ex *Okx) Execute() {
	switch ex.command.Action {
	case "OPEN":
		ex.Open()
	case "CLOSE":
		ex.Close()
	case "INCR":
		ex.Incr()
	case "DECR":
		ex.Decr()
	default:
		log.Printf("exchange: %s, symbol: %s", ex.command.Exchange, ex.command.Symbol)
	}
}

func (ex *Okx) Open() {
	log.Printf("action: %s, exchange: %s, symbol: %s", ex.command.Action, ex.command.Exchange, ex.command.Symbol)
}

func (ex *Okx) Close() {
	log.Printf("action: %s, exchange: %s, symbol: %s", ex.command.Action, ex.command.Exchange, ex.command.Symbol)
}

func (ex *Okx) Incr() {
	log.Printf("action: %s, exchange: %s, symbol: %s", ex.command.Action, ex.command.Exchange, ex.command.Symbol)
}

func (ex *Okx) Decr() {
	log.Printf("action: %s, exchange: %s, symbol: %s", ex.command.Action, ex.command.Exchange, ex.command.Symbol)
}
