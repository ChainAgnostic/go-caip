package caip

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

type AssetID struct {
	ChainID   ChainID `json:"chain_id"`
	Namespace string  `json:"asset_namespace"`
	Reference string  `json:"asset_reference"`
}

var (
	assetNamespaceRegex = regexp.MustCompile("[-a-z0-9]{3,8}")
	assetReferenceRegex = regexp.MustCompile("[-a-zA-Z0-9]{1,64}")
)

func NewAssetID(chainID ChainID, namespace, reference string) (AssetID, error) {
	aID := AssetID{chainID, namespace, reference}
	if err := aID.Validate(); err != nil {
		return AssetID{}, err
	}

	return aID, nil
}

func UnsafeAssetID(chainID ChainID, namespace, reference string) AssetID {
	return AssetID{chainID, namespace, reference}
}

func (a AssetID) Validate() error {
	if ok := assetNamespaceRegex.Match([]byte(a.Namespace)); !ok {
		return errors.New("asset namespace does not match spec")
	}

	if ok := assetReferenceRegex.Match([]byte(a.Reference)); !ok {
		return errors.New("asset reference does not match spec")
	}

	return nil
}

func (a AssetID) String() string {
	if err := a.Validate(); err != nil {
		panic(err)
	}
	return a.ChainID.String() + "/" + a.Namespace + ":" + a.Reference
}

func (a *AssetID) Parse(s string) error {
	components := strings.SplitN(s, "/", 2)
	if len(components) != 2 {
		return fmt.Errorf("invalid asset id: %s", s)
	}

	cID := new(ChainID)
	if err := cID.Parse(components[0]); err != nil {
		return err
	}

	asset := strings.SplitN(components[1], ":", 2)
	if len(asset) != 2 {
		return fmt.Errorf("invalid asset id: %s", s)
	}

	*a = AssetID{*cID, asset[0], asset[1]}
	if err := a.Validate(); err != nil {
		return err
	}

	return nil
}

func (a *AssetID) ParseX(s string) {
	if err := a.Parse(s); err != nil {
		panic(err)
	}
}

func (a *AssetID) UnmarshalJSON(data []byte) error {
	type AssetIDAlias AssetID
	aa := (*AssetIDAlias)(a)
	if err := json.Unmarshal(data, &aa); err != nil {
		return err
	}

	if err := a.Validate(); err != nil {
		return err
	}

	return nil
}

func (a AssetID) MarshalJSON() ([]byte, error) {
	if err := a.Validate(); err != nil {
		return nil, err
	}

	type AssetIDAlias AssetID
	ca := (AssetIDAlias)(a)
	return json.Marshal(ca)
}

func (a AssetID) Value() (driver.Value, error) {
	return a.String(), nil
}

func (a *AssetID) Scan(src interface{}) error {
	var i sql.NullString
	if err := i.Scan(src); err != nil {
		return fmt.Errorf("scanning account id: %w", err)
	}

	if !i.Valid {
		return nil
	}

	if err := a.Parse(i.String); err != nil {
		return err
	}

	return nil
}

type EVMAssetID struct {
	EVMAddressable
	AssetID
}

func NewEVMAssetID(chainID ChainID, namespace, reference string) (EVMAssetID, error) {
	aID := EVMAssetID{AssetID: AssetID{chainID, namespace, reference}}
	if err := aID.Validate(); err != nil {
		return EVMAssetID{}, err
	}
	aID.checksum()
	return aID, nil
}

func UnsafeEVMAssetID(chainID ChainID, namespace, reference string) EVMAssetID {
	aID := EVMAssetID{AssetID: AssetID{chainID, namespace, reference}}
	aID.checksum()
	return aID
}

func (a *EVMAssetID) checksum() {
	split := strings.Split(a.Reference, "/")
	// Make reference checksummed
	split[0] = a.Address().Hex()
	a.Reference = strings.Join(split, "/")
}

func (a EVMAssetID) Validate() error {
	split := strings.Split(a.Reference, "/")
	if ok := common.IsHexAddress(split[0]); !ok {
		return fmt.Errorf("invalid eth address: %s", split[0])
	}

	if a.ChainID.Namespace != "eip155" {
		return fmt.Errorf("invalid chain namespace: %s", a.ChainID.Namespace)
	}

	return a.AssetID.Validate()
}

func (a EVMAssetID) Address() common.Address {
	split := strings.Split(a.Reference, "/")
	return common.HexToAddress(split[0])
}

func (a EVMAssetID) AccountID() EVMAccountID {
	return EVMAccountID{AccountID: AccountID{a.ChainID, a.Reference}}
}
