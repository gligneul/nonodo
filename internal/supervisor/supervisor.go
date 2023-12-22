// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

// This package contains a simple supervisor for goroutine workers.
package supervisor

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

// Timeout when waiting for workers to finish.
const DefaultSupervisorTimeout = time.Second * 5

// Start the workers in order, waiting for each one to be ready before starting the next one.
// When a worker exits, send a cancel signal to all of them and wait for them to finish.
type SupervisorWorker struct {
	Name    string
	Workers []Worker
	Timeout time.Duration
}

func (w SupervisorWorker) String() string {
	return w.Name
}

func (w SupervisorWorker) Start(ctx context.Context, ready chan<- struct{}) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	timeout := w.Timeout
	if timeout == 0 {
		timeout = DefaultSupervisorTimeout
	}

	// Start workers
	var wg sync.WaitGroup
Loop:
	for _, worker := range w.Workers {
		worker := worker
		wg.Add(1)
		innerReady := make(chan struct{})
		go func() {
			defer wg.Done()
			defer cancel()
			err := worker.Start(ctx, innerReady)
			if err != nil && !errors.Is(err, context.Canceled) {
				log.Printf("%v: %v exitted with error: %v", w, worker, err)
			} else {
				log.Printf("%v: %v exitted with success", w, worker)
			}
		}()
		select {
		case <-innerReady:
			log.Printf("%v: %v is ready", w, worker)
		case <-time.After(timeout):
			log.Printf("%v: %v timed out", w, worker)
			cancel()
			break Loop
		case <-ctx.Done():
			break Loop
		}
	}

	// Wait for context to be done
	ready <- struct{}{}
	<-ctx.Done()

	// Wait for all workers
	wait := make(chan struct{})
	go func() {
		wg.Wait()
		wait <- struct{}{}
	}()
	select {
	case <-wait:
		log.Printf("%v: all workers exitted", w)
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("%v: timed out waiting for workers", w)
	}
}
