package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danielmbirochi/trustwallet-assignment/internal/txparser"
)

const (
	Endpoint     = "https://cloudflare-eth.com"
	InitialBlock = 0
)

var (
	ScanInterval = time.Second * 10
)

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := txparser.New(ctx, Endpoint, InitialBlock)
	if !service.StartScan(ScanInterval) {
		return fmt.Errorf("error starting scan")
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	sig := <-shutdown
	fmt.Println("shutdown started - received signal: ", sig)
	cancel()

	ticker := time.NewTicker(time.Second * 2)
	<-ticker.C
	fmt.Println("shutdown completed")

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}
}
