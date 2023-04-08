package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"

	pb "ordermanagement/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type server struct {
	pb.UnimplementedOrderManagementServer

	orderMap map[string]*pb.Order
}

func NewServerDeafult() *server {
	svc := server{}
	svc.orderMap = make(map[string]*pb.Order, 0)
	svc.orderMap["abcd"] = &pb.Order{Id: "abcd", Items: []string{"a", "b", "c", "d"}, Name: "abcd", Price: 123.4, Description: "Description"}
	svc.orderMap["abcd-2"] = &pb.Order{Id: "abcd", Items: []string{"b", "c", "d"}, Name: "bcd", Price: 123.4, Description: "Description"}
	svc.orderMap["abcd-3"] = &pb.Order{Id: "abcd", Items: []string{"c", "d"}, Name: "cd", Price: 123.4, Description: "Description"}
	svc.orderMap["abcd-4"] = &pb.Order{Id: "abcd", Items: []string{"d"}, Name: "d", Price: 123.4, Description: "Description"}

	return &svc
}

func (s *server) GetOrder(ctx context.Context, req *wrapperspb.StringValue) (*pb.Order, error) {
	if s.orderMap == nil {
		return nil, status.Errorf(codes.Internal, "Error while get order")
	}
	order, exists := s.orderMap[req.Value]
	if !exists {
		return nil, status.Errorf(codes.NotFound, codes.NotFound.String())
	}
	return order, status.New(codes.OK, "").Err()
}

func (s *server) SearchOrders(req *wrapperspb.StringValue, stream pb.OrderManagement_SearchOrdersServer) error {
	for key, order := range s.orderMap {
		log.Print(key, order)
		for _, itemStr := range order.Items {
			log.Print(itemStr)
			if strings.Contains(itemStr, req.Value) {
				err := stream.Send(order)
				if err != nil {
					return fmt.Errorf("error sending message to stream err: %v", err)
				}
				log.Printf("Matching Order Found: %v", key)
				break
			}
		}
	}
	return nil
}

func (s *server) UpdateOrders(stream pb.OrderManagement_UpdateOrdersServer) error {
	ordersStr := "Update Order Ids: "
	for {
		order, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&wrapperspb.StringValue{Value: "Orders process " + ordersStr})
		}

		s.orderMap[order.Id] = order

		log.Printf("Order ID %s Updated", order.Id)
		ordersStr += order.Id + ", "
	}
}
