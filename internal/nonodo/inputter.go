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
	model              *model.NonodoModel
	rpcEndpoint        string
	inputBoxAddress    common.Address
	applicationAddress common.Address
}

func (i inputterService) Start(ctx context.Context, ready chan<- struct{}) error {
	client, err := ethclient.DialContext(ctx, i.rpcEndpoint)
	if err != nil {
		return err
	}

	inputBox, err := contracts.NewInputBox(i.inputBoxAddress, client)
	if err != nil {
		return err
	}

	logs := make(chan *contracts.InputBoxInputAdded)
	startingBlock := findGenesis(i.inputBoxAddress)
	opts := bind.WatchOpts{
		Start:   &startingBlock,
		Context: ctx,
	}
	filter := []common.Address{i.applicationAddress}
	sub, err := inputBox.WatchInputAdded(&opts, logs, filter, nil)
	if err != nil {
		return fmt.Errorf("failed to watch inputs: %v", err)
	}

	log.Print("watching inputs")
	ready <- struct{}{}
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

func findGenesis(contract common.Address) uint64 {
	return 1
}
