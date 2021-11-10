package caip

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

type ERC721AssetID struct {
	EVMAddressable
	AssetID
}

func NewERC721AssetID(chainID ChainID, namespace, reference string) (ERC721AssetID, error) {
	aID := ERC721AssetID{AssetID: AssetID{chainID, namespace, reference}}
	if err := aID.Validate(); err != nil {
		return ERC721AssetID{}, err
	}

	return aID, nil
}

func UnsafeERC721AssetID(chainID ChainID, namespace, reference string) ERC721AssetID {
	aID := AssetID{chainID, namespace, reference}
	return ERC721AssetID{AssetID: aID}
}

func (a ERC721AssetID) Validate() error {
	split := strings.Split(a.Reference, "/")
	if ok := common.IsHexAddress(split[0]); !ok {
		return fmt.Errorf("invalid eth address: %s", split[0])
	}

	if a.ChainID.Namespace != "eip155" {
		return fmt.Errorf("invalid chain namespace: %s", a.ChainID.Namespace)
	}

	if a.AssetID.Namespace != "erc721" {
		return fmt.Errorf("invalid asset namespace: %s", a.AssetID.Namespace)
	}

	if len(split) > 1 {
		if _, ok := new(big.Int).SetString(split[1], 10); !ok {
			return fmt.Errorf("invalid token id: %s", split[1])
		}
	}

	return a.AssetID.Validate()
}

func (a ERC721AssetID) Address() common.Address {
	split := strings.Split(a.Reference, "/")
	return common.HexToAddress(split[0])
}
