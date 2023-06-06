package main

import (
	"context"
	"mxshs/ticker_candles/protos"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	conn, err := grpc.Dial("", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	c := protos.NewTickerServiceClient(conn)

	ctx := context.TODO()

	for {
		for i := 0; i < 20; i++ {

			_, err := c.GetTickerCandles(ctx, &protos.TickerId{TickerId: []string{
				"a9eb4238-eba9-488c-b102-b6140fd08e38",
				"de367235-e3c7-4c6c-84b7-1836e5fba34a",
				"81575098-df8a-45c4-82dc-1b64374dcfdb",
				"6afa6f80-03a7-4d83-9cf0-c19d7d021f76",
				"0ac3afbc-ddf7-486d-a8e6-96cca53eeae7",
				"19dfb3f6-2222-4eb4-93b6-f2500aa97e82"}})
			if err != nil {
				panic(err)
			}
		}

		time.Sleep(1 * time.Second)
	}
}
