all: gen build

build:
	@go build -v -x -ldflags="-s -w" .

gen:
	@protoc -I api/ api/service.proto --go_out=plugins=grpc:api
