package protocol

import (
	"encoding/binary"
	"math/rand"
	"time"
)

type peerMessage interface {
	toBytes() []byte
}

type metadata struct {
	DisableMempool bool
	PrivateNode    bool
}

func (message *metadata) fromBytes(data []byte) {
	message.DisableMempool = data[0] != 0
	message.PrivateNode = data[1] != 0
}

func (message *metadata) toBytes() []byte {
	arr := make([]byte, 2)
	if message.DisableMempool {
		arr[0] = 1
	} else {
		arr[0] = 0
	}

	if message.PrivateNode {
		arr[1] = 1
	} else {
		arr[1] = 0
	}
	return arr
}

type ackMessage struct {
	IsNack bool
}

func (message *ackMessage) fromBytes(data []byte) {
	message.IsNack = data[0] != 0
}

func (message *ackMessage) toBytes() (bytes []byte) {
	if message.IsNack {
		return []byte{255}
	}
	return []byte{0}
}

// ConnectionMessage -
type ConnectionMessage struct {
	Port     uint16
	Versions []Version

	PublicKey        []byte
	ProofOfWorkStamp []byte
	MessageNonce     []byte
}

func generateRandomNonce() []byte {
	const nonceSize = 24
	rand.Seed(time.Now().UnixNano())
	nonce := make([]byte, nonceSize)
	rand.Read(nonce)

	return nonce
}

// NewConnectionMessage -
func NewConnectionMessage(port uint16, versions []Version, publicKey, proofOfWorkStamp []byte) (msg ConnectionMessage) {
	return ConnectionMessage{
		Port:             port,
		Versions:         versions,
		PublicKey:        publicKey,
		ProofOfWorkStamp: proofOfWorkStamp,
		MessageNonce:     generateRandomNonce(),
	}
}

const (
	connectionMessageLenSize       = 2
	connectionMessagePortSize      = 2
	connectionMessagePublicKeySize = 32
	connectionMessageNonceSize     = 24
	connectionMessageProofSize     = 24
)

func (msg *ConnectionMessage) toBytes() (bytes []byte) {
	bytes = make([]byte, connectionMessagePortSize+connectionMessagePublicKeySize+connectionMessageNonceSize+connectionMessageProofSize)

	index := 0
	binary.BigEndian.PutUint16(bytes[index:index+connectionMessagePortSize], msg.Port)
	index += connectionMessagePortSize

	copy(bytes[index:index+connectionMessagePublicKeySize], msg.PublicKey)
	index += connectionMessagePublicKeySize

	copy(bytes[index:index+connectionMessageProofSize], msg.ProofOfWorkStamp)
	index += connectionMessageProofSize

	copy(bytes[index:index+connectionMessageNonceSize], msg.MessageNonce)
	index += connectionMessageNonceSize

	versions := versionsToBytes(msg.Versions)
	bytes = append(bytes, versions...)
	bytes = addSize(bytes)
	return
}

func newMessage(bytes []byte) (msg *ConnectionMessage) {
	if len(bytes) < connectionMessageLenSize+connectionMessagePortSize+connectionMessagePublicKeySize+connectionMessageProofSize+connectionMessageNonceSize {
		return nil
	}

	index := 0
	sz := binary.BigEndian.Uint16(bytes[index : index+connectionMessageLenSize])
	if uint16(len(bytes)) != sz+connectionMessageLenSize {
		return nil
	}

	msg = new(ConnectionMessage)
	index += connectionMessageLenSize

	msg.Port = binary.BigEndian.Uint16(bytes[index : index+connectionMessagePortSize])
	index += connectionMessagePortSize

	msg.PublicKey = append(msg.PublicKey, bytes[index:index+connectionMessagePublicKeySize]...)
	index += connectionMessagePublicKeySize

	msg.ProofOfWorkStamp = append(msg.ProofOfWorkStamp, bytes[index:index+connectionMessageProofSize]...)
	index += connectionMessageProofSize

	msg.MessageNonce = append(msg.MessageNonce[:], bytes[index:index+connectionMessageNonceSize]...)
	index += connectionMessageNonceSize

	msg.Versions = bytesToVersions(bytes[index:])
	return
}
