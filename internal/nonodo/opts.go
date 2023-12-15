// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package nonodo

import "github.com/gligneul/nonodo/internal/foundry"

// Options to nonodo.
type NonodoOpts struct {
	AnvilPort       int
	AnvilBlockTime  int
	AnvilVerbose    bool
	BuiltInDApp     bool
	HttpAddress     string
	HttpPort        int
	InputBoxAddress string
	DAppAddress     string
}

// Create the options struct with default values.
func NewNonodoOpts() NonodoOpts {
	return NonodoOpts{
		AnvilPort:       foundry.AnvilDefaultPort,
		AnvilBlockTime:  1,
		AnvilVerbose:    false,
		BuiltInDApp:     false,
		HttpAddress:     "127.0.0.1",
		HttpPort:        8080,
		InputBoxAddress: foundry.InputBoxAddress,
		DAppAddress:     foundry.DAppAddress,
	}
}
