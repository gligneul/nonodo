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

// Return true when the given input should be filtered.
func (f InputFilter) Filter(i *AdvanceInput) bool {
	return (f.IndexGreaterThan != nil && i.Index <= *f.IndexGreaterThan) ||
		(f.IndexLowerThan != nil && i.Index >= *f.IndexLowerThan)
}

// Interface implemented by vouchers, notices, and reports.
type Output interface {
	GetInputIndex() int
}

// Filter outputs (vouchers, notices, and reports).
type OutputFilter struct {
	InputIndex *int
}

// Return true when the given output should be filtered.
func (f OutputFilter) Filter(o Output) bool {
	return f.InputIndex != nil && o.GetInputIndex() != *f.InputIndex
}
