package main

import (
	"context"
	"fmt"
	"github.com/tb0hdan/openva-server/api"
	"github.com/tb0hdan/openva-server/tts"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path"
	"strings"
)

const port = ":50001"

type server struct {
}

func (s *server) TTSStringToMP3(ctx context.Context, request *api.TTSRequest) (reply *api.TTSReply, err error) {
	cacheDir := path.Join("cache", "tts")
	err = os.MkdirAll(cacheDir, 0755)
	if err != nil {
		fmt.Println("Cache dir exists")
	}

	cachedFile := path.Join(cacheDir, strings.ToLower(strings.Replace(request.Text, " ", "_", -1)+".mp3"))
	_, err = os.Open(cachedFile)
	if os.IsNotExist(err) {
		fname := tts.Say(request.Text)
		os.Rename(fname, cachedFile)
	}

	result, err := ioutil.ReadFile(cachedFile)
	if err != nil {
		return nil, err
	}
	reply = &api.TTSReply{MP3Response: result}
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
