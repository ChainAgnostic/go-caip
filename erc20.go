package caip

import (
	"fmt"
)

type ERC20AssetID struct {
	EVMAssetID
}

func NewERC20AssetID(chainID ChainID, namespace, reference string) (ERC20AssetID, error) {
	aID := ERC20AssetID{EVMAssetID{AssetID: AssetID{chainID, namespace, reference}}}
	if err := aID.Validate(); err != nil {
		return ERC20AssetID{}, err
	}

	return aID, nil
}

func UnsafeERC20AssetID(chainID ChainID, namespace, reference string) ERC20AssetID {
	aID := AssetID{chainID, namespace, reference}
	return ERC20AssetID{EVMAssetID{AssetID: aID}}
}

func (a ERC20AssetID) Validate() error {
	if err := a.EVMAssetID.Validate(); err != nil {
		return err
	}

	if a.AssetID.Namespace != "erc20" {
		return fmt.Errorf("invalid asset namespace: %s", a.AssetID.Namespace)
	}

	return nil
}
