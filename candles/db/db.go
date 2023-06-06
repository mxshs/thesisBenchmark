package db

import (
	"database/sql"
	"mxshs/candles/candlesproto"

	"github.com/lib/pq"
)

var DB *sql.DB

func QueryDB(input *candlesproto.TickerId, db *sql.DB) <-chan (*candlesproto.CandlesResponse) {

	ch := make(chan *candlesproto.CandlesResponse)

	go func() {

		defer close(ch)

		rows, err := db.Query("SELECT * FROM candles WHERE ticker_id = any($1)", pq.Array(&input.TickerId))
		if err != nil {
			panic(err)
		}

		prices := []*candlesproto.Candle{}

		for rows.Next() {

			var r candlesproto.Candle

			if err := rows.Scan(&r.Ticker, &r.Low, &r.High, &r.Open, &r.Close); err != nil {
				panic(err)
			}

			prices = append(prices, &r)
		}

		ch <- &candlesproto.CandlesResponse{Prices: prices}
	}()

	return ch
}

func CallDB(inp *candlesproto.TickerId) <-chan (*candlesproto.CandlesResponse) {

	if DB == nil {

		connDef := ""

		db, err := sql.Open("postgres", connDef)
		if err != nil {
			panic(err)
		}

		DB = db
	}

	return QueryDB(inp, DB)
}
