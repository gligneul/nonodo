// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

// This package contains the nonodo run function.
// This is separate from the main package to facilitate testing.
package nonodo

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
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

// Options to nonodo.
type NonodoOpts struct {
	AnvilPort          int
	AnvilVerbose       bool
	BuiltInEcho        bool
	HttpAddress        string
	HttpPort           int
	InputBoxAddress    string
	ApplicationAddress string
}

// Create the options struct with default values.
func NewNonodoOpts() NonodoOpts {
	return NonodoOpts{
		AnvilPort:          foundry.AnvilDefaultPort,
		AnvilVerbose:       false,
		BuiltInEcho:        false,
		HttpAddress:        "127.0.0.1",
		HttpPort:           8080,
		InputBoxAddress:    foundry.InputBoxAddress,
		ApplicationAddress: foundry.ApplicationAddress,
	}
}

// Start nonodo.
func NewService(opts NonodoOpts) supervisor.Service {
	model := model.NewNonodoModel()
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	rollup.Register(e, model)
	inspect.Register(e, model)
	reader.Register(e, model)

	var services []supervisor.Service
	services = append(services, foundry.AnvilService{
		Port:    opts.AnvilPort,
		Verbose: opts.AnvilVerbose,
	})
	services = append(services, inputter.InputterService{
		Model:              model,
		Provider:           fmt.Sprintf("ws://127.0.0.1:%v", opts.AnvilPort),
		InputBoxAddress:    common.HexToAddress(opts.InputBoxAddress),
		ApplicationAddress: common.HexToAddress(opts.ApplicationAddress),
	})
	services = append(services, supervisor.HttpService{
		Address: fmt.Sprintf("%v:%v", opts.HttpAddress, opts.HttpPort),
		Handler: e,
	})
	if opts.BuiltInEcho {
		services = append(services, echoService{
			rollupEndpoint: fmt.Sprintf("http://127.0.0.1:%v/rollup", opts.HttpPort),
		})
	}

	return supervisor.SupervisorService{
		Name:     "nonodo",
		Services: services,
	}
}
