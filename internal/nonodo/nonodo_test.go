// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package nonodo

import (
	"context"
	"crypto/rand"
	"fmt"
	"testing"
	"time"

	"github.com/Khan/genqlient/graphql"
	"github.com/gligneul/nonodo/internal/foundry"
	"github.com/stretchr/testify/suite"
)

type NonodoSuite struct {
	suite.Suite
	ctx           context.Context
	timeoutCancel context.CancelFunc
	serviceCancel context.CancelFunc
	serviceResult chan error
	client        graphql.Client
}

func (s *NonodoSuite) SetupTest() {
	const testTimeout = 5 * time.Second
	s.ctx, s.timeoutCancel = context.WithTimeout(context.Background(), testTimeout)
	s.serviceResult = make(chan error)

	var serviceCtx context.Context
	serviceCtx, s.serviceCancel = context.WithCancel(s.ctx)

	opts := NewNonodoOpts()
	opts.BuiltInEcho = true
	service := NewService(opts)

	ready := make(chan struct{})
	go func() {
		s.serviceResult <- service.Start(serviceCtx, ready)
	}()
	select {
	case <-s.ctx.Done():
		s.Fail(s.ctx.Err().Error())
	case err := <-s.serviceResult:
		if s.Error(err) {
			s.Fail(err.Error())
		} else {
			s.Fail("service exited without an erro")
		}
	case <-ready:
	}

	endpoint := fmt.Sprintf("http://%v:%v/graphql", opts.HttpAddress, opts.HttpPort)
	s.client = graphql.NewClient(endpoint, nil)
}

func (s *NonodoSuite) TearDownTest() {
	s.serviceCancel()
	select {
	case <-s.ctx.Done():
		s.Fail(s.ctx.Err().Error())
	case err := <-s.serviceResult:
		s.Nil(err)
	}
}

func (s *NonodoSuite) TestItProcessesInputs() {
	// send inputs
	const n = 1
	payloads := make([][]byte, n)
	for i := 0; i < n; i++ {
		payloads[i] = make([]byte, 32)
		_, err := rand.Read(payloads[i])
		s.Nil(err)

		err = foundry.AddInput(payloads[i])
		s.Nil(err)
	}

	// wait for input to be ready
	const pollRetry = 100
	const pollInterval = 10 * time.Millisecond
	for i := 0; i < pollRetry; i++ {
		//readerclient.
	}
}

func TestNonodoSuite(t *testing.T) {
	suite.Run(t, &NonodoSuite{})
}
