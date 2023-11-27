hello-service:
	go run testserver/hello-service/main.go

auth-service:
	go run testserver/auth-service/main.go

gateway:
	go run main.go

.PHONY: hello-service auth-service gateway