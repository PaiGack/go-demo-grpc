package main

import (
	"context"

	pb "productinfo/proto"

	uuid "github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type server struct {
	pb.UnimplementedProductInfoServer

	productMap map[string]*pb.Product
}

func (s *server) AddProduct(ctx context.Context, req *pb.Product) (*wrapperspb.StringValue, error) {
	out, err := uuid.NewUUID()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error while generating Product ID, err: %v", err)
	}
	req.Id = out.String()
	if s.productMap == nil {
		s.productMap = make(map[string]*pb.Product)
	}
	s.productMap[req.Id] = req
	return &wrapperspb.StringValue{Value: req.Id}, status.New(codes.OK, "").Err()
}

func (s *server) GetProduct(ctx context.Context, req *wrapperspb.StringValue) (*pb.Product, error) {
	value, exists := s.productMap[req.Value]
	if exists {
		return value, status.New(codes.OK, "").Err()
	}
	return nil, status.Errorf(codes.NotFound, "Product does not exists")
}
