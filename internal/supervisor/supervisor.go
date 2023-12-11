// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

// This package contains a simple supervisor for goroutines.
package supervisor

import (
	"context"
	"log"
	"sync"
	"time"
)

// Timeout when waiting for services to finish.
const ServiceTimeout = time.Second * 5

// Service managed by the supervisor function.
type Service interface {

	// Start the service.
	Start(ctx context.Context) error

	// The service should send a message when it is ready.
	Ready() <-chan struct{}
}

// Start the services in order, waiting for each one to be ready before starting the next one.
// When a service exits, send a cancel signal to all of them and wait for them to finish.
func Start(ctx context.Context, services []Service) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Start services
	var wg sync.WaitGroup
	for _, service := range services {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer cancel()
			err := service.Start(ctx)
			if err != nil {
				log.Print(err)
			}
		}()
		select {
		case <-service.Ready():
		case <-time.After(ServiceTimeout):
			log.Print("service timed out")
			cancel()
			break
		case <-ctx.Done():
			break
		}
	}

	// Wait for context to be done
	if ctx.Err() == nil {
		log.Print("ready")
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
		log.Print("all services were shutdown")
	case <-time.After(ServiceTimeout):
		log.Print("exited after a timeout")
	}
}
