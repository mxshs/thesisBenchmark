package main

import (
	"context"
	"mxshs/candles/candlesproto"
	"mxshs/candles/db"

	"net"

	"google.golang.org/grpc"
)

type server struct {
	candlesproto.UnimplementedStockServiceServer
}

func (*server) GetPrices(ctx context.Context, inp *candlesproto.TickerId) (*candlesproto.CandlesResponse, error) {

	p := db.CallDB(inp)

	return <-p, nil
}

func main() {

	listener, err := net.Listen("tcp", ":9002")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()

	candlesproto.RegisterStockServiceServer(grpcServer, &server{})

	grpcServer.Serve(listener)
}
