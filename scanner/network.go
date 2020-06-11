package scanner

import (
	"log"
	"net"
	"net/url"
	"sync"
	"time"
)

// Network -
type Network struct {
	Nodes []*Node

	checked map[string]struct{}
	chainID string
}

// NewNetwork -
func NewNetwork(chainID string) *Network {
	return &Network{
		Nodes:   make([]*Node, 0),
		checked: make(map[string]struct{}),
		chainID: chainID,
	}
}

// Init -
func (network *Network) Init(bootstrap []string) {
	for i := range bootstrap {
		URL, err := url.Parse(bootstrap[i])
		if err != nil {
			log.Printf("[ERROR] Invalid bootstrap URL: %s", bootstrap[i])
			continue
		}

		var ip string
		addr, err := net.LookupIP(URL.Host)
		if err != nil {
			log.Printf("[WARNING] Can't resolve %s", bootstrap[i])
			continue
		}

		ip = addr[0].String()

		network.Nodes = append(network.Nodes,
			NewNode(
				ip,
				WithRPCURL(bootstrap[i]),
			),
		)
	}
}

// Scan -
func (network *Network) Scan() error {
	start := time.Now()
	for i := range network.Nodes {
		if err := network.findNeighbors(network.Nodes[i]); err != nil {
			log.Printf("[ERROR] %s", err)
		}
	}

	log.Printf("Spent: %s", time.Since(start))
	return nil
}

func (network *Network) findNeighbors(node *Node) error {
	if _, ok := network.checked[node.IP]; ok {
		return nil
	}
	if err := node.findNeighbors(); err != nil {
		network.checked[node.IP] = struct{}{}
		return err
	}

	network.pingNodes(node.Neighbors)

	network.checked[node.IP] = struct{}{}
	network.Nodes = append(network.Nodes, node)

	for i := range node.Neighbors {
		if err := network.findNeighbors(node.Neighbors[i]); err != nil {
			return err
		}
	}
	return nil
}

func (network *Network) pingNodes(nodes []*Node) {
	var wg sync.WaitGroup

	for _, node := range nodes {
		wg.Add(1)
		go network.pingNode(node, &wg)
	}
	wg.Wait()
}

func (network *Network) pingNode(node *Node, wg *sync.WaitGroup) {
	defer wg.Done()
	if err := node.checkHead(network.chainID); err != nil {
		// log.Printf("[WARNING] check head: (%s) %s", node.ip, err)
		return
	}

	if err := node.pingListener(); err != nil {
		// log.Printf("[WARNING] ping listener: (%s) %s", node.ip, err)
		return
	}
}
