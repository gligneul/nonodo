// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

// This package contains the nonodo options.
package opts

// Default port for the Ethereum node.
const EthDefaultPort = 8545

// Options to nonodo.
type NonodoOpts struct {
	AnvilPort      int
	AnvilBlockTime int
	AnvilVerbose   bool
}

// Create the options struct with default values.
func NewNonodoOpts() NonodoOpts {
	return NonodoOpts{
		AnvilPort:      EthDefaultPort,
		AnvilBlockTime: 1,
		AnvilVerbose:   false,
	}
}
