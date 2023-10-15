// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: MIT (see LICENSE)

package eggeth

import "github.com/ethereum/go-ethereum/common"

var (
	AddressCartesiDAppFactory  common.Address
	AddressDAppAddressRelay    common.Address
	AddressERC1155BatchPortal  common.Address
	AddressERC1155SinglePortal common.Address
	AddressERC20Portal         common.Address
	AddressERC721Portal        common.Address
	AddressEtherPortal         common.Address
	AddressInputBox            common.Address
	AddressSunodoToken         common.Address
)

func init() {
	AddressCartesiDAppFactory = common.HexToAddress("0x7122cd1221C20892234186facfE8615e6743Ab02")
	AddressDAppAddressRelay = common.HexToAddress("0xF5DE34d6BbC0446E2a45719E718efEbaaE179daE")
	AddressERC1155BatchPortal = common.HexToAddress("0xedB53860A6B52bbb7561Ad596416ee9965B055Aa")
	AddressERC1155SinglePortal = common.HexToAddress("0x7CFB0193Ca87eB6e48056885E026552c3A941FC4")
	AddressERC20Portal = common.HexToAddress("0x9C21AEb2093C32DDbC53eEF24B873BDCd1aDa1DB")
	AddressERC721Portal = common.HexToAddress("0x237F8DD094C0e47f4236f12b4Fa01d6Dae89fb87")
	AddressEtherPortal = common.HexToAddress("0xFfdbe43d4c855BF7e0f105c400A50857f53AB044")
	AddressInputBox = common.HexToAddress("0x59b22D57D4f067708AB0c00552767405926dc768")
	AddressSunodoToken = common.HexToAddress("0xae7f61eCf06C65405560166b259C54031428A9C4")
}
