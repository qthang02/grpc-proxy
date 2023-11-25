package main

import (
	"context"
	"fmt"
	"github.com/mwitkow/grpc-proxy/proxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"strings"
)

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
	//grpc.UnknownServiceHandler(proxy.TransparentHandler(director)),
	)

	proxy.RegisterService(s, director, "helloworld.Greeter", "SayHello")
	proxy.RegisterService(s, director, "auth.AuthService", "Login", "SignUp")

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
		conn, err = grpc.DialContext(ctx, "localhost:8080", grpc.WithCodec(proxy.Codec()), grpc.WithInsecure())
	case strings.HasPrefix(fullMethodName, "/auth.AuthService/"):
		conn, err = grpc.DialContext(ctx, "localhost:8081", grpc.WithCodec(proxy.Codec()), grpc.WithInsecure())
	default:
		err = status.Errorf(codes.Unimplemented, "Unknown method")
	}
	return ctx, conn, err
}
