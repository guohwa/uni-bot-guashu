package exchange

import (
	"bot/log"
	"bot/models"
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/uncle-gua/gobinance/futures"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	RETRY = 5
	SLEEP = 3
)

var zero = regexp.MustCompile(`^[\+\-]?0\.?0*$`)

var (
	ErrApi   = errors.New("API Error")
	ErrEmpty = errors.New("position empty")
	ErrHold  = errors.New("position hold")
	ErrRisk  = errors.New("position risk")
)

type Action func(customer models.Customer, command models.Command) error

type CreateOrderService interface {
	Do(ctx context.Context, opts ...futures.RequestOption) (res *futures.CreateOrderResponse, err error)
}

var actions = map[string]Action{
	"OPEN":  open,
	"CLOSE": close,
	"INCR":  incr,
	"DECR":  decr,
}

func Execute(customer models.Customer, command models.Command) {
	action, ok := actions[command.Action]
	if !ok {
		log.Errorf("unsupported action: %s", command.Action)
		return
	}

	update := func(err error) bson.M {
		if err == nil {
			return bson.M{"$set": bson.M{
				"status": "SUCCESS",
			}}
		}
		if err == ErrEmpty || err == ErrHold || err == ErrRisk {
			return bson.M{"$set": bson.M{
				"status": "FAILED",
				"reason": err.Error(),
			}}
		}
		log.Error(err)
		return bson.M{"$set": bson.M{
			"status": "FAILED",
			"reason": "Api Error",
		}}
	}(action(customer, command))

	filter := bson.M{"_id": command.ID}
	if err := models.CommandCollection.FindOneAndUpdate(
		context.TODO(),
		filter,
		update,
		options.FindOneAndUpdate(),
	).Err(); err != nil {
		log.Error(err)
	}
}

func FormatSize(symbol string, size float64) string {
	precision := 0
	for _, s := range ExchangeInfo.Symbols {
		if s.Symbol == symbol {
			precision = s.QuantityPrecision
		}
	}

	return strconv.FormatFloat(size, 'f', precision, 64)
}

func createOrder(service CreateOrderService, action string) error {
	retry := 0
	for {
		_, err := service.Do(context.Background())
		if err == nil {
			return nil
		}
		if retry == 0 {
			log.Errorf("error occurred while %s position, message: %s", action, err.Error())
		} else if retry < RETRY {
			log.Errorf("error occurred while %s position, message: %s, retry: %d", action, err.Error(), retry)
		} else {
			return fmt.Errorf("fatal error while %s position, message: %s", action, err.Error())
		}
		time.Sleep(SLEEP * time.Millisecond)
		retry++
	}
}

func getAccount(client *futures.Client) (*futures.Account, error) {
	retry := 0
	for {
		account, err := client.NewGetAccountService().Do(context.Background())
		if err == nil {
			return account, nil
		}
		if retry == 0 {
			log.Errorf("error occurred while get account, message: %s", err.Error())
		} else if retry < RETRY {
			log.Errorf("error occurred while get account, message: %s, retry: %d", err.Error(), retry)
		} else {
			return nil, fmt.Errorf("fatal error while get account, message: %s", err.Error())
		}
		time.Sleep(SLEEP * time.Millisecond)
		retry++
	}
}

func open(customer models.Customer, command models.Command) error {
	client := futures.NewClient(customer.ApiKey, customer.ApiSecret)
	account, err := getAccount(client)
	if err != nil {
		return err
	}

	amount1 := getAmount(account.Positions, command.Symbol, futures.PositionSideType(command.Side))
	if !isZero(amount1) {
		return ErrHold
	}

	otherSide := func(s futures.PositionSideType) futures.PositionSideType {
		if s == futures.PositionSideTypeLong {
			return futures.PositionSideTypeShort
		}

		return futures.PositionSideTypeLong
	}(command.Side)

	amount2 := getAmount(account.Positions, command.Symbol, otherSide)
	if !isZero(amount2) {
		side := func(ps futures.PositionSideType) futures.SideType {
			if ps == futures.PositionSideTypeLong {
				return futures.SideTypeSell
			}
			return futures.SideTypeBuy
		}(otherSide)
		service := client.NewCreateOrderService().
			Symbol(command.Symbol).
			Type(futures.OrderTypeMarket).
			Side(side).
			PositionSide(otherSide).
			Quantity(abs(amount2))
		if err := createOrder(service, "reverse"); err != nil {
			return err
		}
	}

	account, err = getAccount(client)
	if err != nil {
		return err
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

	service := client.NewCreateOrderService().
		Symbol(command.Symbol).
		Type(futures.OrderTypeMarket).
		Side(side).
		PositionSide(command.Side).
		Quantity(format(command.Symbol, command.Quantity))
	return createOrder(service, "open")
}

func close(customer models.Customer, command models.Command) error {
	client := futures.NewClient(customer.ApiKey, customer.ApiSecret)
	account, err := getAccount(client)
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

	service := client.NewCreateOrderService().
		Symbol(command.Symbol).
		Type(futures.OrderTypeMarket).
		Side(side).
		PositionSide(command.Side).
		Quantity(abs(amount))
	return createOrder(service, "close")
}

func incr(customer models.Customer, command models.Command) error {
	client := futures.NewClient(customer.ApiKey, customer.ApiSecret)
	account, err := getAccount(client)
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

	service := client.NewCreateOrderService().
		Symbol(command.Symbol).
		Type(futures.OrderTypeMarket).
		Side(side).
		PositionSide(command.Side).
		Quantity(format(command.Symbol, command.Quantity))
	return createOrder(service, "incr")
}

func decr(customer models.Customer, command models.Command) error {
	client := futures.NewClient(customer.ApiKey, customer.ApiSecret)
	account, err := getAccount(client)
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

	service := client.NewCreateOrderService().
		Symbol(command.Symbol).
		Type(futures.OrderTypeMarket).
		Side(side).
		PositionSide(command.Side).
		Quantity(quantity)
	return createOrder(service, "decr")
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
