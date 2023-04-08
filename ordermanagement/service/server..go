package main

import (
	"context"

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
