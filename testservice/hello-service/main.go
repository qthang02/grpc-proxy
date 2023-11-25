package main

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	helloworld "private-gateway/testservice/hello-service/pb/proto"
)

type server struct {
}

func (s *server) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		log.Println("md: ", md)

		if val, exists := md["authorization"]; exists {
			fmt.Println("authorization: ", val)

			token := val[0]

			sub, err := ValidateToken(token, "secret")
			if err != nil {
				return &helloworld.HelloReply{Message: err.Error()}, nil
			}

			name, ok := sub.(string)
			if !ok {
				return &helloworld.HelloReply{Message: "invalid sub"}, nil
			}

			return &helloworld.HelloReply{Message: "Hello " + name + "!"}, nil
		}
	}

	return &helloworld.HelloReply{Message: "invalid token"}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	s := grpc.NewServer()
	helloworld.RegisterGreeterServer(s, &server{})
	reflection.Register(s)

	log.Println("Serving gRPC on 0.0.0.0:8080")
	log.Fatal(s.Serve(lis))
}

func ValidateToken(token string, signedJWTKey string) (interface{}, error) {
	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}

		return []byte(signedJWTKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalidate token: %w", err)
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return nil, fmt.Errorf("invalid token claim")
	}

	return claims["sub"], nil
}
