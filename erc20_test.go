package caip

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

func TestERC20AssetID(t *testing.T) {
	for _, tc := range []struct {
		id string
	}{{
		// Ethereum mainnet
		id: "eip155:1/erc20:0x6b175474e89094c44da98b954eedeac495271d0f",
	}} {
		a := ERC20AssetID{}
		if err := a.Parse(tc.id); err != nil {
			t.Errorf("Failed to parse asset id")
		}

		if a.String() != tc.id {
			t.Errorf("Failed to serialize asset id to string")
		}

		if _, err := NewERC20AssetID(a.ChainID, a.AssetID.Namespace, a.AssetID.Reference); err != nil {
			t.Errorf("Failed to create asset id from address")
		}

		b, err := json.Marshal(a)
		if err != nil {
			t.Errorf("Failed to marshal to json")
		}

		a = ERC20AssetID{}
		if err := json.Unmarshal(b, &a); err != nil {
			t.Errorf("Failed to unmarshal to json")
		}

		if a.String() != tc.id {
			t.Errorf("Unmarshalled asset id invalid")
		}

		a2 := ERC20AssetID{}
		if err := a2.Scan(a.String()); err != nil {
			t.Errorf("Scanning value from sql.NullString")
		}

		if a2.String() != a.String() {
			t.Errorf("Scanned value not valid")
		}
	}
}

func TestInvalidERC20AssetID(t *testing.T) {
	for _, tc := range []struct {
		id  string
		err error
	}{{
		id:  "eip155:1/erc20:0x6b175474e89094c44da98b954eedeac495271d0x",
		err: fmt.Errorf("invalid eth address: %s", "0x6b175474e89094c44da98b954eedeac495271d0x"),
	}, {
		id:  "eip155:1/erc721:0x6b175474e89094c44da98b954eedeac495271d0a",
		err: fmt.Errorf("invalid asset namespace: %s", "erc721"),
	}, {
		id:  "eip155:1/erc20:0x6b175474e",
		err: fmt.Errorf("invalid eth address: %s", "0x6b175474e"),
	}, {
		id:  "cosmos:1/erc20:0xab16a96d359ec26a11e2c2b3d8f8b8942d5bfcdd",
		err: fmt.Errorf("invalid chain namespace: %s", "cosmos"),
	}} {
		a := ERC20AssetID{}
		if err := a.Parse(tc.id); err != nil {
			t.Errorf("Failed to parse asset id")
		}

		if a.String() != tc.id {
			t.Errorf("Failed to serialize asset id to string")
		}

		err := a.Validate()
		if err == nil {
			t.Errorf("Validate asset id should error")
		}

		if errors.Is(err, tc.err) {
			t.Errorf("expected error: %s", tc.err)
		}

		_, err = NewERC20AssetID(a.ChainID, a.AssetID.Namespace, a.AssetID.Reference)
		if err == nil {
			t.Errorf("Create asset id should error")
		}

		if err.Error() != tc.err.Error() {
			t.Errorf("expected error: %s, got: %s", tc.err, err)
		}
	}
}
