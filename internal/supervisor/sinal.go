// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package supervisor

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// This service listens signals and exit when it receives them.
type SignalListenerService struct {
	ready chan struct{}
}

func NewSignalListenerService() *SignalListenerService {
	return &SignalListenerService{
		ready: make(chan struct{}),
	}
}

func (s *SignalListenerService) Start(ctx context.Context) error {
	log.Print("press Ctrl+C to exit")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	s.ready <- struct{}{}
	select {
	case sig := <-sigs:
		return fmt.Errorf("received signal: %v", sig)
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *SignalListenerService) Ready() <-chan struct{} {
	return s.ready
}
