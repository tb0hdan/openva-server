package http

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/tb0hdan/openva-server/auth"
)

func NoIndexMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

type MusicHTTPServer struct {
	MusicDir     string
	AuthFileName string
	HTTPAddress  string
}

func NewMusicHTTPServer(musicDir, authFileName, address string) *MusicHTTPServer {
	return &MusicHTTPServer{
		MusicDir:     musicDir,
		AuthFileName: authFileName,
		HTTPAddress:  address,
	}
}

func (mh *MusicHTTPServer) Run() {
	handler := http.NewServeMux()

	dir, err := filepath.EvalSymlinks(mh.MusicDir)

	if err != nil {
		log.Fatal("Please create ./music symlink")
	}

	authenticator, err := auth.NewAuthenticator("../passwd")
	if err != nil {
		log.Fatalf("Auth file error: %s", err)
	}

	fs := http.FileServer(http.Dir(dir))

	handler.Handle("/music/", authenticator.AuthenticationMiddleware(
		NoIndexMiddleware(http.StripPrefix("/music/", fs))),
	)

	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Nothing here")
	})

	srv := http.Server{Addr: mh.HTTPAddress, Handler: handler}
	log.Printf("Library server started at %s\n", mh.HTTPAddress)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("failed to serve http: %v", err)
	}
	return
}
