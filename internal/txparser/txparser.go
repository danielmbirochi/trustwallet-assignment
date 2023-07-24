package txparser

import (
	"context"
	"encoding/json"
	"fmt"

	svc "github.com/danielmbirochi/trustwallet-assignment/internal/service"
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

func (s *Service) Subscribe(address string) bool {
	if err := s.kvstate.Put(address, [][]byte{}); err != nil {
		fmt.Printf("error subscribing address: %v", err)
		return false
	}
	return true
}

func (s *Service) GetTransactions(address string) []svc.Transaction {
	txs, err := s.kvstate.Get(address)
	if err != nil {
		fmt.Printf("error getting transactions: %v", err)
		return nil
	}
	transactions := make([]svc.Transaction, len(txs))
	for i, v := range txs {
		var tx svc.Transaction
		if err := json.Unmarshal(v, &tx); err != nil {
			fmt.Printf("error unmarshaling transaction: %v", err)
			continue
		}
		transactions[i] = tx
	}
	return transactions
}
