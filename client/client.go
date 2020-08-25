package main

import (
	"context"
	"fmt"
	"io"
	"log"

	grpcserver "github.com/aanufriev/SoftProTest/grpcServer"
	"google.golang.org/grpc"
)

func main() {

	grcpConn, err := grpc.Dial(
		"localhost:9001",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}
	defer grcpConn.Close()

	sub := grpcserver.NewSubscribeOnSportsLinesClient(grcpConn)

	ctx := context.Background()
	client, err := sub.Subscribe(ctx)
	if err != nil {
		return
	}

	sports := []string{"baseball", "football", "soccer"}

	err = client.Send(&grpcserver.Request{
		Sports:   sports,
		Interval: 3,
	})
	if err != nil {
		return
	}

	for i := 0; i < 2; i++ {
		out, err := client.Recv()
		if err == io.EOF {
			fmt.Println("stream closed")
			return
		} else if err != nil {
			fmt.Println("error happed", err)
			return
		}
		fmt.Println(out.Lines)
	}

	err = client.Send(&grpcserver.Request{
		Sports:   sports,
		Interval: 2,
	})
	if err != nil {
		return
	}

	for i := 0; i < 2; i++ {
		out, err := client.Recv()
		if err == io.EOF {
			fmt.Println("stream closed")
			return
		} else if err != nil {
			fmt.Println("error happed", err)
			return
		}
		fmt.Println(out.Lines)
	}

	newSports := []string{"baseball", "soccer"}

	err = client.Send(&grpcserver.Request{
		Sports:   newSports,
		Interval: 1,
	})
	if err != nil {
		return
	}

	for {
		out, err := client.Recv()
		if err == io.EOF {
			fmt.Println("stream closed")
			return
		} else if err != nil {
			fmt.Println("error happed", err)
			return
		}
		fmt.Println(out.Lines)
	}
}
