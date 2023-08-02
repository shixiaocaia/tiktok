package main

import (
	"context"
	"github.com/shixiaocaia/tiktok/pkg/pb"
	"google.golang.org/grpc"
	"log"
)

var (
	ctx = context.Background()
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:8000", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	client := pb.NewGreeterClient(conn)
	resp, err := client.SayHello(ctx, &pb.HelloRequest{Name: "hello"})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp.GetMessage())
}
