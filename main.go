// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

// This package contains the main function that executes the nonodo command.
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/carlmjohnson/versioninfo"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gligneul/nonodo/internal/nonodo"
	"github.com/spf13/cobra"
)

var startTime = time.Now()

var cmd = &cobra.Command{
	Use:     "nonodo [flags] [-- application [args]...]",
	Short:   "Nonodo is a development node for Cartesi Rollups",
	Run:     run,
	Version: versioninfo.Short(),
}

var opts = nonodo.NewNonodoOpts()

func init() {
	cmd.Flags().IntVar(&opts.AnvilPort, "anvil-port", opts.AnvilPort,
		"HTTP port used by Anvil")
	cmd.Flags().BoolVar(&opts.AnvilVerbose, "anvil-verbose", opts.AnvilVerbose,
		"If set, prints Anvil's output")
	cmd.Flags().BoolVar(&opts.BuiltInEcho, "built-in-echo", opts.BuiltInEcho,
		"If set, nonodo starts a built-in echo application")
	cmd.Flags().StringVar(&opts.HttpAddress, "http-address", opts.HttpAddress,
		"HTTP address used by nonodo to serve its APIs")
	cmd.Flags().IntVar(&opts.HttpPort, "http-port", opts.HttpPort,
		"HTTP port used by nonodo to serve its APIs")
	cmd.Flags().StringVar(&opts.InputBoxAddress, "address-input-box", opts.InputBoxAddress,
		"InputBox contract address")
	cmd.Flags().StringVar(&opts.ApplicationAddress, "address-application", opts.ApplicationAddress,
		"Application contract address")
}

func run(cmd *cobra.Command, args []string) {
	checkEthAddress(cmd, "address-input-box")
	checkEthAddress(cmd, "address-application")
	opts.ApplicationArgs = args

	ctx, cancel := signal.NotifyContext(cmd.Context(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	ready := make(chan struct{}, 1)
	go func() {
		select {
		case <-ready:
			duration := time.Since(startTime)
			log.Printf("nonodo: ready after %v", duration)
		case <-ctx.Done():
		}
	}()

	w, err := nonodo.NewNonodoWorker(opts)
	cobra.CheckErr(err)
	err = w.Start(ctx, ready)
	cobra.CheckErr(err)
}

func main() {
	cobra.CheckErr(cmd.Execute())
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
