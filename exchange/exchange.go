package exchange

import (
	"bot/exchange/binance"
	"bot/exchange/common"
	"bot/exchange/okx"
	"bot/models"
)

var exchanges = map[string]common.Constructor{
	"BINANCE": binance.New,
	"OKX":     okx.New,
}

func New(name string, customer models.Customer, command models.Command) common.Exchange {
	if constructor, ok := exchanges[name]; ok {
		return constructor(customer, command)
	}

	return nil
}

func Support(name string) bool {
	return exchanges[name] != nil
}