package scanner

import (
	"fmt"
	"net"
	"time"

	"github.com/romanserikov/tzgo"
)

// Node -
type Node struct {
	IP        string
	Neighbors []*Node

	RPCURI      string
	ListenerURI string
}

// NewNode -
func NewNode(ip string, opts ...NodeOption) *Node {
	node := &Node{
		IP:        ip,
		Neighbors: make([]*Node, 0),
	}

	for _, opt := range opts {
		opt(node)
	}

	return node
}

// String -
func (n *Node) String() string {
	alive := "alive"
	if !n.IsAlive() {
		alive = "not alive"
	}
	return fmt.Sprintf("[%s] %s | neighbors: %d", n.RPCURI, alive, len(n.Neighbors))
}

// HasNeighbors -
func (n *Node) HasNeighbors() bool {
	return len(n.Neighbors) > 0
}

// IsAlive -
func (n *Node) IsAlive() bool {
	return n.RPCURI != ""
}

func (n *Node) findNeighbors() error {
	if !n.IsAlive() {
		return nil
	}

	node := tzgo.NewTezosNode(n.RPCURI, 5*time.Second)
	points, err := node.NetworkPoints()
	if err != nil {
		return fmt.Errorf("%s: %s", n, err)
	}

	n.Neighbors = make([]*Node, 0)
	for _, point := range points {
		if point.State.EventKind != "running" {
			continue
		}

		host, _, err := net.SplitHostPort(point.URI)
		if err != nil {
			return fmt.Errorf("%s: %s", n, err)
		}
		n.Neighbors = append(n.Neighbors, NewNode(
			host,
			WithListenerURL(point.URI),
		))
	}

	return nil
}

func (n *Node) checkHead(chainID string) error {
	ports := []int{8732, 80}
	for _, port := range ports {
		baseURL := fmt.Sprintf("http://%s:%d", n.IP, port)
		node := tzgo.NewTezosNode(baseURL, 2*time.Second)

		block, err := node.Head()
		if err != nil {
			// log.Printf("[WARNING] %s: %s", baseURL, err)
			continue
		}

		if block.ChainID != chainID {
			return fmt.Errorf("Invalid chain ID: %s != %s", block.ChainID, chainID)
		}

		n.RPCURI = baseURL
		break
	}
	return nil
}

func (n *Node) pingListener() error {
	conn, err := net.DialTimeout("tcp", n.ListenerURI, time.Second)
	if err != nil {
		n.ListenerURI = ""
		return err
	}

	if conn == nil {
		n.ListenerURI = ""
		return fmt.Errorf("Listener connection is nil: %s", n.ListenerURI)
	}
	conn.Close()
	return nil
}
