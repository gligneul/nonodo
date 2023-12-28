// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package inputter

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gligneul/nonodo/internal/contracts"
)

type Model interface {
	AddAdvanceInput(
		sender common.Address,
		payload []byte,
		blockNumber uint64,
		timestamp time.Time,
	)
}

// This worker reads inputs from Ethereum and puts them in the model.
type InputterWorker struct {
	Model              Model
	Provider           string
	InputBoxAddress    common.Address
	ApplicationAddress common.Address
}

func (w InputterWorker) String() string {
	return "inputter"
}

func (w InputterWorker) Start(ctx context.Context, ready chan<- struct{}) error {
	client, err := ethclient.DialContext(ctx, w.Provider)
	if err != nil {
		return fmt.Errorf("inputter: failed to dial: %w", err)
	}

	inputBox, err := contracts.NewInputBox(w.InputBoxAddress, client)
	if err != nil {
		return fmt.Errorf("inputter: failed to bind input box: %w", err)
	}

	logs := make(chan *contracts.InputBoxInputAdded)
	var startingBlock uint64 = 0
	opts := bind.WatchOpts{
		Start:   &startingBlock,
		Context: ctx,
	}
	filter := []common.Address{w.ApplicationAddress}
	sub, err := inputBox.WatchInputAdded(&opts, logs, filter, nil)
	if err != nil {
		return fmt.Errorf("inputter: failed to watch inputs: %w", err)
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
				return fmt.Errorf("inputter: failed to get tx header: %w", err)
			}
			slog.Debug("inputter: read event log", "log", log)
			w.Model.AddAdvanceInput(
				log.Sender,
				log.Input,
				log.Raw.BlockNumber,
				time.Unix(int64(header.Time), 0),
			)
		}
	}
}
