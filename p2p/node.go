package p2p

import (
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/aopoltorzhicky/tezos-scanner/p2p/ffi"
	"github.com/aopoltorzhicky/tezos-scanner/p2p/protocol"
)

// Node -
type Node struct {
	Peer *protocol.Peer

	connection       net.Conn
	attemptsDuration time.Duration
	nextRetryTime    time.Time
	syncedTime       int64
	attemptsCount    int64
	maxAttemptsCount int64
}

// NewNode -
func NewNode(peer *protocol.Peer, attemptsDuration time.Duration, maxAttemptsCount, syncedTime int64) *Node {
	return &Node{
		Peer:             peer,
		attemptsDuration: attemptsDuration,
		nextRetryTime:    time.Now(),
		maxAttemptsCount: maxAttemptsCount,
		syncedTime:       syncedTime,
	}
}

func (node *Node) close() {
	if node.connection != nil {
		node.connection.Close()
	}
}

func (node *Node) getPeer(identity ffi.Identity) (*protocol.Peer, error) {
	if err := node.connect(); err != nil {
		return nil, fmt.Errorf("connection error: %s", err)
	}

	peer, err := node.handshaking(identity)
	if err != nil {
		return peer, fmt.Errorf("handshaking error: %s", err)
	}

	err = peer.UpdateSyncState(node.syncedTime)
	if err != nil {
		err = fmt.Errorf("UpdateSyncState error: %s", err)
	}
	return peer, err
}

func (node *Node) connect() error {
	if time.Now().Before(node.nextRetryTime) {
		return fmt.Errorf("%s is waiting until %v", node.Peer.Address.IP.String(), node.nextRetryTime)
	}

	tcpType := "tcp"
	if node.Peer.Address.IP.To4() == nil {
		tcpType = "tcp6"
	}
	port := strconv.Itoa(node.Peer.Address.Port)

	connection, err := net.DialTimeout(
		tcpType,
		net.JoinHostPort(node.Peer.Address.IP.String(), port),
		time.Second*8,
	)
	if err != nil {
		node.nextRetryTime = time.Now().Add(node.attemptsDuration)
		node.attemptsCount++
		return err
	}

	node.connection = connection
	return nil
}

func (node *Node) handshaking(identity ffi.Identity) (*protocol.Peer, error) {
	secretKey, err := hex.DecodeString(identity.SecretKey)
	if err != nil {
		return nil, err
	}
	pubKey, err := hex.DecodeString(identity.PublicKey)
	if err != nil {
		return nil, err
	}
	bytePow, err := hex.DecodeString(identity.ProofOfWorkStamp)
	if err != nil {
		return nil, err
	}

	port := uint16(node.Peer.Address.Port)
	connMessage := protocol.NewConnectionMessage(port, node.getVersions(), pubKey, bytePow)

	peer := protocol.NewPeer(node.connection, node.Peer.Address)
	if err := peer.Init(connMessage, secretKey); err != nil {
		node.incrementAttemptsWithPeer(peer)
		return nil, err
	}

	if err = peer.Connect(false, false); err != nil {
		node.incrementAttemptsWithPeer(peer)
		return nil, err
	}
	return peer, nil
}

func (node *Node) getVersions() []protocol.Version {
	if node.Peer.Versions != nil {
		return node.Peer.Versions
	}
	return []protocol.Version{
		{
			Name:  "TEZOS_MAINNET",
			Major: 0,
			Minor: 0,
		},
	}
}

func (node *Node) hasAttempts() bool {
	return node.maxAttemptsCount == 0 || node.maxAttemptsCount > node.attemptsCount
}

func (node *Node) setErrorState(err error) {
	node.Peer.Synced = false
	node.Peer.PrivateNode = true
	node.Peer.Error = err
}

func (node *Node) incrementAttemptsWithPeer(peer *protocol.Peer) {
	node.attemptsCount++
	node.Peer = peer
}
