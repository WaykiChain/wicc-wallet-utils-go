package wicc_wallet_utils_go

import (
	"github.com/WaykiChain/wicc-wallet-utils-go/common"
	"testing"
)


func TestBTCGenerateAddressFromMnemonic(t *testing.T){

	mnemonic:= "wreck bullet carpet nerve belt border often trust exchange believe defense rebel"

	BTCAddress , err := BTCW.GenerateAddressFromMnemonic(mnemonic,common.English)
	BTCAddressSegwit , err := BTCSegwitW.GenerateAddressFromMnemonic(mnemonic,common.English)
	BTCAddressTestnet , err :=  BTCTestnetW.GenerateAddressFromMnemonic(mnemonic,common.English)
	BTCAddressTestnetSegwit , err :=  BTCTestnetSegwitW.GenerateAddressFromMnemonic(mnemonic,common.English)

	if err != nil{
		t.Errorf("Failed to TestBTCGenerateAddressFromMnemonic: %v", err)
	}

	//1AG89FCfPQvtVa7DSCUJjEagH13uTs28Zs
	t.Log("TestImportWalletFromMnemonic , BTCAddress=",BTCAddress)
	//3Ku5bUrN1gXM4fH8WVyRHbKkyGyydqJZ6F
	t.Log("TestImportWalletFromMnemonic , BTCAddressSegwit=",BTCAddressSegwit)
	//mxeBxFWLFAY3G1RKijr91B3kzsX2mTvnYv
	t.Log("TestImportWalletFromMnemonic , BTCAddressTestnet=",BTCAddressTestnet)
	//2MsNRcbHbMgwbbkfzx86Z4FdHkRp29NPjmD
	t.Log("TestImportWalletFromMnemonic , BTCAddressTestnetSegwit=",BTCAddressTestnetSegwit)
}

func TestBTCGenerateAddressFromPrivateKey(t *testing.T){

	BTCPrivateKey := "KyUg2abSHhZYP7bZFXKNDw6TnQoHLyJwDbbDvaNfsBsFxMbFCz4g"
	BTCSegwitPrivateKey := "KwLJcTWXgB6a14VaxPoJsZPe2o9GPdp2PLcEnrZnmPeRrAtmckgb"
	BTC_TestnetPrivateKey := "cUQyhR3BbeFMwtqFRrMppFtPPcx6DMQNFEXm8C1yNSzMkEoRoGYM"
	BTC_TestnetSegwitPrivateKey := "cT7214EqFAbtpfuMfg36EDHBmbrkErXLb27ERSKTLX1RUr92LFSE"

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

	//1AG89FCfPQvtVa7DSCUJjEagH13uTs28Zs
	t.Log("TestGenerateAddressFromPrivateKey: BTCAddress=",BTCAddress)
	//3Ku5bUrN1gXM4fH8WVyRHbKkyGyydqJZ6F
	t.Log("TestGenerateAddressFromPrivateKey: BTCAddressSegwit=",BTCAddressSegwit)
	//mxeBxFWLFAY3G1RKijr91B3kzsX2mTvnYv
	t.Log("TestGenerateAddressFromPrivateKey: BTCAddressTestnet=",BTCAddressTestnet)
	//2MsNRcbHbMgwbbkfzx86Z4FdHkRp29NPjmD
	t.Log("TestGenerateAddressFromPrivateKey: BTCAddressTestnetSegwit=",BTCAddressTestnetSegwit)
}

//Test Export PrivateKey by Mnemonics
func TestBTCExportPrivateFromMnemonic(t *testing.T){
	mnemonic := "wreck bullet carpet nerve belt border often trust exchange believe defense rebel"

	BTCPrivateKey ,err := BTCW.ExportPrivateKeyFromMnemonic(mnemonic,common.English)
	if err != nil {
		t.Errorf("Failed to TestExportPrivateFromMnemonic: %v",err)
	}
	BTCPrivateKeySegwit ,err := BTCSegwitW.ExportPrivateKeyFromMnemonic(mnemonic,common.English)
	if err != nil {
		t.Errorf("Failed to TestExportPrivateFromMnemonic: %v",err)
	}
	BTC_TESTNETPrivateKey ,err := BTCTestnetW.ExportPrivateKeyFromMnemonic(mnemonic,common.English)
	if err != nil {
		t.Errorf("Failed to TestExportPrivateFromMnemonic: %v",err)
	}
	BTC_TESTNETPrivateKeySegwit ,err := BTCTestnetSegwitW.ExportPrivateKeyFromMnemonic(mnemonic,common.English)
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
}

func TestCheckBTCAddress(t *testing.T) {

	BTCAddress := "1AG89FCfPQvtVa7DSCUJjEagH13uTs28Zs"
	isValid, err := BTCW.CheckAddress(BTCAddress)
	if err != nil {
		t.Errorf("Failed to Check BTCAddress: %v", err)
	}
	t.Log("TestCheckAddress: BTCAddress=", isValid)

	BTCAddressSegwit := "3Ku5bUrN1gXM4fH8WVyRHbKkyGyydqJZ6F"
	isValid, err = BTCSegwitW.CheckAddress(BTCAddressSegwit)
	if err != nil {
		t.Errorf("Failed to Check BTCAddressSegwit: %v", err)
	}
	t.Log("TestCheckAddress: BTCAddressSegwit=", isValid)

	BTCAddressTestnet := "mxeBxFWLFAY3G1RKijr91B3kzsX2mTvnYv"
	isValid, err = BTCTestnetW.CheckAddress(BTCAddressTestnet)
	if err != nil {
		t.Errorf("Failed to Check BTCAddressTestnet: %v", err)
	}
	t.Log("TestCheckAddress: BTCAddressTestnet=", isValid)

	BTCAddressTestnetSegwit := "2MsNRcbHbMgwbbkfzx86Z4FdHkRp29NPjmD"
	isValid, err = BTCTestnetSegwitW.CheckAddress(BTCAddressTestnetSegwit)
	if err != nil {
		t.Errorf("Failed to Check BTCAddressTestnetSegwit: %v", err)
	}
	t.Log("TestCheckAddress: BTCAddressTestnetSegwit=", isValid)

}

func TestCheckBTCPrivateKey(t *testing.T) {

	BTCPrivateKey := "KyUg2abSHhZYP7bZFXKNDw6TnQoHLyJwDbbDvaNfsBsFxMbFCz4g"
	isValid, err := BTCW.CheckPrivateKey(BTCPrivateKey)
	if err != nil {
		t.Errorf("Failed to Check BTCPrivateKey: %v", err)
	}
	t.Log("TestCheckAddress: BTCPrivateKey=", isValid)

	BTCSegwitPrivateKey := "KwLJcTWXgB6a14VaxPoJsZPe2o9GPdp2PLcEnrZnmPeRrAtmckgb"
	isValid, err = BTCSegwitW.CheckPrivateKey(BTCSegwitPrivateKey)
	if err != nil {
		t.Errorf("Failed to Check BTCSegwitPrivateKey: %v", err)
	}
	t.Log("TestCheckAddress: BTCSegwitPrivateKey=", isValid)

	BTC_TestnetPrivateKey := "cUQyhR3BbeFMwtqFRrMppFtPPcx6DMQNFEXm8C1yNSzMkEoRoGYM"
	isValid, err = BTCTestnetW.CheckPrivateKey(BTC_TestnetPrivateKey)
	if err != nil {
		t.Errorf("Failed to Check BTC_TestnetPrivateKey: %v", err)
	}
	t.Log("TestCheckAddress: BTC_TestnetPrivateKey=", isValid)

	BTC_TestnetSegwitPrivateKey := "cT7214EqFAbtpfuMfg36EDHBmbrkErXLb27ERSKTLX1RUr92LFSE"
	isValid, err = BTCTestnetSegwitW.CheckPrivateKey(BTC_TestnetSegwitPrivateKey)
	if err != nil {
		t.Errorf("Failed to Check BTC_TestnetSegwitPrivateKey: %v", err)
	}
	t.Log("TestCheckAddress: BTC_TestnetSegwitPrivateKey=", isValid)
}

func TestBTCGetPubKeyFromPrivateKey(t *testing.T) {

	BTCPrivateKey := "KyUg2abSHhZYP7bZFXKNDw6TnQoHLyJwDbbDvaNfsBsFxMbFCz4g"
	publicKey, err := BTCW.GetPubKeyFromPrivateKey(BTCPrivateKey)
	if err != nil {
		t.Errorf("Failed to GetPubKeyFromPrivateKey BTCPrivateKey: %v", err)
	}
	t.Log("GetPubKeyFromPrivateKey BTCPrivateKey=", publicKey)

	BTCSegwitPrivateKey := "KwLJcTWXgB6a14VaxPoJsZPe2o9GPdp2PLcEnrZnmPeRrAtmckgb"
	publicKey, err = BTCSegwitW.GetPubKeyFromPrivateKey(BTCSegwitPrivateKey)
	if err != nil {
		t.Errorf("Failed to GetPubKeyFromPrivateKey : %v", err)
	}
	t.Log("GetPubKeyFromPrivateKey BTCSegwitPrivateKey=", publicKey)

	BTC_TestnetPrivateKey := "cUQyhR3BbeFMwtqFRrMppFtPPcx6DMQNFEXm8C1yNSzMkEoRoGYM"
	publicKey, err = BTCTestnetW.GetPubKeyFromPrivateKey(BTC_TestnetPrivateKey)
	if err != nil {
		t.Errorf("Failed to GetPubKeyFromPrivateKey BTC_TestnetPrivateKey: %v", err)
	}
	t.Log("GetPubKeyFromPrivateKey BTC_TestnetPrivateKey=", publicKey)

	BTC_TestnetSegwitPrivateKey := "cT7214EqFAbtpfuMfg36EDHBmbrkErXLb27ERSKTLX1RUr92LFSE"
	publicKey, err = BTCTestnetSegwitW.GetPubKeyFromPrivateKey(BTC_TestnetSegwitPrivateKey)
	if err != nil {
		t.Errorf("Failed to BTC_TestnetSegwitPrivateKey: %v", err)
	}
	t.Log("GetPubKeyFromPrivateKey BTC_TestnetSegwitPrivateKey=", publicKey)
}

