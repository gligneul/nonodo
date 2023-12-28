// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package devnet

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gligneul/nonodo/internal/contracts"
)

// Send an input using the devnet sender.
// This function should be used in the devnet environment.
func AddInput(ctx context.Context, rpcUrl string, payload []byte) error {
	if len(payload) == 0 {
		return fmt.Errorf("cannot send empty payload")
	}

	client, err := ethclient.DialContext(ctx, rpcUrl)
	if err != nil {
		return fmt.Errorf("dial to %v: %w", rpcUrl, err)
	}

	privateKey, err := crypto.ToECDSA(common.Hex2Bytes(SenderPrivateKey[2:]))
	if err != nil {
		return fmt.Errorf("create private key: %w", err)
	}

	chainId, err := client.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("get chain id: %w", err)
	}

	txOpts, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		return fmt.Errorf("create transactor: %w", err)
	}
	nonce, err := client.PendingNonceAt(ctx, common.HexToAddress(SenderAddress))
	if err != nil {
		return fmt.Errorf("get nonce: %w", err)
	}
	txOpts.Nonce = big.NewInt(int64(nonce))
	txOpts.Value = big.NewInt(0)
	txOpts.GasLimit = GasLimit
	txOpts.GasPrice, err = client.SuggestGasPrice(ctx)
	if err != nil {
		return fmt.Errorf("get gas price: %w", err)
	}

	inputBox, err := contracts.NewInputBox(common.HexToAddress(InputBoxAddress), client)
	if err != nil {
		return fmt.Errorf("bind input box: %w", err)
	}

	tx, err := inputBox.AddInput(txOpts, common.HexToAddress(ApplicationAddress), payload)
	if err != nil {
		return fmt.Errorf("add input: %w", err)
	}

	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		return fmt.Errorf("wait mined: %w", err)
	}
	if receipt.Status == 0 {
		return fmt.Errorf("transaction was not accepted")
	}
	return nil
}

// Get all input added events from the input box.
func GetInputAdded(ctx context.Context, rpcUrl string) ([]*contracts.InputBoxInputAdded, error) {
	client, err := ethclient.DialContext(ctx, rpcUrl)
	if err != nil {
		return nil, fmt.Errorf("dial to %v: %w", rpcUrl, err)
	}
	inputBox, err := contracts.NewInputBox(common.HexToAddress(InputBoxAddress), client)
	if err != nil {
		return nil, fmt.Errorf("bind input box: %w", err)
	}
	it, err := inputBox.FilterInputAdded(nil, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to filter input added: %v", err)
	}
	defer it.Close()
	var inputs []*contracts.InputBoxInputAdded
	for it.Next() {
		inputs = append(inputs, it.Event)
	}
	return inputs, nil
}
