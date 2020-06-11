package main

import (
	"log"

	"github.com/aopoltorzhicky/tezos-scanner/scanner"
)

func main() {
	bootstrap := []string{
		"https://api.tez.ie",
		"https://rpc.tzkt.io/mainnet",
		"https://mainnet-tezos.giganode.io",
		"https://mainnet.smartpy.io",
	}

	network := scanner.NewNetwork("NetXdQprcVkpaWU")
	network.Init(bootstrap)

	if err := network.Scan(); err != nil {
		panic(err)
	}

	log.Printf("Found %d nodes", len(network.Nodes))
}
