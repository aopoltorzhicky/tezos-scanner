package protocol

import (
	"fmt"
	"log"
	"net"
)

// GetPeersAddresses -
func (peer *Peer) GetPeersAddresses() ([]*Peer, error) {
	if err := peer.SendMessage(BootstrapMsg{}); err != nil {
		return nil, err
	}
	msg, _, err := peer.ReceivePeerMessage()
	if err != nil {
		return nil, err
	}
	if err := peer.parseMessage(msg); err != nil {
		return nil, err
	}

	if err := peer.SendMessage(BootstrapMsg{}); err != nil {
		return nil, err
	}
	msg, _, err = peer.ReceivePeerMessage()
	if err != nil {
		return nil, err
	}
	if err := peer.parseMessage(msg); err != nil {
		return nil, err
	}

	return getNeighborsTCPAddress(peer.Neighbors), nil
}

func (peer *Peer) parseMessage(msg interface{}) error {
	switch message := msg.(type) {
	case AdvertiseMsg:
		peer.Neighbors = message.Addresses
	case CurrentHeadMsg:
	case GetCurrentBranchMsg:
	default:
		return fmt.Errorf("Unknown message type: %T", msg)
	}
	return nil
}

func getNeighborsTCPAddress(peersAddress []string) []*Peer {
	neighbors := make([]*Peer, len(peersAddress))
	for i, peer := range peersAddress {
		neighbors[i] = &Peer{}
		addr, err := net.ResolveTCPAddr("tcp", peer)
		if err != nil {
			log.Printf("Resolve tcp address error: %s", err)
		}
		neighbors[i].Address = *addr
	}
	return neighbors
}
