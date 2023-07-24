package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/danielmbirochi/trustwallet-assignment/internal/txparser"
)

const (
	Endpoint = "https://cloudflare-eth.com"

	// DefaultInitialBlock will start scanning from the latest block.
	DefaultInitialBlock = 0
)

var (
	ScanInterval = time.Second * 10
)

func main() {
	if err := run(); err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	initialBlock := flag.Int("block", DefaultInitialBlock, "block number to start scanning from")
	flag.Parse()

	service := txparser.New(ctx, Endpoint, *initialBlock)
	service.StartScan(ScanInterval)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// ========================================
	// Initialize the CLI and start the main loop
	{
		reader := bufio.NewReader(os.Stdin)

	Exit:
		for {
			select {
			case sig := <-shutdown:
				fmt.Println("shutdown started - received signal: ", sig)
				cancel()
				break Exit
			default:
				fmt.Print("> ")
				input, err := reader.ReadString('\n')
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					continue
				}

				input = strings.TrimSuffix(input, "\n")
				args := strings.Split(input, " ")

				if len(args) == 0 {
					continue
				}

				operation := args[0]

				if operation == "exit" {
					break Exit
				}

				if operation == "help" {
					help()
					continue
				}

				if operation == "stats" {
					fmt.Println("Current block:", service.GetCurrentBlock())
					fmt.Println()
					continue
				}

				if len(args) < 2 {
					help()
					continue
				}

				switch operation {
				case "subscribe":
					address := args[1]
					if ok := service.Subscribe(address); !ok {
						fmt.Fprintln(os.Stderr, err)
					}
					fmt.Printf("Address [%s] subscribed successfully\n", address)
					fmt.Println()
				case "transactions":
					address := args[1]
					txs := service.GetTransactions(address)
					fmt.Println("Transactions:")
					for _, tx := range txs {
						fmt.Printf("%+v\n", tx)
					}
					fmt.Println()
				}
			}
		}
	}

	ticker := time.NewTicker(time.Second * 2)
	<-ticker.C
	fmt.Println("shutdown completed")

	return nil
}

func help() {
	fmt.Println("Usage: <operation> <input>")
	fmt.Println("Available commands:")
	fmt.Println("  subscribe <address>")
	fmt.Println("  transactions <address>")
	fmt.Println("  stats")
	fmt.Println("  exit")
	fmt.Println("  help")
	fmt.Println()
}
