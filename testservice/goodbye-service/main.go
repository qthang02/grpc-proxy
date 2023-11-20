package main

import (
	"context"
	goodbye "grpc-proxy/testservice/goodbye-service/pb/proto"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

func (s *server) SayBye(ctx context.Context, in *goodbye.ByeRequest) (*goodbye.ByeReply, error) {
	return &goodbye.ByeReply{Message: "Goodbye " + in.Name}, nil
}


func main() {
	lis, err := net.Listen("tcp", ":8000")

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	goodbye.RegisterGreeterServer(s, &server{})
	reflection.Register(s)

	log.Println("Starting server on port :8000")
	log.Fatal(s.Serve(lis))
}