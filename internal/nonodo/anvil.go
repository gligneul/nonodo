// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package nonodo

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/gligneul/nonodo/internal/supervisor"
)

// Generate the devnet state and embed it in the Go binary.
//
//go:generate go run ./gen-devnet-state
//go:embed anvil_state.json
var anvilState []byte

// Start the anvil process in the host machine.
type anvilService struct {
	port      int
	blockTime int
	verbose   bool
}

func (s anvilService) Start(ctx context.Context, ready chan<- struct{}) error {
	// create temp dir
	tempDir, err := os.MkdirTemp("", "")
	if err != nil {
		return fmt.Errorf("anvil: failed to create temp dir: %v", err)
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
	err = os.WriteFile(stateFile, anvilState, permissions)
	if err != nil {
		return fmt.Errorf("anvil: failed to write state file: %v", err)
	}

	// start command
	args := []string{
		"--port", fmt.Sprint(s.port),
		"--block-time", fmt.Sprint(s.blockTime),
		"--load-state", stateFile,
	}
	if !s.verbose {
		args = append(args, "--silent")
	}
	command := supervisor.CommandService{
		Name: "anvil",
		Args: args,
		Port: s.port,
	}
	return command.Start(ctx, ready)
}
