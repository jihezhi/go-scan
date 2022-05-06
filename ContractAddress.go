package main

import (
	"github.com/ethereum/go-ethereum/common"
)

func Returncontractaddress() []common.Address {
	var AddressSet []common.Address
	AddressSet = []common.Address{
		common.HexToAddress("0x74b23882a30290451A17c44f4F05243b6b58C76d"), //Fantom
		common.HexToAddress("0x841FAD6EAe12c286d1Fd18d1d525DFfA75C7EFFE"), //BOO
	}
	return AddressSet
}
