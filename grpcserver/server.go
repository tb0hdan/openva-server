package grpcserver

import (
	"context"
	"errors"
	"fmt"
	"github.com/tb0hdan/openva-server/auth"
	"github.com/tb0hdan/openva-server/fileutils"
	"io"
	"io/ioutil"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	speech "cloud.google.com/go/speech/apiv1"
	"google.golang.org/api/option"
	"google.golang.org/api/transport"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
	"google.golang.org/grpc/peer"

	"github.com/tb0hdan/openva-server/api"
	"github.com/tb0hdan/openva-server/library"
	"github.com/tb0hdan/openva-server/netutils"
	"github.com/tb0hdan/openva-server/stt"
	"github.com/tb0hdan/openva-server/tts"
)

type NodeState struct {
	PlayerState       api.PlayerStateMessage
	SystemInformation api.SystemInformationMessage
}
type GRPCServer struct {
	// UUID - NodeState
	NodeStates        map[string]*NodeState
	MusicDir          string
	HTTPServerAddress string
	Authenticator *auth.Authenticator
}

func (s *GRPCServer) TTSStringToMP3(ctx context.Context, request *api.TTSRequest) (reply *api.TTSReply, err error) {
	cacheDir := path.Join("cache", "tts")
	err = os.MkdirAll(cacheDir, 0755)
	if err != nil {
		fmt.Println("Cache dir exists")
	}

	cachedFile := path.Join(cacheDir, strings.ToLower(strings.Replace(request.Text, " ", "_", -1)+".mp3"))
	_, err = os.Open(cachedFile)
	if os.IsNotExist(err) {
		fname := tts.Say(request.Text)
		err = fileutils.MoveFile(fname, cachedFile)
		if err != nil {
			return nil, err
		}
	}

	result, err := ioutil.ReadFile(cachedFile)
	if err != nil {
		return nil, err
	}
	reply = &api.TTSReply{MP3Response: result}
	return
}

func (s *GRPCServer) STT(stream api.OpenVAService_STTServer) (err error) {
	fmt.Println("Send config...")
	ctx := stream.Context()

	speechStream := getStream()
	defer speechStream.CloseSend()

	go func() {

		for {

			resp, err := speechStream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Println("Cannot speech stream results: %v", err)
				break
			}

			replies := stt.GoogleSTTToOpenVASTT(resp)
			stream.Send(&replies)
		}

	}()

	for {

		// exit if context is done
		// or continue
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		req, err := stream.Recv()
		if err == io.EOF {
			log.Println("completed")
			break
		}
		if err != nil {
			log.Println(err)
			break
		}

		if err = speechStream.Send(&speechpb.StreamingRecognizeRequest{
			StreamingRequest: &speechpb.StreamingRecognizeRequest_AudioContent{
				AudioContent: req.STTBuffer,
			},
		}); err != nil {
			log.Printf("Could not send audio: %v", err)
		}

	}

	return
}

func (s *GRPCServer) HeartBeat(stream api.OpenVAService_HeartBeatServer) (err error) {
	log.Println("HeartBeat stream started...")
	ctx := stream.Context()
	for {

		// exit if context is done
		// or continue
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		req, err := stream.Recv()
		if err == io.EOF {
			log.Println("HeartBeat stream completed...")
			break
		}
		if err != nil {
			log.Println(err)
			continue
		}
		// Update local representaion
		s.NodeStates[req.SystemInformation.SystemUUID] = &NodeState{
			PlayerState:       *req.PlayerState,
			SystemInformation: *req.SystemInformation,
		}

		// Send message back
		// To be used as a latency measurement later
		err = stream.Send(req)
		if err != nil {
			log.Println(err)
			continue
		}
		fmt.Println(req)
	}
	return
}

func (s *GRPCServer) ClientConfig(ctx context.Context, request *api.ClientMessage) (reply *api.ClientConfigMessage, err error) {
	reply = &api.ClientConfigMessage{
		Locale: &api.LocaleMessage{
			LocaleName:                "en-US",
			LocaleLanguage:            "English",
			VolumeMessage:             "volume",
			PauseMessage:              "pause",
			ResumeMessage:             "resume",
			StopMessage:               "stop",
			NextMessage:               "next",
			PreviousMessage:           "previous",
			RebootMessage:             "reboot",
			CouldNotUnderstandMessage: "I could not understand you",
		},
	}
	return
}

func (s *GRPCServer) HandleServerSideCommand(ctx context.Context, request *api.TTSRequest) (reply *api.OpenVAServerResponse, err error) {
	var (
		textResponse = "Unknown command"
		isError      = false
		noCmdMatches = false
		items        = make([]*api.LibraryItem, 0)
	)
	peerInfo, ok := peer.FromContext(ctx)
	if !ok {
		log.Println("Peer not ok")
		return nil, errors.New("Go away")
	}

	if peerInfo.Addr.Network() != "tcp" {
		log.Println("Peer tried using something else than TCP")
		return nil, errors.New("Go away. No, really")
	}

	serverIP := netutils.ServerIPForClientHostPort(peerInfo.Addr.String())

	localLibrary := library.LocalLibrary{
		MusicDir:          s.MusicDir,
		HTTPServerAddress: s.HTTPServerAddress,
	}

	log.Println(peerInfo.Addr.String())


	token, err := s.Authenticator.GetTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}

	cmd := request.Text
	first := strings.ToLower(strings.Split(cmd, " ")[0])
	switch first {
	case "play":
		textResponse, isError, noCmdMatches, items = handlePlayCommand(cmd, token, serverIP, s)
		if noCmdMatches {
			// unknown play command -> forward
			textResponse, isError, items = PlayForward(cmd, token)
		}
	case "shuffle":
		textResponse = "Shuffling your library"
		items, err = localLibrary.Library("", token, serverIP)
		if err != nil {
			isError = true
		}
	default:
		// 3rd-party tools like AVS and GHA
		textResponse, isError, items = UnknownCmdForward(cmd, token)
	}


	reply = &api.OpenVAServerResponse{
		TextResponse: textResponse,
		IsError:      isError,
		Items:        items,
	}
	return
}

func NewGRPCServer(MusicDir, HTTPServerAddress string, authenticator *auth.Authenticator) (s *GRPCServer) {
	s = &GRPCServer{
		MusicDir:          MusicDir,
		HTTPServerAddress: HTTPServerAddress,
		Authenticator: authenticator,
	}
	s.NodeStates = make(map[string]*NodeState)
	return s
}

var PlayRegs = map[string]func(cmd, token, serverIP string, srv *GRPCServer) (textResponse string, isError bool, items []*api.LibraryItem){
	`^play (.*) from my library$`: handlePlayLibraryCommand,
	`^play some music by (.*)$`:   handlePlayLibraryCommand,
	`^play (.*) by (.*)$`:         handlePlayLibraryCommand,
}

func handlePlayLibraryCommand(what, token, serverIP string, srv *GRPCServer) (textResponse string, isError bool, items []*api.LibraryItem) {
	var err error
	textResponse = fmt.Sprintf("Playing %s from your library", what)
	localLibrary := &library.LocalLibrary{
		MusicDir:          srv.MusicDir,
		HTTPServerAddress: srv.HTTPServerAddress,
	}
	items, err = localLibrary.Library(what, token, serverIP)
	if err != nil {
		isError = true
	}
	if len(items) == 0 {
		textResponse = fmt.Sprintf("I could not find %s", what)
	}
	return
}

func handlePlayCommand(cmd, token, serverIP string, srv *GRPCServer) (textResponse string, isError bool, noCmdMatches bool, items []*api.LibraryItem) {
	matches := 0
	for reg, fn := range PlayRegs {
		var what string
		re := regexp.MustCompile(reg)
		submatch := re.FindStringSubmatch(strings.ToLower(cmd))
		for i := len(submatch) - 1; i > 0; i-- {
			if i > 1 {
				what += strings.TrimSpace(submatch[i]) + " - "
			} else {
				what += strings.TrimSpace(submatch[i])
			}
		}
		what = strings.TrimSpace(what)
		log.Println(what)
		if len(what) == 0 {
			// Command starts with play but didn't match regexp
			continue
		}
		matches++
		textResponse, isError, items = fn(what, token, serverIP, srv)
		break
	}
	if matches == 0 {
		noCmdMatches = true
	}
	return
}

func getStream() (stream speechpb.Speech_StreamingRecognizeClient) {
	// connect to Google for a set duration to avoid running forever
	// and charge the user a lot of money.
	runDuration := 240 * time.Second
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


func PlayForward(cmd, token string) (textResponse string, isError bool, items []*api.LibraryItem){
	log.Debug(cmd, token)
	return
}

func UnknownCmdForward(cmd, token string) (textResponse string, isError bool, items []*api.LibraryItem) {
	log.Debug(cmd, token)
	return
}
