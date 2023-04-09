package main

import (
	"log"
	"net"

	pb "ordermanagement/proto"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(
		grpc.UnaryInterceptor(orderUnaryServerInterceptor),
		grpc.StreamInterceptor(orderServerStreamInterceptor),
	)
	pb.RegisterOrderManagementServer(s, NewServerDeafult())

	log.Printf("Starting gRPC listener on port: %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
