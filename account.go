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

type AccountID struct {
	ChainID ChainID `json:"chain_id"`
	Address string  `json:"account_address"`
}

var (
	accountRegex = regexp.MustCompile("[a-zA-Z0-9]{1,64}")
)

func NewAccountID(chainID ChainID, address string) (AccountID, error) {
	cID := AccountID{chainID, address}
	if err := cID.validate(); err != nil {
		return AccountID{}, err
	}

	return AccountID{chainID, address}, nil
}

func (c AccountID) validate() error {
	if err := c.ChainID.validate(); err != nil {
		return err
	}

	if ok := accountRegex.Match([]byte(c.Address)); !ok {
		return errors.New("namespace does not match spec")
	}

	return nil
}

func (c AccountID) String() string {
	if err := c.validate(); err != nil {
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
	if err := c.validate(); err != nil {
		return err
	}

	return nil
}

func (c *AccountID) UnmarshalJSON(data []byte) error {
	type AccountIDAlias AccountID
	ca := (*AccountIDAlias)(c)
	if err := json.Unmarshal(data, &ca); err != nil {
		return err
	}

	if err := c.validate(); err != nil {
		return err
	}

	return nil
}

func (c AccountID) MarshalJSON() ([]byte, error) {
	if err := c.validate(); err != nil {
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
