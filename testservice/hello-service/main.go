package main

import (
	"context"
	helloworld "grpc-proxy/testservice/hello-service/pb/proto"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

func (s *server) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8001")

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	helloworld.RegisterGreeterServer(s, &server{})
	reflection.Register(s)

	log.Println("Starting server on port :8000")
	log.Fatal(s.Serve(lis))
}