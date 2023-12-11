// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

// This package contains the nonodo run function.
// This is separate from the main package to facilitate testing.
package nonodo

import (
	"context"
	"fmt"

	"github.com/gligneul/nonodo/internal/opts"
)

func Run(ctx context.Context, opts opts.NonodoOpts) {
	fmt.Println("vim-go")
}
