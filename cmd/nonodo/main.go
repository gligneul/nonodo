// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package main

import (
	"context"

	"github.com/gligneul/nonodo/internal/nonodo"
	"github.com/gligneul/nonodo/internal/opts"
)

func main() {
	opts := opts.NewNonodoOpts()
	nonodo.Run(context.Background(), opts)
}
