// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package supervisor

import (
	"context"
	"fmt"
)

// Worker managed by the supervisor.
type Worker interface {
	fmt.Stringer

	// Start the worker.
	// The worker should send a message to the ready channel when it is ready.
	Start(ctx context.Context, ready chan<- struct{}) error
}
