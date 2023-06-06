package db

import (
	"database/sql"
	"mxshs/tickers/protos/tickersproto"

	"github.com/lib/pq"
)

var DB *sql.DB

func QueryDB(input *tickersproto.TickerId, db *sql.DB) <-chan ([]*tickersproto.TickerData) {

	ch := make(chan []*tickersproto.TickerData)

	go func() {

		defer close(ch)

		rows, err := db.Query("SELECT * FROM tickers WHERE ticker_id = any($1)", pq.Array(&input.TickerId))
		if err != nil {
			panic(err)
		}

		tickers := []*tickersproto.TickerData{}

		for rows.Next() {

			var r tickersproto.TickerData

			if err := rows.Scan(&r.TickerId, &r.Ticker, &r.ClassCode, &r.InstrumentType, &r.Name); err != nil {
				panic(err)
			}

			tickers = append(tickers, &r)
		}

		ch <- tickers
	}()

	return ch
}

func CallDB(input *tickersproto.TickerId) <-chan ([]*tickersproto.TickerData) {

	if DB == nil {
		connDef := "host=postgres-tickers user=test password=test dbname=test sslmode=disable"

		db, err := sql.Open("postgres", connDef)
		if err != nil {
			panic(err)
		}

		DB = db
	}

	return QueryDB(input, DB)
}
