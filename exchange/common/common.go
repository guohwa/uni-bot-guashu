package common

import "bot/models"

type Exchange interface {
	Execute(command models.Command)
	FormatSize(symbol string, size float64) string
}

type Constructor func(customer models.Customer) Exchange
