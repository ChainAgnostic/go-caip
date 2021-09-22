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

type ChainID struct {
	Namespace string `json:"namespace"`
	Reference string `json:"reference"`
}

var (
	chainNamespaceRegex = regexp.MustCompile("[-a-z0-9]{3,8}")
	chainReferenceRegex = regexp.MustCompile("[-a-zA-Z0-9]{1,32}")
)

func (c *ChainID) validate() error {
	if ok := chainNamespaceRegex.Match([]byte(c.Namespace)); !ok {
		return errors.New("chain namespace does not match spec")
	}

	if ok := chainReferenceRegex.Match([]byte(c.Reference)); !ok {
		return errors.New("chain reference does not match spec")
	}

	return nil
}

func (c *ChainID) String() string {
	if err := c.validate(); err != nil {
		panic(err)
	}
	return c.Namespace + ":" + c.Reference
}

func (c *ChainID) Parse(s string) (*ChainID, error) {
	split := strings.SplitN(s, ":", 2)
	if len(split) != 2 {
		return nil, fmt.Errorf("invalid chain id: %s", s)
	}

	c = &ChainID{split[0], split[1]}
	if err := c.validate(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *ChainID) UnmarshalJSON(data []byte) error {
	type ChainIDAlias ChainID
	ca := (*ChainIDAlias)(c)
	if err := json.Unmarshal(data, &ca); err != nil {
		return err
	}

	if err := c.validate(); err != nil {
		return err
	}

	return nil
}

func (c *ChainID) MarshalJSON() ([]byte, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}

	type ChainIDAlias ChainID
	ca := (*ChainIDAlias)(c)
	return json.Marshal(ca)
}

func (c *ChainID) Format(namespace, reference string) (*ChainID, error) {
	cID := &ChainID{namespace, reference}
	if err := cID.validate(); err != nil {
		return nil, err
	}

	return cID, nil
}

func (c *ChainID) Value() (driver.Value, error) {
	return c.String(), nil
}

func (c *ChainID) Scan(src interface{}) error {
	var i sql.NullString
	if err := i.Scan(src); err != nil {
		return fmt.Errorf("scanning chain id: %w", err)
	}

	if !i.Valid {
		return nil
	}

	if _, err := c.Parse(i.String); err != nil {
		return err
	}

	return nil
}
