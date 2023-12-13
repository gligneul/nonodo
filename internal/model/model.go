// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

// The nonodo model uses a shared-memory paradigm to synchronize between threads.
package model

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// -------------------------------------------------------------------------------------------------
// Rollups types
// -------------------------------------------------------------------------------------------------

// Rollups voucher type.
type Voucher struct {
	Index       int
	InputIndex  int
	Destination common.Address
	Payload     []byte
}

// Rollups notice type.
type Notice struct {
	Index      int
	InputIndex int
	Payload    []byte
}

// Rollups report type.
type Report struct {
	Index      int
	InputIndex int
	Payload    []byte
}

// Completion status for inputs.
type CompletionStatus int

const (
	CompletionStatusUnprocessed CompletionStatus = iota
	CompletionStatusAccepted
	CompletionStatusRejected
	CompletionStatusException
)

// Rollups input, which can be advance or inspect.
type Input interface{}

// Rollups advance input type.
type AdvanceInput struct {
	Index       int
	Status      CompletionStatus
	MsgSender   common.Address
	Payload     []byte
	BlockNumber uint64
	Timestamp   time.Time
	Vouchers    []Voucher
	Notices     []Notice
	Reports     []Report
	Exception   []byte
}

// Rollups inspect input type.
type InspectInput struct {
	Index                int
	Status               CompletionStatus
	Payload              []byte
	ProccessedInputCount int
	Reports              []Report
	Exception            []byte
}

// -------------------------------------------------------------------------------------------------
// Query filters
// -------------------------------------------------------------------------------------------------

// Filter inputs.
type InputFilter struct {
	IndexGreaterThan *int
	IndexLowerThan   *int
}

// Filter outputs (vouchers, notices, and reports).
type OutputFilter struct {
	InputIndex *int
}

// -------------------------------------------------------------------------------------------------
// State
// -------------------------------------------------------------------------------------------------

// Interface that represents the state of the nonodo model.
type modelState interface{}

// In the idle state, the model waits for an finish request from the rollups API.
type modelStateIdle struct {
	Reports []Report
}

// In the advance state, the model accumulates the outputs from an advance.
type modelStateAdvance struct {
	inputIndex int
	vouchers   []Voucher
	notices    []Notice
	reports    []Report
}

// In the inspect state, the model accumulates the reports from an inspect.
type modelStateInspect struct {
	inputIndex int
	reports    []Report
}

// -------------------------------------------------------------------------------------------------
// NonodoModel
// -------------------------------------------------------------------------------------------------

// Nonodo model shared among the internal services.
type NonodoModel struct {
	mutex    sync.Mutex
	advances []AdvanceInput
	inspects []InspectInput
	state    modelState
}

// Create a new model.
func NewNonodoModel() *NonodoModel {
	return &NonodoModel{
		state: &modelStateIdle{},
	}
}

// -------------------------------------------------------------------------------------------------
// Methods for Inputter
// -------------------------------------------------------------------------------------------------

// Add an advance input to the model.
func (m *NonodoModel) AddAdvanceInput(
	inputIndex int,
	sender common.Address,
	payload []byte,
	blockNumber uint64,
	timestamp time.Time,
) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if inputIndex != len(m.advances) {
		panic(fmt.Sprintf("invalid input index: %v; want %v", inputIndex, len(m.advances)))
	}

	input := AdvanceInput{
		Index:       inputIndex,
		Status:      CompletionStatusUnprocessed,
		MsgSender:   sender,
		Payload:     payload,
		Timestamp:   timestamp,
		BlockNumber: blockNumber,
	}
	m.advances = append(m.advances, input)
	log.Printf("nonodo: added advance input: index=%v sender=%v payload=0x%x",
		input.Index, input.MsgSender, input.Payload)
}

// -------------------------------------------------------------------------------------------------
// Methods for Inspector
// -------------------------------------------------------------------------------------------------

// Add an inspect input to the model.
// Return the inspect input index that should be used for polling.
func (m *NonodoModel) AddInspectInput(payload []byte) int {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	index := len(m.inspects)
	input := InspectInput{
		Index:   index,
		Status:  CompletionStatusUnprocessed,
		Payload: payload,
	}
	m.inspects = append(m.inspects, input)
	log.Printf("nonodo: added inspect input: index=%v payload=0x%x", input.Index, input.Payload)

	return index
}

// Get the inspect input from the model.
func (m *NonodoModel) GetInspectInput(index int) InspectInput {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if index >= len(m.inspects) {
		panic(fmt.Sprintf("invalid inspect input index: %v", index))
	}
	return m.inspects[index]
}

// -------------------------------------------------------------------------------------------------
// Methods for Rollups
// -------------------------------------------------------------------------------------------------

// Finish the current input and get the next one.
// If there is no input to be processed return nil.
func (m *NonodoModel) FinishAndGetNext(accepted bool) Input {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	//
	// finish current input
	//

	var status CompletionStatus
	if accepted {
		status = CompletionStatusAccepted
	} else {
		status = CompletionStatusRejected
	}

	switch state := m.state.(type) {
	case *modelStateIdle:
		// nothing to do

	case *modelStateAdvance:
		m.advances[state.inputIndex].Status = status
		if accepted {
			m.advances[state.inputIndex].Vouchers = state.vouchers
			m.advances[state.inputIndex].Notices = state.notices
		}
		m.advances[state.inputIndex].Reports = state.reports
		log.Printf("nonodo: finished advance input %v", state.inputIndex)

	case *modelStateInspect:
		m.inspects[state.inputIndex].Status = status
		m.inspects[state.inputIndex].ProccessedInputCount = m.getProccessedInputCount()
		m.inspects[state.inputIndex].Reports = state.reports
		log.Printf("nonodo: finished inspect input %v", state.inputIndex)

	default:
		panic("invalid state")
	}

	//
	// get next input
	//

	inspectInput, ok := m.getUnprocessedInspectInput()
	if ok {
		m.state = &modelStateInspect{
			inputIndex: inspectInput.Index,
		}
		log.Printf("nonodo: processing inspect input %v", inspectInput.Index)
		return inspectInput
	}

	advanceInput, ok := m.getUnprocessedAdvanceInput()
	if ok {
		m.state = &modelStateAdvance{
			inputIndex: advanceInput.Index,
		}
		log.Printf("nonodo: processing advance input %v", advanceInput.Index)
		return advanceInput
	}

	m.state = &modelStateIdle{}
	return nil
}

// Add a voucher to the model.
// Return the voucher index within the input.
// Return an error if the state isn't advance.
func (m *NonodoModel) AddVoucher(destination common.Address, payload []byte) (int, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	state, ok := m.state.(*modelStateAdvance)
	if !ok {
		return 0, fmt.Errorf("cannot add voucher in current state")
	}

	index := len(state.vouchers)
	voucher := Voucher{
		Index:       index,
		InputIndex:  state.inputIndex,
		Destination: destination,
		Payload:     payload,
	}
	state.vouchers = append(state.vouchers, voucher)
	log.Printf("nonodo: added voucher %v", index)

	return index, nil
}

// Add a notice to the model.
// Return the notice index within the input.
// Return an error if the state isn't advance.
func (m *NonodoModel) AddNotice(payload []byte) (int, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	state, ok := m.state.(*modelStateAdvance)
	if !ok {
		return 0, fmt.Errorf("cannot add notice in current state")
	}

	index := len(state.notices)
	notice := Notice{
		Index:      index,
		InputIndex: state.inputIndex,
		Payload:    payload,
	}
	state.notices = append(state.notices, notice)
	log.Printf("nonodo: added notice %v", index)

	return index, nil
}

// Add a report to the model.
// Return an error if the state isn't advance or inspect.
func (m *NonodoModel) AddReport(payload []byte) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var inputIndex int
	var reports *[]Report
	switch state := m.state.(type) {
	case *modelStateIdle:
		return fmt.Errorf("cannot add report in current state")
	case *modelStateAdvance:
		inputIndex = state.inputIndex
		reports = &state.reports
	case *modelStateInspect:
		inputIndex = state.inputIndex
		reports = &state.reports
	default:
		panic("invalid state")
	}

	index := len(*reports)
	report := Report{
		Index:      index,
		InputIndex: inputIndex,
		Payload:    payload,
	}
	*reports = append(*reports, report)
	log.Printf("nonodo: added report %v", index)

	return nil
}

// Finish the current input with an exception.
// Return an error if the state isn't advance or inspect.
func (m *NonodoModel) RaiseException(payload []byte) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	switch state := m.state.(type) {
	case *modelStateIdle:
		return fmt.Errorf("cannot raise exception in current state")

	case *modelStateAdvance:
		m.advances[state.inputIndex].Status = CompletionStatusException
		m.advances[state.inputIndex].Reports = state.reports
		log.Printf("nonodo: finished advance input %v with exception", state.inputIndex)

	case *modelStateInspect:
		m.inspects[state.inputIndex].Status = CompletionStatusException
		m.inspects[state.inputIndex].ProccessedInputCount = m.getProccessedInputCount()
		m.inspects[state.inputIndex].Reports = state.reports
		log.Printf("nonodo: finished inspect input %v with exception", state.inputIndex)

	default:
		panic("invalid state")
	}

	// set state to idle
	m.state = &modelStateIdle{}
	return nil
}

// -------------------------------------------------------------------------------------------------
// Methods for Reader
// -------------------------------------------------------------------------------------------------

// Get the advance input for the given index.
// Return false if not found.
func (m *NonodoModel) GetAdvanceInput(index int) (AdvanceInput, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if index >= len(m.advances) {
		var input AdvanceInput
		return input, false
	}
	return m.advances[index], true
}

// Get the voucher for the given index and input index.
// Return false if not found.
func (m *NonodoModel) GetVoucher(voucherIndex, inputIndex int) (Voucher, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if inputIndex >= len(m.advances) ||
		voucherIndex >= len(m.advances[inputIndex].Vouchers) {
		var voucher Voucher
		return voucher, false
	}
	return m.advances[inputIndex].Vouchers[voucherIndex], true
}

// Get the notice for the given index and input index.
// Return false if not found.
func (m *NonodoModel) GetNotice(noticeIndex, inputIndex int) (Notice, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if inputIndex >= len(m.advances) ||
		noticeIndex >= len(m.advances[inputIndex].Notices) {
		var notice Notice
		return notice, false
	}
	return m.advances[inputIndex].Notices[noticeIndex], true
}

// Get the report for the given index and input index.
// Return false if not found.
func (m *NonodoModel) GetReport(reportIndex, inputIndex int) (Report, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if inputIndex >= len(m.advances) ||
		reportIndex >= len(m.advances[inputIndex].Reports) {
		var report Report
		return report, false
	}
	return m.advances[inputIndex].Reports[reportIndex], true
}

// Get the inputs given the filter and pagination parameters.
func (m *NonodoModel) GetInputs(filter InputFilter, offset int, limit int) []AdvanceInput {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var inputs []AdvanceInput
	for _, input := range m.advances {
		if filter.IndexGreaterThan != nil &&
			input.Index <= *filter.IndexGreaterThan {
			continue
		}
		if filter.IndexLowerThan != nil &&
			input.Index >= *filter.IndexLowerThan {
			continue
		}
		inputs = append(inputs, input)
	}
	return paginate(inputs, offset, limit)
}

// Get the vouchers given the filter and pagination parameters.
func (m *NonodoModel) GetVouchers(filter OutputFilter, limit int, offset int) []Voucher {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var vouchers []Voucher
	for _, input := range m.advances {
		for _, voucher := range input.Vouchers {
			if filter.InputIndex != nil &&
				voucher.InputIndex != *filter.InputIndex {
				continue
			}
			vouchers = append(vouchers, voucher)
		}
	}
	return paginate(vouchers, offset, limit)
}

// Get the notices given the filter and pagination parameters.
func (m *NonodoModel) GetNotices(filter OutputFilter, limit int, offset int) []Notice {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var notices []Notice
	for _, input := range m.advances {
		for _, notice := range input.Notices {
			if filter.InputIndex != nil &&
				notice.InputIndex != *filter.InputIndex {
				continue
			}
			notices = append(notices, notice)
		}
	}
	return paginate(notices, offset, limit)
}

// Get the reports given the filter and pagination parameters.
func (m *NonodoModel) GetReports(filter OutputFilter, limit int, offset int) []Report {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var reports []Report
	for _, input := range m.advances {
		for _, report := range input.Reports {
			if filter.InputIndex != nil &&
				report.InputIndex != *filter.InputIndex {
				continue
			}
			reports = append(reports, report)
		}
	}
	return paginate(reports, offset, limit)
}

// -------------------------------------------------------------------------------------------------
// Auxiliary Methods
// -------------------------------------------------------------------------------------------------

func (m *NonodoModel) getProccessedInputCount() int {
	n := 0
	for _, input := range m.advances {
		if input.Status == CompletionStatusUnprocessed {
			break
		}
		n++
	}
	return n
}

func (m *NonodoModel) getUnprocessedAdvanceInput() (AdvanceInput, bool) {
	for _, input := range m.advances {
		if input.Status == CompletionStatusUnprocessed {
			return input, true
		}
	}
	var input AdvanceInput
	return input, false
}

func (m *NonodoModel) getUnprocessedInspectInput() (InspectInput, bool) {
	for _, input := range m.inspects {
		if input.Status == CompletionStatusUnprocessed {
			return input, true
		}
	}
	var input InspectInput
	return input, false
}

func paginate[T any](slice []T, offset int, limit int) []T {
	if offset >= len(slice) {
		return nil
	}
	upperBound := offset + limit
	if upperBound > len(slice) {
		upperBound = len(slice)
	}
	return slice[offset:upperBound]
}
