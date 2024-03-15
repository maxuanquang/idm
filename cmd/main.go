package main

import (
	"context"
	"log"

	"github.com/maxuanquang/idm/internal/configs"
	"github.com/maxuanquang/idm/internal/wiring"
)

func main() {
	configFilePath := configs.ConfigFilePath("config.yml")

	grpcServer, cleanupFunc, err := wiring.InitializeGRPCServer(configFilePath)
	if err != nil {
		log.Fatalf("error initializing gprc server: %v", err)
	}

	err = grpcServer.Start(context.Background())
	if err != nil {
		cleanupFunc()
		log.Fatalf("error starting gprc server: %v", err)
	}
}
