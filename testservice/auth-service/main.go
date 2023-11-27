package main

import (
	"context"
	"fmt"
	"log"
	"net"
	auth "private-gateway/testservice/auth-service/pb/proto"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
}

type User struct {
	Username string
	Password string
}

var UserStore = map[string]User{}

const secretJWTKey = "secret"

func (s *server) Login(ctx context.Context, in *auth.LoginRequest) (*auth.LoginResponse, error) {

	user, ok := UserStore[in.Username]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}

	err := VerifyPassword(user.Password, in.Password)
	if err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	token, err := GenerateToken(time.Hour*24, user.Username, secretJWTKey)
	if err != nil {
		return nil, fmt.Errorf("could not generate token: %w", err)
	}

	return &auth.LoginResponse{Token: token}, nil
}

func (s *server) SignUp(ctx context.Context, in *auth.SignUpRequest) (*auth.SignUpResponse, error) {

	passwordHash, err := HashPassword(in.Password)
	if err != nil {
		return nil, fmt.Errorf("could not hash password: %w", err)
	}

	user := User{
		Username: in.Username,
		Password: passwordHash,
	}

	fmt.Println(user)

	UserStore[user.Username] = user

	return &auth.SignUpResponse{IsSuccess: true}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	s := grpc.NewServer()
	auth.RegisterAuthServiceServer(s, &server{})
	reflection.Register(s)

	log.Println("Serving gRPC on 0.0.0.0:8081")
	log.Fatal(s.Serve(lis))
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("could not hash password %w", err)
	}
	return string(hashedPassword), nil
}

func VerifyPassword(hashedPassword string, candidatePassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(candidatePassword))
}

func GenerateToken(ttl time.Duration, payload interface{}, secretJWTKey string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	now := time.Now().UTC()
	claims := token.Claims.(jwt.MapClaims)

	claims["sub"] = payload
	claims["exp"] = now.Add(ttl).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	tokenString, err := token.SignedString([]byte(secretJWTKey))

	if err != nil {
		return "", fmt.Errorf("generating JWT Token failed: %w", err)
	}

	return tokenString, nil
}
