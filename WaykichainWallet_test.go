package wiccwallet

import (
	"fmt"
	_ "fmt"
	//"github.com/btcsuite/btcutil"
	"testing"
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
	address,err := GetAddressFromMnemonic(mnemonic, WAYKI_MAINTNET)
	if err != nil {
		t.Error("GenerateAddress err!",err)
	}
	t.Log("address: " + address)
}


func TestMnemonicWIF(t *testing.T) {
	mnemonic := "empty regular curve turtle student prize toy accuse develop spike scatter ginger"
	privateKey,err := GetPrivateKeyFromMnemonic(mnemonic, WAYKI_MAINTNET)
	fmt.Println("私钥" + privateKey)
	if err != nil {
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
	//地址:wX7cC6qK6RQCLShCevpeciqQaQNEtqLRa8
	//钱包地址对应的私钥:Y8WXc3RYw4TRxdGEpTLPd5GR7VrsAvRgCdiZMZakwFyVST1P7NnC
	//公钥:034edcac8efda301a0919cdf2feeb0376bfcd2a1a29b5d094e5e9ce7a580c82fcc (压缩后)
	msg := "WaykiChain" //原始数据,由开发者后台生成传给前端,生成规则由开发者自己决定
	privateKey := "Y8WXc3RYw4TRxdGEpTLPd5GR7VrsAvRgCdiZMZakwFyVST1P7NnC"
	signResult, _ := SignMessage(privateKey, msg) //签名结果，包含签名后信息 + 签名者公钥

	fmt.Println("signResult: \n\tpublicKey=", signResult.PublicKey, "\n\tsignature=", signResult.SignMessage)
}

func TestVerifyMsgSignature(t *testing.T) {

	signature := "3044022024fafdf62a8414ad28c96354cc310daffee04e8ad46276420bdaafe1aa35091e02205b2c1b1a1e7fe97a74f2e3dc16f790a28cafea2ec40911fd40cff856899a851"
	publicKey := "034edcac8efda301a0919cdf2feeb0376bfcd2a1a29b5d094e5e9ce7a580c82fcc"
	msg := "WaykiChain"
	netType := WAYKI_TESTNET
//	netType := WAYKI_MAINTNET

	isValid, address := VerifyMsgSignature(signature,publicKey,msg,netType)
	fmt.Println("VerifyMsgSignature Result:", isValid, ";Sign address：",address)
}
