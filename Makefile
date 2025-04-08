.PHONY: proto

proto:
	@protoc --go_out=. --go-grpc_out=. proto/*.proto

tracker:
	go run cmd/tracker/main.go