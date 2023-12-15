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
type SignalListenerService struct{}

func (s SignalListenerService) Start(ctx context.Context, ready chan<- struct{}) error {
	log.Print("press Ctrl+C to exit")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	ready <- struct{}{}
	select {
	case sig := <-sigs:
		return fmt.Errorf("received signal: %v", sig)
	case <-ctx.Done():
		return ctx.Err()
	}
}
