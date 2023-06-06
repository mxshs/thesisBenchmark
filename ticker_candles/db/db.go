package db

import (
	"database/sql"
	"mxshs/ticker_candles/protos"

	"github.com/lib/pq"
)

var DB *sql.DB

func GetTickerData(input *protos.TickerId, db *sql.DB) <-chan ([]*protos.TickerData) {

	ch := make(chan []*protos.TickerData)

	go func() {

		defer close(ch)

		res := []*protos.TickerData{}

		ts, err := db.Query("SELECT * FROM tickers WHERE ticker_id = any($1)", pq.Array(input.TickerId))
		if err != nil {
			panic(err)
		}

		for ts.Next() {

			var t protos.TickerData

			if err := ts.Scan(&t.TickerId, &t.Ticker, &t.ClassCode, &t.InstrumentType, &t.Name); err != nil {
				panic(err)
			}

			res = append(res, &t)
		}

		ch <- res
	}()

	return ch
}

func GetCandlesData(input *protos.TickerId, db *sql.DB) <-chan ([]*protos.CandleData) {

	ch := make(chan []*protos.CandleData)

	go func() {

		defer close(ch)

		res := []*protos.CandleData{}

		cs, err := db.Query("SELECT * FROM candles WHERE ticker_id = any($1)", pq.Array(input.TickerId))
		if err != nil {
			panic(err)
		}

		for cs.Next() {

			var r protos.CandleData

			if err := cs.Scan(&r.Ticker, &r.Low, &r.High, &r.Open, &r.Close); err != nil {
				panic(err)
			}

			res = append(res, &r)
		}

		ch <- res
	}()

	return ch
}

func CallDB(input *protos.TickerId) (<-chan ([]*protos.TickerData), <-chan ([]*protos.CandleData)) {

	if DB == nil {

		connDef := ""

		db, err := sql.Open("postgres", connDef)
		if err != nil {
			panic(err)
		}

		DB = db
	}

	return GetTickerData(input, DB), GetCandlesData(input, DB)
}
