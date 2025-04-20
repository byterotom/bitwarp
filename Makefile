.PHONY: build proto

build:
	docker build -t bitwarp-node:latest -f Dockerfile.node .
	docker build -t bitwarp-tracker:latest -f Dockerfile.tracker .

proto:
	@protoc --go_out=. --go-grpc_out=. proto/**/*.proto
