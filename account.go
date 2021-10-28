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

type EVMAddressable interface {
	Address() common.Address
}

type AccountID struct {
	ChainID ChainID `json:"chain_id"`
	Address string  `json:"account_address"`
}

var (
	accountRegex = regexp.MustCompile("[a-zA-Z0-9]{1,64}")
)

func NewAccountID(chainID ChainID, address string) (AccountID, error) {
	aID := AccountID{chainID, address}
	if err := aID.Validate(); err != nil {
		return AccountID{}, err
	}

	return aID, nil
}

func UnsafeAccountID(chainID ChainID, address string) AccountID {
	return AccountID{chainID, address}
}

func (c AccountID) Validate() error {
	if err := c.ChainID.Validate(); err != nil {
		return err
	}

	if ok := accountRegex.Match([]byte(c.Address)); !ok {
		return errors.New("namespace does not match spec")
	}

	return nil
}

func (c AccountID) String() string {
	if err := c.Validate(); err != nil {
		panic(err)
	}
	return c.ChainID.String() + ":" + c.Address
}

func (c *AccountID) Parse(s string) error {
	split := strings.SplitN(s, ":", 3)
	if len(split) != 3 {
		return fmt.Errorf("invalid account id: %s", s)
	}

	*c = AccountID{ChainID{split[0], split[1]}, split[2]}
	if err := c.Validate(); err != nil {
		return err
	}

	return nil
}

func (c *AccountID) ParseX(s string) {
	if err := c.Parse(s); err != nil {
		panic(err)
	}
}

func (c *AccountID) UnmarshalJSON(data []byte) error {
	type AccountIDAlias AccountID
	ca := (*AccountIDAlias)(c)
	if err := json.Unmarshal(data, &ca); err != nil {
		return err
	}

	if err := c.Validate(); err != nil {
		return err
	}

	return nil
}

func (c AccountID) MarshalJSON() ([]byte, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	type AccountIDAlias AccountID
	ca := (AccountIDAlias)(c)
	return json.Marshal(ca)
}

func (c AccountID) Value() (driver.Value, error) {
	return c.String(), nil
}

func (c *AccountID) Scan(src interface{}) error {
	var i sql.NullString
	if err := i.Scan(src); err != nil {
		return fmt.Errorf("scanning account id: %w", err)
	}

	if !i.Valid {
		return nil
	}

	if err := c.Parse(i.String); err != nil {
		return err
	}

	return nil
}

type EVMAccountID struct {
	EVMAddressable
	AccountID
}

func NewEVMAccountID(chainID ChainID, address string) (EVMAccountID, error) {
	aID := AccountID{chainID, address}
	if err := aID.Validate(); err != nil {
		return EVMAccountID{}, err
	}

	aID.Address = common.HexToAddress(aID.Address).Hex()
	return EVMAccountID{AccountID: aID}, nil
}

func UnsafeEVMAccountID(chainID ChainID, address string) EVMAccountID {
	aID := AccountID{chainID, common.HexToAddress(address).Hex()}
	return EVMAccountID{AccountID: aID}
}

func (a EVMAccountID) Validate() error {
	if ok := common.IsHexAddress(a.AccountID.Address); !ok {
		return fmt.Errorf("invalid eth address: %s", a.AccountID.Address)
	}

	return a.AccountID.Validate()
}

func (a EVMAccountID) Address() common.Address {
	return common.HexToAddress(a.AccountID.Address)
}
