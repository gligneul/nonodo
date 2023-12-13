// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package nonodo

import (
	"context"
	"fmt"
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
	ready chan struct{}
	port  int
	model *model.NonodoModel
}

// Creates a new inputter from opts.
func newEchoService(model *model.NonodoModel, port int) *echoService {
	return &echoService{
		ready: make(chan struct{}),
		port:  port,
		model: model,
	}
}

func (s *echoService) Start(ctx context.Context) error {
	// setup routes
	e := echo.New()
	e.Use(middleware.CORS())
	rollup.Register(e, s.model)
	inspect.Register(e, s.model)

	// create server
	addr := fmt.Sprintf("127.0.0.1:%d", s.port)
	server := http.Server{
		Addr:    addr,
		Handler: e,
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	log.Printf("listening on %v", addr)
	s.ready <- struct{}{}

	// create goroutine to shutdown server
	go func() {
		<-ctx.Done()
		server.Shutdown(ctx)
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
