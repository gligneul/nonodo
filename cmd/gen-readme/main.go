// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package main

import (
	"bytes"
	"os"

	"github.com/gligneul/nonodo/cmd/nonodo/cmd"
)

func main() {
	var buf bytes.Buffer

	buf.WriteString("# nonodo\n\n")
	buf.WriteString(cmd.Cmd.Long)
	buf.WriteString("\n\n")

	buf.WriteString("## Flags\n\n```\n")
	flags := cmd.Cmd.NonInheritedFlags()
	flags.SetOutput(&buf)
	flags.PrintDefaults()
	buf.WriteString("```\n")

	const permission = 644
	err := os.WriteFile("README.md", buf.Bytes(), permission)
	if err != nil {
		panic(err)
	}
}
