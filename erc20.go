package caip

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

type ERC20AssetID struct {
	EVMAddressable
	AssetID
}

func NewERC20AssetID(chainID ChainID, namespace, reference string) (ERC20AssetID, error) {
	aID := ERC20AssetID{AssetID: AssetID{chainID, namespace, reference}}
	if err := aID.Validate(); err != nil {
		return ERC20AssetID{}, err
	}

	return aID, nil
}

func UnsafeERC20AssetID(chainID ChainID, namespace, reference string) ERC20AssetID {
	aID := AssetID{chainID, namespace, reference}
	return ERC20AssetID{AssetID: aID}
}

func (a ERC20AssetID) Validate() error {
	if ok := common.IsHexAddress(a.AssetID.Reference); !ok {
		return fmt.Errorf("invalid eth address: %s", a.AssetID.Reference)
	}

	if a.AssetID.Namespace != "erc20" {
		return fmt.Errorf("invalid asset namespace: %s", a.AssetID.Namespace)
	}

	if a.ChainID.Namespace != "eip155" {
		return fmt.Errorf("invalid chain namespace: %s", a.ChainID.Namespace)
	}

	return a.AssetID.Validate()
}

func (a ERC20AssetID) Address() common.Address {
	return common.HexToAddress(a.Reference)
}
