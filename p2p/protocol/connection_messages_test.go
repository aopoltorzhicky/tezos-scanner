package protocol

import (
	"encoding/hex"
	"testing"
)

func TestNewMessage(t *testing.T) {
	messages := [...]string{
		"007c2604911787157ac31ac86ba49ae929ec665c6b955eb51cfcddb8dde96c71cc0f9b2536daf0fc3b618be0d6e04396ac77664535295b0c2656721493d6c99cc0d985948a128278b3c02f42499d9fddbc7693020000002254455a4f535f5a45524f4e45545f323031392d30382d30365431353a31383a35365a00000000",
	}

	expectedNonces := [...]string{
		"93d6c99cc0d985948a128278b3c02f42499d9fddbc769302",
	}

	expectedProofOfWorkStamp := [...]string{
		"36daf0fc3b618be0d6e04396ac77664535295b0c26567214",
	}

	expectedPublicKey := [...]string{
		"911787157ac31ac86ba49ae929ec665c6b955eb51cfcddb8dde96c71cc0f9b25",
	}

	expectedVersionNames := [...]string{
		"TEZOS_ZERONET_2019-08-06T15:18:56Z",
	}

	for i, v := range messages {
		data, _ := hex.DecodeString(v)
		con := newMessage(data)

		if hex.EncodeToString(con.MessageNonce) != expectedNonces[i] {
			t.Errorf("%s != %s", hex.EncodeToString(con.MessageNonce), expectedNonces[i])
		}
		if hex.EncodeToString(con.PublicKey) != expectedPublicKey[i] {
			t.Errorf("%s != %s", hex.EncodeToString(con.PublicKey), expectedPublicKey[i])
		}
		if hex.EncodeToString(con.ProofOfWorkStamp) != expectedProofOfWorkStamp[i] {
			t.Errorf("%s != %s", hex.EncodeToString(con.ProofOfWorkStamp), expectedProofOfWorkStamp[i])
		}
		if con.Versions[0].Name != expectedVersionNames[i] {
			t.Errorf("%s != %s", con.Versions[0].Name, expectedVersionNames[i])
		}
	}
}

func TestMessageToBytes(t *testing.T) {
	expectedMessages := [...]string{
		"007c2604911787157ac31ac86ba49ae929ec665c6b955eb51cfcddb8dde96c71cc0f9b2536daf0fc3b618be0d6e04396ac77664535295b0c2656721493d6c99cc0d985948a128278b3c02f42499d9fddbc7693020000002254455a4f535f5a45524f4e45545f323031392d30382d30365431353a31383a35365a00000000",
	}

	nonces := [...]string{
		"93d6c99cc0d985948a128278b3c02f42499d9fddbc769302",
	}

	proofOfWorkStamp := [...]string{
		"36daf0fc3b618be0d6e04396ac77664535295b0c26567214",
	}

	publicKey := [...]string{
		"911787157ac31ac86ba49ae929ec665c6b955eb51cfcddb8dde96c71cc0f9b25",
	}

	versionNames := [...]string{
		"TEZOS_ZERONET_2019-08-06T15:18:56Z",
	}
	ports := [...]uint16{
		9732,
	}

	for i, expected := range expectedMessages {
		pubKey, _ := hex.DecodeString(publicKey[i])
		bytePow, _ := hex.DecodeString(proofOfWorkStamp[i])
		nonce, _ := hex.DecodeString(nonces[i])
		connMessage := ConnectionMessage{
			Port: ports[i],
			Versions: []Version{
				{
					Name:  versionNames[i],
					Major: 0,
					Minor: 0,
				},
			},
			PublicKey:        pubKey,
			ProofOfWorkStamp: bytePow,
			MessageNonce:     nonce,
		}

		data := connMessage.toBytes()

		if hex.EncodeToString(data) != expected {
			t.Errorf("%s != %s", hex.EncodeToString(data), expected)
		}
	}
}

func TestAddSize(t *testing.T) {
	messages := [...]string{
		"000000000000000000000000000000000000",
		"2604911787157ac31ac86ba49ae929ec665c6b955eb51cfcddb8dde96c71cc0f9b2536daf0fc3b618be0d6e04396ac77664535295b0c2656721493d6c99cc0d985948a128278b3c02f42499d9fddbc7693020000002254455a4f535f5a45524f4e45545f323031392d30382d30365431353a31383a35365a00000000",
	}

	expectedMessages := [...]string{
		"0012000000000000000000000000000000000000",
		"007c2604911787157ac31ac86ba49ae929ec665c6b955eb51cfcddb8dde96c71cc0f9b2536daf0fc3b618be0d6e04396ac77664535295b0c2656721493d6c99cc0d985948a128278b3c02f42499d9fddbc7693020000002254455a4f535f5a45524f4e45545f323031392d30382d30365431353a31383a35365a00000000",
	}

	for i, msg := range messages {
		data, _ := hex.DecodeString(msg)
		buff := addSize(data)
		res := hex.EncodeToString(buff)

		if res != expectedMessages[i] {
			t.Errorf("%s != %s", expectedMessages[i], res)
		}
	}
}
