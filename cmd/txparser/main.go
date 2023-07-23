package main

import (
	"fmt"
	"os"

	svc "github.com/danielmbirochi/trustwallet-assignment/internal/service"
	"github.com/danielmbirochi/trustwallet-assignment/internal/txparser"
)

const (
	Endpoint     = "https://cloudflare-eth.com"
	InitialBlock = 0
)

func run(service svc.Parser) error {
	return nil
}

func main() {

	service := txparser.New(Endpoint, InitialBlock)

	if err := run(service); err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}
}
