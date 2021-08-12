# go-caip

[CAIP standard utils](https://github.com/ChainAgnostic/CAIPs)

## ChainID (CAIP-2)

```go
// From namespace + reference
cid := NewChainID("eip155", "1")
cid.String() // "eip155:1"

// Parse CAIP-2 ChainID
cid, err := new(ChainID).Parse("eip155:1")
if err != nil {
    panic(err)
}
```
