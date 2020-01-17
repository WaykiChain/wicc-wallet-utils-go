package waykichain

import (
	"fmt"
	"github.com/WaykiChain/wicc-wallet-utils-go/common"
	"testing"
)

func TestGenerateAddressFromMnemonic(t *testing.T){

	mnemonic:= "wreck bullet carpet nerve belt border often trust exchange believe defense rebel"

	WICCAddress , err :=  WICCW.GenerateAddressFromMnemonic(mnemonic,common.English)
	WICCAddressTestnet , err :=  WICCTestnetW.GenerateAddressFromMnemonic(mnemonic,common.English)

	if err != nil{
		t.Errorf("Failed to TestImportWalletFromMnemonic: %v", err)
	}

	//WhHGDnhL3ny5VES9CFuA38WqDhAQ4VNGuo
	t.Log("TestImportWalletFromMnemonic , WICCAddress=",WICCAddress)
	//wLYLCxsBDjbRiPVEzvbX2bgFftqnWuQxB7
	t.Log("TestImportWalletFromMnemonic , WICCAddressTestnet=",WICCAddressTestnet)
}

//Test Generate Address From PrivateKey
func TestGenerateAddressFromPrivateKey(t *testing.T){

	WICCPrivateKey := "PemqPzcsCJXjU4PovGSC9zBv89YZSisrWePF9N1skxdEbftbdkDo"
	WICC_TestnetPrivateKey := "Y6amwxjHqUM37UrquokPsbCXTRNughoM27gDUGfbXhJikS39i9h1"

	WICCAddress, err := WICCW.GenerateAddressFromPrivateKey(WICCPrivateKey)
	if err != nil {
		t.Errorf("Failed to Import WICCPrivateKey: %v",err)
	}
	WICCAddressTestnet, err := WICCTestnetW.GenerateAddressFromPrivateKey(WICC_TestnetPrivateKey)
	if err != nil {
		t.Errorf("Failed to Import WICC_TestnetPrivateKey: %v",err)
	}

	//WhHGDnhL3ny5VES9CFuA38WqDhAQ4VNGuo
	t.Log("TestGenerateAddressFromPrivateKey: WICCAddress=",WICCAddress)
	//wLYLCxsBDjbRiPVEzvbX2bgFftqnWuQxB7
	t.Log("TestGenerateAddressFromPrivateKey: WICCAddressTestnet=",WICCAddressTestnet)
}

//Test Export PrivateKey by Mnemonics
func TestExportPrivateFromMnemonic(t *testing.T){
	mnemonic := "wreck bullet carpet nerve belt border often trust exchange believe defense rebel"

	WICCPrivateKey ,err := WICCW.ExportPrivateKeyFromMnemonic(mnemonic,common.English)

	WICC_TESTNETPrivateKey ,err := WICCTestnetW.ExportPrivateKeyFromMnemonic(mnemonic,common.English)
	if err != nil {
		t.Errorf("Failed to TestExportPrivateFromMnemonic: %v",err)
	}

	//PemqPzcsCJXjU4PovGSC9zBv89YZSisrWePF9N1skxdEbftbdkDo
	t.Log("WICCPrivateKey=",WICCPrivateKey)
	//Y6amwxjHqUM37UrquokPsbCXTRNughoM27gDUGfbXhJikS39i9h1
	t.Log("WICC_TESTNETPrivateKey=",WICC_TESTNETPrivateKey)
}

func TestCheckAddress(t *testing.T){

	WICCAddress := "WhHGDnhL3ny5VES9CFuA38WqDhAQ4VNGuo"
	isValid,err := WICCW.CheckAddress(WICCAddress)
	if err != nil {
		t.Errorf("Failed to Check WICCAddress: %v",err)
	}
	t.Log("TestCheckAddress: WICCAddress=",isValid)

	WICCAddressTestnet := "wLYLCxsBDjbRiPVEzvbX2bgFftqnWuQxB7"
	isValid,err = WICCW.CheckAddress(WICCAddressTestnet)
	if err != nil {
		t.Errorf("Failed to Check WICCAddressTestnet: %v",err)
	}
	t.Log("TestCheckAddress: WICCAddressTestnet=",isValid)
}

func TestCheckWICCPrivateKey(t *testing.T){

	WICCPrivateKey := "PemqPzcsCJXjU4PovGSC9zBv89YZSisrWePF9N1skxdEbftbdkDo"
	isValid,err := WICCW.CheckPrivateKey(WICCPrivateKey)
	if err != nil {
		t.Errorf("Failed to Check WICCPrivateKey: %v",err)
	}
	t.Log("CheckWICCPrivateKey BTCSegwitPrivateKey=",isValid)

	WICC_TestnetPrivateKey := "Y6amwxjHqUM37UrquokPsbCXTRNughoM27gDUGfbXhJikS39i9h1"
	isValid,err = WICCW.CheckPrivateKey(WICC_TestnetPrivateKey)
	if err != nil {
		t.Errorf("Failed to Check WICC_TestnetPrivateKey: %v",err)
	}
	t.Log("CheckWICCPrivateKey WICC_TestnetPrivateKey=",isValid)
}

/*
获得公钥
get publicKey hex String
*/
func TestGetPubKeyFromPrivateKey(t *testing.T) {

	WICCPrivateKey := "PemqPzcsCJXjU4PovGSC9zBv89YZSisrWePF9N1skxdEbftbdkDo"
	publicKey, err := WICCW.GetPubKeyFromPrivateKey(WICCPrivateKey)
	if err != nil {
		t.Errorf("Failed to GetPubKeyFromPrivateKey WICCPrivateKey: %v", err)
	}
	t.Log("GetPubKeyFromPrivateKey WICCPrivateKey=", publicKey)

	WICCTestnetPrivateKey := "Y6amwxjHqUM37UrquokPsbCXTRNughoM27gDUGfbXhJikS39i9h1"
	publicKey, err = WICCW.GetPubKeyFromPrivateKey(WICCTestnetPrivateKey)
	if err != nil {
		t.Errorf("Failed to GetPubKeyFromPrivateKey WICCTestnetPrivateKey: %v", err)
	}
	t.Log("GetPubKeyFromPrivateKey WICCTestnetPrivateKey=", publicKey)
}

func TestSignMessage(t *testing.T) {
	//地址:wX7cC6qK6RQCLShCevpeciqQaQNEtqLRa8
	//钱包地址对应的私钥:Y8WXc3RYw4TRxdGEpTLPd5GR7VrsAvRgCdiZMZakwFyVST1P7NnC
	//公钥:034edcac8efda301a0919cdf2feeb0376bfcd2a1a29b5d094e5e9ce7a580c82fcc (压缩后)
	input := SignMsgInput{}
	input.Data = "WaykiChain" //原始数据,由开发者后台生成传给前端,生成规则由开发者自己决定
	input.PrivateKey = "Y6amwxjHqUM37UrquokPsbCXTRNughoM27gDUGfbXhJikS39i9h1"
	signResult, _ := input.SignMessage() //签名结果，包含签名后信息 + 签名者公钥

	fmt.Println("signResult: \n\tpublicKey=", signResult.PublicKey, "\n\tsignature=", signResult.Signature)
}

func TestVerifyMsgSignature(t *testing.T) {

	input := VerifySignInput{}
	input.Signature = "3045022100cbd99999466e02d7c824e39769c1c553776084d73106fa5cd980696cdf279ce9022050bd31a13e02e29674b8303a58451a41e2f4fed875237db9dbb361ec9856b8cc"
	input.PublicKey = "031b27286c65b81ac13cfd4067b030398a19eb147f439c094fbb19a2f3ab9ec10b"
	input.Data = "WaykiChain"
	input.NetParams = common.WICCTestnetParams //testnet
//	netParams := &common.WICCParams  //mainnet

	isValid, address := input.VerifyMsgSignature()
	fmt.Println("VerifyMsgSignature Result:", isValid, ";Sign address：",address)
}
