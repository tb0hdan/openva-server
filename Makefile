all: gen lint build

lint:
	@golangci-lint run

regen:
	@echo 'module github.com/tb0hdan/openva-server' > ./go.mod
	@rm -f ./go.sum
	@go mod why

build:
	@go build -v -x -ldflags="-s -w" .

gen:
	@protoc -I api/ api/service.proto --go_out=plugins=grpc:api
