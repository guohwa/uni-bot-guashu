package binance

import (
	"bot/log"
	"context"

	"github.com/robfig/cron/v3"
	"github.com/uncle-gua/gobinance/futures"
)

var ExchangeInfo *futures.ExchangeInfo

func init() {
	client := futures.NewClient("", "")
	info, err := client.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		panic(err)
	}
	ExchangeInfo = info

	c := cron.New(cron.WithSeconds())
	c.AddFunc("@every 10m", func() {
		info, err := client.NewExchangeInfoService().Do(context.Background())
		if err != nil {
			log.Error(err)
		}

		ExchangeInfo = info
	})
	c.Start()
}
