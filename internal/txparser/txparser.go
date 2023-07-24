package txparser

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	svc "github.com/danielmbirochi/trustwallet-assignment/internal"
	"github.com/danielmbirochi/trustwallet-assignment/internal/state"
	db "github.com/danielmbirochi/trustwallet-assignment/internal/state/inmemorydb"
	"github.com/danielmbirochi/trustwallet-assignment/pkg/ethclient"
)

type Service struct {
	kvstate state.KeyValueStorer
	*Blockscan
}

func New(ctx context.Context, endpoint string, startAtBlock int) *Service {
	datastore := db.New()
	ethclt := ethclient.New(endpoint)
	scan := NewScan(ctx, datastore, ethclt, startAtBlock)
	return &Service{
		kvstate:   datastore,
		Blockscan: scan,
	}
}

// Subscribe adds the address to the list of addresses to be scanned
// for transactions. Returns true if the address was added successfully.
// It will return true if the address is already subscribed.
func (s *Service) Subscribe(address string) bool {
	if err := s.kvstate.Put(strings.ToLower(address), [][]byte{}); err != nil {
		fmt.Println("error subscribing address: ", err)
		return false
	}
	return true
}

// GetTransactions return a list of scanned transactions for the given address.
func (s *Service) GetTransactions(address string) []svc.Transaction {
	txs, err := s.kvstate.Get(strings.ToLower(address))
	if err != nil {
		fmt.Println("error getting transactions: ", err)
		return nil
	}
	transactions := make([]svc.Transaction, len(txs))
	for i, v := range txs {
		var tx svc.Transaction
		if err := json.Unmarshal(v, &tx); err != nil {
			fmt.Println("error unmarshaling transaction: ", err)
			continue
		}
		transactions[i] = tx
	}
	return transactions
}
