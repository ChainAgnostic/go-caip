# go-caip

[CAIP standard utils](https://github.com/ChainAgnostic/CAIPs)

## Installation

```
go get github.com/ChainAgnostic/go-caip
```

## ChainID (CAIP-2)

```go
// From namespace + reference
c, err := NewChainID("eip155", "1")
if err != nil {
    panic(err)
}
c.String() // "eip155:1"

// Parse CAIP-2 ChainID
c := ChainID{}
if err := c.Parse("eip155:1"); err != nil {
    panic(err)
}

b, err := json.Marshal(c) // {"namespace": "eip155", "reference": "1"}

c := ChainID{}
err := json.Unmarshal(b, &c)
```

## AccountID (CAIP-10)

```go
// From chain id + address
a, err := NewAccountID{ChainID{"eip155", "1"}, "0xab16a96d359ec26a11e2c2b3d8f8b8942d5bfcdb"}
if err != nil {
    panic(err)
}
a.String() // "eip155:1:0xab16a96d359ec26a11e2c2b3d8f8b8942d5bfcdb"

// Parse CAIP-10 AccountID
a := AccountID{}
if err := a.Parse("eip155:1:0xab16a96d359ec26a11e2c2b3d8f8b8942d5bfcdb"); err != nil {
    panic(err)
}

b, err := json.Marshal(a) // { "ahain_id": {"namespace": "eip155", "reference": "1"}, "account_address": 0xab16a96d359ec26a11e2c2b3d8f8b8942d5bfcdb" }

a := new(AccountID)
err := json.Unmarshal(b, &a)
```

## AssetID (CAIP-19)

```go
// From namespace + reference
a, err := NewAssetID(ChainID{"eip155", "1"}, "erc20", "0x6b175474e89094c44da98b954eedeac495271d0f")
if err != nil {
    panic(err)
}
a.String() // "eip155:1/erc20:0x6b175474e89094c44da98b954eedeac495271d0f"

// Parse CAIP-19 AssetID
a := AssetID{}
if err := a.Parse("eip155:1/erc20:0x6b175474e89094c44da98b954eedeac495271d0f"); err != nil {
    panic(err)
}

b, err := json.Marshal(a) // { "chain_id": {"namespace": "eip155", "reference": "1"}, "asset_namespace": erc20", "asset_reference": "0x6b175474e89094c44da98b954eedeac495271d0f" }

a := AssetID{}
err := json.Unmarshal(b, a)
```
