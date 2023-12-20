// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

// This package contains the nonodo command.
package cmd

import (
	"fmt"
	"os"

	"github.com/carlmjohnson/versioninfo"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gligneul/nonodo/internal/nonodo"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "nonodo",
	Short:   "Development Node for Cartesi Rollups",
	Long:    removeFirst(long),
	Run:     run,
	Version: versioninfo.Short(),
}

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

var opts = nonodo.NewNonodoOpts()

func init() {
	Cmd.Flags().IntVar(&opts.AnvilBlockTime, "anvil-block-time", opts.AnvilBlockTime,
		"Time in seconds between Anvil blocks")
	Cmd.Flags().IntVar(&opts.AnvilPort, "anvil-port", opts.AnvilPort,
		"HTTP port used by Anvil")
	Cmd.Flags().BoolVar(&opts.AnvilVerbose, "anvil-verbose", opts.AnvilVerbose,
		"If set, prints Anvil's output")
	Cmd.Flags().BoolVar(&opts.BuiltInEcho, "built-in-echo", opts.BuiltInEcho,
		"If set, nonodo starts a built-in echo application")
	Cmd.Flags().StringVar(&opts.HttpAddress, "http-address", opts.HttpAddress,
		"HTTP address used by nonodo to serve its APIs")
	Cmd.Flags().IntVar(&opts.HttpPort, "http-port", opts.HttpPort,
		"HTTP port used by nonodo to serve its APIs")
	Cmd.Flags().StringVar(&opts.InputBoxAddress, "address-input-box", opts.InputBoxAddress,
		"InputBox contract address")
	Cmd.Flags().StringVar(&opts.ApplicationAddress, "address-application", opts.ApplicationAddress,
		"Application contract address")
}

func run(cmd *cobra.Command, args []string) {
	checkEthAddress(cmd, "address-input-box")
	checkEthAddress(cmd, "address-application")
	nonodo.Run(cmd.Context(), opts)
}

func exitf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func checkEthAddress(cmd *cobra.Command, varName string) {
	if cmd.Flags().Changed(varName) {
		value, err := cmd.Flags().GetString(varName)
		cobra.CheckErr(err)
		bytes, err := hexutil.Decode(value)
		if err != nil {
			exitf("invalid address for --%v: %v\n", varName, err)
		}
		if len(bytes) != common.AddressLength {
			exitf("invalid address for --%v: wrong length\n", varName)
		}
	}
}

func removeFirst(s string) string {
	return s[1:]
}
