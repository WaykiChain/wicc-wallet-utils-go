package waykichain

import (
	"encoding/hex"
	"fmt"
	"github.com/WaykiChain/wicc-wallet-utils-go/common"
	"github.com/WaykiChain/wicc-wallet-utils-go/common/hash"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
)

var(
	WICCW 			  *WICCWallet
	WICCTestnetW       *WICCWallet
)
func init(){
	WICCW = NewWICCWallet(NewWalletConfig(WICCMainnetConf))
	WICCTestnetW = NewWICCWallet(NewWalletConfig(WICCTestnetConf))
}

type WICCWallet struct{
	wallet *common.Wallet
	mnemonicLen int
}
func NewWICCWallet(wc *WICCWalletConfig) *WICCWallet{
	return &WICCWallet{wallet:common.NewWallet(wc.coinType,false,false,wc.netParam),mnemonicLen:12}
}

func (WICCw *WICCWallet) GenerateAddressFromMnemonic(mnemonic,language string) (string, error) {
	return WICCw.wallet.GenerateAddressFromMnemonic(mnemonic,language)
}

func (WICCw *WICCWallet) GenerateAddressFromPrivateKey(privateKey string) (string,error){
	return WICCw.wallet.GenerateAddressFromPrivateKey(privateKey)
}

func (WICCw *WICCWallet) ExportPrivateKeyFromMnemonic(mnemonic ,language string ) (string,error) {
	return WICCw.wallet.ExportPrivateKeyFromMnemonic(mnemonic ,language)
}

func (WICCw *WICCWallet) CheckAddress(address string) (bool, error) {
	return common.CheckAddress(address, WICCw.wallet.NetParam)
}

func (WICCw *WICCWallet) CheckPrivateKey(address string) (bool, error) {
	return common.CheckPrivateKey(address, WICCw.wallet.NetParam)
}

// get publickey from privatekey
func (WICCw *WICCWallet) GetPubKeyFromPrivateKey(privateKey string) (string, error) {

	isValid, err := WICCw.CheckPrivateKey(privateKey)
	if isValid == false || err != nil{
		return "", err
	}

	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return "", common.ERR_INVALID_PRIVATEKEY
	}
	pubHex := hex.EncodeToString(wifKey.SerializePubKey())
	return pubHex, nil
}


//Sign message by private Key
func (input SignMsgInput) SignMessage() (SignMsgResult, error) {
	//Use sha256_sha160 instead of sha256_twice
	hash := hash.Hash256(hash.Hash160([]byte(input.Data)))

	hashStr := hex.EncodeToString(hash)
	fmt.Println("hashStr=",hashStr)

	wifKey, errorDecode := btcutil.DecodeWIF(input.PrivateKey)
	if (errorDecode != nil) {
		return SignMsgResult{}, common.ERR_INVALID_PRIVATEKEY
	}
	key := wifKey.PrivKey
	signature, errorSign := key.Sign(hash)
	if (errorSign != nil) {
		return SignMsgResult{}, common.ERR_SIGNATURE_ERROR
	}
	signResult := SignMsgResult{hex.EncodeToString(wifKey.SerializePubKey()),
		hex.EncodeToString(signature.Serialize())}
	return signResult, nil
}

func (input VerifySignInput) VerifyMsgSignature() (isRight bool,addr string){

	if  (len(input.PublicKey) != 66) || (len(input.Signature) % 2 != 0) {
		fmt.Println("The length of publicKey or signature error")
		return false,""
	}

	publicKeyBytes, err := hex.DecodeString(input.PublicKey)
	if err != nil {
		fmt.Println("publicKey err:", err)
		return false,""
	}

	//check publicKey
	pubKey, err := btcec.ParsePubKey(publicKeyBytes, btcec.S256())
	if (err != nil) {
		fmt.Println("PublicKey invaild")
		return false,""
	}

	//get address from public
	address, err := btcutil.NewAddressPubKey(publicKeyBytes,&input.NetParams)
	if (err != nil) {
		fmt.Println("Failed to generate address")
		return false,""
	}

	//get signature hash
	sigBytes, err := hex.DecodeString(input.Signature)
	if err != nil {
		fmt.Println("signature err:", err)
		return false,""
	}
	sig, err := btcec.ParseDERSignature(sigBytes, btcec.S256())
	if err != nil {
		fmt.Println("sigBytes err:", err)
		return false,""
	}

	//Verify
	//Use sha256_sha160 instead of sha256_twice
	isValid := sig.Verify(hash.Hash256(btcutil.Hash160([]byte(input.Data))), pubKey)
	if (isValid){
		return true,address.EncodeAddress()
	}else{
		return false,""
	}
}
