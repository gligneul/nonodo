// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

// This package contains the nonodo run function.
// This is separate from the main package to facilitate testing.
package nonodo

import (
	"context"

	"github.com/gligneul/nonodo/internal/opts"
	"github.com/gligneul/nonodo/internal/supervisor"
)

func Run(ctx context.Context, opts opts.NonodoOpts) {
	var services []supervisor.Service

	services = append(services, supervisor.NewSignalListenerService())

	anvil, cleanup := newAnvil(opts)
	defer cleanup()
	services = append(services, anvil)

	supervisor.Start(ctx, services)
}
