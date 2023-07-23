package ethclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const (
	ApiVersion           = "2.0"
	GetBlocknumberMethod = "eth_blockNumber"
	GetBlockByNumber     = "eth_getBlockByNumber"
)

type Client struct {
	endpoint string
}

type RequestBody struct {
	Jsonrpc string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

type Block struct {
	Number       string        `json:"number"`
	Hash         string        `json:"hash"`
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	ChainID          string            `json:"chainId"`
	BlockNumber      string            `json:"blockNumber"`
	BlockHash        string            `json:"-"`
	Hash             string            `json:"hash"`
	Nonce            string            `json:"nonce"`
	From             string            `json:"from"`
	To               string            `json:"to"`
	Value            string            `json:"value"`
	Gas              string            `json:"gas"`
	GasPrice         string            `json:"gasPrice"`
	Input            string            `json:"input"`
	Type             string            `json:"-"`
	R                string            `json:"-"`
	S                string            `json:"-"`
	V                string            `json:"-"`
	TransactionIndex string            `json:"-"`
	AccessList       []AccessListEntry `json:"-"`
}

type AccessListEntry struct {
	Address     string   `json:"address"`
	StorageKeys []string `json:"storageKeys"`
}

func New(endpoint string) *Client {
	return &Client{
		endpoint: endpoint,
	}
}

// BlockNumber returns the current block number. It will call
// the eth_blockNumber method of the JSON-RPC API in the given endpoint.
func (c Client) BlockNumber() (int, error) {
	body, err := json.Marshal(makeRequestBody(GetBlocknumberMethod, []string{}))
	if err != nil {
		return 0, fmt.Errorf("error marshaling json: %v", err)
	}

	r, err := http.Post(c.endpoint, "application/json", bytes.NewReader(body))
	if err != nil {
		return 0, fmt.Errorf("error making request: %v", err)
	}
	if r.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("error response status code: %v", r.StatusCode)
	}

	var responseBody struct {
		Jsonrpc string `json:"jsonrpc"`
		ID      int    `json:"id"`
		Result  string `json:"result"`
	}
	if err := json.NewDecoder(r.Body).Decode(&responseBody); err != nil {
		return 0, fmt.Errorf("error decoding response body: %v", err)
	}

	blocknumber, err := strconv.ParseInt(responseBody.Result[2:], 16, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing response body: %v", err)
	}

	return int(blocknumber), nil
}

func (c Client) BlockByNumber(blocknumber int) (Block, error) {
	body, err := json.Marshal(makeRequestBody(GetBlockByNumber, []string{fmt.Sprintf("0x%x", blocknumber), "true"}))
	if err != nil {
		return Block{}, fmt.Errorf("error marshaling json: %v", err)
	}

	r, err := http.Post(c.endpoint, "application/json", bytes.NewReader(body))
	if err != nil {
		return Block{}, fmt.Errorf("error making request: %v", err)
	}
	if r.StatusCode != http.StatusOK {
		return Block{}, fmt.Errorf("error response status code: %v", r.StatusCode)
	}

	var responseBody struct {
		Jsonrpc string `json:"jsonrpc"`
		ID      int    `json:"id"`
		Result  Block  `json:"result"`
	}
	if err := json.NewDecoder(r.Body).Decode(&responseBody); err != nil {
		return Block{}, fmt.Errorf("error decoding response body: %v", err)
	}

	return responseBody.Result, nil
}

func makeRequestBody(method string, params interface{}) RequestBody {
	rand.Seed(time.Now().UnixNano())
	return RequestBody{
		Jsonrpc: ApiVersion,
		ID:      rand.Int(),
		Method:  method,
		Params:  params,
	}
}
