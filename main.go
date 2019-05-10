package main

import (
	"flag"

	log "github.com/sirupsen/logrus"
	"github.com/tb0hdan/openva-server/server"
)

const (
	GRPCPort     = ":50001"
	HTTPPort     = ":50002"
	MusicDir     = "./music"
	AuthFileName = "../passwd"
)

var Debug = flag.Bool("debug", false, "Enable debug")

func main() {
	flag.Parse()

	if *Debug {
		log.SetLevel(log.DebugLevel)
	}

	server.Run(MusicDir, AuthFileName, HTTPPort, GRPCPort)
}
