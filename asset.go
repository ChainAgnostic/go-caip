package caip

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type AssetID struct {
	ChainID   *ChainID `json:"chain_id"`
	Namespace string   `json:"asset_namespace"`
	Reference string   `json:"asset_reference"`
}

var (
	assetNamespaceRegex = regexp.MustCompile("[-a-z0-9]{3,8}")
	assetReferenceRegex = regexp.MustCompile("[-a-zA-Z0-9]{1,64}")
)

func (a *AssetID) validate() error {
	if ok := assetNamespaceRegex.Match([]byte(a.Namespace)); !ok {
		return errors.New("asset namespace does not match spec")
	}

	if ok := assetReferenceRegex.Match([]byte(a.Reference)); !ok {
		return errors.New("asset reference does not match spec")
	}

	return nil
}

func (a *AssetID) String() string {
	if err := a.validate(); err != nil {
		panic(err)
	}
	return a.ChainID.String() + "/" + a.Namespace + ":" + a.Reference
}

func (a *AssetID) Parse(s string) (*AssetID, error) {
	components := strings.SplitN(s, "/", 2)
	if len(components) != 2 {
		return nil, fmt.Errorf("invalid asset id: %s", s)
	}

	cID, err := new(ChainID).Parse(components[0])
	if err != nil {
		return nil, err
	}

	asset := strings.SplitN(components[1], ":", 2)
	if len(asset) != 2 {
		return nil, fmt.Errorf("invalid asset id: %s", s)
	}

	a = &AssetID{cID, asset[0], asset[1]}
	if err := a.validate(); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *AssetID) UnmarshalJSON(data []byte) error {
	type AssetIDAlias AssetID
	aa := (*AssetIDAlias)(a)
	if err := json.Unmarshal(data, &aa); err != nil {
		return err
	}

	if err := a.validate(); err != nil {
		return err
	}

	return nil
}

func (a *AssetID) MarshalJSON() ([]byte, error) {
	if err := a.validate(); err != nil {
		return nil, err
	}

	type AssetIDAlias AssetID
	ca := (*AssetIDAlias)(a)
	return json.Marshal(ca)
}

func (a *AssetID) Format(chainID *ChainID, namespace, reference string) (*AssetID, error) {
	aID := &AssetID{chainID, namespace, reference}
	if err := aID.validate(); err != nil {
		return nil, err
	}

	return aID, nil
}

func (a *AssetID) Value() (driver.Value, error) {
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

	if _, err := a.Parse(i.String); err != nil {
		return err
	}

	return nil
}
