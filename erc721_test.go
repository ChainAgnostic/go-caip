package caip

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

func TestERC721AssetID(t *testing.T) {
	for _, tc := range []struct {
		id string
	}{{
		// CryptoKitties Collectible
		id: "eip155:1/erc721:0x06012c8cf97BEaD5deAe237070F9587f8E7A266d",
	}, {
		// CryptoKitties Collectible ID
		id: "eip155:1/erc721:0x06012c8cf97BEaD5deAe237070F9587f8E7A266d/771769",
	}} {
		a := ERC721AssetID{}
		if err := a.Parse(tc.id); err != nil {
			t.Errorf("Failed to parse asset id")
		}

		if a.String() != tc.id {
			t.Errorf("Failed to serialize asset id to string")
		}

		if _, err := NewERC721AssetID(a.ChainID, a.AssetID.Namespace, a.AssetID.Reference); err != nil {
			t.Errorf("Failed to create asset id from address")
		}

		b, err := json.Marshal(a)
		if err != nil {
			t.Errorf("Failed to marshal to json")
		}

		a = ERC721AssetID{}
		if err := json.Unmarshal(b, &a); err != nil {
			t.Errorf("Failed to unmarshal to json")
		}

		if a.String() != tc.id {
			t.Errorf("Unmarshalled asset id invalid")
		}

		a2 := ERC721AssetID{}
		if err := a2.Scan(a.String()); err != nil {
			t.Errorf("Scanning value from sql.NullString")
		}

		if a2.String() != a.String() {
			t.Errorf("Scanned value not valid")
		}
	}
}

func TestInvalidERC721AssetID(t *testing.T) {
	for _, tc := range []struct {
		id  string
		err error
	}{{
		id:  "eip155:1/erc721:0x06012c8cf97BEaD5deAe237070F9587f8E7A266x",
		err: fmt.Errorf("invalid eth address: %s", "0x06012c8cf97BEaD5deAe237070F9587f8E7A266x"),
	}, {
		id:  "eip155:1/erc20:0x06012c8cf97BEaD5deAe237070F9587f8E7A266a",
		err: fmt.Errorf("invalid asset namespace: %s", "erc20"),
	}, {
		id:  "eip155:1/erc721:0x06012c8cf97BEaD5deA",
		err: fmt.Errorf("invalid eth address: %s", "0x06012c8cf97BEaD5deA"),
	}, {
		id:  "cosmos:1/erc721:0x06012c8cf97BEaD5deAe237070F9587f8E7A266d",
		err: fmt.Errorf("invalid chain namespace: %s", "cosmos"),
	}, {
		id:  "eip155:1/erc721:0x06012c8cf97BEaD5deAe237070F9587f8E7A266d/cat",
		err: fmt.Errorf("invalid token id: %s", "cat"),
	}} {
		a := ERC721AssetID{}
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

		_, err = NewERC721AssetID(a.ChainID, a.AssetID.Namespace, a.AssetID.Reference)
		if err == nil {
			t.Errorf("Create asset id should error")
		}

		if err.Error() != tc.err.Error() {
			t.Errorf("expected error: %s, got: %s", tc.err, err)
		}
	}
}
