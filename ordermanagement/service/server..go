package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	pb "ordermanagement/proto"

	"google.golang.org/grpc"
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

	svc.orderMap["102"] = &pb.Order{Id: "102", Items: []string{"Google Pixel 3A", "Mac Book Pro"}, Destination: "Mountain View, CA", Price: 1800.00}
	svc.orderMap["103"] = &pb.Order{Id: "103", Items: []string{"Apple Watch S4"}, Destination: "San Jose, CA", Price: 400.00}
	svc.orderMap["104"] = &pb.Order{Id: "104", Items: []string{"Google Home Mini", "Google Nest Hub"}, Destination: "Mountain View, CA", Price: 400.00}
	svc.orderMap["105"] = &pb.Order{Id: "105", Items: []string{"Amazon Echo"}, Destination: "San Jose, CA", Price: 30.00}
	svc.orderMap["106"] = &pb.Order{Id: "106", Items: []string{"Amazon Echo", "Apple iPhone XS"}, Destination: "Mountain View, CA", Price: 300.00}

	return &svc
}

// Simple RPC
func (s *server) AddOrder(ctx context.Context, req *pb.Order) (*wrapperspb.StringValue, error) {
	time.Sleep(time.Second * 5)
	log.Printf("Order Added. ID : %v", req.Id)
	s.orderMap[req.Id] = req
	return &wrapperspb.StringValue{Value: "Order Added: " + req.Id}, nil
}

// Simple RPC
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

// Server-side Streaming RPC
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

// Client-side Streaming RPC
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

// Bi-directional Streaming RPC
func (s *server) ProcessOrders(stream pb.OrderManagement_ProcessOrdersServer) error {
	batchMarker := 1
	orderBatchSize := 3

	var combinedShipmentMap = make(map[string]pb.CombinedShipment)
	for {
		orderId, err := stream.Recv()
		log.Printf("Reading Proc order : %s", orderId)
		if err == io.EOF {
			// Client has sent all the messages
			// Send remaining shipments
			log.Printf("EOF : %s", orderId)
			for _, shipment := range combinedShipmentMap {
				if err := stream.Send(&shipment); err != nil {
					return err
				}
			}
			return nil
		}
		if err != nil {
			log.Println(err)
			return err
		}

		destination := s.orderMap[orderId.GetValue()].Destination
		shipment, found := combinedShipmentMap[destination]

		if found {
			ord := s.orderMap[orderId.GetValue()]
			shipment.OrdersList = append(shipment.OrdersList, ord)
			combinedShipmentMap[destination] = shipment
		} else {
			comShip := pb.CombinedShipment{Id: "cmb - " + (s.orderMap[orderId.GetValue()].Destination), Status: "Processed!"}
			ord := s.orderMap[orderId.GetValue()]
			comShip.OrdersList = append(shipment.OrdersList, ord)
			combinedShipmentMap[destination] = comShip
			log.Print(len(comShip.OrdersList), comShip.GetId())
		}

		if batchMarker == orderBatchSize {
			for _, comb := range combinedShipmentMap {
				log.Printf("Shipping : %v -> %v", comb.Id, len(comb.OrdersList))
				if err := stream.Send(&comb); err != nil {
					return err
				}
			}
			batchMarker = 0
			combinedShipmentMap = make(map[string]pb.CombinedShipment)
		} else {
			batchMarker++
		}
	}
}

// 服务端一元拦截器
func orderUnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// 前置处理逻辑
	// 通过检查传入的参数，获取关于当前 RPC 的信息
	log.Println("====== [Server Interceptor] ", info.FullMethod)

	// 调用 handler 完成一元 RPC 的正常执行
	m, err := handler(ctx, req)

	// 后置处理逻辑
	log.Printf(" Post Proc Message: %s", m)
	return m, err
}

// 服务端流拦截器
// wrappedStream 包装嵌入 grpc.ServerStream
// 并拦截对 RecvMsg 和 SendMsg 函数的调用
type wrappedStream struct {
	grpc.ServerStream
}

func (w *wrappedStream) RecvMsg(m interface{}) error {
	log.Printf("===== [Server Stream Interceptor Wrapper] Receive a message (Type: %T) at %s", m, time.Now().Format(time.RFC3339))

	return w.ServerStream.RecvMsg(m)
}

func (w *wrappedStream) SendMsg(m interface{}) error {
	log.Printf("===== [Server Stream Interceptor Wrapper] Send a message (Type: %T) at %s", m, time.Now().Format(time.RFC3339))

	return w.ServerStream.SendMsg(m)
}

func newWrappedStream(s grpc.ServerStream) grpc.ServerStream {
	return &wrappedStream{s}
}

func orderServerStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	log.Println("===== [Server Stream Interceptor] ", info.FullMethod)

	err := handler(srv, newWrappedStream(ss))
	if err != nil {
		log.Printf("RPC failed with error %v", err)
	}
	return err
}
