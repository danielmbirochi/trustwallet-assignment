package main

import (
	"context"
	"fmt"
	"os"
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

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}
}
