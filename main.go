package main

import (
	"context"
	"fmt"
	"github.com/tb0hdan/openva-server/api"
	"github.com/tb0hdan/openva-server/tts"
	"google.golang.org/api/option"
	"google.golang.org/api/transport"
	"google.golang.org/grpc"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/speech/apiv1"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
)

const (
	port     = ":50001"
	HTTPPort = ":50002"
)

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

func (s *server) STT(srv api.OpenVAService_STTServer) (err error) {
	stream := getStream()

	for {

		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("Cannot stream results: %v", err)
			break
		}
		if err := resp.Error; err != nil {
			log.Println("Could not recognize: %v", err)
			break
		}

		results := make([]*api.StreamingRecognitionResult, 0)
		for _, res := range resp.Results {
			alternatives := make([]*api.SpeechRecognitionAlternative, 0)
			for _, alt := range res.Alternatives {
				words := make([]*api.WordInfo, 0)
				for _, word := range alt.Words {
					wrd := &api.WordInfo{
						StartTime: word.StartTime,
						EndTime:   word.EndTime,
						Word:      word.Word,
					}
					words = append(words, wrd)
				}

				alternative := &api.SpeechRecognitionAlternative{
					Transcript: alt.Transcript,
					Confidence: alt.Confidence,
					Words:      words,
				}
				alternatives = append(alternatives, alternative)
			}

			result := &api.StreamingRecognitionResult{
				Alternatives: alternatives,
				IsFinal:      res.IsFinal,
				Stability:    res.Stability,
			}
			results = append(results, result)
		}

		response := api.StreamingRecognizeResponse{
			Results:         results,
			SpeechEventType: api.StreamingRecognizeResponse_SpeechEventType(resp.SpeechEventType),
		}

		err = srv.Send(&response)
		if err != nil {
			log.Printf("%+v", err)
		}

	}

	return
}

func getStream() (stream speechpb.Speech_StreamingRecognizeClient) {
	// connect to Google for a set duration to avoid running forever
	// and charge the user a lot of money.
	runDuration := 70 * time.Second
	bgctx := context.Background()
	ctx, _ := context.WithDeadline(bgctx, time.Now().Add(runDuration))
	conn, err := transport.DialGRPC(ctx,
		option.WithEndpoint("speech.googleapis.com:443"),
		option.WithScopes("https://www.googleapis.com/auth/cloud-platform"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client, err := speech.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	stream, err = client.StreamingRecognize(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// Send the initial configuration message.
	if err := stream.Send(&speechpb.StreamingRecognizeRequest{
		StreamingRequest: &speechpb.StreamingRecognizeRequest_StreamingConfig{
			StreamingConfig: &speechpb.StreamingRecognitionConfig{
				Config: &speechpb.RecognitionConfig{
					// Uncompressed 16-bit signed little-endian samples (Linear PCM).
					Encoding:        speechpb.RecognitionConfig_LINEAR16,
					SampleRateHertz: 16000,
					LanguageCode:    "en-US",
				},
			},
		},
	}); err != nil {
		log.Fatal(err)
	}
	return
}

func main() {
	handler := http.NewServeMux()

	dir, err := filepath.EvalSymlinks("./music")

	if err != nil {
		log.Fatal("Please create ./music symlink")
	}

	fs := http.FileServer(http.Dir(dir))

	handler.Handle("/music/", http.StripPrefix("/music/", fs))

	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Nothing here")
	})

	srv := http.Server{Addr: HTTPPort, Handler: handler}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("failed to serve http: %v", err)
		}
	}()

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
