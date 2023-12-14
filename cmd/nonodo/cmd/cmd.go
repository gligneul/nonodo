// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

// This package contains the nonodo command.
package cmd

import (
	"github.com/gligneul/nonodo/internal/nonodo"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "nonodo",
	Short: "Development Node for Cartesi Rollups",
	Long:  removeFirst(long),
	Run:   run,
}

const long = `
Nonodo is a development node for Cartesi Rollups. It was designed to work with DApps running in the
host machine instead of the Cartesi machine. The DApp back-end should call nonodo's Rollup API to
advance the rollups state and to perform inspect requests.

Nonodo uses the Anvil as the underlying Ethereum node. To install Anvil, read the instructions in
the Foundry book: https://book.getfoundry.sh/getting-started/installation.

To start nonodo with default configuration, run the following command.

	nonodo

In the default configuration, nonodo starts an Anvil node with the Cartesi Rollups contracts
deployed. This is the same deployment used by sunodo, so the contract addresses are the same.
Nonodo offer some flags to configure Anvil; these flags start with --anvil-*.

Nonodo exposes the Cartesi Rollups GraphQL (/graphql) and Inspect (/inspect) APIs for the DApp
front-end, and the Rollup (/rollup) API for the DApp back-end. Nonodo uses the HTTP address and port
set by the --http-* flags.

To start nonodo with a built-in echo DApp, use the --built-in-dapp flag. This flag is useful when
testing the DApp front-end without a working back-end.

	nonodo --built-in-dapp

The flag usage is described below.`

var opts = nonodo.NewNonodoOpts()

func init() {
	Cmd.Flags().IntVar(&opts.AnvilBlockTime, "anvil-block-time", opts.AnvilBlockTime,
		"Time in seconds between Anvil blocks")
	Cmd.Flags().IntVar(&opts.AnvilPort, "anvil-port", opts.AnvilPort,
		"HTTP port used by Anvil")
	Cmd.Flags().BoolVar(&opts.AnvilVerbose, "anvil-verbose", opts.AnvilVerbose,
		"If true, prints Anvil's output")
	Cmd.Flags().IntVar(&opts.HttpPort, "http-port", opts.HttpPort,
		"HTTP port used by nonodo to serve its APIs")
	Cmd.Flags().BoolVar(&opts.BuiltInDApp, "built-in-dapp", opts.BuiltInDApp,
		"If true, nonodo starts a built-in echo DApp")
}

func run(Cmd *cobra.Command, args []string) {
	nonodo.Run(Cmd.Context(), opts)
}

func removeFirst(s string) string {
	return s[1:]
}
