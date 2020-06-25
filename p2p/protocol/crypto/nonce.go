package crypto

import (
	"math/big"

	"golang.org/x/crypto/blake2b"
)

// Nonce -
type Nonce = []byte

const nonceSize = 24

// GenerateNonces -
func GenerateNonces(sentMsg []byte, recvMsg []byte, incoming bool) (local Nonce, remote Nonce) {

	local = make([]byte, nonceSize)
	remote = make([]byte, nonceSize)

	var initMsg []byte
	var respMsg []byte
	if incoming {
		initMsg = recvMsg
		respMsg = sentMsg
	} else {
		initMsg = sentMsg
		respMsg = recvMsg
	}

	initToResp := append(initMsg, respMsg...)
	respToInit := make([]byte, len(initToResp))
	copy(respToInit, initToResp)

	initToResp = append(initToResp, []byte("Init -> Resp")...)
	respToInit = append(respToInit, []byte("Resp -> Init")...)

	irHash := blake2b.Sum256(initToResp)
	riHash := blake2b.Sum256(respToInit)

	if incoming {
		copy(local, irHash[:nonceSize])
		copy(remote, riHash[:nonceSize])
	} else {
		copy(local, riHash[:nonceSize])
		copy(remote, irHash[:nonceSize])
	}
	return
}

// NonceIncrement -
func NonceIncrement(nonce Nonce) Nonce {
	a := new(big.Int).SetBytes(nonce)
	b := big.NewInt(1)
	c := big.NewInt(0)
	c = c.Add(a, b)
	d := make([]byte, 24)
	tmp := c.Bytes()
	copy(d[len(d)-len(tmp):], tmp)
	nonce = d
	return nonce
}
