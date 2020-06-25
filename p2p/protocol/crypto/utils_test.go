package crypto

import (
	"encoding/hex"
	"testing"
)

func TestCalcPeerId(t *testing.T) {

	type testCaseType struct {
		expectedID string
		publicKey  string
	}

	cases := [...]testCaseType{
		{"idr9VuS3wKx7CQdhAGSR6Z2e9hjD4b", "9e83ee6e795be9551bb11d026034923ff3361fad0b3108b996797c5640588341"},
		{"idtm6jM1Fz1eFu9D2Sn34AruzAxVbW", "87ae33eef29911bbf1a2cee6386efe1879d5589c85cd2b156d5efed7b1290826"},
		{"idsMpA9XmaY8hMx1pypAzv3eisWtud", "7a8bc8486d5279c4bd16f720b1f5fc55b82e7507d7cc5bb4a368b31f11c4e334"},
		{"idqthtAP6gEYCMJbFBvxKiND8HSzoJ", "4e33c71c0265e2dc91e028bdf950052315f354709d6a21df98bb3710ae629649"},
	}

	for _, testCase := range cases {

		publicKeyBuf, err := hex.DecodeString(testCase.publicKey)
		if err != nil {
			t.Error(err)
		}

		peerID, err := CalcPeerID(publicKeyBuf)
		if err != nil {
			t.Error(err)
		}

		if testCase.expectedID != peerID {
			t.Errorf("Expected %s, but %s", testCase.expectedID, peerID)
		}
	}
}

func TestGeneratePrecomputedKey(t *testing.T) {

	pk, _ := hex.DecodeString("96678b88756dd6cfd6c129980247b70a6e44da77823c3672a2ec0eae870d8646")
	sk, _ := hex.DecodeString("a18dc11cb480ebd31081e1541df8bd70c57da0fa419b5036242f8619d605e75a")

	expectedKey := "5228751a6f5a6494e38e1042f578e3a64ae3462b7899356f49e50be846c9609c"

	sharedKeyBuff := PrecomputeSharedKey(pk, sk)
	sharedKey := hex.EncodeToString(sharedKeyBuff[:])

	if expectedKey != sharedKey {
		t.Errorf("Expected %s, but key = %s", expectedKey, sharedKey)
	}

}

func TestEncryptMessage(t *testing.T) {
	nonce, _ := hex.DecodeString("8dde158c55cff52f4be9352787d333e616a67853640d72c5")
	msg, _ := hex.DecodeString("00874d1b98317bd6efad8352a7144c9eb0b218c9130e0a875973908ddc894b764ffc0d7f176cf800b978af9e919bdc35122585168475096d0ebcaca1f2a1172412b91b363ff484d1c64c03417e0e755e696c386a0000002d53414e44424f5845445f54455a4f535f414c5048414e45545f323031382d31312d33305431353a33303a35365a00000000")
	shKey := new([32]byte)
	key, _ := hex.DecodeString("5228751a6f5a6494e38e1042f578e3a64ae3462b7899356f49e50be846c9609c")
	copy(shKey[:], key)

	encrypted, _ := EncryptMessage(msg, nonce, shKey)
	encryptedMsg := hex.EncodeToString(encrypted)
	expected := "45d82d5c4067f5c32748596c1bbc93a9f87b5b1f2058ddd82b6f081ca484b672395c7473ab897c64c01c33878ac1ccb6919a75c9938d8bcf0e7917ddac13a787cfb5c9a5aea50d24502cf86b5c9b000358c039334ec077afe98936feec0dabfff35f14cafd2cd3173bbd56a7c6e5bf6f5f57c92b59b129918a5895e883e7d999b191aad078c4a5b164144c1beaed58b49ba9be094abf3a3bd9"

	if expected != encryptedMsg {
		t.Errorf("Expected %s, but %s", expected, encryptedMsg)
	}
}

func TestDecryptMessage(t *testing.T) {
	nonce, _ := hex.DecodeString("8dde158c55cff52f4be9352787d333e616a67853640d72c5")
	enc, _ := hex.DecodeString("45d82d5c4067f5c32748596c1bbc93a9f87b5b1f2058ddd82b6f081ca484b672395c7473ab897c64c01c33878ac1ccb6919a75c9938d8bcf0e7917ddac13a787cfb5c9a5aea50d24502cf86b5c9b000358c039334ec077afe98936feec0dabfff35f14cafd2cd3173bbd56a7c6e5bf6f5f57c92b59b129918a5895e883e7d999b191aad078c4a5b164144c1beaed58b49ba9be094abf3a3bd9")
	shKey := new([32]byte)
	key, _ := hex.DecodeString("5228751a6f5a6494e38e1042f578e3a64ae3462b7899356f49e50be846c9609c")
	copy(shKey[:], key)

	decrypted, successed := DecryptMessage(enc, nonce, shKey)
	msg := hex.EncodeToString(decrypted)
	expected := "00874d1b98317bd6efad8352a7144c9eb0b218c9130e0a875973908ddc894b764ffc0d7f176cf800b978af9e919bdc35122585168475096d0ebcaca1f2a1172412b91b363ff484d1c64c03417e0e755e696c386a0000002d53414e44424f5845445f54455a4f535f414c5048414e45545f323031382d31312d33305431353a33303a35365a00000000"

	if expected != msg || !successed {
		t.Errorf("Expected %s, but %s", expected, msg)
	}

	shKey = new([32]byte)
	_, successed = DecryptMessage(enc, nonce, shKey)
	if successed {
		t.Errorf("Seriously? Did you decrypte the message?")
	}
}
