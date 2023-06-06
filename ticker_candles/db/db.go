package db

import (
	"database/sql"
	"mxshs/ticker_candles/protos"

	"github.com/lib/pq"
)

type DB struct {
	db *sql.DB
}

var DBInstance *DB

func (db DB) GetTickerData(input *protos.TickerId) <-chan ([]*protos.TickerData) {

	ch := make(chan []*protos.TickerData)

	go func() {

		defer close(ch)

		res := []*protos.TickerData{}

		ts, err := db.db.Query("SELECT * FROM tickers WHERE ticker_id = any($1)", pq.Array(input.TickerId))
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

func (db DB) GetCandlesData(input *protos.TickerId) <-chan ([]*protos.CandleData) {

	ch := make(chan []*protos.CandleData)

	go func() {

		defer close(ch)

		res := []*protos.CandleData{}

		cs, err := db.db.Query("SELECT * FROM candles WHERE ticker_id = any($1)", pq.Array(input.TickerId))
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

func SetDB() (*DB, error) {

	if DBInstance == nil {

		connDef := "host=postgres-ticker-candles user=test password=test dbname=test sslmode=disable"

		db, err := sql.Open("postgres", connDef)
		if err != nil {
			return nil, err
		}

		DBInstance = &DB{db: db}
	}

	return DBInstance, nil
}

func CallDB(input *protos.TickerId) (<-chan ([]*protos.TickerData), <-chan ([]*protos.CandleData)) {

	db, err := SetDB()
	if err != nil {
		panic(err)
	}

	return db.GetTickerData(input), db.GetCandlesData(input)
}
