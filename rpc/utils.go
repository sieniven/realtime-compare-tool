package rpc

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strings"
	"testing"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/sieniven/realtime-compare-tool/rpc/rtclient"
	"github.com/stretchr/testify/require"
)

var (
	erc20ABI, _ = abi.JSON(strings.NewReader(erc20ABIJson))
)

// GetAuth configures and returns an auth object.
func GetAuth(privateKeyStr string, chainID uint64) (*bind.TransactOpts, error) {
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(privateKeyStr, "0x"))
	if err != nil {
		return nil, err
	}

	return bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(0).SetUint64(chainID))
}

func GetTestChainConfig(chainID uint64) *params.ChainConfig {
	return &params.ChainConfig{
		ChainID:                 big.NewInt(int64(chainID)),
		HomesteadBlock:          big.NewInt(0),
		DAOForkBlock:            nil,
		DAOForkSupport:          false,
		EIP150Block:             big.NewInt(0),
		EIP155Block:             big.NewInt(0),
		EIP158Block:             big.NewInt(0),
		ByzantiumBlock:          big.NewInt(0),
		ConstantinopleBlock:     big.NewInt(0),
		PetersburgBlock:         big.NewInt(0),
		IstanbulBlock:           big.NewInt(0),
		MuirGlacierBlock:        big.NewInt(0),
		BerlinBlock:             big.NewInt(0),
		LondonBlock:             big.NewInt(0),
		ArrowGlacierBlock:       big.NewInt(0),
		GrayGlacierBlock:        big.NewInt(0),
		MergeNetsplitBlock:      nil,
		ShanghaiTime:            nil,
		CancunTime:              nil,
		PragueTime:              nil,
		OsakaTime:               nil,
		VerkleTime:              nil,
		TerminalTotalDifficulty: big.NewInt(math.MaxInt64),
		Ethash:                  new(params.EthashConfig),
		Clique:                  nil,
	}
}

func deployERC20Contract(
	t *testing.T,
	ctx context.Context,
	privateKey *ecdsa.PrivateKey,
	client *rtclient.RealtimeClient,
) common.Address {
	publicKey := privateKey.Public()
	tmpPublicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	require.True(t, ok)
	fromAddress := crypto.PubkeyToAddress(*tmpPublicKeyECDSA)
	fmt.Printf("Sender: %s\n", fromAddress)

	nonce, err := client.RealtimeGetTransactionCount(ctx, fromAddress)
	require.NoError(t, err)

	// Define gas parameters
	gasPrice, err := client.SuggestGasPrice(ctx)
	require.NoError(t, err)

	// Set up transaction options
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(int64(DefaultL2ChainID)))
	require.NoError(t, err)

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(3000000)
	auth.GasPrice = gasPrice

	// Deploy the contract
	erc20Bytecode, err := hex.DecodeString(erc20BytecodeStr)
	require.NoError(t, err)
	erc20Address, tx, _, err := bind.DeployContract(auth, erc20ABI, erc20Bytecode, client)
	require.NoError(t, err)

	// Wait for contract deployment to be mined
	err = WaitTxToBeMined(ctx, client, tx, DefaultTimeoutTxToBeMined)
	require.NoError(t, err)
	fmt.Printf("ERC20 Contract deployed at: %s, transaction hash: %s\n", erc20Address.Hex(), tx.Hash().Hex())

	return erc20Address
}

func erc20TransferTx(
	t *testing.T,
	ctx context.Context,
	privateKey *ecdsa.PrivateKey,
	client *rtclient.RealtimeClient,
	amount *big.Int,
	toAddress common.Address,
	erc20Address common.Address,
	nonce uint64,
) *types.Transaction {

	// Prepare transfer data
	data, err := erc20ABI.Pack("transfer", toAddress, amount)
	require.NoError(t, err)

	gasPrice, err := client.SuggestGasPrice(ctx)
	require.NoError(t, err)
	transferERC20TokenTx := types.NewTransaction(
		nonce,
		erc20Address,
		big.NewInt(0),
		60000,
		gasPrice,
		data,
	)

	signer := types.MakeSigner(GetTestChainConfig(DefaultL2ChainID), big.NewInt(1), 0)
	signedTx, err := types.SignTx(transferERC20TokenTx, signer, privateKey)
	require.NoError(t, err)

	err = client.SendTransaction(ctx, signedTx)
	require.NoError(t, err)
	fmt.Printf("Sent signedTx: %s\n", signedTx.Hash().String())

	return signedTx
}

// WaitMined waits for tx to be mined on the blockchain.
// It stops waiting when the context is canceled.
func WaitMined(ctx context.Context, b bind.DeployBackend, txHash common.Hash) (*types.Receipt, error) {
	queryTicker := time.NewTicker(10 * time.Millisecond)
	defer queryTicker.Stop()

	for {
		receipt, err := b.TransactionReceipt(ctx, txHash)
		if err == nil {
			return receipt, nil
		}

		// Wait for the next round.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-queryTicker.C:
		}
	}
}

type ethClienter interface {
	ethereum.TransactionReader
	ethereum.ContractCaller
	bind.DeployBackend
}

// RevertReason returns the revert reason for a tx that has a receipt with failed status
func RevertReason(ctx context.Context, c ethClienter, tx *types.Transaction, blockNumber *big.Int) (string, error) {
	if tx == nil {
		return "", nil
	}

	signer := types.MakeSigner(GetTestChainConfig(DefaultL2ChainID), big.NewInt(1), 0)
	from, err := types.Sender(signer, tx)
	if err != nil {
		return "", err
	}

	msg := ethereum.CallMsg{
		From: from,
		To:   tx.To(),
		Gas:  tx.Gas(),

		Value: tx.Value(),
		Data:  tx.Data(),
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

// RevertReasonRealtime returns the revert reason for a tx that has a receipt with failed status
func RevertReasonRealtime(
	ctx context.Context,
	client *rtclient.RealtimeClient,
	tx *types.Transaction,
) (string, error) {
	if tx == nil {
		return "", nil
	}

	signer := types.MakeSigner(GetTestChainConfig(DefaultL2ChainID), big.NewInt(1), 0)
	from, err := types.Sender(signer, tx)
	if err != nil {
		return "", err
	}

	hex, err := client.RealtimeCall(ctx, from, *tx.To(), fmt.Sprintf("0x%x", tx.Gas()), fmt.Sprintf("0x%x", tx.GasPrice()), fmt.Sprintf("0x%x", tx.Value()), fmt.Sprintf("0x%x", tx.Data()))
	if err != nil {
		return "", err
	}

	unpackedMsg, err := abi.UnpackRevert(common.FromHex(hex))
	if err != nil {
		fmt.Printf("failed to get the revert message for tx %v: %v\n", tx.Hash(), err)
		return "", errors.New("execution reverted")
	}

	return unpackedMsg, nil
}

// WaitTxToBeMined waits until a tx has been mined or the given timeout expires.
func WaitTxToBeMined(parentCtx context.Context, client *rtclient.RealtimeClient, tx *types.Transaction, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(parentCtx, timeout)
	defer cancel()
	receipt, err := WaitMined(ctx, client, tx.Hash())
	if errors.Is(err, context.DeadlineExceeded) {
		return err
	} else if err != nil {
		fmt.Printf("error waiting tx %s to be mined: %v\n", tx.Hash(), err)
		return err
	}
	if receipt.Status == types.ReceiptStatusFailed {
		// Get revert reason
		reason, reasonErr := RevertReasonRealtime(ctx, client, tx)
		if reasonErr != nil {
			reason = reasonErr.Error()
		}
		return fmt.Errorf("transaction has failed, reason: %s, receipt: %+v. tx: %+v, gas: %v", reason, receipt, tx, tx.Gas())
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		currNum, err := client.BlockNumber(ctx)
		if err != nil {
			return err
		}

		if currNum >= receipt.BlockNumber.Uint64() {
			break
		}
	}

	fmt.Printf("Transaction successfully mined: %v\n", tx.Hash())
	return nil
}
