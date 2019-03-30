package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"strings"

	"google.golang.org/grpc"

	"github.com/tb0hdan/openva-server/api"
	"github.com/tb0hdan/openva-server/auth"
	"github.com/tb0hdan/openva-server/grpcserver"
)

const (
	GRPCPort = ":50001"
	HTTPPort = ":50002"
	MusicDir = "./music"
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

func main() {
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
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
