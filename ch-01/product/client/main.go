package main

import (
	"context"
	"flag"
	"log"
	"time"

	uuid "github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "demo-grpc/ch-01/product/proto"
)

var (
	addr = flag.String("add", "localhost:50051", "the address to connect to")
)

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewProductInfoClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	for i := 0; i < 100; i++ {
		id := uuid.NewString()

		r, err := c.AddProduct(ctx, &pb.Product{Id: id, Name: "p-" + id, Description: "p-" + id + ", is very good goods."})
		if err != nil {
			log.Fatalf("could not add product: %v", err)
		}
		log.Printf("add product success: %s", r.GetValue())

		r2, err := c.GetProduct(ctx, &pb.ProductID{Value: id})
		log.Printf("get product success: %v, err: %v", r2, err)
	}
}
