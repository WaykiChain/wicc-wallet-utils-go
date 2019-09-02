package wiccwallet

import (
	_ "fmt"
	"testing"
	"fmt"
)

/*
生成助记词
generate Mnemonics
*/
func TestGenerateMnemonics(t *testing.T) {
	mnemonic := GenerateMnemonics()
	if mnemonic == "" {
		t.Error("GenerateMnemonics err!")
	}

	t.Log("mnemonic=", mnemonic)
}

/*
助记词生成地址
Mnemonics to Address
*/
func TestGetAddressFromMnemonic(t *testing.T) {
	mnemonic := "empty regular curve turtle student prize toy accuse develop spike scatter ginger"
	address := GetAddressFromMnemonic(mnemonic, WAYKI_MAINTNET)
	if address == "" {
		t.Error("GenerateAddress err!")
	}
	t.Log("address: " + address)
}


func TestMnemonicWIF(t *testing.T) {
	mnemonic := "empty regular curve turtle student prize toy accuse develop spike scatter ginger"
	privateKey := GetPrivateKeyFromMnemonic(mnemonic, WAYKI_MAINTNET)
	fmt.Println("私钥" + privateKey)
	if privateKey == "" {
		t.Error("MnemonicWIF error!")
		return
	}
	address := GetAddressFromPrivateKey(privateKey, WAYKI_MAINTNET)
	t.Log("地址", address)
}

/*
获得公钥
get publicKey hex String
*/
func TestGetPubKey(t *testing.T) {
	str,_:=GetPubKeyFromPrivateKey("Y6J4aK6Wcs4A3Ex4HXdfjJ6ZsHpNZfjaS4B9w7xqEnmFEYMqQd13")
	println(str)
}
