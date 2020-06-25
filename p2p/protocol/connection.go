package protocol

import (
	"encoding/binary"
	"errors"
	"net"
	"time"

	"github.com/aopoltorzhicky/tezos-scanner/p2p/protocol/crypto"
)

func sendConnectionMessage(conn net.Conn, message ConnectionMessage) (err error) {
	data := message.toBytes()

	err = conn.SetWriteDeadline(time.Now().Add(time.Second * 6))
	if err != nil {
		return
	}
	_, err = conn.Write(data)
	if err != nil {
		return
	}
	return
}

func receiveConnectionMessage(conn net.Conn) (message *ConnectionMessage, err error) {
	err = conn.SetReadDeadline(time.Now().Add(time.Second * 6))
	if err != nil {
		return
	}

	buf := make([]byte, 1024)
	received, err := conn.Read(buf)
	if err != nil {
		return
	}

	message = newMessage(buf[:received])
	return
}

func sendEncryptedMessage(conn net.Conn, message []byte, nonce nonceType, precomputedKey *[32]byte) (err error) {
	// TODO: Split by chunks
	err = conn.SetWriteDeadline(time.Now().Add(time.Second * 10))
	if err != nil {
		return
	}

	encrypted, err := crypto.EncryptMessage(message, nonce, precomputedKey)
	if err != nil {
		return
	}

	data := addSize(encrypted)
	_, err = conn.Write(data)
	if err != nil {
		return
	}
	return
}

func receiveEncryptedMessage(conn net.Conn, nonce nonceType, precomputedKey *[32]byte) (message []byte, err error) {
	err = conn.SetReadDeadline(time.Now().Add(time.Second * 10))
	if err != nil {
		return
	}

	buff := make([]byte, 65536)
	_, err = conn.Read(buff)
	if err != nil {
		return
	}
	size := binary.BigEndian.Uint16(buff[:2])

	message, success := crypto.DecryptMessage(buff[2:size+2], nonce, precomputedKey)
	if !success {
		return nil, errors.New("can't decrypt the message")
	}

	return
}

func receiveMetaMessage(conn net.Conn, nonce nonceType, precomputedKey *[32]byte) (message *metadata, err error) {
	data, err := receiveEncryptedMessage(conn, nonce, precomputedKey)
	if err != nil {
		return
	}

	if len(data) < 2 {
		return nil, errors.New("received wrong data")
	}

	message = new(metadata)
	message.fromBytes(data)
	return
}

func receiveAckMessage(conn net.Conn, nonce nonceType, precomputedKey *[32]byte) (message *ackMessage, err error) {
	data, err := receiveEncryptedMessage(conn, nonce, precomputedKey)
	if err != nil {
		return
	}

	if len(data) < 1 {
		return nil, errors.New("received wrong data")
	}

	message = new(ackMessage)
	message.fromBytes(data)
	return
}
