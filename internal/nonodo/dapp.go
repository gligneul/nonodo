// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package nonodo

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gligneul/nonodo/internal/rollup"
)

// The dappService uses the rollup API to implement an echo dapp.
// This services uses the API rather than talking directly to the model so it can be used in
// integration tests.
type dappService struct {
	ready          chan struct{}
	rollupEndpoint string
}

// Creates a new echo dapp.
func newDappService(rollupEndpoint string) *dappService {
	return &dappService{
		ready:          make(chan struct{}),
		rollupEndpoint: rollupEndpoint,
	}
}

func (s *dappService) Start(ctx context.Context) error {
	client, err := rollup.NewClientWithResponses(s.rollupEndpoint)
	if err != nil {
		return fmt.Errorf("dapp: %v", err)
	}

	s.ready <- struct{}{}
	log.Print("starting built-in echo dapp")

	finishReq := rollup.Finish{
		Status: rollup.Accept,
	}
	for {
		finishResp, err := client.FinishWithResponse(ctx, finishReq)
		if err != nil {
			return fmt.Errorf("dapp: %v", err)
		}
		if finishResp.StatusCode() == http.StatusAccepted {
			continue
		}
		if finishResp.StatusCode() != http.StatusOK {
			return fmt.Errorf("dapp: invalid finish response: status=%v body=`%v`",
				finishResp.StatusCode(), string(finishResp.Body))
		}
		finishBody := finishResp.JSON200
		if finishBody == nil {
			return fmt.Errorf("dapp: missing finish response body")
		}
		switch finishBody.RequestType {
		case rollup.AdvanceState:
			advance, err := finishBody.Data.AsAdvance()
			if err != nil {
				return fmt.Errorf("dapp: failed to parser advance: %v", err)
			}
			if err := handleAdvance(ctx, client, advance); err != nil {
				return err
			}
		case rollup.InspectState:
			inspect, err := finishBody.Data.AsInspect()
			if err != nil {
				return fmt.Errorf("dapp: failed to parser inspect: %v", err)
			}
			if err := handleInspect(ctx, client, inspect); err != nil {
				return err
			}
		default:
			return fmt.Errorf("dapp: invalid request type: %v", finishBody.RequestType)
		}
	}
}

func (s *dappService) Ready() <-chan struct{} {
	return s.ready
}

func handleAdvance(
	ctx context.Context,
	client *rollup.ClientWithResponses,
	advance rollup.Advance,
) error {
	log.Printf("dapp: handling advance with payload %v", advance.Payload)

	// add voucher
	voucherReq := rollup.Voucher{
		Destination: advance.Metadata.MsgSender,
		Payload:     advance.Payload,
	}
	voucherResp, err := client.AddVoucher(ctx, voucherReq)
	if err != nil {
		return fmt.Errorf("dapp: %v", err)
	}
	if voucherResp.StatusCode != http.StatusOK {
		return fmt.Errorf("dapp: failed to add report")
	}

	// add notice
	noticeReq := rollup.Notice{
		Payload: advance.Payload,
	}
	noticeResp, err := client.AddNotice(ctx, noticeReq)
	if err != nil {
		return fmt.Errorf("dapp: %v", err)
	}
	if noticeResp.StatusCode != http.StatusOK {
		return fmt.Errorf("dapp: failed to add notice")
	}

	// add report
	reportReq := rollup.Report{
		Payload: advance.Payload,
	}
	reportResp, err := client.AddReport(ctx, reportReq)
	if err != nil {
		return fmt.Errorf("dapp: %v", err)
	}
	if reportResp.StatusCode != http.StatusOK {
		return fmt.Errorf("dapp: failed to add report")
	}

	return nil
}

func handleInspect(
	ctx context.Context,
	client *rollup.ClientWithResponses,
	inspect rollup.Inspect,
) error {
	log.Printf("dapp: handling inspect with payload %v", inspect.Payload)

	// add report
	reportReq := rollup.Report{
		Payload: inspect.Payload,
	}
	reportResp, err := client.AddReport(ctx, reportReq)
	if err != nil {
		return fmt.Errorf("dapp: %v", err)
	}
	if reportResp.StatusCode != http.StatusOK {
		return fmt.Errorf("dapp: failed to add report")
	}

	return nil
}
