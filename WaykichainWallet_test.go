package wiccwallet

import (
	_ "fmt"
	"testing"
	"fmt"
	"encoding/hex"
	"crypto/ecdsa"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcd/btcec"
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
	mnemonic := "vote despair mind rescue crumble choice garden elite venture cattle oxygen voyage"//"empty regular curve turtle student prize toy accuse develop spike scatter ginger"
	address := GetAddressFromMnemonic(mnemonic, WAYKI_TESTNET)
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
	checkPriv,_:=CheckPrivateKey("Y6J4aK6Wcs4A3Ex4HXdfjJ6ZsHpNZfjaS4B9w7xqEnmFEYMqQd13",2)
	fmt.Println("公钥",str,"测试私钥？",checkPriv)
}


func TestSignMessage(t *testing.T) {
	msg := "WaykiChain"
	privateKey := "Y9dJaHVk7Rs4sVq1Uk8TGLW4PQzzesA7Lss2Xz1inZY9KMfHBSPE"
	signMsg,_:=SignMessage(privateKey,msg)
	fmt.Println(signMsg)

	wifKey, _ := btcutil.DecodeWIF(privateKey)
	key := wifKey.PrivKey
	publicKey := key.PubKey().ToECDSA()
	decode,_:=	hex.DecodeString(signMsg.SignMessage)
	fmt.Println(signMsg.PublicKey)
	sign,_:=btcec.ParseDERSignature(decode, btcec.S256())
	success := ecdsa.Verify(publicKey, btcutil.Hash160([]byte("WaykiChain")), sign.R, sign.S)
	fmt.Println("验证签名成功？", success)
}
