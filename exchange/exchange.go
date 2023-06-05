package exchange

import (
	"bot/models"
)

type Exchange interface {
	Execute(customer models.Customer, command models.Command)
}

type Constructor func() Exchange

var exchanges map[string]Constructor = map[string]Constructor{
	"BINANCE": NewBinance,
	"OKX":     NewOkx,
}

func New(name string) Exchange {
	if constructor, ok := exchanges[name]; ok {
		return constructor()
	}

	return nil
}
