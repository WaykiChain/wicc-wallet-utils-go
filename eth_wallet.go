package wicc_wallet_utils_go

import (
	"encoding/hex"
	wicc_common "github.com/WaykiChain/wicc-wallet-utils-go/common"
	"github.com/WaykiChain/wicc-wallet-utils-go/common/ec"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var(
	ETH 		 ="ETH"
	ETHLedger 	 ="ETHLedger"

	ETHW 			  *ETHWallet
	ETHLedgerW  	  *ETHWallet
)

//init ETH multi type address
const (
	keystoreDir = "/Users/wujinquan/workspace/eth/"
)

func init(){
	ETHW = NewETHWallet(ETH)
	ETHLedgerW = NewETHWallet(ETHLedger)
}

// KeyStore manages a key storage directory on disk.
type KeyStore struct{ k *keystore.KeyStore }
// NewKeyStore creates a keystore for the given directory.
func NewKeyStore(keydir string, scryptN, scryptP int) *KeyStore {
	return &KeyStore{k: keystore.NewKeyStore(keydir, scryptN, scryptP)}
}

//ETHWallet
type ETHWallet struct{
	wallet *wicc_common.Wallet
	keystore *KeyStore
	mnemonicLen int
}

func NewETHWallet(wc string) *ETHWallet{

	newWallet := ETHWallet{}
	switch wc {
	case ETH:
		newWallet =  ETHWallet{
			wallet : wicc_common.NewWallet(wicc_common.ETH,false, false,nil),
			keystore : NewKeyStore(keystoreDir, keystore.StandardScryptN, keystore.StandardScryptP),
			mnemonicLen : 12,
		}
	case ETHLedger:
		newWallet =  ETHWallet{
			wallet : wicc_common.NewWallet(wicc_common.ETH,false, false,nil),
			keystore : NewKeyStore(keystoreDir, keystore.StandardScryptN, keystore.StandardScryptP),
			mnemonicLen : 12,
		}
	default:
		newWallet =  ETHWallet{
			wallet : wicc_common.NewWallet(wicc_common.ETH,false, false,nil),
			keystore : NewKeyStore(keystoreDir, keystore.StandardScryptN, keystore.StandardScryptP),
			mnemonicLen : 12,
		}
	}
	return  &newWallet
}

func (ETHw *ETHWallet) GenerateAddressFromMnemonic(mnemonic,language string) (string, error) {
	return ETHw.wallet.GenerateAddressFromMnemonic(mnemonic,language)
}

func (ETHw *ETHWallet) GenerateAddressFromPrivateKey(privateKey string) (string, error){

	privateKeyBytes,_ := hex.DecodeString(privateKey)
	// (1) new  *PrivateKey、 *PublicKey
	_, publicKey:= ec.PrivKeyFromBytes(privateKeyBytes)
	if  publicKey == nil{
		return "", nil
	}
	/*publicKey Compressed 33 bytes
	publicKeyStr := hex.EncodeToString(publicKey.SerializeCompressed())*/

	// (2) pubBytes为04 开头的65字节公钥,去掉04后剩下64字节进行Keccak256运算
	pubBytes := crypto.Keccak256(publicKey.SerializeUnCompressed()[1:])
	// (3) 经过Keccak256运算后变成32字节，最终取这32字节的后20字节作为真正的地址
	address := common.BytesToAddress(pubBytes[12:])

	return address.Hex(), nil
}

func (ETHw *ETHWallet) ExportPrivateKeyFromMnemonic(mnemonic ,language string ) (string,error) {
	return ETHw.wallet.ExportPrivateKeyFromMnemonic(mnemonic ,language)
}

func (ETHw *ETHWallet) CheckAddress(address string) (bool,error){
	//去掉0x（如有）
	rm0xaddr := wicc_common.RemoveOxFromHex(address)
	//判断长度
	if len(rm0xaddr) != wicc_common.ETHADDRESSLEN{
		return false, wicc_common.ERR_ADDRESS_LEN
	}
	//判断stringTohex是否成功
	_,err := hex.DecodeString(rm0xaddr)
	if err != nil{
		return false, wicc_common.ERR_INVALID_ADDRESS
	}
	return true, nil
}

func (ETHw *ETHWallet) CheckPrivateKey(privateKey string) (bool,error){
	//去掉0x（如有）
	rm0xaddr := wicc_common.RemoveOxFromHex(privateKey)
	//判断长度
	if len(rm0xaddr) != wicc_common.ETHPRIVATEKEYLEN{
		return false, wicc_common.ERR_INVALID_PRIVATEKEY_LEN
	}
	//判断stringTohex是否成功
	_,err := hex.DecodeString(rm0xaddr)
	if err != nil{
		return false, wicc_common.ERR_INVALID_PRIVATEKEY
	}
	return true, nil
}

func (ETHw *ETHWallet) GetPubKeyFromPrivateKey(privateKey string) (string, error) {

	isValid, err := ETHw.CheckPrivateKey(privateKey)
	if isValid == false || err != nil{
		return "", err
	}

	privateKeyBytes, _ := hex.DecodeString(privateKey)
	// (1) new  *PrivateKey、 *PublicKey
	_, publicKey := ec.PrivKeyFromBytes(privateKeyBytes)
	if publicKey == nil {
		return "", nil
	}
	//publicKey Compressed 33 bytes
	return hex.EncodeToString(publicKey.SerializeCompressed()), nil
}


/**********************Keystore********************/

//Import Wallet By Keystore  //wjq 增加去重判断:keystore.Find()
func (ETHw *ETHWallet) ImportKeystore(password string, keyjson string) (string, error){

	// Import back the account we've exported (and then deleted) above with yet
	// again a fresh passphrase
	newAccount, err := ETHw.keystore.k.Import(common.CopyBytes([]byte(keyjson)), password, password)
	if err != nil {
		return  "",err
	}
	return newAccount.Address.Hex(),nil
}

//Import PrivateKey, Save as keystore
func (ETHw *ETHWallet) GenerateAddressFromPrivateKeySaveAsKeystore(password, privateKey string ) (string,error){

	privKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return  "",err
	}

	newAccount,err := ETHw.keystore.k.ImportECDSA(privKey,password)
	if err != nil {
		return "",err
	}
	return newAccount.Address.Hex(), nil
}

//Export keystore by Mnemonic
func (ETHw *ETHWallet) ExportKeystoreFromMnemonic(mnemonic ,language, password string,isledger bool) (string ,error){

	privateKey ,err := ETHw.ExportPrivateKeyFromMnemonic(mnemonic,wicc_common.English)
	if err != nil {
		return "", err
	}

	//import privatekey save as keystore
	address,err := ETHw.GenerateAddressFromPrivateKeySaveAsKeystore(password,privateKey)
	if err != nil {
		return "", err
	}

	//get keystore
	keyjson , err := ETHw.ExportKeystore(address,password)
	if err != nil {
		return "", err
	}

	//delete keystore
	err = ETHw.DeleteKeystoreByAddress(address,password)
	if err != nil {
		return "", err
	}

	return keyjson,nil
}

//Export keystore by keystore
func (ETHw *ETHWallet) ExportKeystore(address, password string )(string,error) {

	account, err := utils.MakeAddress(ETHw.keystore.k, address)
	if err != nil {
		return "", err
	}

	// Export the newly created account with a different passphrase. The returned
	// data from this method invocation is a JSON encoded, encrypted key-file
	keyJson,err := ETHw.keystore.k.Export(account,password,password)
	if err != nil {
		return "", err
	}

	return string(keyJson),nil
}


//Export PrivateKey by keystore
func (ETHw *ETHWallet) ExportPrivateKeyByKeystore(address, password string )(string,error) {

	keyJson , err := ETHw.ExportKeystore(address,password)
	if err != nil {
		return "", err
	}

	key, err := keystore.DecryptKey([]byte(keyJson), password)
	//seckey := math.PaddedBigBytes(key.PrivateKey.D, key.PrivateKey.Params().BitSize/8)
	//fmt.Println("seckey="+ hex.EncodeToString(seckey))
	if err != nil {
		return "", err
	}
	privateKey := crypto.FromECDSA(key.PrivateKey)

	return hex.EncodeToString(privateKey), nil
}


//Delete Keystore By Address
func (ETHw *ETHWallet) DeleteKeystoreByAddress(address, password string ) error {

	account, err := utils.MakeAddress(ETHw.keystore.k, address)
	if err != nil {
		utils.Fatalf("Could not list accounts: %v", err)
	}

	// Delete the account updated above from the local keystore
	if err := ETHw.keystore.k.Delete(account, password); err != nil {
		return err
	}

	return nil
}


//Update Keystore password
func (ETHw *ETHWallet) UpdateKeystorePassword(address, oldpass , newpass string ) error {

	account, err := utils.MakeAddress(ETHw.keystore.k ,address)
	if err != nil {
		return err
	}

	// Update the passphrase on the account created above inside the local keystore
	if err := ETHw.keystore.k.Update(account,oldpass,newpass); err != nil {
		return err
	}

	return nil
}

