package grpc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tb0hdan/openva-server/auth"
	httpclient "github.com/tb0hdan/openva-server/client/http"
	"github.com/tb0hdan/openva-server/fileutils"
	"github.com/tb0hdan/openva-server/node"
	"github.com/tb0hdan/openva-server/stringutil"

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

const ICouldNotProcessYourRequestEN = "I could not process your request"

type Server struct {
	// UUID - State
	NodeStates        *node.StateType
	MusicDir          string
	HTTPServerAddress string
	Authenticator     *auth.Authenticator
	Library           *library.Library
	IndexTicker       *time.Ticker
}

func (s *Server) TTSStringToMP3(ctx context.Context, request *api.TTSRequest) (reply *api.TTSReply, err error) {
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

func (s *Server) STT(stream api.OpenVAService_STTServer) (err error) { // nolint gocyclo
	log.Debug("STT Send config...")
	ctx := stream.Context()

	speechStream, cancelFunc, err := getStream()
	defer cancelFunc()
	if err != nil {
		return err
	}
	defer speechStream.CloseSend() // nolint errcheck

	go func() {

		for {

			resp, err := speechStream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Debugf("STT Cannot speech stream results: %v", err)
				break
			}

			replies := stt.GoogleSTTToOpenVASTT(resp)
			err = stream.Send(&replies)
			if err != nil {
				log.Debugf("STT Stream send error: %+v\n", err)
			}
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
			log.Debug("STT completed")
			break
		}
		if err != nil {
			log.Debug("STT", err)
			break
		}

		if err = speechStream.Send(&speechpb.StreamingRecognizeRequest{
			StreamingRequest: &speechpb.StreamingRecognizeRequest_AudioContent{
				AudioContent: req.STTBuffer,
			},
		}); err != nil {
			log.Debugf("STT Could not send audio: %v", err)
		}

	}

	return nil
}

func (s *Server) HeartBeat(stream api.OpenVAService_HeartBeatServer) (err error) {
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
		// Update local representation
		s.NodeStates.Set(req.SystemInformation.SystemUUID, &node.State{
			PlayerState:       *req.PlayerState,
			SystemInformation: *req.SystemInformation,
			LastUpdatedTS:     time.Now().Unix(),
		})

		// Send message back
		// To be used as a latency measurement later
		err = stream.Send(req)
		if err != nil {
			log.Println(err)
			continue
		}
		fmt.Println(req)

	}
	return nil
}

func (s *Server) ClientConfig(ctx context.Context, request *api.ClientMessage) (reply *api.ClientConfigMessage, err error) {
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

func (s *Server) HandleServerSideCommand(ctx context.Context, request *api.TTSRequest) (reply *api.OpenVAServerResponse, err error) {
	var (
		textResponse string
		isError,
		noCmdMatches bool
		items []*api.LibraryItem
	)
	peerInfo, ok := peer.FromContext(ctx)
	if !ok {
		log.Println("Peer not ok")
		return nil, errors.New("go away")
	}

	if peerInfo.Addr.Network() != "tcp" {
		log.Println("Peer tried using something else than TCP")
		return nil, errors.New("go away. No, really")
	}

	serverIP := netutils.ServerIPForClientHostPort(peerInfo.Addr.String())

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
		items, err = s.Library.Library("", token, serverIP)
		if err != nil {
			isError = true
		}
	default:
		// 3rd-party tools like AVS and GHA
		textResponse, isError, items = UnknownCmdForward(cmd, token)
	}

	return &api.OpenVAServerResponse{
		TextResponse: textResponse,
		IsError:      isError,
		Items:        items,
	}, nil
}

func (s *Server) StartPeriodicIndexUpdater() {
	s.IndexTicker = time.NewTicker(60 * time.Second)
	// Update index immediately
	go func() {
		s.Library.UpdateIndex()
	}()
	// Update index every minute
	go func() {
		for t := range s.IndexTicker.C {
			s.Library.UpdateIndex()
			log.Debug("Library update completed at ", t)
		}
	}()
}

func (s *Server) Stop() {
	s.IndexTicker.Stop()
	log.Debug("Stop called...")
}

func NewGRPCServer(musicDir, httpServerAddress string, authenticator *auth.Authenticator) (s *Server) {
	s = &Server{
		MusicDir:          musicDir,
		HTTPServerAddress: httpServerAddress,
		Authenticator:     authenticator,
		Library: &library.Library{
			MusicDir:          musicDir,
			HTTPServerAddress: httpServerAddress,
		},
	}
	s.NodeStates = node.New()
	return s
}

var PlayRegs = map[string]func(cmd, token, serverIP string, srv *Server) (textResponse string, isError bool, items []*api.LibraryItem){
	`^play (.*) from my library$`: handlePlayLibraryCommand,
	`^play some music by (.*)$`:   handlePlayLibraryCommand,
	`^play (.*) by (.*)$`:         handlePlayLibraryCommand,
}

func handlePlayLibraryCommand(what, token, serverIP string, srv *Server) (textResponse string, isError bool, items []*api.LibraryItem) {
	var err error
	textResponse = fmt.Sprintf("Playing %s from your library", what)

	for _, query := range []string{what, stringutil.SplitWords(what)} {
		items, err = srv.Library.Library(query, token, serverIP)
		if err != nil {
			isError = true
		}
		if len(items) > 0 {
			break
		}
	}

	if len(items) == 0 {
		textResponse = fmt.Sprintf("I could not find %s", what)
	}

	return
}

func handlePlayCommand(cmd, token, serverIP string,
	srv *Server) (textResponse string, isError, noCmdMatches bool, items []*api.LibraryItem) {
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
		log.Debug(what)
		if what == "" {
			// Command starts with play but didn't match regexp
			continue
		}
		matches++
		textResponse, isError, items = fn(what, token, serverIP, srv)
		break
	}
	if matches == 0 || len(items) == 0 {
		noCmdMatches = true
	}
	return
}

func getStream() (stream speechpb.Speech_StreamingRecognizeClient, cancelFunc context.CancelFunc, err error) {
	log.Debug("STT Getstream started....")
	// connect to Google for a set duration to avoid running forever
	// and charge the user a lot of money.

	runDuration := 240 * time.Second
	bgctx := context.Background()
	ctx, cancel := context.WithDeadline(bgctx, time.Now().Add(runDuration))

	conn, err := transport.DialGRPC(ctx,
		option.WithEndpoint("speech.googleapis.com:443"),
		option.WithScopes("https://www.googleapis.com/auth/cloud-platform"),
	)

	if err != nil {
		log.Printf("getStream DialGRPC %+v\n", err)
		return nil, cancel, err
	}

	defer conn.Close()

	client, err := speech.NewClient(ctx)
	if err != nil {
		log.Printf("getStream SpeachNewclient %+v\n", err)
		return nil, cancel, err
	}
	stream, err = client.StreamingRecognize(ctx)
	if err != nil {
		log.Printf("getStream StreamingRecognize %+v\n", err)
		return nil, cancel, err
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
		log.Printf("getStream Send %+v\n", err)
		return nil, cancel, err
	}
	return stream, cancel, nil
}

func PlayForward(cmd, token string) (textResponse string, isError bool, items []*api.LibraryItem) {
	items, err := httpclient.PlayForward(cmd, token)
	if err != nil {
		log.Error(err)
		return ICouldNotProcessYourRequestEN, true, nil
	}
	return
}

func UnknownCmdForward(cmd, token string) (textResponse string, isError bool, items []*api.LibraryItem) {
	items, err := httpclient.UnknownForward(cmd, token)
	if err != nil {
		log.Error(err)
		return ICouldNotProcessYourRequestEN, true, nil
	}
	return
}
