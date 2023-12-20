// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

// This package contains the strings that are used to generate the nonodo docs.
package docs

import "bytes"

// Use go generate to update the README file
//go:generate go run ./gen-readme

// Get the short string for the CLI command.
func Short() string {
	return short
}

// Get the long string for the CLI command.
func Long() string {
	return long[1:]
}

// Get the contents for the read-me.
func Readme() []byte {
	var buf bytes.Buffer
	buf.WriteString("# NoNodo\n\n")
	buf.WriteString(long)
	buf.WriteString("\n\n")
	return buf.Bytes()
}

const short = "Development Node for Cartesi Rollups"

const long = `
Nonodo is a development node for Cartesi Rollups.

Nonodo was designed to work with applications running in the host machine instead of the Cartesi
machine. The application back-end should call the Rollup HTTP API to advance the rollups state and
to process inspect inputs.

Nonodo uses the Anvil as the underlying Ethereum node. To install Anvil, read the instructions in
the Foundry book: https://book.getfoundry.sh/getting-started/installation.

To start nonodo with default configuration, run the command below.

	nonodo

With the default configuration, nonodo starts an Anvil node with the Cartesi Rollups contracts
deployed. This is the same deployment used by sunodo; so, the contract addresses are the same.
Nonodo offer some flags to configure Anvil; these flags start with --anvil-*.

To send an input to the Cartesi application, you may use cast; a command-line tool from the foundry
package. For instance, the invocation below sends an input with contents 0xdeadbeef to the running
application.

	INPUT=0xdeadbeef; \
	INPUT_BOX_ADDRESS=0x59b22D57D4f067708AB0c00552767405926dc768; \
	APPLICATION_ADDRESS=0x70ac08179605AF2D9e75782b8DEcDD3c22aA4D0C; \
	MNEMONIC="test test test test test test test test test test test junk"; \
	cast send --mnemonic $MNEMONIC --rpc-url "http://localhost:8545" $INPUT_BOX_ADDRESS \
		"addInput(address,bytes)(bytes32)" $APPLICATION_ADDRESS $INPUT

Nonodo exposes the Cartesi Rollups GraphQL (/graphql) and Inspect (/inspect) APIs for the
application front-end, and the Rollup (/rollup) API for the application back-end. Nonodo uses the
HTTP address and port set by the --http-* flags.

To start nonodo with a built-in echo application, use the --built-in-echo flag. This flag is useful
when testing the application front-end without a working back-end.

	nonodo --built-in-app`
