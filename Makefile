hello-service:
	go run testserver/hello-service/main.go

auth-service:
	go run testserver/auth-service/main.go

gateway:
	go run main.go

test-auth-service-login:
	grpcurl -proto testservice/auth-service/proto/auth.proto -plaintext -d '{"username": "thang", "password": "123456"}' localhost:50052 auth.AuthService.Login

test-auth-service-signup:
	grpcurl -proto testservice/auth-service/proto/auth.proto -plaintext -d '{"username": "thang", "password": "123456"}' localhost:50052 auth.AuthService.SignUp

test-hello-service:
	grpcurl -H "authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDA3OTg3MTUsImlhdCI6MTcwMDcxMjMxNSwibmJmIjoxNzAwNzEyMzE1LCJzdWIiOiJ0aGFuZyJ9.62dDopqahZLu9DLNZRjHrAgu2ZhtYhgr7vxWYXdL9ZE" -proto testservice/hello-service/proto/hello_world.proto -plaintext -d '{"name": "thang"}' localhost:50052 helloworld.Greeter.SayHello

.PHONY: hello-service test-hello-service auth-service test-auth-service-login test-auth-service-signup gateway