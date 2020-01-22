package wicc_wallet_utils_go

import (
	"encoding/hex"
	"github.com/WaykiChain/wicc-wallet-utils-go/common"
	"github.com/WaykiChain/wicc-wallet-utils-go/common/hash"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
)

var(
	WICCMainnet 		= "WICCMainnet"
	WICCTestnet 		= "WICCTestnet"
	WICCW 			  *WICCWallet
	WICCTestnetW       *WICCWallet
)

type WICCWallet struct{
	wallet *common.Wallet
	mnemonicLen int
}

func init(){
	WICCW = NewWICCWallet(WICCMainnet)
	WICCTestnetW = NewWICCWallet(WICCTestnet)
}

func NewWICCWallet(wc string) *WICCWallet{

	newWallet := WICCWallet{}
	switch wc {
	case WICCMainnet:
		newWallet =  WICCWallet{wallet:common.NewWallet(common.WICC,false,false,&common.WICCParams),mnemonicLen:12}
	case WICCTestnet:
		newWallet =  WICCWallet{wallet:common.NewWallet(common.WICC_TESTNET,false,false,&common.WICCTestnetParams),mnemonicLen:12}
	default:
		newWallet =  WICCWallet{wallet:common.NewWallet(common.WICC,false,false,&common.WICCParams),mnemonicLen:12}
	}
	return  &newWallet
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
func  SignMessage(input *SignMsgInput) (*SignMsgResult, error) {
	//Use sha256_sha160 instead of sha256_twice
	hash := hash.Hash256(hash.Hash160([]byte(input.Data)))

	wifKey, errorDecode := btcutil.DecodeWIF(input.PrivateKey)
	if (errorDecode != nil) {
		return nil, common.ERR_INVALID_PRIVATEKEY
	}
	key := wifKey.PrivKey
	signature, errorSign := key.Sign(hash)
	if (errorSign != nil) {
		return nil, common.ERR_SIGNATURE_ERROR
	}
	signResult := SignMsgResult{hex.EncodeToString(wifKey.SerializePubKey()),
		hex.EncodeToString(signature.Serialize())}
	return &signResult, nil
}

func VerifyMsgSignature(input *VerifySignInput) (*VerifyMsgSignResult,error){

	if  (len(input.PublicKey) != 66) || (len(input.Signature) % 2 != 0) {
		return nil, common.ERR_PUBLICKEY_SIGNATURE_ERROR
	}

	publicKeyBytes, err := hex.DecodeString(input.PublicKey)
	if err != nil {
		return nil,err
	}

	//check publicKey
	pubKey, err := btcec.ParsePubKey(publicKeyBytes, btcec.S256())
	if (err != nil) {
		return nil,err
	}

	//get address from public
	address, err := btcutil.NewAddressPubKey(publicKeyBytes,&input.NetParams)
	if (err != nil) {
		return nil,err
	}

	//get signature hash
	sigBytes, err := hex.DecodeString(input.Signature)
	if err != nil {
		return nil,err
	}
	sig, err := btcec.ParseDERSignature(sigBytes, btcec.S256())
	if err != nil {
		return nil,err
	}

	//Verify
	//Use sha256_sha160 instead of sha256_twice
	isValid := sig.Verify(hash.Hash256(btcutil.Hash160([]byte(input.Data))), pubKey)
	if (isValid){
		return &VerifyMsgSignResult{true,address.EncodeAddress()},nil
	}else{
		return &VerifyMsgSignResult{false,""},nil
	}
}
