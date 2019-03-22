all: gen lint build

lint:
	@golangci-lint run

build:
	@go build -v -x -ldflags="-s -w" .

gen:
	@protoc -I api/ api/service.proto --go_out=plugins=grpc:api
