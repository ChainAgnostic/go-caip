package caip

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

// See: https://github.com/ChainAgnostic/CAIPs/blob/master/CAIPs/caip-19.md#test-cases
func TestAssetID(t *testing.T) {
	for _, tc := range []struct {
		id string
	}{{
		// Ether Token
		id: "eip155:1/slip44:60",
	}, {
		// Bitcoin Token
		id: "bip122:000000000019d6689c085ae165831e93/slip44:0",
	}, {
		// ATOM Token
		id: "cosmos:cosmoshub-3/slip44:118",
	}, {
		// Litecoin Token
		id: "bip122:12a765e31ffd4059bada1e25190f6e98/slip44:2",
	}, {
		// Binance Token
		id: "cosmos:Binance-Chain-Tigris/slip44:714",
	}, {
		// IOV Token
		id: "cosmos:iov-mainnet/slip44:234",
	}, {
		// Lisk Token
		id: "lip9:9ee11e9df416b18b/slip44:134",
	}, {
		// DAI Token
		id: "eip155:1/erc20:0x6b175474e89094c44da98b954eedeac495271d0f",
	}, {
		// CryptoKitties Collectible
		id: "eip155:1/erc721:0x06012c8cf97BEaD5deAe237070F9587f8E7A266d",
	}, {
		// CryptoKitties Collectible ID
		id: "eip155:1/erc721:0x06012c8cf97BEaD5deAe237070F9587f8E7A266d/771769",
	}} {
		a := AssetID{}
		if err := a.Parse(tc.id); err != nil {
			t.Fatalf("Failed to parse asset id: %v", err)
		}

		if a.String() != tc.id {
			t.Fatalf("Failed to serialize asset id to string")
		}

		if _, err := NewAssetID(a.ChainID, a.Namespace, a.Reference); err != nil {
			t.Fatalf("Failed to create asset id from namespace and reference")
		}

		b, err := json.Marshal(a)
		if err != nil {
			t.Fatalf("Failed to marshal to json")
		}

		a = AssetID{}
		if err := json.Unmarshal(b, &a); err != nil {
			t.Fatalf("Failed to unmarshal to json")
		}

		if a.String() != tc.id {
			t.Fatalf("Unmarshalled asset id invalid")
		}

		a2 := AssetID{}
		if err := a2.Scan(a.String()); err != nil {
			t.Errorf("Scanning value from sql.NullString")
		}

		if a2.String() != a.String() {
			t.Errorf("Scanned value not valid")
		}
	}
}

func TestEVMAssetID(t *testing.T) {
	for _, tc := range []struct {
		id string
	}{{
		id: "eip155:1/erc20:0x6b175474e89094c44da98b954eedeac495271d0f",
	}, {
		id: "eip155:1/erc720:0x6b175474e89094c44da98b954eedeac495271d0f",
	}, {
		id: "eip155:1/erc720:0x6b175474e89094c44da98b954eedeac495271d0f/1",
	}, {
		id: "eip155:1/erc1155:0x6b175474e89094c44da98b954eedeac495271d0f",
	}, {
		id: "eip155:1/erc1155:0x6b175474e89094c44da98b954eedeac495271d0f/1",
	}} {
		a := EVMAssetID{}
		if err := a.Parse(tc.id); err != nil {
			t.Errorf("Failed to parse asset id")
		}

		if a.String() != tc.id {
			t.Errorf("Failed to serialize asset id to string")
		}

		if _, err := NewEVMAssetID(a.ChainID, a.AssetID.Namespace, a.AssetID.Reference); err != nil {
			t.Errorf("Failed to create asset id from address")
		}

		b, err := json.Marshal(a)
		if err != nil {
			t.Errorf("Failed to marshal to json")
		}

		a = EVMAssetID{}
		if err := json.Unmarshal(b, &a); err != nil {
			t.Errorf("Failed to unmarshal to json")
		}

		if a.String() != tc.id {
			t.Errorf("Unmarshalled asset id invalid")
		}

		a2 := EVMAssetID{}
		if err := a2.Scan(a.String()); err != nil {
			t.Errorf("Scanning value from sql.NullString")
		}

		if a2.String() != a.String() {
			t.Errorf("Scanned value not valid")
		}
	}
}

func TestInvalidEVMAssetID(t *testing.T) {
	for _, tc := range []struct {
		id  string
		err error
	}{{
		id:  "eip155:1/erc20:0x6b175474e89094c44da98b954eedeac495271d0x",
		err: fmt.Errorf("invalid eth address: %s", "0x6b175474e89094c44da98b954eedeac495271d0x"),
	}, {
		id:  "eip155:1/erc20:0x6b175474e",
		err: fmt.Errorf("invalid eth address: %s", "0x6b175474e"),
	}, {
		id:  "cosmos:1/erc20:0xab16a96d359ec26a11e2c2b3d8f8b8942d5bfcdd",
		err: fmt.Errorf("invalid chain namespace: %s", "cosmos"),
	}} {
		a := EVMAssetID{}
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

		_, err = NewEVMAssetID(a.ChainID, a.AssetID.Namespace, a.AssetID.Reference)
		if err == nil {
			t.Errorf("Create asset id should error")
		}

		if err.Error() != tc.err.Error() {
			t.Errorf("expected error: %s, got: %s", tc.err, err)
		}
	}
}
