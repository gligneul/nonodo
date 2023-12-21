// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package nonodo

import (
	"context"
	"crypto/rand"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Khan/genqlient/graphql"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gligneul/nonodo/internal/foundry"
	"github.com/gligneul/nonodo/internal/readerclient"
	"github.com/stretchr/testify/suite"
)

type NonodoSuite struct {
	suite.Suite
	ctx           context.Context
	timeoutCancel context.CancelFunc
	serviceCancel context.CancelFunc
	serviceResult chan error
	graphqlClient graphql.Client
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
		s.Fail("context error: %v", s.ctx.Err())
	case err := <-s.serviceResult:
		s.Fail("service exited before being ready: %v", err)
	case <-ready:
	}

	endpoint := fmt.Sprintf("http://%v:%v/graphql", opts.HttpAddress, opts.HttpPort)
	s.graphqlClient = graphql.NewClient(endpoint, nil)
}

func (s *NonodoSuite) TearDownTest() {
	s.serviceCancel()
	select {
	case <-s.ctx.Done():
		s.Fail("context error: %v", s.ctx.Err())
	case err := <-s.serviceResult:
		s.Nil(err)
	}
	s.timeoutCancel()
}

// Wait for the given input to be ready.
func (s *NonodoSuite) waitForAdvanceInput(inputIndex int) error {
	const pollRetries = 100
	const pollInterval = 10 * time.Millisecond
	for i := 0; i < pollRetries; i++ {
		result, err := readerclient.InputStatus(s.ctx, s.graphqlClient, inputIndex)
		if err != nil && !strings.Contains(err.Error(), "input not found") {
			return fmt.Errorf("failed to get input status: %w", err)
		}
		if result.Input.Status == readerclient.CompletionStatusAccepted {
			return nil
		}
		select {
		case <-s.ctx.Done():
			return s.ctx.Err()
		case <-time.After(pollInterval):
		}
	}
	return fmt.Errorf("input never got ready")
}

// Decode the hex string into bytes.
func (s *NonodoSuite) decodeHex(value string) []byte {
	bytes, err := hexutil.Decode(value)
	s.Require().Nil(err)
	return bytes
}

func (s *NonodoSuite) TestItProcessesInputs() {
	s.T().Log("sending inputs")
	const n = 3
	payloads := make([][]byte, n)
	for i := 0; i < n; i++ {
		payloads[i] = make([]byte, 32)
		_, err := rand.Read(payloads[i])
		s.Require().Nil(err)

		err = foundry.AddInput(s.ctx, payloads[i])
		s.Require().Nil(err)
	}

	s.T().Log("waiting until last input is ready")
	err := s.waitForAdvanceInput(n - 1)
	s.Require().Nil(err)

	s.T().Log("verifying node state")
	response, err := readerclient.State(s.ctx, s.graphqlClient)
	s.Require().Nil(err)
	for i := 0; i < n; i++ {
		input := response.Inputs.Edges[i].Node
		s.Equal(i, input.Index)
		s.Equal(payloads[i], s.decodeHex(input.Payload))
		s.Equal(payloads[i], s.decodeHex(input.Payload))
		s.Equal(foundry.SenderAddress[:], s.decodeHex(input.MsgSender))
		voucher := input.Vouchers.Edges[0].Node
		s.Equal(payloads[i], s.decodeHex(voucher.Payload))
		s.Equal(foundry.SenderAddress[:], s.decodeHex(voucher.Destination))
		s.Equal(payloads[i], s.decodeHex(input.Notices.Edges[0].Node.Payload))
		s.Equal(payloads[i], s.decodeHex(input.Reports.Edges[0].Node.Payload))
	}
}

func TestNonodoSuite(t *testing.T) {
	suite.Run(t, &NonodoSuite{})
}
