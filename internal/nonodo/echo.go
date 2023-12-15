// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package nonodo

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/gligneul/nonodo/internal/inspect"
	"github.com/gligneul/nonodo/internal/model"
	"github.com/gligneul/nonodo/internal/rollup"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// The inputter reads inputs from anvil and puts them in the model.
type echoService struct {
	ready   chan struct{}
	address string
	model   *model.NonodoModel
}

// Creates a new inputter from opts.
func newEchoService(model *model.NonodoModel, address string) *echoService {
	return &echoService{
		ready:   make(chan struct{}),
		address: address,
		model:   model,
	}
}

func (s *echoService) Start(ctx context.Context) error {
	// setup routes
	e := echo.New()
	e.Use(middleware.CORS())
	rollup.Register(e, s.model)
	inspect.Register(e, s.model)

	// create server
	server := http.Server{
		Addr:    s.address,
		Handler: e,
	}
	ln, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}
	log.Printf("listening on %v", ln.Addr())
	s.ready <- struct{}{}

	// create goroutine to shutdown server
	go func() {
		<-ctx.Done()
		err := server.Shutdown(ctx)
		if err != nil {
			log.Printf("error shutting down http server: %v", err)
		}
	}()

	// serve
	err = server.Serve(ln)
	if err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *echoService) Ready() <-chan struct{} {
	return s.ready
}
