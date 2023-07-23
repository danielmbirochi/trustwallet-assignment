package txparser

import (
	"encoding/json"
	"fmt"

	svc "github.com/danielmbirochi/trustwallet-assignment/internal/service"
	"github.com/danielmbirochi/trustwallet-assignment/internal/state"
	"github.com/danielmbirochi/trustwallet-assignment/pkg/ethclient"
)

type Blockscan struct {
	kvstate          state.KeyValueStorer
	clt              *ethclient.Client
	lastScannedBlock int
}

func NewScan(kvstate state.KeyValueStorer, clt *ethclient.Client, startAt int) Blockscan {
	return Blockscan{
		kvstate:          kvstate,
		clt:              clt,
		lastScannedBlock: startAt,
	}
}

// ParseTx converts an ethclient.Transaction into a the domain
// type service.Transaction.
func ParseTx(tx ethclient.Transaction) svc.Transaction {
	return svc.Transaction{
		ChainID:     tx.ChainID,
		BlockNumber: tx.BlockNumber,
		Hash:        tx.Hash,
		Nonce:       tx.Nonce,
		From:        tx.From,
		To:          tx.To,
		Value:       tx.Value,
		Gas:         tx.Gas,
		GasPrice:    tx.GasPrice,
		Input:       tx.Input,
	}
}

// GetCurrentBlock returns the last scanned block.
func (b *Blockscan) GetCurrentBlock() int {
	return b.lastScannedBlock
}

// Run starts the block scanning process. It will return the number
// of the last scanned block and an error if any. In case of no pending
// blocks to be scanned it will return 0.
func (b *Blockscan) Run() (int, error) {
	headBlock, err := b.clt.BlockNumber()
	if err != nil {
		fmt.Println("error querying head block number: ", err)
		return 0, err
	}

	nextBlock := nextBlock(b.lastScannedBlock, headBlock)
	if nextBlock == 0 {
		return 0, nil
	}

	txs, err := b.ScanBlock(nextBlock)
	if err != nil {
		fmt.Println("error scanning block: ", err)
		return 0, err
	}

	b.saveTxs(txs)
	b.lastScannedBlock = nextBlock

	return b.lastScannedBlock, nil
}

// ScanBlock retrieves the block with the given block number and
// returns a map containing the ingoing/outgoing transactions for
// the addresses subscribed.
func (b *Blockscan) ScanBlock(blockNumber int) (map[string][]svc.Transaction, error) {
	block, err := b.clt.BlockByNumber(blockNumber)
	if err != nil {
		fmt.Println("error querying block: ", err)
		return nil, err
	}

	newTxs := b.Pull(parseTxs(block.Transactions))
	if len(newTxs) == 0 {
		return nil, nil
	}

	return newTxs, nil
}

// Pull retrieves ingoing/outgoing transactions for the given list
// of address. This method does not check for internal transactions
// from smart contract executions.
func (b *Blockscan) Pull(txs []svc.Transaction) map[string][]svc.Transaction {
	result := make(map[string][]svc.Transaction)
	for _, tx := range txs {
		if exist, _ := b.kvstate.Has(tx.From); exist {
			result[tx.From] = append(result[tx.From], tx)
		}
		if exist, _ := b.kvstate.Has(tx.To); exist {
			result[tx.To] = append(result[tx.To], tx)
		}
	}
	return result
}

// saveTxs saves the given transactions into the key value store.
func (b *Blockscan) saveTxs(newTxs map[string][]svc.Transaction) {
	for address, txs := range newTxs {
		if err := b.kvstate.Put(address, encodeTxBatch(txs)); err != nil {
			fmt.Printf("error saving transactions: %v", err)
			continue
		}
	}
}

func encodeTxBatch(batch []svc.Transaction) [][]byte {
	var txs [][]byte
	for _, v := range batch {
		tx, err := json.Marshal(v)
		if err != nil {
			fmt.Printf("error marshaling transaction: %v", err)
			continue
		}
		txs = append(txs, tx)
	}
	return txs
}

// nextBlock returns the next block to be scanned. It will return
// 0 if there is any pending block to be scanned. If the last scanned
// block is 0 it will return the head block number.
func nextBlock(lastScannedBlock, headBlock int) int {
	if lastScannedBlock == headBlock {
		return 0
	}
	if lastScannedBlock == 0 {
		return headBlock
	}
	next := lastScannedBlock + 1
	return next
}

// parseTxs converts a list of ethclient.Transaction into a list of
// service.Transaction.
func parseTxs(txs []ethclient.Transaction) []svc.Transaction {
	transactions := make([]svc.Transaction, len(txs))
	for i, v := range txs {
		transactions[i] = ParseTx(v)
	}
	return transactions
}
