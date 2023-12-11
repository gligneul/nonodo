// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package nonodo

import "github.com/ethereum/go-ethereum/common"

var (
	CartesiDAppFactoryAddress  common.Address
	DAppAddressRelayAddress    common.Address
	ERC1155BatchPortalAddress  common.Address
	ERC1155SinglePortalAddress common.Address
	ERC20PortalAddress         common.Address
	ERC721PortalAddress        common.Address
	EtherPortalAddress         common.Address
	InputBoxAddress            common.Address
	CartesiDAppAddress         common.Address
)

func init() {
	CartesiDAppFactoryAddress = common.HexToAddress("0x7122cd1221C20892234186facfE8615e6743Ab02")
	DAppAddressRelayAddress = common.HexToAddress("0xF5DE34d6BbC0446E2a45719E718efEbaaE179daE")
	ERC1155BatchPortalAddress = common.HexToAddress("0xedB53860A6B52bbb7561Ad596416ee9965B055Aa")
	ERC1155SinglePortalAddress = common.HexToAddress("0x7CFB0193Ca87eB6e48056885E026552c3A941FC4")
	ERC20PortalAddress = common.HexToAddress("0x9C21AEb2093C32DDbC53eEF24B873BDCd1aDa1DB")
	ERC721PortalAddress = common.HexToAddress("0x237F8DD094C0e47f4236f12b4Fa01d6Dae89fb87")
	EtherPortalAddress = common.HexToAddress("0xFfdbe43d4c855BF7e0f105c400A50857f53AB044")
	InputBoxAddress = common.HexToAddress("0x59b22D57D4f067708AB0c00552767405926dc768")
	CartesiDAppAddress = common.HexToAddress("0x70ac08179605AF2D9e75782b8DEcDD3c22aA4D0C")
}
