package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	ethereum "github.com/ledgerwatch/erigon"
	"github.com/ledgerwatch/erigon/accounts/abi"
	"github.com/ledgerwatch/erigon/accounts/abi/bind"
	"github.com/ledgerwatch/erigon/core/types"
)

var (
	erc20ABI, _ = abi.JSON(strings.NewReader(erc20ABIJson))
)

type ethClienter interface {
	ethereum.TransactionReader
	ethereum.ContractCaller
	bind.DeployBackend
}

// RevertReason returns the revert reason for a tx that has a receipt with failed status
func RevertReason(
	ctx context.Context,
	c ethClienter,
	tx types.Transaction,
	blockNumber *big.Int,
) (string, error) {
	if tx == nil {
		return "", nil
	}

	from, _ := tx.GetSender()
	msg := ethereum.CallMsg{
		From: from,
		To:   tx.GetTo(),
		Gas:  tx.GetGas(),

		Value: tx.GetValue(),
		Data:  tx.GetData(),
	}
	hex, err := c.CallContract(ctx, msg, blockNumber)
	if err != nil {
		return "", err
	}

	unpackedMsg, err := abi.UnpackRevert(hex)
	if err != nil {
		fmt.Printf("failed to get the revert message for tx %v: %v\n", tx.Hash(), err)
		return "", errors.New("execution reverted")
	}

	return unpackedMsg, nil
}

func transHexToUint64(hex json.RawMessage) (uint64, error) {
	var result string
	err := json.Unmarshal(hex, &result)
	if err != nil {
		return 0, err
	}

	if len(result) > 1 && (result[:2] == "0x" || result[:2] == "0X") {
		result = result[2:]
	}

	result1, err := strconv.ParseUint(result, 16, 64)
	if err != nil {
		return 0, err
	}

	return result1, nil
}
