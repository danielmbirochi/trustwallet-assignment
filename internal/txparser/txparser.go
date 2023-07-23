package txparser

import (
	"encoding/json"
	"fmt"

	svc "github.com/danielmbirochi/trustwallet-assignment/internal/service"
	"github.com/danielmbirochi/trustwallet-assignment/internal/state"
	"github.com/danielmbirochi/trustwallet-assignment/pkg/ethclient"
)

type TxParser struct {
	kvstate          state.KeyValueStorer
	clt              *ethclient.Client
	lastScannedBlock int
}

func New(kvstate state.KeyValueStorer, clt *ethclient.Client, startAt int) *TxParser {
	return &TxParser{
		kvstate:          kvstate,
		clt:              clt,
		lastScannedBlock: startAt,
	}
}

func (a *TxParser) GetCurrentBlock() int {
	return a.lastScannedBlock
}

func (a *TxParser) Subscribe(address string) bool {
	if err := a.kvstate.Put(address, [][]byte{}); err != nil {
		fmt.Printf("error subscribing address: %v", err)
		return false
	}
	return true
}

func (a *TxParser) GetTransactions(address string) []svc.Transaction {
	txs, err := a.kvstate.Get(address)
	if err != nil {
		fmt.Printf("error getting transactions: %v", err)
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
