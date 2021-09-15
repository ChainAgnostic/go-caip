package caip

import (
	"encoding/json"
	"testing"
)

// See: https://github.com/ChainAgnostic/CAIPs/blob/master/CAIPs/caip-10.md#test-cases
func TestAccountID(t *testing.T) {
	for _, tc := range []struct {
		id string
	}{{
		// Ethereum mainnet
		id: "eip155:1:0xab16a96d359ec26a11e2c2b3d8f8b8942d5bfcdb",
	}, {
		// Bitcoin mainnet
		id: "bip122:000000000019d6689c085ae165831e93:128Lkh3S7CkDTBZ8W7BbpsN3YYizJMp8p6",
	}, {
		// Cosmos Hub
		id: "cosmos:cosmoshub-3:cosmos1t2uflqwqe0fsj0shcfkrvpukewcw40yjj6hdc0",
	}, {
		// Kusama network
		id: "polkadot:b0a8d493285c2df73290dfb7e61f870f:5hmuyxw9xdgbpptgypokw4thfyoe3ryenebr381z9iaegmfy",
	}, {
		// Dummy max length (64+1+8+1+32 = 106 chars/bytes)
		id: "chainstd:8c3444cf8970a9e41a706fab93e7a6c4:6d9b0b4b9994e8a6afbd3dc3ed983cd51c755afb27cd1dc7825ef59c134a39f7",
	}} {
		a, err := new(AccountID).Parse(tc.id)
		if err != nil {
			t.Errorf("Failed to parse account id")
		}

		if a.String() != tc.id {
			t.Errorf("Failed to serialize account id to string")
		}

		if _, err := new(AccountID).Format(a.ChainID, a.Address); err != nil {
			t.Errorf("Failed to create account id from namespace and reference")
		}

		b, err := json.Marshal(a)
		if err != nil {
			t.Errorf("Failed to marshal to json")
		}

		a = new(AccountID)
		if err := json.Unmarshal(b, a); err != nil {
			t.Errorf("Failed to unmarshal to json")
		}

		if a.String() != tc.id {
			t.Errorf("Unmarshalled client id invalid")
		}
	}
}
