// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package model

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
)

// Interface that represents the state of the rollup.
type rollupsState interface {

	// Finish the current state, saving the result to the model.
	finish(status CompletionStatus)

	// Add voucher to current state.
	addVoucher(destination common.Address, payload []byte) (int, error)

	// Add notice to current state.
	addNotice(payload []byte) (int, error)

	// Add report to current state.
	addReport(payload []byte) error

	// Register exception in current state.
	registerException(payload []byte) error
}

//
// Idle
//

// In the idle state, the model waits for an finish request from the rollups API.
type rollupsStateIdle struct{}

func newRollupsStateIdle() *rollupsStateIdle {
	return &rollupsStateIdle{}
}

func (s *rollupsStateIdle) finish(status CompletionStatus) {
	// Do nothing
}

func (s *rollupsStateIdle) addVoucher(destination common.Address, payload []byte) (int, error) {
	return 0, fmt.Errorf("cannot add voucher in current state")
}

func (s *rollupsStateIdle) addNotice(payload []byte) (int, error) {
	return 0, fmt.Errorf("cannot add notice in current state")
}

func (s *rollupsStateIdle) addReport(payload []byte) error {
	return fmt.Errorf("cannot add report in current state")
}

func (s *rollupsStateIdle) registerException(payload []byte) error {
	return fmt.Errorf("cannot register exception in current state")
}

//
// Advance
//

// In the advance state, the model accumulates the outputs from an advance.
type rollupsStateAdvance struct {
	input    *AdvanceInput
	vouchers []Voucher
	notices  []Notice
	reports  []Report
}

func newRollupsStateAdvance(input *AdvanceInput) *rollupsStateAdvance {
	log.Printf("nonodo: processing advance input %v", input.Index)
	return &rollupsStateAdvance{
		input: input,
	}
}

func (s *rollupsStateAdvance) finish(status CompletionStatus) {
	s.input.Status = status
	if status == CompletionStatusAccepted {
		s.input.Vouchers = s.vouchers
		s.input.Notices = s.notices
	}
	s.input.Reports = s.reports
	log.Printf("nonodo: finished advance")
}

func (s *rollupsStateAdvance) addVoucher(destination common.Address, payload []byte) (int, error) {
	index := len(s.vouchers)
	voucher := Voucher{
		Index:       index,
		InputIndex:  s.input.Index,
		Destination: destination,
		Payload:     payload,
	}
	s.vouchers = append(s.vouchers, voucher)
	log.Printf("nonodo: added voucher index=%v destination=%v payload=0x%x",
		index, destination, payload)
	return index, nil
}

func (s *rollupsStateAdvance) addNotice(payload []byte) (int, error) {
	index := len(s.notices)
	notice := Notice{
		Index:      index,
		InputIndex: s.input.Index,
		Payload:    payload,
	}
	s.notices = append(s.notices, notice)
	log.Printf("nonodo: added notice index=%v payload=0x%x", index, payload)
	return index, nil
}

func (s *rollupsStateAdvance) addReport(payload []byte) error {
	index := len(s.reports)
	report := Report{
		Index:      index,
		InputIndex: s.input.Index,
		Payload:    payload,
	}
	s.reports = append(s.reports, report)
	log.Printf("nonodo: added report index=%v payload=0x%x", index, payload)
	return nil
}

func (s *rollupsStateAdvance) registerException(payload []byte) error {
	s.input.Status = CompletionStatusException
	s.input.Reports = s.reports
	s.input.Exception = payload
	log.Printf("nonodo: finished advance with exception")
	return nil
}

//
// Inspect
//

// In the inspect state, the model accumulates the reports from an inspect.
type rollupsStateInspect struct {
	input                   *InspectInput
	reports                 []Report
	getProccessedInputCount func() int
}

func newRollupsStateInspect(
	input *InspectInput,
	getProccessedInputCount func() int,
) *rollupsStateInspect {
	log.Printf("nonodo: processing inspect input %v", input.Index)
	return &rollupsStateInspect{
		input:                   input,
		getProccessedInputCount: getProccessedInputCount,
	}
}

func (s *rollupsStateInspect) finish(status CompletionStatus) {
	s.input.Status = status
	s.input.ProccessedInputCount = s.getProccessedInputCount()
	s.input.Reports = s.reports
	log.Printf("nonodo: finished inspect")
}

func (s *rollupsStateInspect) addVoucher(destination common.Address, payload []byte) (int, error) {
	return 0, fmt.Errorf("cannot add voucher in current state")
}

func (s *rollupsStateInspect) addNotice(payload []byte) (int, error) {
	return 0, fmt.Errorf("cannot add notice in current state")
}

func (s *rollupsStateInspect) addReport(payload []byte) error {
	index := len(s.reports)
	report := Report{
		Index:      index,
		InputIndex: s.input.Index,
		Payload:    payload,
	}
	s.reports = append(s.reports, report)
	log.Printf("nonodo: added report index=%v payload=0x%x", index, payload)
	return nil
}

func (s *rollupsStateInspect) registerException(payload []byte) error {
	s.input.Status = CompletionStatusException
	s.input.ProccessedInputCount = s.getProccessedInputCount()
	s.input.Reports = s.reports
	s.input.Exception = payload
	log.Printf("nonodo: finished inspect with exception")
	return nil
}
