hello-service:
	go run testservice/hello-service/main.go

auth-service:
	go run testservice/auth-service/main.go

gateway:
	go run main.go

.PHONY: hello-service auth-service gateway