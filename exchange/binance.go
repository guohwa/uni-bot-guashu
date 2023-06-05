package exchange

import (
	"bot/log"
	"bot/models"
)

func NewBinance() Exchange {
	return &Binance{}
}

type Binance struct {
}

func (ex *Binance) Execute(customer models.Customer, command models.Command) {
	log.Printf("exchange: %s, symbol: %s", command.Exchange, command.Symbol)
}
