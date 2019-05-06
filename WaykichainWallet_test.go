package wiccwallet

import (
	"encoding/hex"
	_ "fmt"
	"io/ioutil"
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
// 	mnemonic := "empty regular curve turtle student prize toy accuse develop spike scatter ginger"
// 	privateKey := GetAddressFromPrivateKey(mnemonic, WAYKI_MAINTNET)
// 	fmt.Println("私钥" + privateKey)
	mnemonic := "Y9sx4Y8sBAbWDAqAWytYuUnJige3ZPwKDZp1SCDqqRby1YMgRG9c"
	privateKey := GetAddressFromPrivateKey(mnemonic, WAYKI_MAINTNET)
	if privateKey == "" {
		t.Error("MnemonicWIF error!")
		return
	}
	t.Log("私钥",privateKey)
}

func TestSignCallContractTx(t *testing.T) {

	privateKey := "Y9sx4Y8sBAbWDAqAWytYuUnJige3ZPwKDZp1SCDqqRby1YMgRG9c"

	var txParam CallContractTxParam
	txParam.ValidHeight = 22365
	txParam.SrcRegId = "7849-1"
	txParam.AppId = "20988-1"
	txParam.Fees = 100000
	txParam.Values = 10000
	txParam.ContractHex = "f017"

	hash, err := SignCallContractTx(privateKey, &txParam)
	if err != nil {
		t.Error("SignCallContractTx err: ", err)
	}
	println(hash)
}

func TestSignDelegateTx(t *testing.T) {
	privateKey := "Y9sx4Y8sBAbWDAqAWytYuUnJige3ZPwKDZp1SCDqqRby1YMgRG9c"
	var txParams DelegateTxParam
	txParams.ValidHeight = 95728
	txParams.SrcRegId = "0-1"
	txParams.Fees = 10000
	txParams.Votes = NewOperVoteFunds()
	pubKey, _ := hex.DecodeString("025a37cb6ec9f63bb17e562865e006f0bafa9afbd8a846bd87fc8ff9e35db1252e")
	vote:=OperVoteFund{PubKey:pubKey,VoteValue:10000}
	txParams.Votes.Add(&vote)

	hash, err := SignDelegateTx(privateKey, &txParams)
	if err != nil {
		t.Error("SignDelegateTx err: ", err)
	}
	println(hash)
}

func TestSignRegisterAccountTx(t *testing.T) {

	privateKey := "Y9sx4Y8sBAbWDAqAWytYuUnJige3ZPwKDZp1SCDqqRby1YMgRG9c"
	var txParam RegisterAccountTxParam
	txParam.ValidHeight = 7783
	txParam.Fees = 10000

	hash, err := SignRegisterAccountTx(privateKey, &txParam)
	if err != nil {
		t.Error("SignRegisterAccountTx err: ", err)
	}
	println(hash)
}

func TestSignCommonTx(t *testing.T) {

	privateKey := "Y7V1jwCRr8D3tyPTkcsjgBTHwZN45b1U3ueZfJ5oWVJqwcKpArou"
	var txParams CommonTxParam
	txParams.ValidHeight = 14897
	txParams.SrcRegId = "158-1"
	txParams.DestAddr = "wSSbTePArv6BkDsQW9gpGCTX55AXVxVKbd"
	txParams.Values = 10000
	txParams.Fees = 10000

	hash, err := SignCommonTx(privateKey, &txParams)
	if err != nil {
		t.Error("SignCommonTx err: ", err)
	}
	println(hash)
}

func TestSignRegisterContractTx(t *testing.T) {

	privateKey := "Y9sx4Y8sBAbWDAqAWytYuUnJige3ZPwKDZp1SCDqqRby1YMgRG9c"

	script, err := ioutil.ReadFile("./demo/data/hello.lua")
	if err != nil {
		t.Error("Read contract script file err: ", err)
	}

	var txParam RegisterContractTxParam
	txParam.ValidHeight = 20999
	txParam.SrcRegId = "7849-1"
	txParam.Fees = 110000000
	txParam.Script = script
	txParam.Description = "My hello contract!!!"

	hash, err := SignRegisterContractTx(privateKey, &txParam)
	if err != nil {
		t.Error("SignRegisterContractTx err: ", err)
	}
	println(hash)
}
