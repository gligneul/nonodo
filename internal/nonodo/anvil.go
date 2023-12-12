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
	ready   chan struct{}
	command *supervisor.CommandService
	cleanup func()
}

func newAnvilService(opts NonodoOpts) *anvilService {
	stateFile, cleanup := loadStateFile()
	args := []string{
		"--port", fmt.Sprint(opts.AnvilPort),
		"--block-time", fmt.Sprint(opts.AnvilBlockTime),
		"--load-state", stateFile,
	}
	if !opts.AnvilVerbose {
		args = append(args, "--silent")
	}
	return &anvilService{
		ready:   make(chan struct{}),
		command: supervisor.NewCommandService("anvil", args, nil, opts.AnvilPort),
		cleanup: cleanup,
	}
}

func (s *anvilService) Start(ctx context.Context) error {
	go func() {
		// cleanup after the internal service is ready
		defer s.cleanup()
		select {
		case <-s.command.Ready():
			s.ready <- struct{}{}
		case <-ctx.Done():
		}
	}()
	return s.command.Start(ctx)
}

func (s *anvilService) Ready() <-chan struct{} {
	return s.ready
}

// Create a temporary file with the anvil state.
// Return a function to delete this file.
func loadStateFile() (string, func()) {
	tempDir, err := os.MkdirTemp("", "")
	if err != nil {
		panic(err)
	}
	path := path.Join(tempDir, "anvil_state.json")
	const permissions = 0644
	err = os.WriteFile(path, anvilState, permissions)
	if err != nil {
		panic(err)
	}
	cleanup := func() {
		err := os.RemoveAll(tempDir)
		if err != nil {
			log.Printf("failed to remove temp file: %v", err)
		}
	}
	return path, cleanup
}
