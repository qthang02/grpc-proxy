hello:
	go run testservice/hello-service/main.go

auth:
	go run testservice/auth-service/main.go

gateway:
	go run main.go

.PHONY: hello-service auth-service gateway