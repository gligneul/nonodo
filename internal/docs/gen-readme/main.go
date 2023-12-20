// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package main

import (
	"os"

	"github.com/gligneul/nonodo/internal/docs"
)

func main() {
	const permission = 644
	err := os.WriteFile("../../README.md", docs.Readme(), permission)
	if err != nil {
		panic(err)
	}
}
