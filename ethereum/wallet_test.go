package ethereum

import (
	"github.com/WaykiChain/wicc-wallet-utils-go/common"
	"testing"
)

func TestGenerateAddressFromMnemonic(t *testing.T){

	mnemonic:= "wreck bullet carpet nerve belt border often trust exchange believe defense rebel"

	ETHAddress , err :=  ETHW.GenerateAddressFromMnemonic(mnemonic,common.English)
	if err != nil{
		t.Errorf("Failed to TestGenerateAddressFromMnemonic: %v", err)
	}
	ETHAddressLedger , err :=  ETHLedgerW.GenerateAddressFromMnemonic(mnemonic,common.English)
	if err != nil{
		t.Errorf("Failed to TestGenerateAddressFromMnemonic: %v", err)
	}

	//0x81FD1F7aE91041aAc5fCF7d8Ed3e1dd88Cc1359a
	t.Log("TestImportWalletFromMnemonic , ETHAddress=",ETHAddress)
	//0x791893c14f0a8dCa4ADB0A8297F8d12063865cd2
	t.Log("TestImportWalletFromMnemonic , ETHAddressLedger=",ETHAddressLedger)
}

func TestGenerateAddressFromPrivateKey(t *testing.T) {

	privkey := "6B93D965D9981F9066CCC44B9DBF32B50F411C0DCEDF4A41CA4E7424ABDB6112"

	//0x81FD1F7aE91041aAc5fCF7d8Ed3e1dd88Cc1359a
	address,err := ETHW.GenerateAddressFromPrivateKey(privkey)
	if err != nil{
		t.Error("GenerateAddressFromPrivateKey err:",err)
	}
	t.Log("address :",address)
}

func TestExportPrivateKeyFromMnemonic(t *testing.T){

	mnemonic:= "wreck bullet carpet nerve belt border often trust exchange believe defense rebel"

	ETHPrivateKey ,err := ETHW.ExportPrivateKeyFromMnemonic(mnemonic,common.English)
	if err != nil {
		t.Errorf("Failed to TestExportPrivateKeyFromMnemonic: %v",err)
	}
	ETHPrivateKeyLedger ,err := ETHLedgerW.ExportPrivateKeyFromMnemonic(mnemonic,common.English)
	if err != nil {
		t.Errorf("Failed to TestExportPrivateKeyFromMnemonic: %v",err)
	}

	//494f8228ae5b6fda6bee1f44eb2c4ed120f210e06acaa8053763efb65638b315
	t.Log("ETHPrivateKey=",ETHPrivateKey)
	//0b98e389e449fa5f388f94bf702066e9ad373e19c2119076f0c276cdd50d776a
	t.Log("ETHPrivateKeyLedger=",ETHPrivateKeyLedger)

}

func TestCheckETHAddress(t *testing.T) {

	ETHAddress := "0x96b4213eD85031b02A1bE101FfA3F82ee929285a"
	isValid,err := ETHW.CheckAddress(ETHAddress)
	if err != nil{
		t.Error("CheckETHAddress err:",err)
	}
	t.Log("CheckETHAddress :",isValid)
}

func TestCheckETHPrivateKey(t *testing.T) {

	ETHPrivateKey := "6B93D965D9981F9066CCC44B9DBF32B50F411C0DCEDF4A41CA4E7424ABDB6112"
	isValid, err := ETHW.CheckPrivateKey(ETHPrivateKey)
	if err != nil {
		t.Errorf("Failed to Check ETHPrivateKey: %v", err)
	}
	t.Log("TestCheckAddress: ETHPrivateKey=", isValid)
}

func TestGetPubKeyFromPrivateKey(t *testing.T) {

	ETHPrivateKey := "0b98e389e449fa5f388f94bf702066e9ad373e19c2119076f0c276cdd50d776a"
	publicKey, err := ETHW.GetPubKeyFromPrivateKey(ETHPrivateKey)
	if err != nil {
		t.Errorf("Failed to GetPubKeyFromPrivateKey BTCPrivateKey: %v", err)
	}
	t.Log("GetPubKeyFromPrivateKey BTCPrivateKey=", publicKey)
}

//TestImportKeystore
func TestImportKeystore(t *testing.T) {

	//ks := NewKeyStore(keystoreDir, keystore.StandardScryptN, keystore.StandardScryptP)

	password := "12345678"
	keyjsonStr := "{\"address\":\"ec9c88fc291ddc0e18dc321d82e29aa5454efb9d\",\"crypto\":{\"cipher\":\"aes-128-ctr\",\"ciphertext\":\"525b496910610bcc48c5488edc2a1daf19cf57f553d747fb214994fd32145096\",\"cipherparams\":{\"iv\":\"5c49b2e7f13f6afc5321a43d91ee55d0\"},\"kdf\":\"scrypt\",\"kdfparams\":{\"dklen\":32,\"n\":262144,\"p\":1,\"r\":8,\"salt\":\"783f5d34ebb98b7a06b39f65f172b5e1e11a0c126a0a8e621dabe5906de2d307\"},\"mac\":\"3fc4239a550687f4b038a5740802371be876e4bb9b68a174850220b6abb140bc\"},\"id\":\"299ff3d4-14cd-49c9-aec3-561fd8ce88a8\",\"version\":3}"

	err , address := ETHW.ImportKeystore(password,keyjsonStr)
	if err != nil {
		t.Errorf("Failed to import account: %v", err)
	}

	t.Log("TestImportKeystore: new account=",address)
}

//Test Import PrivateKey, Save as keystore
func TestGenerateAddressFromPrivateKeySaveAsKeystore(t *testing.T){

	//ks := NewKeyStore(keystoreDir, keystore.StandardScryptN, keystore.StandardScryptP)

	privkey := "a7c0cf50fecdf99c570b987e21c9f63f3f152e5b74885e5f7a16dd3bbebe4d7b"
	password := "87654321"

	err , address := ETHW.GenerateAddressFromPrivateKeySaveAsKeystore(password,privkey)
	if err != nil {
		t.Errorf("Failed to import account: %v",err)
	}

	t.Log("TestImportWalletByPrivateKey: new account=",address)
}

//Test Export keystore from Mnemonic
func TestExportKeystoreFromMnemonic(t *testing.T){

	mnemonic := "wreck bullet carpet nerve belt border often trust exchange believe defense rebel"
//	ks := NewKeyStore(keystoreDir, keystore.StandardScryptN, keystore.StandardScryptP)
	password := "12345678"

	keyJson, err := ETHW.ExportKeystoreFromMnemonic(mnemonic, common.English, password, true)
	if err != nil {
		t.Errorf("Failed to TestExportkeystoreFromMnemonic: %v",err)
	}

	t.Log("TestExportkeystoreFromMnemonic: keyJson=",keyJson)
}

//Test Export Keystore
func TestExportKeystore(t *testing.T){

//	ks := NewKeyStore(keystoreDir, keystore.StandardScryptN, keystore.StandardScryptP)
	address := "ec9c88fc291ddc0e18dc321d82e29aa5454efb9d"
	password := "12345678"
	keyJson , err := ETHW.ExportKeystore(address,password)
	if err != nil {
		t.Errorf("Failed to export keystore: %v",err)
	}

	t.Log("keyJson=",keyJson)
}

//Test Export PrivateKey by Keystore
func TestExportPrivateKeyByKeystore(t *testing.T){

//	ks := NewKeyStore(keystoreDir, keystore.StandardScryptN, keystore.StandardScryptP)
	address := "1c05cb077d2e2d28bfffb73a51cca25af22bc355"
	password := "12345678"
	privateKey , err := ETHW.ExportPrivateKeyByKeystore(address,password)
	if err != nil {
		t.Errorf("Failed to export keystore: %v",err)
	}

	t.Log("privateKey=",privateKey)
}

//Test DeleteAccountByKeystore
func TestDeleteKeystoreByAddress(t *testing.T){

//	ks := NewKeyStore(keystoreDir, keystore.StandardScryptN, keystore.StandardScryptP)
	password := "12345678"
	address := "51c9180cf26dd1e26b481cac0f555677bfe51d95"

	if err := ETHW.DeleteKeystoreByAddress(address,password); err != nil {
		t.Errorf("Failed DeleteAccountByKeystore: %v",err)
	}

	t.Log("TestDeleteAccountByKeystore successful")
}

//TestUpdateKeystorePassword
func TestUpdateKeystorePassword(t *testing.T){

//	ks := NewKeyStore(keystoreDir, keystore.StandardScryptN, keystore.StandardScryptP)
	address := "ec9c88fc291ddc0e18dc321d82e29aa5454efb9d"
	oldpass := "12345678"
	newpass := "12345678"

	if err := ETHW.UpdateKeystorePassword(address,oldpass,newpass); err != nil {
		t.Errorf("Failed to UpdateKeystorePassword: %v",err)
	}

	t.Log("TestUpdateKeystorePassword successful")
}







