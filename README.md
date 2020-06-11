# Tezos scanner

Scan tezos network via RPC

## Usage

```go
network := scanner.NewNetwork("NetXdQprcVkpaWU")  // Set chain ID for checking

bootstrap := []string{
  // Known nodes URLS (https://some-api-address.com)
}
network.Init(bootstrap)

if err := network.Scan(); err != nil {
  panic(err)
}
```

After calling `Scan` method network nodes `network.Nodes` will fill. `Node` structure described below.

```go

// Node -
type Node struct {
  IP        string
  Neighbors []*Node
  RPCURI      string
  ListenerURI string
}
```
