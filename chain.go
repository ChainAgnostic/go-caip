package caip

import (
	"database/sql"
	"database/sql/driver"
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
	namespaceRegex = regexp.MustCompile("[-a-z0-9]{3,8}")
	referenceRegex = regexp.MustCompile("[-a-zA-Z0-9]{1,32}")
)

func NewChainID(namespace, reference string) (*ChainID, error) {
	cID := &ChainID{namespace, reference}
	if err := cID.valid(); err != nil {
		return nil, err
	}

	return &ChainID{namespace, reference}, nil
}

func (c *ChainID) valid() error {
	if ok := namespaceRegex.Match([]byte(c.Namespace)); !ok {
		return errors.New("namespace does not match spec")
	}

	if ok := referenceRegex.Match([]byte(c.Reference)); !ok {
		return errors.New("namespace does not match spec")
	}

	return nil
}

func (c *ChainID) String() string {
	return c.Namespace + ":" + c.Reference
}

func (c *ChainID) Parse(s string) (*ChainID, error) {
	split := strings.Split(s, ":")
	if len(split) != 2 {
		return nil, fmt.Errorf("invalid chain id: %s", s)
	}

	c = &ChainID{split[0], split[1]}
	if err := c.valid(); err != nil {
		return nil, err
	}

	return c, nil
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
