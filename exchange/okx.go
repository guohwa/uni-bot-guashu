package exchange

import (
	"bot/log"
	"bot/models"
)

func NewOkx() Exchange {
	return &Okx{}
}

type Okx struct {
}

func (ex *Okx) Execute(customer models.Customer, command models.Command) {
	log.Printf("exchange: %s, symbol: %s", command.Exchange, command.Symbol)
}
