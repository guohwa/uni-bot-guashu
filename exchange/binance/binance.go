package binance

import (
	"bot/exchange/common"
	"bot/log"
	"bot/models"
	"context"
	"errors"
	"regexp"
	"strconv"

	"github.com/uncle-gua/gobinance/futures"
)

var zero = regexp.MustCompile(`^[\+\-]?0\.?0*$`)

type Action func(customer models.Customer, command models.Command) error

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
	if !isZero(amount1) {
		return errors.New("position hold")
	}

	oppositeSide := opposite(command.Side)

	amount2 := amount(account.Positions, command.Symbol, oppositeSide)
	if !isZero(amount2) {
		side := func(ps futures.PositionSideType) futures.SideType {
			if ps == futures.PositionSideTypeLong {
				return futures.SideTypeSell
			}
			return futures.SideTypeBuy
		}(oppositeSide)
		if _, err := client.NewCreateOrderService().
			Symbol(command.Symbol).
			Type(futures.OrderTypeMarket).
			Side(side).
			PositionSide(oppositeSide).
			Quantity(abs(amount2)).
			Do(context.Background()); err != nil {
			return err
		}
	}

	side := func(ps futures.PositionSideType) futures.SideType {
		if ps == futures.PositionSideTypeShort {
			return futures.SideTypeSell
		}
		return futures.SideTypeBuy
	}(command.Side)

	if _, err := client.NewCreateOrderService().
		Symbol(command.Symbol).
		Type(futures.OrderTypeMarket).
		Side(side).
		PositionSide(command.Side).
		Quantity(format(command.Symbol, command.Quantity)).
		Do(context.Background()); err != nil {
		return err
	}

	return nil
}

func close(customer models.Customer, command models.Command) error {
	client := futures.NewClient(customer.ApiKey, customer.ApiSecret)
	account, err := client.NewGetAccountService().Do(context.Background())
	if err != nil {
		return err
	}

	amount := amount(account.Positions, command.Symbol, futures.PositionSideType(command.Side))
	if isZero(amount) {
		return errors.New("position empty")
	}

	side := func(ps futures.PositionSideType) futures.SideType {
		if ps == futures.PositionSideTypeShort {
			return futures.SideTypeSell
		}
		return futures.SideTypeBuy
	}(command.Side)
	if _, err := client.NewCreateOrderService().
		Symbol(command.Symbol).
		Type(futures.OrderTypeMarket).
		Side(side).
		PositionSide(command.Side).
		Quantity(abs(amount)).
		Do(context.Background()); err != nil {
		return err
	}

	return nil
}

func incr(customer models.Customer, command models.Command) error {
	client := futures.NewClient(customer.ApiKey, customer.ApiSecret)
	account, err := client.NewGetAccountService().Do(context.Background())
	if err != nil {
		return err
	}

	amount := amount(account.Positions, command.Symbol, futures.PositionSideType(command.Side))
	if isZero(amount) {
		return errors.New("position empty")
	}

	side := func(ps futures.PositionSideType) futures.SideType {
		if ps == futures.PositionSideTypeShort {
			return futures.SideTypeSell
		}
		return futures.SideTypeBuy
	}(command.Side)
	if _, err := client.NewCreateOrderService().
		Symbol(command.Symbol).
		Type(futures.OrderTypeMarket).
		Side(side).
		PositionSide(command.Side).
		Quantity(format(command.Symbol, command.Quantity)).
		Do(context.Background()); err != nil {
		return err
	}

	return nil
}

func decr(customer models.Customer, command models.Command) error {
	client := futures.NewClient(customer.ApiKey, customer.ApiSecret)
	account, err := client.NewGetAccountService().Do(context.Background())
	if err != nil {
		return err
	}

	amount := amount(account.Positions, command.Symbol, futures.PositionSideType(command.Side))
	if isZero(amount) {
		return errors.New("position empty")
	}

	side := func(ps futures.PositionSideType) futures.SideType {
		if ps == futures.PositionSideTypeShort {
			return futures.SideTypeBuy
		}
		return futures.SideTypeSell
	}(command.Side)

	quantity := format(command.Symbol, command.Quantity)
	f, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return err
	}
	if command.Quantity > f {
		quantity = abs(amount)
	}

	if _, err := client.NewCreateOrderService().
		Symbol(command.Symbol).
		Type(futures.OrderTypeMarket).
		Side(side).
		PositionSide(command.Side).
		Quantity(quantity).
		Do(context.Background()); err != nil {
		return err
	}

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

func abs(s string) string {
	if s[0:1] == "-" {
		return s[1:]
	}
	return s
}

func isZero(s string) bool {
	return zero.MatchString(s)
}

func opposite(s futures.PositionSideType) futures.PositionSideType {
	if s == futures.PositionSideTypeLong {
		return futures.PositionSideTypeShort
	}

	return futures.PositionSideTypeLong
}

func format(symbol string, f float64) string {
	precision := 0
	for _, s := range ExchangeInfo.Symbols {
		if s.Symbol == symbol {
			precision = s.QuantityPrecision
		}
	}

	return strconv.FormatFloat(f, 'f', precision, 64)
}
