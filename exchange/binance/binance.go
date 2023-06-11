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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var zero = regexp.MustCompile(`^[\+\-]?0\.?0*$`)

var (
	ErrApi   = errors.New("API Error")
	ErrEmpty = errors.New("position empty")
	ErrHold  = errors.New("position hold")
	ErrRisk  = errors.New("position risk")
)

type Action func(customer models.Customer, command models.Command) error

func New(customer models.Customer) common.Exchange {
	return &Binance{
		customer,
	}
}

type Binance struct {
	customer models.Customer
}

var actions = map[string]Action{
	"OPEN":  open,
	"CLOSE": close,
	"INCR":  incr,
	"DECR":  decr,
}

func (ex *Binance) FormatSize(symbol string, size float64) string {
	return format(symbol, size)
}

func (ex *Binance) Execute(command models.Command) {
	if action, ok := actions[command.Action]; ok {
		err := action(ex.customer, command)

		filter := bson.M{"_id": command.ID}

		var update bson.M
		if err != nil {
			if err == ErrEmpty || err == ErrHold || err == ErrRisk {
				update = bson.M{"$set": bson.M{
					"status": "FAILED",
					"reason": err.Error(),
				}}
			} else {
				log.Error(err)
				update = bson.M{"$set": bson.M{
					"status": "FAILED",
					"reason": "Api Error",
				}}
			}
		} else {
			update = bson.M{"$set": bson.M{
				"status": "SUCCESS",
			}}
		}

		if err := models.CommandCollection.FindOneAndUpdate(
			context.TODO(),
			filter,
			update,
			options.FindOneAndUpdate(),
		).Err(); err != nil {
			log.Error(err)
		}
	} else {
		log.Errorf("unsupported action: %s", command.Action)
	}
}

func open(customer models.Customer, command models.Command) error {
	client := futures.NewClient(customer.ApiKey, customer.ApiSecret)
	account, err := client.NewGetAccountService().Do(context.Background())
	if err != nil {
		return err
	}

	amount1 := getAmount(account.Positions, command.Symbol, futures.PositionSideType(command.Side))
	if !isZero(amount1) {
		return ErrHold
	}

	oppositeSide := opposite(command.Side)

	amount2 := getAmount(account.Positions, command.Symbol, oppositeSide)
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

	price, err := getPrice(command.Symbol)
	if err != nil {
		return err
	}

	risk1, err := risk(account, command.Quantity, price)
	if err != nil {
		return err
	}
	if risk1 > customer.Level1 {
		return ErrRisk
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

	amount := getAmount(account.Positions, command.Symbol, futures.PositionSideType(command.Side))
	if isZero(amount) {
		return ErrEmpty
	}

	side := func(ps futures.PositionSideType) futures.SideType {
		if ps == futures.PositionSideTypeShort {
			return futures.SideTypeBuy
		}
		return futures.SideTypeSell
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

	amount := getAmount(account.Positions, command.Symbol, futures.PositionSideType(command.Side))
	if isZero(amount) {
		return ErrEmpty
	}

	price, err := getPrice(command.Symbol)
	if err != nil {
		return err
	}

	risk2, err := risk(account, command.Quantity, price)
	if err != nil {
		return err
	}
	if risk2 > customer.Level2 {
		return ErrRisk
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

	amount := getAmount(account.Positions, command.Symbol, futures.PositionSideType(command.Side))
	if isZero(amount) {
		return ErrEmpty
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

func getAmount(positions []*futures.AccountPosition, symbol string, side futures.PositionSideType) string {
	for _, p := range positions {
		if p.Symbol == symbol && p.PositionSide == side {
			return p.PositionAmt
		}
	}

	return "0"
}

func getPrice(symbol string) (float64, error) {
	client := futures.NewClient("", "")
	prices, err := client.NewListPricesService().Do(context.Background())
	if err != nil {
		return 0.0, err
	}

	for _, p := range prices {
		if p.Symbol == symbol {
			return strconv.ParseFloat(p.Price, 64)
		}
	}

	return 0.0, ErrApi
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

func format(symbol string, size float64) string {
	precision := 0
	for _, s := range ExchangeInfo.Symbols {
		if s.Symbol == symbol {
			precision = s.QuantityPrecision
		}
	}

	return strconv.FormatFloat(size, 'f', precision, 64)
}

func risk(account *futures.Account, quantity, price float64) (float64, error) {
	totalMaintMargin, err := strconv.ParseFloat(account.TotalMaintMargin, 64)
	if err != nil {
		return totalMaintMargin, err
	}
	totalMarginBalance, err := strconv.ParseFloat(account.TotalWalletBalance, 64)
	if err != nil {
		return totalMarginBalance, err
	}

	return (totalMaintMargin*100 + quantity*price) / totalMarginBalance, nil
}
