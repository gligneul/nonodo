// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

// This package contains the nonodo binary.
// To read the usage, run the binary passing the -h flag.
package main

import (
	"os"

	"github.com/gligneul/nonodo/internal/nonodo"
	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "nonodo",
	Short: "Development Node for Cartesi Rollups",
	Run:   run,
}

var opts = nonodo.NewNonodoOpts()

func init() {
	cmd.Flags().IntVar(&opts.AnvilBlockTime, "anvil-block-time", opts.AnvilBlockTime,
		"Time in seconds between Anvil blocks")
	cmd.Flags().IntVar(&opts.AnvilPort, "anvil-port", opts.AnvilPort,
		"HTTP port used by Anvil")
	cmd.Flags().BoolVar(&opts.AnvilVerbose, "anvil-verbose", opts.AnvilVerbose,
		"If true, prints Anvil's output")
	cmd.Flags().IntVar(&opts.HttpPort, "http-port", opts.HttpPort,
		"HTTP port used by nonodo to serve its APIs")
	cmd.Flags().BoolVar(&opts.BuiltInDApp, "built-in-dapp", opts.BuiltInDApp,
		"If true, nonodo starts a built-in echo DApp")
}

func run(cmd *cobra.Command, args []string) {
	nonodo.Run(cmd.Context(), opts)
}

func main() {
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
