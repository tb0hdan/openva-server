package main

import (
	"context"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"google.golang.org/grpc"
	"github.com/shkh/lastfm-go/lastfm"
	"github.com/spf13/viper"

	"github.com/tb0hdan/openva-server/api"
	"github.com/tb0hdan/openva-server/auth"
	"github.com/tb0hdan/openva-server/grpcserver"
)

const (
	GRPCPort = ":50001"
	HTTPPort = ":50002"
	MusicDir = "./music"
)

var Debug = flag.Bool("debug", false, "Enable debug")

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
	flag.Parse()

	if *Debug {
		log.SetLevel(log.DebugLevel)
	}

	handler := http.NewServeMux()

	dir, err := filepath.EvalSymlinks(MusicDir)

	if err != nil {
		log.Fatal("Please create ./music symlink")
	}

	fs := http.FileServer(http.Dir(dir))

	handler.Handle("/music/", auth.AuthenticationMiddleware(
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
	server := grpc.NewServer()

	openVAServer := grpcserver.NewGRPCServer(MusicDir, HTTPPort)
	api.RegisterOpenVAServiceServer(server, openVAServer)

	log.Printf("gRPC server started at %s\n", GRPCPort)

	viper.AutomaticEnv()

	nodeWatcher := &NodeWatcher{
		LastFMAPIKey: viper.GetString("LASTFM_API_KEY"),
		LastFMAPISecret: viper.GetString("LASTFM_API_SECRET"),
		LastFMUsername: viper.GetString("LASTFM_USERNAME"),
		LastFMPassword: viper.GetString("LASTFM_PASSWORD"),
	}

	nodeWatcher.LoginLastFm()

	go nodeWatcher.ServerWatcher(openVAServer.NodeStates)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type NodeWatcher struct {
	LastFMAPIKey string
	LastFMAPISecret string
	LastFMUsername string
	LastFMPassword string
	lastFMLoggedIn bool
	lastFMAPI *lastfm.Api
}

func (nw *NodeWatcher) LoginLastFm() {
	lfmAPI := lastfm.New(nw.LastFMAPIKey, nw.LastFMAPISecret)

	err := lfmAPI.Login(nw.LastFMUsername, nw.LastFMPassword)
	if err != nil {
		fmt.Println(err)
		return
	}
	nw.lastFMLoggedIn = true
	nw.lastFMAPI = lfmAPI
}

func (nw *NodeWatcher) Scrobble(nowPlaying string) {
	var (
		artist,
		track string
	)
	if ! nw.lastFMLoggedIn {
		log.Debug("Scrobble exit, not logged in")
		return
	}
	if len(strings.Split(nowPlaying, " - ")) < 2 {
		artist = nowPlaying
		track = nowPlaying
	} else {
		artist = strings.Split(nowPlaying, " - ")[0]
		track = strings.Split(nowPlaying, " - ")[1]
	}
	p := lastfm.P{"artist": artist, "track": track}
	start := time.Now().Unix()

	p["timestamp"] = start
	_, err := nw.lastFMAPI.Track.Scrobble(p)
	if err != nil {
		log.Debug(err)
		return
	}

}

func (nw *NodeWatcher) SetNowPlaying(nowPlaying string) {
	var (
		artist,
		track string
	)
	if ! nw.lastFMLoggedIn {
		log.Debug("SetNowPlaying exit, not logged in")
		return
	}
	if len(strings.Split(nowPlaying, " - ")) < 2 {
		artist = nowPlaying
		track = nowPlaying
	} else {
		artist = strings.Split(nowPlaying, " - ")[0]
		track = strings.Split(nowPlaying, " - ")[1]
	}
	p := lastfm.P{"artist": artist, "track": track}
	_, err := nw.lastFMAPI.Track.UpdateNowPlaying(p)
	if err != nil {
		log.Debug(err)
		return
	}
}

func (nw *NodeWatcher) ClientWatcher(ctx context.Context, key string, nodeState map[string]*grpcserver.NodeState) {
	passed := 0
	previousTrack := nodeState[key].PlayerState.NowPlaying
	trackChanged := false

	if len(previousTrack) > 0 {
		nw.SetNowPlaying(previousTrack)
	}

	for {
		select {
		case <-ctx.Done():
			break
		default:
		}

		if nodeState[key].PlayerState.NowPlaying != previousTrack && len(nodeState[key].PlayerState.NowPlaying) > 0{
			previousTrack = nodeState[key].PlayerState.NowPlaying
			passed = 0
			trackChanged = true
			// track change
			nw.SetNowPlaying(previousTrack)
		}
		log.Println(key, nodeState[key].PlayerState.NowPlaying)
		if passed >= 35 && len(nodeState[key].PlayerState.NowPlaying) > 0 && trackChanged {
			log.Println("Can scrobble: ", nodeState[key].PlayerState.NowPlaying)
			nw.Scrobble(nodeState[key].PlayerState.NowPlaying)
			trackChanged = false
		}

		passed++
		time.Sleep(time.Second)
	}
}

func (nw *NodeWatcher) ServerWatcher(nodeStates map[string]*grpcserver.NodeState) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	oldmap := make(map[string]*grpcserver.NodeState)
	for {
		time.Sleep(time.Second)
		if len(nodeStates) == 0 {
			log.Println("Waiting for first client...")
			continue
		}
		for key := range nodeStates {
			if oldmap[key] == nil && nodeStates[key] != nil {
				oldmap[key] = nodeStates[key]
				log.Println("Got new client ", key)
				go nw.ClientWatcher(ctx, key, nodeStates)
			}
		}
	}
}
