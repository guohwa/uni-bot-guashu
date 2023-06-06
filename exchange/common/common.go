package common

import "bot/models"

type Exchange interface {
	Execute()
}

type Constructor func(customer models.Customer, command models.Command) Exchange
