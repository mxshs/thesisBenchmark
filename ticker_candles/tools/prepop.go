package main

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	"github.com/tinkoff/invest-api-go-sdk/investgo"
	pb "github.com/tinkoff/invest-api-go-sdk/proto"
)

type Ticker struct {
	Id             string
	Ticker         string
	ClassCode      string
	InstrumentType string
	Name           string
}

type Candle struct {
	Ticker string
	Low    float64
	High   float64
	Open   float64
	Close  float64
}

type DB struct {
	db *sql.DB
}

func (db *DB) getWriteCandles(client *investgo.Client, tickers []string) {

	instrumentsService := client.NewInstrumentsServiceClient()
	quotes := client.NewMarketDataServiceClient()

	from := time.Date(2022, time.February, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2022, time.April, 31, 0, 0, 0, 0, time.UTC)

	for _, ticker := range tickers {

		instr, err := instrumentsService.FindInstrument(ticker)
		if err != nil {
			panic(err)
		}

		id_quote := instr.GetInstruments()[0].GetUid()

		candlesResp, err := quotes.GetCandles(id_quote, pb.CandleInterval_CANDLE_INTERVAL_DAY, from, to)
		if err != nil {
			panic(err)
		}

		candles := candlesResp.GetCandles()

		for _, candle := range candles {
			_, err := db.db.Exec("INSERT INTO candles VALUES ($1, $2, $3, $4, $5);",
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

func (db *DB) getWriteTickers(client *investgo.Client, str_tickers []string) {

	instrumentsService := client.NewInstrumentsServiceClient()

	for i := range str_tickers {

		instrResp, err := instrumentsService.FindInstrument(str_tickers[i])
		if err != nil {
			panic(err)
		}

		ticker := instrResp.GetInstruments()[0]

		_, err = db.db.Exec("INSERT INTO tickers VALUES ($1, $2, $3, $4, $5);",
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

func setDB() (*DB, error) {

	connDef := ""

	db, err := sql.Open("postgres", connDef)
	if err != nil {
		return nil, err
	}

	return &DB{db: db}, nil
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

	db, err := setDB()
	if err != nil {
		panic(err)
	}

	db.getWriteCandles(client, str_tickers)
	db.getWriteTickers(client, str_tickers)
}
