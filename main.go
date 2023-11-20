package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/mwitkow/grpc-proxy/proxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.UnknownServiceHandler(proxy.TransparentHandler(director)),
	)

	log.Println("Serving gRPC on 0.0.0.0:50052")
	log.Fatal(s.Serve(lis))
}

func director(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {

	fmt.Println("fullName: ", fullMethodName)

	var conn *grpc.ClientConn
	var err error
	switch {
	case strings.HasPrefix(fullMethodName, "/helloworld.Greeter/"):
		conn, err = grpc.DialContext(ctx, "localhost:8000", grpc.WithCodec(proxy.Codec()), grpc.WithInsecure())
	case strings.HasPrefix(fullMethodName, "/goodbye.Greeter/"):
		conn, err = grpc.DialContext(ctx, "localhost:8001", grpc.WithCodec(proxy.Codec()), grpc.WithInsecure())
	default:
		err = status.Errorf(codes.Unimplemented, "Unknown method")
	}
	return ctx, conn, err
}
