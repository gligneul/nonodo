// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

// This package contains the nonodo run function.
// This is separate from the main package to facilitate testing.
package nonodo

import (
	"context"
	"fmt"

	"github.com/gligneul/nonodo/internal/model"
	"github.com/gligneul/nonodo/internal/opts"
	"github.com/gligneul/nonodo/internal/supervisor"
)

func Run(ctx context.Context, opts opts.NonodoOpts) {
	var services []supervisor.Service

	services = append(services, supervisor.NewSignalListenerService())

	anvil, cleanup := newAnvilService(opts)
	defer cleanup()
	services = append(services, anvil)

	model := model.NewNonodoModel()
	rpcEndpoint := fmt.Sprintf("ws://127.0.0.1:%v", opts.AnvilPort)
	services = append(services, newInputterService(model, rpcEndpoint))

	services = append(services, newEchoService(model, opts.HttpPort))

	supervisor.Start(ctx, services)
}
