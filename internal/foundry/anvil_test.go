// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package foundry

import (
	"context"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

const testTimeout = 5 * time.Second

func TestAnvilService(t *testing.T) {
	ctx, timeoutCancel := context.WithTimeout(context.Background(), testTimeout)
	defer timeoutCancel()

	command := AnvilService{
		Port:    AnvilDefaultPort,
		Verbose: false,
	}

	// start service in goroutine
	serviceCtx, serviceCancel := context.WithCancel(ctx)
	defer serviceCancel()
	ready := make(chan struct{})
	result := make(chan error)
	go func() {
		result <- command.Start(serviceCtx, ready)
	}()

	// wait until service is ready
	select {
	case <-ready:
	case <-ctx.Done():
		t.Error(ctx.Err())
	}

	// send input
	err := AddInput(ctx, common.Hex2Bytes("deadbeef"))
	assert.Nil(t, err)

	// stop service
	serviceCancel()
	select {
	case err := <-result:
		assert.Equal(t, context.Canceled, err)
	case <-ctx.Done():
		t.Error(ctx.Err())
	}
}
