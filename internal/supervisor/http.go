// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package supervisor

import (
	"context"
	"log"
	"net"
	"net/http"
)

// The HTTP service starts a HTTP server.
type HttpService struct {
	Address string
	Handler http.Handler
}

func (s HttpService) String() string {
	return "http"
}

func (s HttpService) Start(ctx context.Context, ready chan<- struct{}) error {
	// create server
	server := http.Server{
		Addr:    s.Address,
		Handler: s.Handler,
	}
	ln, err := net.Listen("tcp", s.Address)
	if err != nil {
		return err
	}
	log.Printf("%s: listening on %v", s, ln.Addr())
	ready <- struct{}{}

	// create goroutine to shutdown server
	go func() {
		<-ctx.Done()
		err := server.Shutdown(ctx)
		if err != nil {
			log.Printf("%v: error shutting down http server: %v", s, err)
		}
	}()

	// serve
	err = server.Serve(ln)
	if err != http.ErrServerClosed {
		return err
	}
	return nil
}
