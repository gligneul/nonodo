// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

// This package contains the nonodo binary.
// To read the usage, run the binary passing the -h flag.
package main

import (
	"os"

	"github.com/gligneul/nonodo/internal/nonodo"
	"github.com/gligneul/nonodo/internal/opts"
	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "nonodo",
	Short: "Development Node for Cartesi Rollups",
	Run:   run,
}

var nonodoOpts = opts.NewNonodoOpts()

func init() {
	cmd.Flags().IntVar(&nonodoOpts.AnvilBlockTime, "anvil-block-time",
		nonodoOpts.AnvilBlockTime, "Time in seconds between Anvil blocks")
	cmd.Flags().IntVar(&nonodoOpts.AnvilPort, "anvil-port", nonodoOpts.AnvilPort,
		"HTTP port used by Anvil")
	cmd.Flags().BoolVar(&nonodoOpts.AnvilVerbose, "anvil-verbose", nonodoOpts.AnvilVerbose,
		"If true, prints Anvil's output")
	cmd.Flags().IntVar(&nonodoOpts.HttpPort, "http-port", nonodoOpts.HttpPort,
		"HTTP port used by nonodo to serve its APIs")
}

func run(cmd *cobra.Command, args []string) {
	nonodo.Run(cmd.Context(), nonodoOpts)
}

func main() {
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
