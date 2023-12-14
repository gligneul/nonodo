// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package main

import (
	"os"

	"github.com/gligneul/nonodo/cmd/nonodo/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	file, err := os.Create("README.md")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	err = doc.GenMarkdown(cmd.Cmd, file)
	if err != nil {
		panic(err)
	}
}
