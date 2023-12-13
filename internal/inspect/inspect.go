// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

// This package contains the bindings for the inspect OpenAPI spec.
package inspect

import (
	"log"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gligneul/nonodo/internal/model"
	"github.com/labstack/echo/v4"
)

//go:generate oapi-codegen -config=oapi.yaml ../../api/inspect.yaml

const InspectRetries = 50
const InspectPollInterval = time.Millisecond * 100

// Register the rollup API to echo
func Register(e *echo.Echo, model *model.NonodoModel) {
	inspectAPI := &inspectAPI{model}
	RegisterHandlersWithBaseURL(e, inspectAPI, "inspect")
	log.Print("added /inspect to HTTP server")
}

// Shared struct for request handlers.
type inspectAPI struct {
	model *model.NonodoModel
}

// Handle requests to /.
func (a *inspectAPI) InspectPost(c echo.Context) error {
	contentType := c.Request().Header["Content-Type"]
	if len(contentType) != 1 || contentType[0] != "application/octet-stream" {
		return c.String(http.StatusUnsupportedMediaType, "invalid content type")
	}
	var payload []byte
	err := (&echo.DefaultBinder{}).BindBody(c, &payload)
	if err != nil {
		return nil
	}
	return a.inspect(c, payload)
}

// Handle requests to /{payload}.
func (a *inspectAPI) Inspect(c echo.Context, payload string) error {
	payloadBytes, err := hexutil.Decode(payload)
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid hex payload")
	}
	return a.inspect(c, payloadBytes)
}

// Send the inspect input to the model and wait until it is completed.
func (a *inspectAPI) inspect(c echo.Context, payload []byte) error {
	index := a.model.AddInspectInput(payload)
	for i := 0; i < InspectRetries; i++ {
		input := a.model.GetInspectInput(index)
		if input.Status != model.CompletionStatusUnprocessed {
			resp := convertInput(input)
			return c.JSON(http.StatusOK, &resp)
		}
		ctx := c.Request().Context()
		select {
		case <-ctx.Done():
			return c.String(http.StatusInternalServerError, ctx.Err().Error())
		case <-time.After(InspectPollInterval):
		}
	}
	return c.String(http.StatusRequestTimeout, "inspect timed out")
}

// Convert model input to API type.
func convertInput(input model.InspectInput) InspectResult {
	var status CompletionStatus
	switch input.Status {
	case model.CompletionStatusUnprocessed:
		panic("impossible")
	case model.CompletionStatusAccepted:
		status = Accepted
	case model.CompletionStatusRejected:
		status = Rejected
	case model.CompletionStatusException:
		status = Exception
	default:
		panic("invalid completion status")
	}

	var reports []Report
	for _, report := range input.Reports {
		reports = append(reports, Report{
			Payload: hexutil.Encode(report.Payload),
		})
	}

	return InspectResult{
		Status:              status,
		Reports:             reports,
		ExceptionPayload:    hexutil.Encode(input.Exception),
		ProcessedInputCount: input.ProccessedInputCount,
	}
}
