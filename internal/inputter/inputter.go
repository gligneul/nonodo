// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package inputter

import (
	"context"
	"fmt"
	"time"

	"github.com/cartesi/rollups-node/pkg/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Model interface {
	AddAdvanceInput(
		sender common.Address,
		payload []byte,
		blockNumber uint64,
		timestamp time.Time,
	)
}

// The inputterService reads inputs from Ethereum and puts them in the model.
type InputterService struct {
	Model              Model
	Provider           string
	InputBoxAddress    common.Address
	ApplicationAddress common.Address
}

func (s InputterService) String() string {
	return "inputter"
}

func (s InputterService) Start(ctx context.Context, ready chan<- struct{}) error {
	client, err := ethclient.DialContext(ctx, s.Provider)
	if err != nil {
		return err
	}

	inputBox, err := contracts.NewInputBox(s.InputBoxAddress, client)
	if err != nil {
		return err
	}

	logs := make(chan *contracts.InputBoxInputAdded)
	startingBlock := findGenesis(s.InputBoxAddress)
	opts := bind.WatchOpts{
		Start:   &startingBlock,
		Context: ctx,
	}
	filter := []common.Address{s.ApplicationAddress}
	sub, err := inputBox.WatchInputAdded(&opts, logs, filter, nil)
	if err != nil {
		return fmt.Errorf("failed to watch inputs: %w", err)
	}

	ready <- struct{}{}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-sub.Err():
			return err
		case log := <-logs:
			header, err := client.HeaderByHash(ctx, log.Raw.BlockHash)
			if err != nil {
				return fmt.Errorf("failed to get tx header: %w", err)
			}
			s.Model.AddAdvanceInput(
				log.Sender,
				log.Input,
				log.Raw.BlockNumber,
				time.Unix(int64(header.Time), 0),
			)
		}
	}
}

func findGenesis(contract common.Address) uint64 {
	return 1
}
