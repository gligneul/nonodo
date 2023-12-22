// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package supervisor

import (
	"bytes"
	"errors"
	"io"
	"log"
)

// Log the command output with the name as a prefix.
type commandLogger struct {
	name   string
	buffer bytes.Buffer
}

func (w *commandLogger) Write(data []byte) (int, error) {
	_, err := w.buffer.Write(data)
	if err != nil {
		return 0, err
	}
	for {
		line, err := w.buffer.ReadBytes('\n')
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return 0, err
		}
		log.Printf("%v: %v", w.name, string(line))
	}
	return len(data), nil
}
