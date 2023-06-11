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

func New(customer models.Customer) common.Exchange {
	if constructor, ok := exchanges[customer.Exchange]; ok {
		return constructor(customer)
	}

	return nil
}
