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
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var store = make(map[string]*grpc.ClientConn)

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

	fmt.Println("fullMethodName: ", fullMethodName)

	md, _ := metadata.FromIncomingContext(ctx)

	var header metadata.MD
	if val, exists := md["authorization"]; exists {
		header = metadata.New(map[string]string{
			"authorization": val[0],
		})	
	}

	ctx = metadata.NewOutgoingContext(ctx, header)

	var conn *grpc.ClientConn
	var err error

	switch {
	case strings.HasPrefix(fullMethodName, "/helloworld.Greeter/"):
		conn, err = getConnection(ctx, store, "localhost:8081")
	case strings.HasPrefix(fullMethodName, "/auth.AuthService/"):
		conn, err = getConnection(ctx, store, "localhost:8081")
	default:
		err = status.Errorf(codes.Unimplemented, "Unknown method")
	}

	return ctx, conn, err
}

func getConnection(ctx context.Context, store map[string]*grpc.ClientConn, address string) (*grpc.ClientConn, error) {
	// If connection exists, return it.
	if conn, ok := store[address]; ok {
		return conn, nil
	}

	// If connection does not exist, create a new one and store it.
	newConn, err := grpc.DialContext(ctx, address, grpc.WithCodec(proxy.Codec()), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	store[address] = newConn
	
	return newConn, nil
}