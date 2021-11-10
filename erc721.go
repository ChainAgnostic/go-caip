package caip

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

type ERC721AssetID struct {
	EVMAssetID
}

func NewERC721AssetID(chainID ChainID, namespace, reference string) (ERC721AssetID, error) {
	aID := ERC721AssetID{EVMAssetID{AssetID: AssetID{chainID, namespace, reference}}}
	if err := aID.Validate(); err != nil {
		return ERC721AssetID{}, err
	}

	return aID, nil
}

func UnsafeERC721AssetID(chainID ChainID, namespace, reference string) ERC721AssetID {
	aID := AssetID{chainID, namespace, reference}
	return ERC721AssetID{EVMAssetID{AssetID: aID}}
}

func (a ERC721AssetID) Validate() error {
	if err := a.EVMAssetID.Validate(); err != nil {
		return err
	}

	split := strings.Split(a.Reference, "/")
	if ok := common.IsHexAddress(split[0]); !ok {
		return fmt.Errorf("invalid eth address: %s", split[0])
	}

	if a.AssetID.Namespace != "erc721" {
		return fmt.Errorf("invalid asset namespace: %s", a.AssetID.Namespace)
	}

	if len(split) > 1 {
		if _, ok := new(big.Int).SetString(split[1], 10); !ok {
			return fmt.Errorf("invalid token id: %s", split[1])
		}
	}

	return nil
}

func (a ERC721AssetID) Address() common.Address {
	split := strings.Split(a.Reference, "/")
	return common.HexToAddress(split[0])
}
