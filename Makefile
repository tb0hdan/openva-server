all: gen build

build:
	@go build -v -x -ldflags="-s -w" .
	@strip ./openva-server

gen:
	@protoc -I api/ api/service.proto --go_out=plugins=grpc:api
