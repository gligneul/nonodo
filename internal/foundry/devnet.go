// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package foundry

import _ "embed"

// Generate the devnet state and embed it in the Go binary.
//
//go:generate go run ./gen-devnet-state
//go:embed anvil_state.json
var devnetState []byte

// Input box address in devnet.
const InputBoxAddress = "0x59b22D57D4f067708AB0c00552767405926dc768"

// DApp address in devnet.
const DAppAddress = "0x70ac08179605AF2D9e75782b8DEcDD3c22aA4D0C"

// Foundry test mnemonic.
const TestMnemonic = "test test test test test test test test test test test junk"
