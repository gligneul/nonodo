// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package model

//
// Query filters
//

// Filter inputs.
type InputFilter struct {
	IndexGreaterThan *int
	IndexLowerThan   *int
}

// Filter outputs (vouchers, notices, and reports).
type OutputFilter struct {
	InputIndex *int
}
