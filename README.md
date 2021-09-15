# go-caip

[CAIP standard utils](https://github.com/ChainAgnostic/CAIPs)

## ChainID (CAIP-2)

```go
// From namespace + reference
c, err := new(ChainID).Format("eip155", "1")
if err != nil {
    panic(err)
}
c.String() // "eip155:1"

// Parse CAIP-2 ChainID
c, err := new(ChainID).Parse("eip155:1")
if err != nil {
    panic(err)
}

b, err := json.Marshal(c) // {"namespace": "eip155", "reference": "1"}

c := new(ChainID)
err := json.Unmarshal(b, c)
```

## AccountID (CAIP-10)

```go
// From namespace + reference
c, err := new(AccountID).Format(new(ChainID).Format("eip155", "1"), "0xab16a96d359ec26a11e2c2b3d8f8b8942d5bfcdb")
if err != nil {
    panic(err)
}
c.String() // "eip155:1"

// Parse CAIP-2 AccountID
c, err := new(AccountID).Parse("eip155:1:0xab16a96d359ec26a11e2c2b3d8f8b8942d5bfcdb")
if err != nil {
    panic(err)
}

b, err := json.Marshal(c) // { "chain_id": {"namespace": "eip155", "reference": "1"}, "account_address": 0xab16a96d359ec26a11e2c2b3d8f8b8942d5bfcdb" }

c := new(AccountID)
err := json.Unmarshal(b, c)
```
