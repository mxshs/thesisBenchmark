package main

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	"github.com/tinkoff/invest-api-go-sdk/investgo"
	pb "github.com/tinkoff/invest-api-go-sdk/proto"
)

type Candle struct {
	Ticker string
	Low    float64
	High   float64
	Open   float64
	Close  float64
}

func main() {

	config, err := investgo.LoadConfig("config.yaml")
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	client, err := investgo.NewClient(ctx, config, nil)
	if err != nil {
		panic(err)
	}

	str_quotes := []string{"AAPL", "NFLX", "NVDA", "TCSG", "WMT", "GOOG"}

	from := time.Date(2022, time.February, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2022, time.April, 31, 0, 0, 0, 0, time.UTC)

	instrumentsService := client.NewInstrumentsServiceClient()

	quotes := client.NewMarketDataServiceClient()

	connDef := ""

	db, err := sql.Open("postgres", connDef)
	if err != nil {
		panic(err)
	}

	for _, quote := range str_quotes {

		instrResp, err := instrumentsService.FindInstrument(quote)
		if err != nil {
			panic(err)
		}

		id_quote := instrResp.GetInstruments()[0].GetUid()

		candlesResp, err := quotes.GetCandles(id_quote, pb.CandleInterval_CANDLE_INTERVAL_DAY, from, to)
		if err != nil {
			panic(err.Error())
		}

		candles := candlesResp.GetCandles()

		for _, candle := range candles {
			_, err := db.Exec("INSERT INTO candles VALUES ($1, $2, $3, $4, $5)",
				id_quote,
				candle.GetLow().ToFloat(),
				candle.GetHigh().ToFloat(),
				candle.GetOpen().ToFloat(),
				candle.GetClose().ToFloat())
			if err != nil {
				panic(err)
			}
		}
	}
}
