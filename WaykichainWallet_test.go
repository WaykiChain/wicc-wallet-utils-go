package wiccwallet

import (
	"fmt"
	"testing"
)

func TestGenerateMnemonics(t *testing.T) {
	mnemonic := GenerateMnemonics()
	if mnemonic == "" {
		t.Error("GenerateMnemonics err!")
	}

	t.Log("mnemonic=", mnemonic)
}

func TestGetAddressFromMnemonic(t *testing.T) {
	mnemonic := "empty regular curve turtle student prize toy accuse develop spike scatter ginger"
	//seed := bip.NewSeed(mnemonic, "")
	////fmt.Println(hex.EncodeToString(seed))
	address := GetAddressFromMnemonic(mnemonic, WAYKI_MAINTNET)
	if address == "" {
		t.Error("GenerateAddress err!")
	}
	t.Log("address: " + address)
}

func TestMnemonicWIF(t *testing.T) {
	mnemonic := "empty regular curve turtle student prize toy accuse develop spike scatter ginger"
	privateKey := GetAddressFromPrivateKey(mnemonic, WAYKI_MAINTNET)
	fmt.Println("私钥" + privateKey)
}
