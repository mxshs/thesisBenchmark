package main

import (
	"context"
	"mxshs/ticker_candles/db"
	"mxshs/ticker_candles/metrics"
	"mxshs/ticker_candles/protos"
	"net"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"google.golang.org/grpc"
)

var pusher *push.Pusher

type server struct {
	protos.UnimplementedTickerServiceServer
}

func (*server) GetTickerCandles(ctx context.Context, in *protos.TickerId) (*protos.TickerCandlesResponse, error) {

	start := time.Now()

	prom1, prom2 := db.CallDB(in)

	res := protos.TickerCandlesResponse{Tickers: <-prom1, Candles: <-prom2}

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

	protos.RegisterTickerServiceServer(grpcServer, &server{})

	reg := prometheus.NewRegistry()

	metrics.SetMetrics()

	reg.MustRegister(metrics.TTS, metrics.TotalRequests)

	pusher = push.New("bench-prom-prometheus-pushgateway:9091", "ticker_metrics").Gatherer(reg)

	grpcServer.Serve(listener)
}
