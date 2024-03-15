package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/maxuanquang/idm/internal/configs"
	"github.com/maxuanquang/idm/internal/wiring"
)

func startGRPCServer(configFilePath configs.ConfigFilePath, errorCh chan error) {
	grpcServer, cleanupFunc, err := wiring.InitializeGRPCServer(configFilePath)
	if err != nil {
		errorCh <- fmt.Errorf("error initializing gprc server: %w", err)
	}

	err = grpcServer.Start(context.Background())
	if err != nil {
		cleanupFunc()
		errorCh <- fmt.Errorf("error starting gprc server: %w", err)
	}
}

func startHTTPServer(configFilePath configs.ConfigFilePath, errorCh chan error) {
	httpServer, cleanupFunc, err := wiring.InitializeHTTPServer(configFilePath)
	if err != nil {
		errorCh <- fmt.Errorf("error initializing http gateway server: %w", err)
	}

	err = httpServer.Start(context.Background())
	if err != nil {
		cleanupFunc()
		errorCh <- fmt.Errorf("error starting http gateway server: %w", err)
	}
}

func main() {
	errorCh := make(chan error)

	configFilePath := configs.ConfigFilePath("config.yml")

	go startGRPCServer(configFilePath, errorCh)
	go startHTTPServer(configFilePath, errorCh)

	for {
		select {
		case err := <-errorCh:
			log.Fatal(err)
		default:
			time.Sleep(time.Second)
		}
	}
}
