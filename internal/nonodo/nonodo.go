// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

// This package contains the nonodo run function.
// This is separate from the main package to facilitate testing.
package nonodo

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gligneul/nonodo/internal/echoapp"
	"github.com/gligneul/nonodo/internal/foundry"
	"github.com/gligneul/nonodo/internal/inputter"
	"github.com/gligneul/nonodo/internal/inspect"
	"github.com/gligneul/nonodo/internal/model"
	"github.com/gligneul/nonodo/internal/reader"
	"github.com/gligneul/nonodo/internal/rollup"
	"github.com/gligneul/nonodo/internal/supervisor"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var ApplicationConflictErr = errors.New("can't use built-in echo with custom application")

// Options to nonodo.
type NonodoOpts struct {
	AnvilPort          int
	AnvilVerbose       bool
	EnableEcho         bool
	HttpAddress        string
	HttpPort           int
	InputBoxAddress    string
	ApplicationAddress string
	ApplicationArgs    []string
}

// Create the options struct with default values.
func NewNonodoOpts() NonodoOpts {
	return NonodoOpts{
		AnvilPort:          foundry.AnvilDefaultPort,
		AnvilVerbose:       false,
		EnableEcho:         false,
		HttpAddress:        "127.0.0.1",
		HttpPort:           8080,
		InputBoxAddress:    foundry.InputBoxAddress,
		ApplicationAddress: foundry.ApplicationAddress,
		ApplicationArgs:    nil,
	}
}

// Start nonodo.
func NewNonodoWorker(opts NonodoOpts) (w supervisor.SupervisorWorker, err error) {
	if opts.EnableEcho && len(opts.ApplicationArgs) > 0 {
		return w, ApplicationConflictErr
	}

	w.Name = "nonodo"

	model := model.NewNonodoModel()
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	rollup.Register(e, model)
	inspect.Register(e, model)
	reader.Register(e, model)

	w.Workers = append(w.Workers, foundry.AnvilWorker{
		Port:    opts.AnvilPort,
		Verbose: opts.AnvilVerbose,
	})
	w.Workers = append(w.Workers, inputter.InputterWorker{
		Model:              model,
		Provider:           fmt.Sprintf("ws://127.0.0.1:%v", opts.AnvilPort),
		InputBoxAddress:    common.HexToAddress(opts.InputBoxAddress),
		ApplicationAddress: common.HexToAddress(opts.ApplicationAddress),
	})
	w.Workers = append(w.Workers, supervisor.HttpWorker{
		Address: fmt.Sprintf("%v:%v", opts.HttpAddress, opts.HttpPort),
		Handler: e,
	})
	if len(opts.ApplicationArgs) > 0 {
		w.Workers = append(w.Workers, supervisor.CommandWorker{
			Name:    "app",
			Command: opts.ApplicationArgs[0],
			Args:    opts.ApplicationArgs[1:],
		})
	} else if opts.EnableEcho {
		w.Workers = append(w.Workers, echoapp.EchoAppWorker{
			RollupEndpoint: fmt.Sprintf("http://127.0.0.1:%v/rollup", opts.HttpPort),
		})
	}

	return w, nil
}
