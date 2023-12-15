// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package nonodo

// Default port for the Ethereum node.
const EthDefaultPort = 8545

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
		AnvilPort:       EthDefaultPort,
		AnvilBlockTime:  1,
		AnvilVerbose:    false,
		BuiltInDApp:     false,
		HttpAddress:     "127.0.0.1",
		HttpPort:        8080,
		InputBoxAddress: "0x59b22D57D4f067708AB0c00552767405926dc768",
		DAppAddress:     "0x70ac08179605AF2D9e75782b8DEcDD3c22aA4D0C",
	}
}
