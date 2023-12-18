// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package supervisor

import (
	"context"
	"fmt"
	"log"
	"net"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

// Poll interval when checking whether the command is ready.
const CommandPollInterval = 100 * time.Millisecond

// This service is responsible for a shell command that runs endlessly.
// The service polls the given port to know when it is ready.
type CommandService struct {
	Name string
	Args []string
	Env  []string
	Port int
}

func (s CommandService) String() string {
	return s.Name
}

func (s CommandService) Start(ctx context.Context, ready chan<- struct{}) error {
	cmd := exec.CommandContext(ctx, s.Name, s.Args...)
	cmd.Env = s.Env
	cmd.Stderr = commandLogger{s.Name}
	cmd.Stdout = commandLogger{s.Name}
	cmd.Cancel = func() error {
		err := cmd.Process.Signal(syscall.SIGTERM)
		if err != nil {
			log.Printf("failed to send SIGTERM to %v: %v", s.Name, err)
		}
		return err
	}
	go s.pollTcp(ctx, ready)
	return cmd.Run()
}

// Polls the command tcp port until it is ready.
func (s CommandService) pollTcp(ctx context.Context, ready chan<- struct{}) {
	for {
		conn, err := net.Dial("tcp", fmt.Sprintf("0.0.0.0:%v", s.Port))
		if err == nil {
			conn.Close()
			ready <- struct{}{}
			return
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(CommandPollInterval):
		}
	}
}

type commandLogger struct {
	Name string
}

// Log the command output.
func (s commandLogger) Write(data []byte) (int, error) {
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		log.Printf("%v: %v", s.Name, line)
	}
	return len(data), nil
}
