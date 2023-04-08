package main

import (
	"context"
	"log"
	pb "ordermanagement/proto"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewOrderManagementClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	order, err := c.GetOrder(ctx, &wrapperspb.StringValue{Value: "abcd"})
	if err != nil {
		log.Fatalf("Could not get order: %v", err)
	}
	log.Printf("Order: %s", order.String())
}
