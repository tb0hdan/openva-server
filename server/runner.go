package server

import (
	"log"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/spf13/viper"
	"github.com/tb0hdan/openva-server/api"
	"github.com/tb0hdan/openva-server/auth"
	"github.com/tb0hdan/openva-server/node"
	"github.com/tb0hdan/openva-server/server/http"
	"google.golang.org/grpc"

	grpc2 "github.com/tb0hdan/openva-server/server/grpc"
)

func Run(musicDir, authFileName, httpPort, grpcPort string) {
	httpServer := http.NewMusicHTTPServer(musicDir, authFileName, httpPort)
	go httpServer.Run()

	authenticator, err := auth.NewAuthenticator(authFileName)
	if err != nil {
		log.Fatalf("Auth file error: %s", err)
	}

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_auth.StreamServerInterceptor(authenticator.MyGRPCAuthFunction),
			grpc_recovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_auth.UnaryServerInterceptor(authenticator.MyGRPCAuthFunction),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)

	openVAServer := grpc2.NewGRPCServer(musicDir, httpPort, authenticator)
	api.RegisterOpenVAServiceServer(server, openVAServer)

	log.Printf("gRPC server started at %s\n", grpcPort)

	openVAServer.Library.UpdateIndex()

	viper.AutomaticEnv()

	nodeWatcher := node.NewNodeWatcherWithConnection(
		viper.GetString("LASTFM_API_KEY"),
		viper.GetString("LASTFM_API_SECRET"),
		viper.GetString("LASTFM_USERNAME"),
		viper.GetString("LASTFM_PASSWORD"),
	)

	go nodeWatcher.ServerWatcher(openVAServer.NodeStates)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
