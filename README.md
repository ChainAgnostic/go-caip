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
