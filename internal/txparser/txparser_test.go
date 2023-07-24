package txparser_test

import (
	"context"
	"testing"

	svc "github.com/danielmbirochi/trustwallet-assignment/internal/service"
	"github.com/danielmbirochi/trustwallet-assignment/internal/txparser"
	"github.com/danielmbirochi/trustwallet-assignment/pkg/ethclient"
)

// Success and failure markers.
const (
	Success = "\u2713"
	Failed  = "\u2717"
)

func TestTxParserService(t *testing.T) {

	var (
		Endpoint     = "https://cloudflare-eth.com"
		InitialBlock = 17758854
	)

	ethclt := ethclient.New(Endpoint)
	dataset := makeSampleDataset(t, ethclt, InitialBlock, InitialBlock+30)

	service := txparser.New(context.Background(), Endpoint, InitialBlock)

	t.Run("GetCurrentBlock", func(t *testing.T) {
		testID := 0
		if service.GetCurrentBlock() != InitialBlock {
			t.Fatalf("\t%s\tTest %d:\tShould be able to return the last scanned block : Expected: %d. Got: %d ", Failed, testID, InitialBlock, service.GetCurrentBlock())
		}
		t.Logf("\t%s\tTest %d:\tShould be able to return the last scanned block", Success, testID)
	})

	t.Run("Subscribe", func(t *testing.T) {
		testID := 1
		for addr := range dataset {
			if !service.Subscribe(addr) {
				t.Fatalf("\t%s\tTest %d:\tShould be able to subscribe to an address", Failed, testID)
			}
		}
		t.Logf("\t%s\tTest %d:\tShould be able to subscribe to an address", Success, testID)
	})

	t.Run("ScanBlock", func(t *testing.T) {
		testID := 2
		entries, err := service.ScanBlock(InitialBlock)
		if err != nil {
			t.Fatalf("\t%s\tTest %d:\tShould be able to scan a block: %v", Failed, testID, err)
		}
		if len(entries) == 0 {
			t.Fatalf("\t%s\tTest %d:\tShould be able to pull txs from the given block", Failed, testID)
		}
		service.SaveTxs(entries)
		t.Logf("\t%s\tTest %d:\tShould be able to pull txs from the given block", Success, testID)
	})

	t.Run("TestBlockscan", func(t *testing.T) {
		testID := 3
		for i := InitialBlock + 1; i <= InitialBlock+30; i++ {
			entries, err := service.ScanBlock(i)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to scan a block: %v", Failed, testID, err)
			}
			if len(entries) == 0 {
				t.Fatalf("\t%s\tTest %d:\tShould be able to pull txs from the given block", Failed, testID)
			}
			service.SaveTxs(entries)
		}
		t.Logf("\t%s\tTest %d:\tShould be able to pull txs from the given block", Success, testID)
	})

	t.Run("GetTransactions", func(t *testing.T) {
		testID := 4
		for addr, txs := range dataset {
			txsFromService := service.GetTransactions(addr)
			if len(txsFromService) != len(txs) {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve the expected amount of transactions for an address: Expected: %d. Got: %d", Failed, testID, len(txs), len(txsFromService))
			}
		}
		t.Logf("\t%s\tTest %d:\tShould be able to retrieve the expected amount of transactions for an address. %d addresses scanned in 30 blocks", Success, testID, len(dataset))
	})
}

func makeSampleDataset(t *testing.T, ethclt *ethclient.Client, initialBlock, finalBlock int) map[string][]svc.Transaction {
	dataset := make(map[string][]svc.Transaction)
	for i := initialBlock; i <= finalBlock; i++ {
		block, err := ethclt.BlockByNumber(i)
		if err != nil {
			t.Fatal("error generating sample data: ", err)
		}
		for _, tx := range block.Transactions {
			dataset[tx.To] = append(dataset[tx.To], txparser.ParseTx(tx))
			dataset[tx.From] = append(dataset[tx.From], txparser.ParseTx(tx))
		}
	}
	return dataset
}
