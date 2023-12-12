// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package nonodo

// Default port for the Ethereum node.
const EthDefaultPort = 8545

// Options to nonodo.
type NonodoOpts struct {
	AnvilPort                  int
	AnvilBlockTime             int
	AnvilVerbose               bool
	HttpPort                   int
	BuiltInDApp                bool
	CartesiDAppFactoryAddress  string
	DAppAddressRelayAddress    string
	ERC1155BatchPortalAddress  string
	ERC1155SinglePortalAddress string
	ERC20PortalAddress         string
	ERC721PortalAddress        string
	EtherPortalAddress         string
	InputBoxAddress            string
	CartesiDAppAddress         string
}

// Create the options struct with default values.
func NewNonodoOpts() NonodoOpts {
	return NonodoOpts{
		AnvilPort:                  EthDefaultPort,
		AnvilBlockTime:             1,
		AnvilVerbose:               false,
		HttpPort:                   8080,
		BuiltInDApp:                false,
		CartesiDAppFactoryAddress:  "0x7122cd1221C20892234186facfE8615e6743Ab02",
		DAppAddressRelayAddress:    "0xF5DE34d6BbC0446E2a45719E718efEbaaE179daE",
		ERC1155BatchPortalAddress:  "0xedB53860A6B52bbb7561Ad596416ee9965B055Aa",
		ERC1155SinglePortalAddress: "0x7CFB0193Ca87eB6e48056885E026552c3A941FC4",
		ERC20PortalAddress:         "0x9C21AEb2093C32DDbC53eEF24B873BDCd1aDa1DB",
		ERC721PortalAddress:        "0x237F8DD094C0e47f4236f12b4Fa01d6Dae89fb87",
		EtherPortalAddress:         "0xFfdbe43d4c855BF7e0f105c400A50857f53AB044",
		InputBoxAddress:            "0x59b22D57D4f067708AB0c00552767405926dc768",
		CartesiDAppAddress:         "0x70ac08179605AF2D9e75782b8DEcDD3c22aA4D0C",
	}
}
