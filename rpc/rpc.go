package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	ethereum "github.com/ledgerwatch/erigon"
	"github.com/ledgerwatch/erigon-lib/common"
	"github.com/ledgerwatch/erigon/core/types"
	rpcTypes "github.com/ledgerwatch/erigon/zk/rpcdaemon"
	zktypes "github.com/ledgerwatch/erigon/zk/types"
	"github.com/ledgerwatch/erigon/zkevm/jsonrpc/client"
)

type RealtimeClient struct {
	client ethClienter
	rpcUrl string
}

// RealtimeBlockNumber returns the number of the most recent block in real-time
func (c *RealtimeClient) RealtimeBlockNumber() (uint64, error) {
	response, err := client.JSONRPCCall(c.rpcUrl, "realtime_blockNumber")
	if err != nil {
		return 0, err
	}
	if response.Error != nil {
		return 0, fmt.Errorf("%d - %s", response.Error.Code, response.Error.Message)
	}

	return transHexToUint64(response.Result)
}

// RealtimeGetBlockTransactionCountByNumber returns the number of transactions in a block by number in real-time
func (c *RealtimeClient) RealtimeGetBlockTransactionCountByNumber(blockNumber uint64) (uint64, error) {
	response, err := client.JSONRPCCall(c.rpcUrl, "realtime_getBlockTransactionCountByNumber", blockNumber)
	if err != nil {
		return 0, err
	}
	if response.Error != nil {
		return 0, fmt.Errorf("%d - %s", response.Error.Code, response.Error.Message)
	}

	return transHexToUint64(response.Result)
}

// RealtimeGetTransactionByHash returns the information about a transaction requested by transaction hash in real-time
func (c *RealtimeClient) RealtimeGetTransactionByHash(txHash common.Hash, includeExtraInfo *bool) (rpcTypes.Transaction, error) {
	response, err := client.JSONRPCCall(c.rpcUrl, "realtime_getTransactionByHash", txHash, includeExtraInfo)
	if err != nil {
		return rpcTypes.Transaction{}, err
	}
	if response.Error != nil {
		return rpcTypes.Transaction{}, fmt.Errorf("%d - %s", response.Error.Code, response.Error.Message)
	}

	result := rpcTypes.Transaction{}
	err = json.Unmarshal(response.Result, &result)
	if err != nil {
		return rpcTypes.Transaction{}, err
	}

	return result, nil
}

// RealtimeGetTransactionByHash returns raw information about a transaction requested by transaction hash in real-time
func (c *RealtimeClient) RealtimeGetRawTransactionByHash(txHash common.Hash) ([]byte, error) {
	response, err := client.JSONRPCCall(c.rpcUrl, "realtime_getRawTransactionByHash", txHash)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, fmt.Errorf("%d - %s", response.Error.Code, response.Error.Message)
	}

	var result []byte
	err = json.Unmarshal(response.Result, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// RealtimeGetTransactionReceipt returns the receipt of a transaction by transaction hash in real-time
func (c *RealtimeClient) RealtimeGetTransactionReceipt(txHash common.Hash) (*types.Receipt, error) {
	response, err := client.JSONRPCCall(c.rpcUrl, "realtime_getTransactionReceipt", txHash)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, fmt.Errorf("%d - %s", response.Error.Code, response.Error.Message)
	}

	var result types.Receipt
	err = json.Unmarshal(response.Result, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// RealtimeGetInternalTransactions returns the internal transactions for a given transaction hash in real-time
func (c *RealtimeClient) RealtimeGetInternalTransactions(txHash common.Hash) ([]zktypes.InnerTx, error) {
	response, err := client.JSONRPCCall(c.rpcUrl, "realtime_getInternalTransactions", txHash)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, fmt.Errorf("%d - %s", response.Error.Code, response.Error.Message)
	}

	result := []zktypes.InnerTx{}
	err = json.Unmarshal(response.Result, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// RealtimeGetBalance returns the balance of an account in real-time
func (c *RealtimeClient) RealtimeGetBalance(address common.Address) (*big.Int, error) {
	response, err := client.JSONRPCCall(c.rpcUrl, "realtime_getBalance", address)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, fmt.Errorf("%d - %s", response.Error.Code, response.Error.Message)
	}

	var hexBalance string
	err = json.Unmarshal(response.Result, &hexBalance)
	if err != nil {
		return nil, err
	}

	if len(hexBalance) > 2 && (hexBalance[:2] == "0x" || hexBalance[:2] == "0X") {
		hexBalance = hexBalance[2:]
	}

	balance := new(big.Int)
	balance, ok := balance.SetString(hexBalance, 16)
	if !ok {
		return nil, fmt.Errorf("failed to convert hex to big.Int: %s", hexBalance)
	}

	return balance, nil
}

// RealtimeGetCode returns the code at a given address in real-time
func (c *RealtimeClient) RealtimeGetCode(address common.Address) (string, error) {
	response, err := client.JSONRPCCall(c.rpcUrl, "realtime_getCode", address)
	if err != nil {
		return "", err
	}
	if response.Error != nil {
		return "", fmt.Errorf("%d - %s", response.Error.Code, response.Error.Message)
	}

	var code string
	err = json.Unmarshal(response.Result, &code)
	if err != nil {
		return "", err
	}

	return code, nil
}

// RealtimeGetTransactionCount returns the number of transactions sent from an address in real-time
func (c *RealtimeClient) RealtimeGetTransactionCount(address common.Address) (uint64, error) {
	response, err := client.JSONRPCCall(c.rpcUrl, "realtime_getTransactionCount", address)
	if err != nil {
		return 0, err
	}
	if response.Error != nil {
		return 0, fmt.Errorf("%d - %s", response.Error.Code, response.Error.Message)
	}

	return transHexToUint64(response.Result)
}

// RealtimeGetStorageAt returns the value from a storage position at a given address in real-time
func (c *RealtimeClient) RealtimeGetStorageAt(address common.Address, position string) (string, error) {
	response, err := client.JSONRPCCall(c.rpcUrl, "realtime_getStorageAt", address, position)
	if err != nil {
		return "", err
	}
	if response.Error != nil {
		return "", fmt.Errorf("%d - %s", response.Error.Code, response.Error.Message)
	}

	var result string
	err = json.Unmarshal(response.Result, &result)
	if err != nil {
		return "", err
	}

	return result, nil
}

// RealtimeCall executes a new message call immediately without creating a transaction in real-time
func (c *RealtimeClient) RealtimeCall(from, to common.Address, gas string, gasPrice string, value string, data string) (string, error) {
	txParams := map[string]any{
		"from":     from,
		"to":       to,
		"gas":      gas,
		"gasPrice": gasPrice,
		"value":    value,
		"data":     data,
	}

	response, err := client.JSONRPCCall(c.rpcUrl, "realtime_call", txParams)
	if err != nil {
		return "", err
	}
	if response.Error != nil {
		return "", fmt.Errorf("%d - %s", response.Error.Code, response.Error.Message)
	}

	var result string
	err = json.Unmarshal(response.Result, &result)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (c *RealtimeClient) RealtimeGetTokenBalance(
	ctx context.Context,
	fromAddress common.Address,
	toAddress common.Address,
	erc20Addr common.Address,
) (*big.Int, error) {
	// Pack the balanceOf function call
	data, err := erc20ABI.Pack("balanceOf", toAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to pack balanceOf call: %v", err)
	}

	// Make the realtime_call
	result, err := c.RealtimeCall(fromAddress, erc20Addr, "0x100000", "0x1", "0x0", fmt.Sprintf("0x%x", data))
	if err != nil {
		return nil, fmt.Errorf("failed to call contract: %v", err)
	}

	// Parse the hex result
	if len(result) > 2 && (result[:2] == "0x" || result[:2] == "0X") {
		result = result[2:]
	}

	balance := new(big.Int)
	balance, ok := balance.SetString(result, 16)
	if !ok {
		return nil, fmt.Errorf("failed to convert hex to big.Int: %s", result)
	}

	return balance, nil
}

// RealtimeDumpStateCache dumps the state cache
func (c *RealtimeClient) RealtimeDumpStateCache() error {
	response, err := client.JSONRPCCall(c.rpcUrl, "realtime_dumpStateCache")
	if err != nil {
		return err
	}
	if response.Error != nil {
		return fmt.Errorf("%d - %s", response.Error.Code, response.Error.Message)
	}

	return nil
}

// EthGetBalance returns the balance of an account
func (c *RealtimeClient) EthGetBalance(address common.Address, block string) (*big.Int, error) {
	response, err := client.JSONRPCCall(c.rpcUrl, "eth_getBalance", address, block)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, fmt.Errorf("%d - %s", response.Error.Code, response.Error.Message)
	}

	var hexBalance string
	err = json.Unmarshal(response.Result, &hexBalance)
	if err != nil {
		return nil, err
	}

	if len(hexBalance) > 2 && (hexBalance[:2] == "0x" || hexBalance[:2] == "0X") {
		hexBalance = hexBalance[2:]
	}

	balance := new(big.Int)
	balance, ok := balance.SetString(hexBalance, 16)
	if !ok {
		return nil, fmt.Errorf("failed to convert hex to big.Int: %s", hexBalance)
	}

	return balance, nil
}

// EthGetTransactionCount returns the number of transactions sent from an address
func (c *RealtimeClient) EthGetTransactionCount(address common.Address, block string) (uint64, error) {
	response, err := client.JSONRPCCall(c.rpcUrl, "eth_getTransactionCount", address, block)
	if err != nil {
		return 0, err
	}
	if response.Error != nil {
		return 0, fmt.Errorf("%d - %s", response.Error.Code, response.Error.Message)
	}

	return transHexToUint64(response.Result)
}

func (c *RealtimeClient) EthGetTokenBalance(
	ctx context.Context,
	addr common.Address,
	erc20Addr common.Address,
) (*big.Int, error) {
	// Pack the balanceOf function call
	data, err := erc20ABI.Pack("balanceOf", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to pack balanceOf call: %v", err)
	}

	// Make the eth_call
	result, err := c.client.CallContract(ctx, ethereum.CallMsg{
		To:   &erc20Addr,
		Data: data,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call contract: %v", err)
	}

	// Unpack the result
	var balance *big.Int
	err = erc20ABI.UnpackIntoInterface(&balance, "balanceOf", result)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack result: %v", err)
	}

	return balance, nil
}
