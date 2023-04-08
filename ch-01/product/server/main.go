package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"

	pb "demo-grpc/ch-01/product/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	pb.UnimplementedProductInfoServer
	products sync.Map
}

func (s *server) AddProduct(ctx context.Context, req *pb.Product) (*pb.ProductID, error) {
	log.Printf("AddProduct  %v", req)
	_, ok := s.products.LoadOrStore(req.Id, req)
	if ok {
		return nil, status.Errorf(codes.AlreadyExists, fmt.Sprintf("%s product id: %s", codes.AlreadyExists.String(), req.Id))
	}

	return &pb.ProductID{Value: req.Id}, nil
}

func (s *server) GetProduct(ctx context.Context, req *pb.ProductID) (*pb.Product, error) {
	log.Printf("GetProduct  %v", req)
	val, ok := s.products.Load(req.Value)
	if !ok {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("%s product id: %s", codes.NotFound.String(), req.Value))
	}
	return val.(*pb.Product), nil
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterProductInfoServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
