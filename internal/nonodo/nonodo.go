// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

// This package contains the nonodo run function.
// This is separate from the main package to facilitate testing.
package nonodo

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gligneul/nonodo/internal/model"
	"github.com/gligneul/nonodo/internal/supervisor"
)

func Run(ctx context.Context, opts NonodoOpts) {
	model := model.NewNonodoModel()
	var services []supervisor.Service

	services = append(services, supervisor.NewSignalListenerService())

	services = append(services, newAnvilService(opts))

	rpcEndpoint := fmt.Sprintf("ws://127.0.0.1:%v", opts.AnvilPort)
	inputBoxAddress := common.HexToAddress(opts.InputBoxAddress)
	services = append(services, newInputterService(model, rpcEndpoint, inputBoxAddress))

	services = append(services, newEchoService(model, opts.HttpPort))

	if opts.BuiltInDApp {
		rollupsEndpoint := fmt.Sprintf("http://127.0.0.1:%v/rollup", opts.HttpPort)
		services = append(services, newDappService(rollupsEndpoint))
	}

	supervisor.Start(ctx, services)
}
