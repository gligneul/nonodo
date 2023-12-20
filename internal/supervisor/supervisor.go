// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

// This package contains a simple supervisor for goroutines.
package supervisor

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

// Timeout when waiting for services to finish.
const DefaultSupervisorTimeout = time.Second * 5

// Start the sub-services in order, waiting for each one to be ready before starting the next one.
// When a service exits, send a cancel signal to all of them and wait for them to finish.
type SupervisorService struct {
	Name     string
	Services []Service
	Timeout  time.Duration
}

func (s SupervisorService) String() string {
	return s.Name
}

func (s SupervisorService) Start(ctx context.Context, ready chan<- struct{}) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	timeout := s.Timeout
	if timeout == 0 {
		timeout = DefaultSupervisorTimeout
	}

	// Start services
	var wg sync.WaitGroup
Loop:
	for _, service := range s.Services {
		service := service
		wg.Add(1)
		innerReady := make(chan struct{})
		go func() {
			defer wg.Done()
			defer cancel()
			err := service.Start(ctx, innerReady)
			if err != nil && !errors.Is(err, context.Canceled) {
				log.Printf("%v: %v exitted with error: %v", s, service, err)
			}
		}()
		select {
		case <-innerReady:
			log.Printf("%v: %v is ready", s, service)
		case <-time.After(timeout):
			log.Printf("%v: %v timed out", s, service)
			cancel()
			break Loop
		case <-ctx.Done():
			break Loop
		}
	}

	// Wait for context to be done
	ready <- struct{}{}
	if ctx.Err() == nil {
		log.Printf("%v: all services are ready", s)
	}
	<-ctx.Done()

	// Wait for all services
	wait := make(chan struct{})
	go func() {
		wg.Wait()
		wait <- struct{}{}
	}()
	select {
	case <-wait:
		log.Printf("%v: all services exitted", s)
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("%v: timed out waiting for services", s)
	}
}
