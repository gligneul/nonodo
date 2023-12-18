// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package foundry

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/gligneul/nonodo/internal/supervisor"
)

// Default port for the Ethereum node.
const AnvilDefaultPort = 8545

// Start the anvil process in the host machine.
type AnvilService struct {
	Port      int
	BlockTime int
	Verbose   bool
}

func (s AnvilService) String() string {
	return "anvil"
}

func (s AnvilService) Start(ctx context.Context, ready chan<- struct{}) error {
	// create temp dir
	tempDir, err := os.MkdirTemp("", "")
	if err != nil {
		return fmt.Errorf("anvil: failed to create temp dir: %w", err)
	}
	defer func() {
		err := os.RemoveAll(tempDir)
		if err != nil {
			log.Printf("anvil: failed to remove temp file: %v", err)
		}
	}()

	// create state file in temp dir
	stateFile := path.Join(tempDir, "anvil_state.json")
	const permissions = 0644
	err = os.WriteFile(stateFile, devnetState, permissions)
	if err != nil {
		return fmt.Errorf("anvil: failed to write state file: %w", err)
	}

	// start command
	args := []string{
		"--port", fmt.Sprint(s.Port),
		"--block-time", fmt.Sprint(s.BlockTime),
		"--load-state", stateFile,
	}
	if !s.Verbose {
		args = append(args, "--silent")
	}
	command := supervisor.CommandService{
		Name: "anvil",
		Args: args,
		Port: s.Port,
	}
	return command.Start(ctx, ready)
}
