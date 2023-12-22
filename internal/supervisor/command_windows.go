// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

//go:build windows

package supervisor

import (
	"context"
	"log"
	"os/exec"
)

// This worker is responsible for a shell command that runs endlessly.
type CommandWorker struct {
	Name    string
	Command string
	Args    []string
	Env     []string
}

func (w CommandWorker) String() string {
	return w.Name
}

func (w CommandWorker) Start(ctx context.Context, ready chan<- struct{}) error {
	cmd := exec.CommandContext(ctx, w.Command, w.Args...)
	cmd.Env = w.Env
	cmd.Stderr = &commandLogger{name: w.Name}
	cmd.Stdout = &commandLogger{name: w.Name}
	cmd.Cancel = func() error {
		// Sending Interrupt on Windows is not implemented, so we just kill the process.
		// See: https://pkg.go.dev/os#Process.Signal
		err := cmd.Process.Kill()
		if err != nil {
			log.Printf("%v: failed to send SIGTERM: %v", w, err)
		}
		return err
	}
	ready <- struct{}{}
	err := cmd.Run()
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return err
}
