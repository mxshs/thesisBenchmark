package main

import (
	"context"
	"mxshs/tickers/db"
	"mxshs/tickers/metrics"
	"mxshs/tickers/protos/candlesproto"
	"mxshs/tickers/protos/tickersproto"
	"net"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var pusher *push.Pusher

type server struct {
	tickersproto.UnimplementedTickerServiceServer
}

func GetCandles(ctx context.Context, inp *tickersproto.TickerId) <-chan (*candlesproto.CandlesResponse) {

	ch := make(chan *candlesproto.CandlesResponse)

	go func() {

		conn, err := grpc.Dial("", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(err)
		}

		c := candlesproto.NewStockServiceClient(conn)

		res, err := c.GetPrices(ctx, &candlesproto.TickerId{TickerId: inp.TickerId})
		if err != nil {
			panic(err)
		}

		ch <- res
	}()

	return ch
}

func (*server) GetTickers(ctx context.Context, inp *tickersproto.TickerId) (*tickersproto.TickerResponse, error) {

	start := time.Now()

	prom1, prom2 := db.CallDB(inp), GetCandles(ctx, inp)

	candles := <-prom2

	res := tickersproto.TickerResponse{Tickers: <-prom1,
		Candles: candles.Prices}

	metrics.TTS.Observe(float64(time.Since(start).Milliseconds()))
	metrics.TotalRequests.Add(1)

	pusher.Add()

	return &res, nil
}

func main() {

	listener, err := net.Listen("tcp", ":9002")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()

	tickersproto.RegisterTickerServiceServer(grpcServer, &server{})

	reg := prometheus.NewRegistry()

	metrics.SetMetrics()

	reg.MustRegister(metrics.TTS, metrics.TotalRequests)

	pusher = push.New("", "ticker_metrics").Gatherer(reg)

	grpcServer.Serve(listener)
}
