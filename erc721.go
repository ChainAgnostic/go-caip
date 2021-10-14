package caip

import (
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

type ERC721AssetID struct {
	EVMAddressable
	AssetID
}

func NewERC721AssetID(chainID ChainID, namespace, reference string) (ERC721AssetID, error) {
	aID := AssetID{chainID, namespace, reference}
	if err := aID.Validate(); err != nil {
		return ERC721AssetID{}, err
	}

	return ERC721AssetID{AssetID: aID}, nil
}

func UnsafeERC721AssetID(chainID ChainID, namespace, reference string) ERC721AssetID {
	aID := AssetID{chainID, namespace, reference}
	return ERC721AssetID{AssetID: aID}
}

func (a ERC721AssetID) Address() common.Address {
	split := strings.Split(a.Reference, "/")
	return common.HexToAddress(split[0])
}
