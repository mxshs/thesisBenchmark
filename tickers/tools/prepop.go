package main

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/tinkoff/invest-api-go-sdk/investgo"
)

type Ticker struct {
	Id             string
	Ticker         string
	ClassCode      string
	InstrumentType string
	Name           string
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

	str_tickers := []string{"AAPL", "NFLX", "NVDA", "TCSG", "WMT", "GOOG"}

	instrumentsService := client.NewInstrumentsServiceClient()

	connDef := ""

	db, err := sql.Open("postgres", connDef)
	if err != nil {
		panic(err)
	}

	for i := range str_tickers {
		instrResp, err := instrumentsService.FindInstrument(str_tickers[i])
		if err != nil {
			panic(err)
		} else {
			ticker := instrResp.GetInstruments()[0]
			_, err := db.Exec("INSERT INTO tickers VALUES ($1, $2, $3, $4, $5);",
				ticker.Uid,
				ticker.Ticker,
				ticker.ClassCode,
				ticker.InstrumentType,
				ticker.Name)
			if err != nil {
				panic(err)
			}
		}
	}
}
