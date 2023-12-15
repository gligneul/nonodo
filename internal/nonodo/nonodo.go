// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

// This package contains the nonodo run function.
// This is separate from the main package to facilitate testing.
package nonodo

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gligneul/nonodo/internal/foundry"
	"github.com/gligneul/nonodo/internal/model"
	"github.com/gligneul/nonodo/internal/supervisor"
)

func Run(ctx context.Context, opts NonodoOpts) {
	model := model.NewNonodoModel()
	var services []supervisor.Service
	services = append(services, supervisor.SignalListenerService{})
	services = append(services, foundry.AnvilService{
		Port:      opts.AnvilPort,
		BlockTime: opts.AnvilBlockTime,
		Verbose:   opts.AnvilVerbose,
	})
	services = append(services, inputterService{
		model:           model,
		rpcEndpoint:     fmt.Sprintf("ws://127.0.0.1:%v", opts.AnvilPort),
		inputBoxAddress: common.HexToAddress(opts.InputBoxAddress),
		dappAddress:     common.HexToAddress(opts.DAppAddress),
	})
	services = append(services, echoService{
		model:   model,
		address: fmt.Sprintf("%v:%v", opts.HttpAddress, opts.HttpPort),
	})
	if opts.BuiltInDApp {
		services = append(services, dappService{
			rollupEndpoint: fmt.Sprintf("http://127.0.0.1:%v/rollup", opts.HttpPort),
		})
	}
	supervisor.Start(ctx, services)
}
