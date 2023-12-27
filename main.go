// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

// This package contains the main function that executes the nonodo command.
package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/carlmjohnson/versioninfo"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gligneul/nonodo/internal/nonodo"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:     "nonodo [flags] [-- application [args]...]",
	Short:   "nonodo is a development node for Cartesi Rollups",
	Run:     run,
	Version: versioninfo.Short(),
}

var debug bool
var color bool
var opts = nonodo.NewNonodoOpts()

func init() {
	// address-*
	cmd.Flags().StringVar(&opts.ApplicationAddress, "address-application", opts.ApplicationAddress,
		"Application contract address")
	cmd.Flags().StringVar(&opts.InputBoxAddress, "address-input-box", opts.InputBoxAddress,
		"InputBox contract address")

	// anvil-*
	cmd.Flags().IntVar(&opts.AnvilPort, "anvil-port", opts.AnvilPort,
		"HTTP port used by Anvil")
	cmd.Flags().BoolVar(&opts.AnvilVerbose, "anvil-verbose", opts.AnvilVerbose,
		"If set, prints Anvil's output")

	// enable-*
	cmd.Flags().BoolVarP(&debug, "enable-debug", "d", false, "If set, enable debug output")
	cmd.Flags().BoolVar(&color, "enable-color", true, "If set, enables logs color")
	cmd.Flags().BoolVar(&opts.EnableEcho, "enable-echo", opts.EnableEcho,
		"If set, nonodo starts a built-in echo application")

	// http-*
	cmd.Flags().StringVar(&opts.HttpAddress, "http-address", opts.HttpAddress,
		"HTTP address used by nonodo to serve its APIs")
	cmd.Flags().IntVar(&opts.HttpPort, "http-port", opts.HttpPort,
		"HTTP port used by nonodo to serve its APIs")
}

func run(cmd *cobra.Command, args []string) {
	var startTime = time.Now()

	logOpts := new(tint.Options)
	if debug {
		logOpts.Level = slog.LevelDebug
	}
	logOpts.AddSource = debug
	logOpts.NoColor = !color || !isatty.IsTerminal(os.Stdout.Fd())
	logOpts.TimeFormat = "[15:04:05.000]"
	handler := tint.NewHandler(os.Stdout, logOpts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	if opts.AnvilPort == 0 {
		exitf("anvil port cannot be 0")
	}
	checkEthAddress(cmd, "address-input-box")
	checkEthAddress(cmd, "address-application")
	opts.ApplicationArgs = args

	ctx, cancel := signal.NotifyContext(cmd.Context(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	ready := make(chan struct{}, 1)
	go func() {
		select {
		case <-ready:
			slog.Info("nonodo: ready", "after", time.Since(startTime))
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
	err := fmt.Sprintf(format, args...)
	slog.Error("configuration error", "error", err)
	os.Exit(1)
}

func checkEthAddress(cmd *cobra.Command, varName string) {
	if cmd.Flags().Changed(varName) {
		value, err := cmd.Flags().GetString(varName)
		cobra.CheckErr(err)
		bytes, err := hexutil.Decode(value)
		if err != nil {
			exitf("invalid address for --%v: %v", varName, err)
		}
		if len(bytes) != common.AddressLength {
			exitf("invalid address for --%v: wrong length", varName)
		}
	}
}
