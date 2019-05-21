package server

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

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

// https://gobyexample.com/signals
func configureServerShutdown(stopFunc func()) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		log.Debug(sig)
		stopFunc()
	}()
}

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
	configureServerShutdown(func() {
		openVAServer.Stop()
		server.Stop()
	})

	api.RegisterOpenVAServiceServer(server, openVAServer)

	log.Printf("gRPC server started at %s\n", grpcPort)

	openVAServer.StartPeriodicIndexUpdater()

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
