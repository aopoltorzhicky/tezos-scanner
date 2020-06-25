package crypto

import (
	"encoding/hex"
	"math/big"
	"testing"
)

func TestGenerateNonces(t *testing.T) {
	{
		sentMsg, _ := hex.DecodeString("00874d1b98317bd6efad8352a7144c9eb0b218c9130e0a875973908ddc894b764ffc0d7f176cf800b978af9e919bdc35122585168475096d0ebcaca1f2a1172412b91b363ff484d1c64c03417e0e755e696c386a0000002d53414e44424f5845445f54455a4f535f414c5048414e45545f323031382d31312d33305431353a33303a35365a00000000")
		recvMsg, _ := hex.DecodeString("00874d1ab3845960b32b039fef38ca5c9f8f867df1d522f27a83e07d9dfbe3b296a6c076412d98b369ab015d57247e5380d708b9edfcca0ca2c865346ef9c3d7ed00182cf4f613a6303c9b2a28cda8ff93687bd20000002d53414e44424f5845445f54455a4f535f414c5048414e45545f323031382d31312d33305431353a33303a35365a00000000")

		l, r := GenerateNonces(sentMsg, recvMsg, false)

		ls := hex.EncodeToString(l)
		if ls != "8dde158c55cff52f4be9352787d333e616a67853640d72c5" {
			t.Error("Expected '8dde158c55cff52f4be9352787d333e616a67853640d72c5', but it's", ls)
		}

		rs := hex.EncodeToString(r)
		if rs != "e67481a23cf9b404626a12bd405066e161b32dc53f469153" {
			t.Error("Expected 'e67481a23cf9b404626a12bd405066e161b32dc53f469153', but it's", rs)
		}
	}

	{
		sentMsg, _ := hex.DecodeString("00874d1b98317bd6efad8352a7144c9eb0b218c9130e0a875973908ddc894b764ffc0d7f176cf800b978af9e919bdc35122585168475096d0ebcaca1f2a1172412b91b363ff484d1c64c03417e0e755e696c386a0000002d53414e44424f5845445f54455a4f535f414c5048414e45545f323031382d31312d33305431353a33303a35365a00000000")
		recvMsg, _ := hex.DecodeString("00874d1ab3845960b32b039fef38ca5c9f8f867df1d522f27a83e07d9dfbe3b296a6c076412d98b369ab015d57247e5380d708b9edfcca0ca2c865346ef9c3d7ed00182cf4f613a6303c9b2a28cda8ff93687bd20000002d53414e44424f5845445f54455a4f535f414c5048414e45545f323031382d31312d33305431353a33303a35365a00000000")

		l, r := GenerateNonces(sentMsg, recvMsg, true)

		ls := hex.EncodeToString(l)
		if ls != "ff0451d94af9f75a46d74a2a9f685cff20222a15829f121d" {
			t.Error("Expected 'ff0451d94af9f75a46d74a2a9f685cff20222a15829f121d', but it's", ls)
		}

		rs := hex.EncodeToString(r)
		if rs != "8a09a2c43a61aa6eccee084aa66da9bc94b441b17615be58" {
			t.Error("Expected '8a09a2c43a61aa6eccee084aa66da9bc94b441b17615be58', but it's", rs)
		}
	}
}

func TestIncrementNonce(t *testing.T) {
	// Увеличение nonce на один. Надо перенести в библиотеку и юзать в отправке сообщений.
	bytes, _ := hex.DecodeString("0000000000cff52f4be9352787d333e616a67853640d72c5")

	a := new(big.Int).SetBytes(bytes)
	b := big.NewInt(1)
	c := big.NewInt(0)
	c = c.Add(a, b)

	d := make([]byte, 24)
	tmp := c.Bytes()
	copy(d[len(d)-len(tmp):], tmp)
	s := hex.EncodeToString(d)
	if s != "0000000000cff52f4be9352787d333e616a67853640d72c6" {
		t.Error("wrong answer!", s)
	}
}
