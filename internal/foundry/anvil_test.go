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

func TestAnvilWorker(t *testing.T) {
	ctx, timeoutCancel := context.WithTimeout(context.Background(), testTimeout)
	defer timeoutCancel()

	w := AnvilWorker{
		Port:    AnvilDefaultPort,
		Verbose: false,
	}

	// start worker in goroutine
	workerCtx, workerCancel := context.WithCancel(ctx)
	defer workerCancel()
	ready := make(chan struct{})
	result := make(chan error)
	go func() {
		result <- w.Start(workerCtx, ready)
	}()

	// wait until worker is ready
	select {
	case <-ready:
	case <-ctx.Done():
		t.Error(ctx.Err())
	}

	// send input
	err := AddInput(ctx, common.Hex2Bytes("deadbeef"))
	assert.Nil(t, err)

	// stop worker
	workerCancel()
	select {
	case err := <-result:
		assert.Equal(t, context.Canceled, err)
	case <-ctx.Done():
		t.Error(ctx.Err())
	}
}
