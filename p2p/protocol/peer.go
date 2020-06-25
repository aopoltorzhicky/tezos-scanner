package protocol

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/aopoltorzhicky/tezos-scanner/p2p/protocol/crypto"
)

// Version -
type Version struct {
	Name  string
	Major uint16
	Minor uint16
}

type nonceType = []byte

//
// Peer is a struct for communicating with tezos nodes.
//
// There are several basic methods: SendMessage and ReceivePeerMessage.
// For initialize a Peer object call three methods: New, Init, Connect.
type Peer struct {
	conn           net.Conn
	localNonce     nonceType
	remoteNonce    nonceType
	precomputedKey [32]byte

	ID             string      `json:"id"`
	Versions       []Version   `json:"versions"`
	DisableMempool bool        `json:"disable_mempool"`
	PrivateNode    bool        `json:"private_node"`
	Address        net.TCPAddr `json:"address"`
	Synced         bool        `json:"synced"`
	RPC            bool        `json:"rpc"`
	Error          error       `json:"error"`
	Neighbors      []string    `json:"neighbors"`
}

// SendMessage -
func (peer *Peer) SendMessage(message peerMessage) (err error) {
	err = sendEncryptedMessage(peer.conn, message.toBytes(), peer.localNonce, &peer.precomputedKey)
	if err != nil {
		return
	}
	peer.localNonce = crypto.NonceIncrement(peer.localNonce)
	return
}

// ReceivePeerMessage - Returns the received message (see peer_messages.go) and type of this message.
//
// Typical use of the method:
/*
msg, msgType, err := peer.ReceivePeerMessage()
if err == nil {
	if msgType == CurrentBranchTag {
		log.Print(msg.(CurrentBranchMsg).Branch)
	}
}
*/
func (peer *Peer) ReceivePeerMessage() (msg interface{}, messageType PeerMessageType, err error) {
	messageType = UnknownTag
	msg = nil

	data, err := peer.receiveData()
	if err != nil {
		return
	}
	msg, messageType = parseMessage(data)
	return
}

// NewPeer -
func NewPeer(conn net.Conn, address net.TCPAddr) (peer *Peer) {
	peer = new(Peer)
	peer.conn = conn
	peer.Address = address
	return
}

// Init -
func (peer *Peer) Init(connMessage ConnectionMessage, secretKey []byte) (err error) {
	if err = sendConnectionMessage(peer.conn, connMessage); err != nil {
		return
	}

	receiveMessage, err := receiveConnectionMessage(peer.conn)
	if err != nil || receiveMessage == nil {
		return err
	}

	peer.localNonce, peer.remoteNonce = crypto.GenerateNonces(connMessage.toBytes(), receiveMessage.toBytes(), false)
	peer.precomputedKey = crypto.PrecomputeSharedKey(receiveMessage.PublicKey, secretKey)
	peer.Versions = receiveMessage.Versions
	peer.ID, err = crypto.CalcPeerID(receiveMessage.PublicKey)
	return err
}

// Connect -
func (peer *Peer) Connect(disableMemoryPool bool, privateNode bool) (err error) {
	metaMsg := &metadata{
		DisableMempool: disableMemoryPool,
		PrivateNode:    privateNode,
	}

	if err = peer.SendMessage(metaMsg); err != nil {
		return
	}

	if err = peer.receiveMeta(); err != nil {
		return
	}

	ack := &ackMessage{
		IsNack: false,
	}

	if err = peer.SendMessage(ack); err != nil {
		return
	}

	if err = peer.receiveAck(); err != nil {
		return err
	}

	return nil
}

// String -
func (peer *Peer) String() string {
	var jsonData []byte
	jsonData, err := json.Marshal(peer)
	if err != nil {
		return ""
	}
	return string(jsonData)
}

func (peer *Peer) receiveData() (data []byte, err error) {
	data, err = receiveEncryptedMessage(peer.conn, peer.remoteNonce, &peer.precomputedKey)
	if err != nil {
		return
	}
	peer.remoteNonce = crypto.NonceIncrement(peer.remoteNonce)
	return
}

func (peer *Peer) receiveMeta() (err error) {
	meta, err := receiveMetaMessage(peer.conn, peer.remoteNonce, &peer.precomputedKey)
	if err != nil {
		return
	}
	peer.remoteNonce = crypto.NonceIncrement(peer.remoteNonce)
	peer.DisableMempool = meta.DisableMempool
	peer.PrivateNode = meta.PrivateNode
	return
}

func (peer *Peer) receiveAck() (err error) {
	ack, err := receiveAckMessage(peer.conn, peer.remoteNonce, &peer.precomputedKey)
	if err != nil {
		return
	}
	peer.remoteNonce = crypto.NonceIncrement(peer.remoteNonce)
	if ack.IsNack {
		err = &NackError{ip: peer.Address.IP}
	}
	return
}

// FindRPC -
func (peer *Peer) FindRPC() error {
	for _, port := range []string{"8732", "18732"} {
		host := net.JoinHostPort(peer.Address.IP.String(), port)
		path := fmt.Sprintf("http://%s/chains/main/blocks/head", host)
		resp, err := http.Get(path)
		if err != nil {
			return err
		}
		if resp.StatusCode == http.StatusOK {
			peer.RPC = true
			break
		}
	}

	return nil
}

// NackError -
type NackError struct {
	ip net.IP
}

// Error -
func (obj *NackError) Error() string {
	return fmt.Sprintf("%s - received nack", obj.ip.String())
}
