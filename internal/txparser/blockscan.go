package txparser

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	svc "github.com/danielmbirochi/trustwallet-assignment/internal"
	"github.com/danielmbirochi/trustwallet-assignment/internal/state"
	"github.com/danielmbirochi/trustwallet-assignment/pkg/ethclient"
)

type Scanner interface {
	// Run starts the block scanning process. It will return the number
	// of the last scanned block and an error if any. In case of no pending
	// blocks to be scanned it will return 0.
	Run() (int, error)

	// GetCurrentBlock returns the last scanned block.
	GetCurrentBlock() int
}

type Blockscan struct {
	ctx              context.Context
	kvstate          state.KeyValueStorer
	clt              *ethclient.Client
	lastScannedBlock int
	once             sync.Once
}

// ParseTx converts an ethclient.Transaction into a the domain
// type service.Transaction.
func ParseTx(tx ethclient.Transaction) svc.Transaction {
	return svc.Transaction{
		ChainID:     decodeHexString(tx.ChainID),
		BlockNumber: decodeHexString(tx.BlockNumber),
		Hash:        tx.Hash,
		Nonce:       decodeHexString(tx.Nonce),
		From:        tx.From,
		To:          tx.To,
		Value:       decodeHexString(tx.Value),
		Gas:         decodeHexString(tx.Gas),
		GasPrice:    decodeHexString(tx.GasPrice),
		Input:       tx.Input,
	}
}

func NewScan(ctx context.Context, kvstate state.KeyValueStorer, clt *ethclient.Client, startAt int) *Blockscan {
	fmt.Println("Blockscan set to start at block: ", startAt)
	return &Blockscan{
		ctx:              ctx,
		kvstate:          kvstate,
		clt:              clt,
		lastScannedBlock: startAt,
	}
}

// StartScan spawn a goroutine that will run the block scanning process
// at the given interval. It will stop the process when a signal is received
// on the shutdown channel. It will spawn only one goroutine for each Blockscan
// instance.
func (b *Blockscan) StartScan(interval time.Duration) {
	b.once.Do(func() {
		go func() {
			ticker := time.NewTicker(interval)
			for {
				select {
				case <-b.ctx.Done():
					ticker.Stop()
					fmt.Println("stopping blockscan")
					return
				case <-ticker.C:
					ticker.Stop()
					for scannedBlock, err := b.Run(); scannedBlock != 0 || err != nil; scannedBlock, err = b.Run() {
						if err != nil {
							fmt.Println(fmt.Errorf("error scanning block: %s", err))
							continue
						}
					}
					ticker.Reset(interval)
					fmt.Printf("last scanned block %d\n", b.GetCurrentBlock())
				}
			}
		}()
	})
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

	b.SaveTxs(txs)
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
func (b *Blockscan) SaveTxs(newTxs map[string][]svc.Transaction) {
	for address, txs := range newTxs {
		if err := b.kvstate.Put(address, encodeTxBatch(txs)); err != nil {
			fmt.Println("error saving transactions: ", err)
			continue
		}
	}
}

func encodeTxBatch(batch []svc.Transaction) [][]byte {
	var txs [][]byte
	for _, v := range batch {
		tx, err := json.Marshal(v)
		if err != nil {
			fmt.Println("error marshaling transaction: ", err)
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

func decodeHexString(hexStr string) *big.Int {
	hexStr = strings.TrimPrefix(hexStr, "0x")

	n := new(big.Int)
	n.SetString(hexStr, 16)
	return n
}
