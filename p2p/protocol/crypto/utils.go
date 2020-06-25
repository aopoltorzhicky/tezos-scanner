package crypto

import (
	"crypto/sha256"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/nacl/box"
)

type prefix = []byte

var (
	// For (de)constructing addresses
	publicKeyHash = []byte{153, 103}
)

func encodeBase58Check(payload []byte, prefix prefix) string {
	var data []byte
	data = append(data, prefix...)
	data = append(data, payload...)

	h := sha256.Sum256(data)
	hash := sha256.Sum256(h[:])
	checksum := hash[:4]

	data = append(data, checksum...)

	return base58.Encode(data)
}

// CalcPeerID -
func CalcPeerID(publicKey []byte) (peerID string, err error) {
	cryptoBlake, err := blake2b.New(16, nil)
	if err != nil {
		return
	}
	_, err = cryptoBlake.Write(publicKey)
	if err != nil {
		return
	}

	peerIDBuff := cryptoBlake.Sum(nil)
	peerID = encodeBase58Check(peerIDBuff, publicKeyHash)

	return
}

// PrecomputeSharedKey -
func PrecomputeSharedKey(publicKey []byte, privateKey []byte) (key [32]byte) {
	privateKeyBuff := [32]byte{}
	copy(privateKeyBuff[:], privateKey)
	publicKeyBuff := [32]byte{}
	copy(publicKeyBuff[:], publicKey)
	box.Precompute(&key, &publicKeyBuff, &privateKeyBuff)
	return
}

// EncryptMessage -
func EncryptMessage(msg []byte, nonce Nonce, sharedKey *[32]byte) (encrypted []byte, err error) {
	if nonceSize != len(nonce) {
		err = fmt.Errorf("Nonce's size must be 24 byte")
		return
	}

	var tmp [24]byte
	copy(tmp[:], nonce[0:24])
	//var outBuff []byte
	encrypted = box.SealAfterPrecomputation(nil, msg, &tmp, sharedKey)
	return
}

// DecryptMessage -
func DecryptMessage(enc []byte, nonce Nonce, sharedKey *[32]byte) (msg []byte, successed bool) {
	if nonceSize != len(nonce) {
		return nil, false
	}

	var tmp [24]byte
	copy(tmp[:], nonce)

	msg, successed = box.OpenAfterPrecomputation(nil, enc[:], &tmp, sharedKey)
	return
}
