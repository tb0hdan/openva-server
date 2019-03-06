package main

import (
	"context"
	"github.com/tb0hdan/openva-server/api"
	"github.com/tb0hdan/openva-server/tts"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net"
)

const port = ":50001"

type server struct {

}

func (s *server) TTSStringToMP3(ctx context.Context, request *api.TTSRequest) (reply *api.TTSReply, err error){
	fileName := tts.Say(request.Text)
	result, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	reply.MP3Response = result
	return
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	api.RegisterOpenVAServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
