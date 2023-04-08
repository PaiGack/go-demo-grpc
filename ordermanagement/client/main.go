package main

import (
	"context"
	"io"
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

	// stream
	searchStream, _ := c.SearchOrders(ctx, &wrapperspb.StringValue{Value: "b"})
	for {
		searchOrder, err := searchStream.Recv()
		if err == io.EOF {
			break
		}
		log.Printf("Search Result: %v, err: %v", searchOrder, err)
	}

	// update order
	updateStream, err := c.UpdateOrders(ctx)
	if err != nil {
		log.Fatalf("%v.UpdateOrders(_)= _, %v", c, err)
	}

	updOrder1 := pb.Order{Id: "abcd", Name: "u1"}
	updOrder2 := pb.Order{Id: "abcd-2", Name: "u2"}
	updOrder3 := pb.Order{Id: "abcd-3", Name: "u3"}
	// update order1
	if err := updateStream.Send(&updOrder1); err != nil {
		log.Fatalf("%v.Send(%v) = %v", updateStream, updOrder1, err)
	}

	if err := updateStream.Send(&updOrder2); err != nil {
		log.Fatalf("%v.Send(%v) = %v", updateStream, updOrder2, err)
	}

	if err := updateStream.Send(&updOrder3); err != nil {
		log.Fatalf("%v.Send(%v) = %v", updateStream, updOrder3, err)
	}

	updateRes, err := updateStream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v", updateStream, err, nil)
	}
	log.Printf("Update orders Res: %s", updateRes)
}
