package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/dhowden/tag"
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
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"cloud.google.com/go/speech/apiv1"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"

	"google.golang.org/grpc/peer"
)

const (
	GRPCPort = ":50001"
	HTTPPort = ":50002"
	MusicDir = "./music"
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

func (s *server) STT(stream api.OpenVAService_STTServer) (err error) {
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

			replies := GoogleSTTToOpenVASTT(resp)
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

func (s *server) HeartBeat(stream api.OpenVAService_HeartBeatServer) (err error) {
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
			break
		}
		// Send message back
		// To be used as a latency measurement later
		err = stream.Send(req)
		if err != nil {
			log.Println(err)
			break
		}
		fmt.Println(req)
	}
	return
}

func (s *server) ClientConfig(ctx context.Context, request *api.ClientMessage) (reply *api.ClientConfigMessage, err error) {
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

func (s *server) HandleServerSideCommand(ctx context.Context, request *api.TTSRequest) (reply *api.OpenVAServerResponse, err error) {
	var (
		textResponse = "Unknown command"
		isError      = false
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

	log.Println(peerInfo.Addr.String())
	token := "12345"
	cmd := request.Text
	first := strings.ToLower(strings.Split(cmd, " ")[0])
	switch first {
	case "play":
		textResponse, isError, items = handlePlayCommand(cmd, token)
	case "shuffle":
		textResponse = "Shuffling your library"
		items, err = Library("", token)
		if err != nil {
			isError = true
		}
	default:
		// 3rd-party tools like AVS and GHA
		// ...
	}
	reply = &api.OpenVAServerResponse{
		TextResponse: textResponse,
		IsError:      isError,
		Items:        items,
	}
	return
}

var PlayRegs = map[string]func(cmd, token string) (textResponse string, isError bool, items []*api.LibraryItem){
	`^play (.*) from my library$`: handlePlayLibraryCommand,
	`^play some music by (.*)$`:   handlePlayLibraryCommand,
	`^play (.*) by (.*)$`:         handlePlayLibraryCommand,
}

func handlePlayLibraryCommand(what, token string) (textResponse string, isError bool, items []*api.LibraryItem) {
	var err error
	textResponse = fmt.Sprintf("Playing %s from your library", what)
	items, err = Library(what, token)
	if err != nil {
		isError = true
	}
	if len(items) == 0 {
		textResponse = fmt.Sprintf("I could not find %s", what)
	}
	return
}

func handlePlayCommand(cmd, token string) (textResponse string, isError bool, items []*api.LibraryItem) {
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
		textResponse, isError, items = fn(what, token)
		break
	}
	return
}

func GoogleSTTToOpenVASTT(resp *speechpb.StreamingRecognizeResponse) (response api.StreamingRecognizeResponse) {
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

	response = api.StreamingRecognizeResponse{
		Results:         results,
		SpeechEventType: api.StreamingRecognizeResponse_SpeechEventType(resp.SpeechEventType),
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

func Library(criteria, token string) (libraryItems []*api.LibraryItem, err error) {
	dir, err := filepath.EvalSymlinks(MusicDir)
	if err != nil {
		log.Fatal(err)
	}
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !strings.HasSuffix(path, ".mp3") {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}

		artist := ""
		album := ""
		track := ""

		m, err := tag.ReadFrom(file)
		if err != nil {
			log.Println(path, err)

		} else {
			artist = strings.Map(fixUTF, m.Artist())
			album = strings.Map(fixUTF, m.Album())
			track = strings.Map(fixUTF, m.Title())
		}

		if !libraryFilterPassed(criteria, artist, album, track, pathWords(path)) {
			return nil
		}

		escapedPath := ""
		for _, r := range strings.Split(strings.TrimPrefix(path, dir), "/") {
			escapedPath += "/" + url.PathEscape(r)
		}

		if strings.HasPrefix(escapedPath, "//") {
			escapedPath = strings.TrimPrefix(escapedPath, "/")
		}

		item := &api.LibraryItem{
			URL:    "http://localhost" + HTTPPort + "/music" + escapedPath + fmt.Sprintf("?token=%s", token),
			Artist: artist,
			Album:  album,
			Track:  track,
		}
		libraryItems = append(libraryItems, item)

		return nil
	})

	return
}

func libraryFilterPassed(criteria string, args ...string) bool {
	var (
		artist       string
		//album        string
		track        string
		searchArtist string
		searchTrack  string
	)
	if len(criteria) == 0 {
		return true
	}

	if len(args) == 0 {
		return true
	}
	criteria = strings.ToLower(criteria)

	if len(strings.Split(criteria, " - ")) >= 2 {
		searchArtist = strings.TrimSpace(strings.Split(criteria, " - ")[0])
		searchTrack = strings.TrimSpace(strings.Split(criteria, " - ")[1])
	}

	// artist, album, track
	if len(args) > 3 {
		artist = strings.TrimSpace(args[0])
		// album = args[1]
		track = strings.TrimSpace(args[2])
	}

	for _, arg := range args {
		arg = strings.ToLower(arg)
		if len(arg) > 0 && strings.Contains(arg, criteria) {
			return true
		}
		// Special case: Artist Name - Track Name
		if searchArtist == "" || searchTrack == "" {
			continue
		}
		if artist == "" || track == "" {
			continue
		}
		if strings.ToLower(searchArtist) != strings.ToLower(artist) {
			continue
		}
		if strings.ToLower(searchTrack) == strings.ToLower(track) {
			return true

		}
	}
	return false
}

func pathWords(path string) (newString string) {
	re := regexp.MustCompile(`[/|_|-|-|(|)|\.]`)
	for _, str := range strings.Split(re.ReplaceAllString(path, " "), " ") {
		if strings.TrimSpace(str) == "" {
			continue
		}
		newString += " " + str
	}
	return
}

// https://stackoverflow.com/questions/20401873/remove-invalid-utf-8-characters-from-a-string-go-lang
func fixUTF(r rune) rune {
	if r == utf8.RuneError {
		return -1
	}
	return r
}

func VerifyToken(token string) bool {
	log.Println("Token: ", token)
	return true
}

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		token_values := values["token"]
		tokenValid := false
		if len(token_values) > 0 {
			token := token_values[0]
			tokenValid = VerifyToken(token)
		}
		if !tokenValid {
			http.Error(w, "403 Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func NoIndexMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	handler := http.NewServeMux()

	dir, err := filepath.EvalSymlinks(MusicDir)

	if err != nil {
		log.Fatal("Please create ./music symlink")
	}

	fs := http.FileServer(http.Dir(dir))

	handler.Handle("/music/", AuthenticationMiddleware(
		NoIndexMiddleware(http.StripPrefix("/music/", fs))),
	)

	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Nothing here")
	})

	srv := http.Server{Addr: HTTPPort, Handler: handler}
	go func() {
		log.Printf("Library server started at %s\n", HTTPPort)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("failed to serve http: %v", err)
		}
	}()

	lis, err := net.Listen("tcp", GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	api.RegisterOpenVAServiceServer(s, &server{})

	log.Printf("gRPC server started at %s\n", GRPCPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
