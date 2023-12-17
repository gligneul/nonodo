// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package supervisor

import (
	"context"
	"fmt"
)

// Service managed by the supervisor.
type Service interface {
	fmt.Stringer

	// Start the service.
	Start(ctx context.Context, ready chan<- struct{}) error
}

// Run the service, waiting for it to exit
func Serve(ctx context.Context, service Service) error {
	ready := make(chan struct{}, 1)
	return service.Start(ctx, ready)
}
