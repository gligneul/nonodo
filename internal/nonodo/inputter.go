// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package nonodo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cartesi/rollups-node/pkg/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gligneul/nonodo/internal/model"
)

// The inputterService reads inputs from anvil and puts them in the model.
type inputterService struct {
	ready           chan struct{}
	model           *model.NonodoModel
	rpcEndpoint     string
	inputBoxAddress common.Address
	dappAddress     common.Address
}

// Creates a new inputterService from opts.
func newInputterService(
	model *model.NonodoModel,
	rpcEndpoint string,
	inputBoxAddress common.Address,
	dappAddress common.Address,
) *inputterService {
	return &inputterService{
		ready:           make(chan struct{}),
		model:           model,
		rpcEndpoint:     rpcEndpoint,
		inputBoxAddress: inputBoxAddress,
		dappAddress:     dappAddress,
	}
}

func (i *inputterService) Start(ctx context.Context) error {
	client, err := ethclient.DialContext(ctx, i.rpcEndpoint)
	if err != nil {
		return err
	}

	inputBox, err := contracts.NewInputBox(i.inputBoxAddress, client)
	if err != nil {
		return err
	}

	logs := make(chan *contracts.InputBoxInputAdded)
	startingBlock := uint64(1)
	opts := bind.WatchOpts{
		Start:   &startingBlock,
		Context: ctx,
	}
	dappFilter := []common.Address{i.dappAddress}
	sub, err := inputBox.WatchInputAdded(&opts, logs, dappFilter, nil)
	if err != nil {
		return fmt.Errorf("failed to watch inputs: %v", err)
	}

	log.Print("watching inputs")
	i.ready <- struct{}{}
	for {
		select {
		case err := <-sub.Err():
			return err
		case log := <-logs:
			header, err := client.HeaderByHash(ctx, log.Raw.BlockHash)
			if err != nil {
				return fmt.Errorf("failed to get tx header: %v", err)
			}
			i.model.AddAdvanceInput(
				log.Sender,
				log.Input,
				log.Raw.BlockNumber,
				time.Unix(int64(header.Time), 0),
			)
		}
	}
}

func (i *inputterService) Ready() <-chan struct{} {
	return i.ready
}
