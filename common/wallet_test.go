package common

import (
	"encoding/hex"
	"fmt"
	"github.com/JKinGH/go-hdwallet"
	"testing"
)

func TestGenerateMnemonics(t *testing.T) {

	mnemonic,err := BTCW.GenerateMnemonic(12)
	if err != nil {
		t.Errorf("Failed to GenerateMnemonic: %v",err)
	}

	t.Log("mnemonic=",mnemonic)
}

func TestChangeMnemonicLanguage(t *testing.T){
	//englishMnemonic := "jazz wine firm worth cry dumb glad foam viable knee pride purse"
	englishMnemonic := "wreck bullet carpet nerve belt border often trust exchange believe defense rebel"
	newMnemonic,err :=  BTCW.ChangeMnemonicLanguage(englishMnemonic, hdwallet.ChineseSimplified)
	if err != nil {
		t.Errorf("Failed to ChangeMnemonicLanguage: %v",err)
	}

	t.Log("mnemonic=",newMnemonic)
}

func TestGenerateSeed(t *testing.T){
	mnemonic := "trend memory raccoon escape crush nut arm alley melody spread spin cute"
//	mnemonic := "sinfonia amanita mangiare sugo duplice girone ognuno muovere vicinanza maglia caduco silenzio"
//	mnemonic := "泡 瑞 值 树 苗 什 饼 河 盛 师 划 诗"
//	mnemonic := "みとめる げざい ひろい いせい すぶり といれ ふよう つつじ てきとう いちじ はせる せつだん"
//	mnemonic :=  "的 一 是 在 不 了 有 和 人 这 中 大"
//	mnemonic :=  "abandon ability able about above absent absorb abstract absurd abuse access accident"


	password := hdwallet.DefaultPassword
	seed ,err  := GenerateSeed(mnemonic,password,hdwallet.English)

	if err != nil {
		t.Errorf("Failed to GenerateMnemonics: %v",err)
	}

	t.Log("seed=",hex.EncodeToString(seed))
}

func TestImportMnemonicCreateMultiWallet1(t *testing.T) {
	mnemonic:= "wreck bullet carpet nerve belt border often trust exchange believe defense rebel"

	master, err := hdwallet.NewKey(false,hdwallet.Mnemonic(mnemonic))
	if err != nil {
		panic(err)
	}

	// BTC: 1AwEPfoojHnKrhgt1vfuZAhrvPrmz7Rh4
	//wallet, _ := master.GetWallet(hdwallet.CoinType(hdwallet.BTC), hdwallet.AddressIndex(1))
	wallet, _ := master.GetWallet(hdwallet.CoinType(hdwallet.BTC))
	address, _ := wallet.GetAddress()

	addressP2WPKH, _ := wallet.GetKey().AddressP2WPKH()
	addressP2WPKHInP2SH, _ := wallet.GetKey().AddressP2WPKHInP2SH()
	fmt.Println("BTC: ", address, addressP2WPKH, addressP2WPKHInP2SH)

	//隔离见证
	// BTC_Segwit:
	wallet, _ = master.GetWallet(hdwallet.CoinType(hdwallet.BTC),hdwallet.Purpose(hdwallet.BIP49Purpose))
	address, _ = wallet.GetAddress()

	addressP2WPKH, _ = wallet.GetKey().AddressP2WPKH()
	addressP2WPKHInP2SH, _ = wallet.GetKey().AddressP2WPKHInP2SH()
	fmt.Println("BTC_Segwit: ", address, addressP2WPKH, addressP2WPKHInP2SH)

	//BTC_TESTNET
	wallet, _ = master.GetWallet(hdwallet.CoinType(hdwallet.BTC_TESTNET))
	address, _ = wallet.GetAddress()

	addressP2WPKH, _ = wallet.GetKey().AddressP2WPKH()
	addressP2WPKHInP2SH, _ = wallet.GetKey().AddressP2WPKHInP2SH()
	fmt.Println("BTC_TESTNET: ", address, addressP2WPKH, addressP2WPKHInP2SH)

	//隔离见证
	// BTCTestnet_Segwit:
	wallet, _ = master.GetWallet(hdwallet.CoinType(hdwallet.BTC_TESTNET),hdwallet.Purpose(hdwallet.BIP49Purpose))
	address, _ = wallet.GetAddress()

	addressP2WPKH, _ = wallet.GetKey().AddressP2WPKH()
	addressP2WPKHInP2SH, _ = wallet.GetKey().AddressP2WPKHInP2SH()
	fmt.Println("BTCTestnet_Segwit: ", address, addressP2WPKH, addressP2WPKHInP2SH)

	// BCH: 1CSBT18sjcCwLCpmnnyN5iqLc46Qx7CC91
	wallet, _ = master.GetWallet(hdwallet.CoinType(hdwallet.BCH))
	address, _ = wallet.GetAddress()
	addressBCH, _ := wallet.GetKey().AddressBCH()
	fmt.Println("BCH: ", address, addressBCH)

	// LTC: LLCaMFT8AKjDTvz1Ju8JoyYXxuug4PZZmS
	wallet, _ = master.GetWallet(hdwallet.CoinType(hdwallet.LTC))
	address, _ = wallet.GetAddress()
	fmt.Println("LTC: ", address)

	// DOGE: DHLA3rJcCjG2tQwvnmoJzD5Ej7dBTQqhHK
	wallet, _ = master.GetWallet(hdwallet.CoinType(hdwallet.DOGE))
	address, _ = wallet.GetAddress()
	fmt.Println("DOGE:", address)

	// ETH: 0x37039021cBA199663cBCb8e86bB63576991A28C1
	wallet, _ = master.GetWallet(hdwallet.CoinType(hdwallet.ETH))
	address, _ = wallet.GetAddress()
	fmt.Println("ETH: ", address)

	pubkeyBytes := wallet.GetKey().Public.SerializeCompressed()
	fmt.Println("pubkey=",hex.EncodeToString(pubkeyBytes))

	addressP2PKH, _ := wallet.GetKey().AddressP2PKH()
	addressP2WPKHInP2SH, _ = wallet.GetKey().AddressP2WPKHInP2SH()
	fmt.Println("ETH: ", address, addressP2PKH, addressP2WPKHInP2SH)

	// ETC: 0x480C69E014C7f018dAbF17A98273e90f0b0680cf
	wallet, _ = master.GetWallet(hdwallet.CoinType(hdwallet.ETC))
	address, _ = wallet.GetAddress()
	fmt.Println("ETC: ", address)
}


func TestGenerateAddressFromMnemonic(t *testing.T){

	mnemonic:= "wreck bullet carpet nerve belt border often trust exchange believe defense rebel"


	BTCAddress , err := BTCW.GenerateAddressFromMnemonic(mnemonic,hdwallet.English)
	BTCAddressSegwit , err := BTCSegwitW.GenerateAddressFromMnemonic(mnemonic,hdwallet.English)
	BTCAddressTestnet , err :=  BTCTestnetW.GenerateAddressFromMnemonic(mnemonic,hdwallet.English)
	BTCAddressTestnetSegwit , err :=  BTCTestnetSegwitW.GenerateAddressFromMnemonic(mnemonic,hdwallet.English)
	ETHAddress , err :=  ETHW.GenerateAddressFromMnemonic(mnemonic,hdwallet.English)
	ETHAddressLedger , err :=  ETHLedgerW.GenerateAddressFromMnemonic(mnemonic,hdwallet.English)
	WICCAddress , err :=  WICCW.GenerateAddressFromMnemonic(mnemonic,hdwallet.English)
	WICCAddressTestnet , err :=  WICCTestnetW.GenerateAddressFromMnemonic(mnemonic,hdwallet.English)

	if err != nil{
		t.Errorf("Failed to TestImportWalletFromMnemonic: %v", err)
	}

	//1AG89FCfPQvtVa7DSCUJjEagH13uTs28Zs
	t.Log("TestImportWalletFromMnemonic , BTCAddress=",BTCAddress)
	//3Ku5bUrN1gXM4fH8WVyRHbKkyGyydqJZ6F
	t.Log("TestImportWalletFromMnemonic , BTCAddressSegwit=",BTCAddressSegwit)
	//mxeBxFWLFAY3G1RKijr91B3kzsX2mTvnYv
	t.Log("TestImportWalletFromMnemonic , BTCAddressTestnet=",BTCAddressTestnet)
	//2MsNRcbHbMgwbbkfzx86Z4FdHkRp29NPjmD
	t.Log("TestImportWalletFromMnemonic , BTCAddressTestnetSegwit=",BTCAddressTestnetSegwit)
	//0x81FD1F7aE91041aAc5fCF7d8Ed3e1dd88Cc1359a
	t.Log("TestImportWalletFromMnemonic , ETHAddress=",ETHAddress)
	//0x791893c14f0a8dCa4ADB0A8297F8d12063865cd2
	t.Log("TestImportWalletFromMnemonic , ETHAddressLedger=",ETHAddressLedger)
	//WhHGDnhL3ny5VES9CFuA38WqDhAQ4VNGuo
	t.Log("TestImportWalletFromMnemonic , WICCAddress=",WICCAddress)
	//wLYLCxsBDjbRiPVEzvbX2bgFftqnWuQxB7
	t.Log("TestImportWalletFromMnemonic , WICCAddressTestnet=",WICCAddressTestnet)
}

//Test Generate Address From PrivateKey
func TestGenerateAddressFromPrivateKey(t *testing.T){

	BTCPrivateKey := "KyUg2abSHhZYP7bZFXKNDw6TnQoHLyJwDbbDvaNfsBsFxMbFCz4g"
	BTCSegwitPrivateKey := "KwLJcTWXgB6a14VaxPoJsZPe2o9GPdp2PLcEnrZnmPeRrAtmckgb"
	BTC_TestnetPrivateKey := "cUQyhR3BbeFMwtqFRrMppFtPPcx6DMQNFEXm8C1yNSzMkEoRoGYM"
	BTC_TestnetSegwitPrivateKey := "cT7214EqFAbtpfuMfg36EDHBmbrkErXLb27ERSKTLX1RUr92LFSE"
	//ETHPrivateKey := "494f8228ae5b6fda6bee1f44eb2c4ed120f210e06acaa8053763efb65638b315"
	WICCPrivateKey := "PemqPzcsCJXjU4PovGSC9zBv89YZSisrWePF9N1skxdEbftbdkDo"
	WICC_TestnetPrivateKey := "Y6amwxjHqUM37UrquokPsbCXTRNughoM27gDUGfbXhJikS39i9h1"

	BTCAddress, err := BTCW.GenerateAddressFromPrivateKey(BTCPrivateKey)
	if err != nil {
		t.Errorf("Failed to Import BTCPrivateKey: %v",err)
	}
	BTCAddressSegwit, err := BTCSegwitW.GenerateAddressFromPrivateKey(BTCSegwitPrivateKey)
	if err != nil {
		t.Errorf("Failed to Import BTCSegwitPrivateKey: %v",err)
	}
	BTCAddressTestnet, err := BTCTestnetW.GenerateAddressFromPrivateKey(BTC_TestnetPrivateKey)
	if err != nil {
		t.Errorf("Failed to Import BTC_TestnetPrivateKey: %v",err)
	}
	BTCAddressTestnetSegwit, err := BTCTestnetSegwitW.GenerateAddressFromPrivateKey(BTC_TestnetSegwitPrivateKey)
	if err != nil {
		t.Errorf("Failed to Import BTC_TestnetSegwitPrivateKey: %v",err)
	}
//	ETHAddress, err := GenerateAddressFromPrivateKey(privateKey2,net2)
	WICCAddress, err := WICCW.GenerateAddressFromPrivateKey(WICCPrivateKey)
	if err != nil {
		t.Errorf("Failed to Import WICCPrivateKey: %v",err)
	}
	WICCAddressTestnet, err := WICCTestnetW.GenerateAddressFromPrivateKey(WICC_TestnetPrivateKey)
	if err != nil {
		t.Errorf("Failed to Import WICC_TestnetPrivateKey: %v",err)
	}

	//1AG89FCfPQvtVa7DSCUJjEagH13uTs28Zs
	t.Log("TestGenerateAddressFromPrivateKey: BTCAddress=",BTCAddress)
	//3Ku5bUrN1gXM4fH8WVyRHbKkyGyydqJZ6F
	t.Log("TestGenerateAddressFromPrivateKey: BTCAddressSegwit=",BTCAddressSegwit)
	//mxeBxFWLFAY3G1RKijr91B3kzsX2mTvnYv
	t.Log("TestGenerateAddressFromPrivateKey: BTCAddressTestnet=",BTCAddressTestnet)
	//2MsNRcbHbMgwbbkfzx86Z4FdHkRp29NPjmD
	t.Log("TestGenerateAddressFromPrivateKey: BTCAddressTestnetSegwit=",BTCAddressTestnetSegwit)
	//WhHGDnhL3ny5VES9CFuA38WqDhAQ4VNGuo
	t.Log("TestGenerateAddressFromPrivateKey: WICCAddress=",WICCAddress)
	//wLYLCxsBDjbRiPVEzvbX2bgFftqnWuQxB7
	t.Log("TestGenerateAddressFromPrivateKey: WICCAddressTestnet=",WICCAddressTestnet)
}

//Test Export PrivateKey by Mnemonics
func TestExportPrivateFromMnemonic(t *testing.T){
	mnemonic := "wreck bullet carpet nerve belt border often trust exchange believe defense rebel"

	BTCPrivateKey ,err := BTCW.ExportPrivateKeyFromMnemonic(mnemonic,hdwallet.English)

	BTCPrivateKeySegwit ,err := BTCSegwitW.ExportPrivateKeyFromMnemonic(mnemonic,hdwallet.English)

	BTC_TESTNETPrivateKey ,err := BTCTestnetW.ExportPrivateKeyFromMnemonic(mnemonic,hdwallet.English)

	BTC_TESTNETPrivateKeySegwit ,err := BTCTestnetSegwitW.ExportPrivateKeyFromMnemonic(mnemonic,hdwallet.English)

	ETHPrivateKey ,err := ETHW.ExportPrivateKeyFromMnemonic(mnemonic,hdwallet.English)

	ETHPrivateKeyLedger ,err := ETHLedgerW.ExportPrivateKeyFromMnemonic(mnemonic,hdwallet.English)

	WICCPrivateKey ,err := WICCW.ExportPrivateKeyFromMnemonic(mnemonic,hdwallet.English)

	WICC_TESTNETPrivateKey ,err := WICCTestnetW.ExportPrivateKeyFromMnemonic(mnemonic,hdwallet.English)
	if err != nil {
		t.Errorf("Failed to TestExportPrivateFromMnemonic: %v",err)
	}

	//KyUg2abSHhZYP7bZFXKNDw6TnQoHLyJwDbbDvaNfsBsFxMbFCz4g
	t.Log("BTCPrivateKey=",BTCPrivateKey)
	//KwLJcTWXgB6a14VaxPoJsZPe2o9GPdp2PLcEnrZnmPeRrAtmckgb
	t.Log("BTCPrivateKeySegwit=",BTCPrivateKeySegwit)
	//cUQyhR3BbeFMwtqFRrMppFtPPcx6DMQNFEXm8C1yNSzMkEoRoGYM
	t.Log("BTC_TESTNETPrivateKey=",BTC_TESTNETPrivateKey)
	//cT7214EqFAbtpfuMfg36EDHBmbrkErXLb27ERSKTLX1RUr92LFSE
	t.Log("BTC_TESTNETPrivateKeySegwit=",BTC_TESTNETPrivateKeySegwit)
	//494f8228ae5b6fda6bee1f44eb2c4ed120f210e06acaa8053763efb65638b315
	t.Log("ETHPrivateKey=",ETHPrivateKey)
	//0b98e389e449fa5f388f94bf702066e9ad373e19c2119076f0c276cdd50d776a
	t.Log("ETHPrivateKeyLedger=",ETHPrivateKeyLedger)
	//PemqPzcsCJXjU4PovGSC9zBv89YZSisrWePF9N1skxdEbftbdkDo
	t.Log("WICCPrivateKey=",WICCPrivateKey)
	//Y6amwxjHqUM37UrquokPsbCXTRNughoM27gDUGfbXhJikS39i9h1
	t.Log("WICC_TESTNETPrivateKey=",WICC_TESTNETPrivateKey)
}

func TestCheckAddress(t *testing.T){

	BTCAddress := "1AG89FCfPQvtVa7DSCUJjEagH13uTs28Zs"
	isValid,err := CheckAddress(BTCAddress,&BTCParams)
	if err != nil {
		t.Errorf("Failed to Check BTCAddress: %v",err)
	}
	t.Log("TestCheckAddress: BTCAddress=",isValid)

	BTCAddressSegwit := "3Ku5bUrN1gXM4fH8WVyRHbKkyGyydqJZ6F"
	isValid,err = CheckAddress(BTCAddressSegwit,&BTCParams)
	if err != nil {
		t.Errorf("Failed to Check BTCAddressSegwit: %v",err)
	}
	t.Log("TestCheckAddress: BTCAddressSegwit=",isValid)

	BTCAddressTestnet := "mxeBxFWLFAY3G1RKijr91B3kzsX2mTvnYv"
	isValid,err = CheckAddress(BTCAddressTestnet,&BTCTestnetParams)
	if err != nil {
		t.Errorf("Failed to Check BTCAddressTestnet: %v",err)
	}
	t.Log("TestCheckAddress: BTCAddressTestnet=",isValid)

	BTCAddressTestnetSegwit := "2MsNRcbHbMgwbbkfzx86Z4FdHkRp29NPjmD"
	isValid,err = CheckAddress(BTCAddressTestnetSegwit,&BTCTestnetParams)
	if err != nil {
		t.Errorf("Failed to Check BTCAddressTestnetSegwit: %v",err)
	}
	t.Log("TestCheckAddress: BTCAddressTestnetSegwit=",isValid)

	WICCAddress := "WhHGDnhL3ny5VES9CFuA38WqDhAQ4VNGuo"
	isValid,err = CheckAddress(WICCAddress,&WICCParams)
	if err != nil {
		t.Errorf("Failed to Check WICCAddress: %v",err)
	}
	t.Log("TestCheckAddress: WICCAddress=",isValid)

	WICCAddressTestnet := "wLYLCxsBDjbRiPVEzvbX2bgFftqnWuQxB7"
	isValid,err = CheckAddress(WICCAddressTestnet,&WICCTestnetParams)
	if err != nil {
		t.Errorf("Failed to Check WICCAddressTestnet: %v",err)
	}
	t.Log("TestCheckAddress: WICCAddressTestnet=",isValid)

}

func TestCheckPrivateKey(t *testing.T){

	BTCPrivateKey := "KyUg2abSHhZYP7bZFXKNDw6TnQoHLyJwDbbDvaNfsBsFxMbFCz4g"
	isValid,err := CheckPrivateKey(BTCPrivateKey,&BTCParams)
	if err != nil {
		t.Errorf("Failed to Check BTCPrivateKey: %v",err)
	}
	t.Log("TestCheckAddress: BTCPrivateKey=",isValid)

	BTCSegwitPrivateKey := "KwLJcTWXgB6a14VaxPoJsZPe2o9GPdp2PLcEnrZnmPeRrAtmckgb"
	isValid,err = CheckPrivateKey(BTCSegwitPrivateKey,&BTCParams)
	if err != nil {
		t.Errorf("Failed to Check BTCSegwitPrivateKey: %v",err)
	}
	t.Log("TestCheckAddress: BTCSegwitPrivateKey=",isValid)

	BTC_TestnetPrivateKey := "cUQyhR3BbeFMwtqFRrMppFtPPcx6DMQNFEXm8C1yNSzMkEoRoGYM"
	isValid,err = CheckPrivateKey(BTC_TestnetPrivateKey,&BTCTestnetParams)
	if err != nil {
		t.Errorf("Failed to Check BTC_TestnetPrivateKey: %v",err)
	}
	t.Log("TestCheckAddress: BTC_TestnetPrivateKey=",isValid)

	BTC_TestnetSegwitPrivateKey := "cT7214EqFAbtpfuMfg36EDHBmbrkErXLb27ERSKTLX1RUr92LFSE"
	isValid,err = CheckPrivateKey(BTC_TestnetSegwitPrivateKey,&BTCTestnetParams)
	if err != nil {
		t.Errorf("Failed to Check BTC_TestnetSegwitPrivateKey: %v",err)
	}
	t.Log("TestCheckAddress: BTC_TestnetSegwitPrivateKey=",isValid)

	WICCPrivateKey := "PemqPzcsCJXjU4PovGSC9zBv89YZSisrWePF9N1skxdEbftbdkDo"
	isValid,err = CheckPrivateKey(WICCPrivateKey,&WICCParams)
	if err != nil {
		t.Errorf("Failed to Check WICCPrivateKey: %v",err)
	}
	t.Log("TestCheckAddress: BTCSegwitPrivateKey=",isValid)

	WICC_TestnetPrivateKey := "Y6amwxjHqUM37UrquokPsbCXTRNughoM27gDUGfbXhJikS39i9h1"
	isValid,err = CheckPrivateKey(WICC_TestnetPrivateKey,&WICCTestnetParams)
	if err != nil {
		t.Errorf("Failed to Check WICC_TestnetPrivateKey: %v",err)
	}
	t.Log("TestCheckAddress: WICC_TestnetPrivateKey=",isValid)
}


func TestM(t *testing.T){
	str := "0xXBod"
	content := RemoveOxFromHex(str)
	fmt.Println(content)
}