package binance

import (
	"bot/log"
	"context"

	"github.com/robfig/cron/v3"
	"github.com/uncle-gua/gobinance/futures"
)

var (
	count        = 0
	ExchangeInfo *futures.ExchangeInfo
)

func init() {
	synchronize()
	count++

	c := cron.New(cron.WithSeconds())
	c.AddFunc("@every 10m", synchronize)
	c.Start()
}

func synchronize() {
	client := futures.NewClient("", "")
	info, err := client.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		if count == 0 {
			panic(err)
		} else {
			log.Error(err)
		}

		return
	}

	ExchangeInfo = info
}
