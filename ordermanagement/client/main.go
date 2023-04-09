package main

import (
	"context"
	"io"
	"log"
	pb "ordermanagement/proto"
	"time"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(orderUnaryClientInterceptor),
		grpc.WithStreamInterceptor(orderClientStreamInterceptor),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewOrderManagementClient(conn)

	// 流式的 rpc 也受限
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	// Add Order
	order1 := pb.Order{Id: "101", Items: []string{"iPhone XS", "Mac Book Pro"}, Destination: "San Jose, CA", Price: 2300.00}
	res, err := c.AddOrder(ctx, &order1)
	if err != nil {
		log.Fatalf("Error Occured -> addOrder: %v", err)
	}
	log.Print("AddOrder Response -> ", res.Value)

	order2 := pb.Order{Id: "-1", Items: []string{"iPhone XS", "Mac Book Pro"}, Destination: "San Jose, CA", Price: 2300.00}
	res2, err := c.AddOrder(ctx, &order2)
	if err != nil {
		errorCode := status.Code(err)
		if errorCode == codes.InvalidArgument {
			log.Printf("Invalid Argument -> addOrder: %v", err)
			errorStatus := status.Convert(err)
			for _, d := range errorStatus.Details() {
				switch info := d.(type) {
				case *errdetails.BadRequest_FieldViolation:
					log.Fatalf("Request Field Invalid: %s", info)
				default:
					log.Fatalf("Unexpected error type: %s", info)
				}
			}
		} else {
			log.Fatalf("Error Occured -> addOrder: %v", err)
		}
	}
	log.Print("AddOrder Response -> ", res2.Value)

	order, err := c.GetOrder(ctx, &wrapperspb.StringValue{Value: "101"})
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

	updOrder1 := pb.Order{Id: "102", Items: []string{"Google Pixel 3A", "Google Pixel Book"}, Destination: "Mountain View, CA", Price: 1100.00}
	updOrder2 := pb.Order{Id: "103", Items: []string{"Apple Watch S4", "Mac Book Pro", "iPad Pro"}, Destination: "San Jose, CA", Price: 2800.00}
	updOrder3 := pb.Order{Id: "104", Items: []string{"Google Home Mini", "Google Nest Hub", "iPad Mini"}, Destination: "Mountain View, CA", Price: 2200.00}

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

	// ProcessOrders
	streamProcOrder, err := c.ProcessOrders(ctx)
	if err != nil {
		log.Fatalf("%v.ProcessOrders(_) = _, %v", c, err)
	}

	if err := streamProcOrder.Send(&wrapperspb.StringValue{Value: "102"}); err != nil {
		log.Fatalf("%v.Send(%v) = %v", c, "102", err)
	}

	if err := streamProcOrder.Send(&wrapperspb.StringValue{Value: "103"}); err != nil {
		log.Fatalf("%v.Send(%v) = %v", c, "103", err)
	}

	if err := streamProcOrder.Send(&wrapperspb.StringValue{Value: "104"}); err != nil {
		log.Fatalf("%v.Send(%v) = %v", c, "104", err)
	}

	channel := make(chan struct{})
	go asncClientBidirectionalRPC(streamProcOrder, channel)
	time.Sleep(time.Millisecond * 1000)

	cancel()

	if err := streamProcOrder.Send(&wrapperspb.StringValue{Value: "101"}); err != nil {
		log.Fatalf("%v.Send(%v) = %v", c, "101", err)
	}
	if err := streamProcOrder.CloseSend(); err != nil {
		log.Fatal(err)
	}
	channel <- struct{}{}
}

func asncClientBidirectionalRPC(streamProcOrder pb.OrderManagement_ProcessOrdersClient, c chan struct{}) {
	for {
		combinedShipment, errProcOrder := streamProcOrder.Recv()
		if errProcOrder == io.EOF {
			break
		}
		log.Printf("Combined shipment: %v", combinedShipment.OrdersList)
	}
	<-c
}
