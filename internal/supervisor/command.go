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
	ready chan struct{}
	name  string
	args  []string
	env   []string
	port  int
}

func NewCommandService(name string, args []string, env []string, port int) *CommandService {
	return &CommandService{
		ready: make(chan struct{}),
		name:  name,
		args:  args,
		env:   env,
		port:  port,
	}
}

func (s *CommandService) Start(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, s.name, s.args...)
	cmd.Env = s.env
	cmd.Stderr = s
	cmd.Stdout = s
	cmd.Cancel = func() error {
		err := cmd.Process.Signal(syscall.SIGTERM)
		if err != nil {
			log.Printf("failed to send SIGTERM to %v: %v\n", s.name, err)
		}
		return err
	}
	go s.pollTcp(ctx)
	return cmd.Run()
}

func (s *CommandService) Ready() <-chan struct{} {
	return s.ready
}

// Log the command output.
func (s *CommandService) Write(data []byte) (int, error) {
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		log.Printf("%v: %v", s.name, line)
	}
	return len(data), nil
}

// Polls the command tcp port until it is ready.
func (s *CommandService) pollTcp(ctx context.Context) {
	for {
		conn, err := net.Dial("tcp", fmt.Sprintf("0.0.0.0:%v", s.port))
		if err == nil {
			conn.Close()
			log.Printf("%v is ready", s.name)
			s.ready <- struct{}{}
			return
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(CommandPollInterval):
		}
	}
}
