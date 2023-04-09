package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

func orderUnaryClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// 前置处理阶段
	log.Println("===== [Client Interceptor] Method: " + method)

	// 调用远程方法
	err := invoker(ctx, method, req, reply, cc, opts...)

	// 后置处理阶段
	log.Printf("===== [Client Interceptor] Reply: %v", reply)

	return err
}
