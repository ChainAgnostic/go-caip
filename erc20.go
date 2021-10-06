package caip

import "github.com/ethereum/go-ethereum/common"

type ERC20AssetID struct {
	EVMAddressable
	AssetID
}

func NewERC20AssetID(chainID ChainID, namespace, reference string) (ERC20AssetID, error) {
	aID := AssetID{chainID, namespace, reference}
	if err := aID.validate(); err != nil {
		return ERC20AssetID{}, err
	}

	return ERC20AssetID{AssetID: aID}, nil
}

func (a ERC20AssetID) Address() common.Address {
	return common.HexToAddress(a.Reference)
}
