package bitcoin

import (
	"encoding/hex"
	"github.com/WaykiChain/wicc-wallet-utils-go/common"
	"github.com/btcsuite/btcutil"
)

var(
	BTCW 			  *BTCWallet
	BTCSegwitW  	  *BTCWallet
	BTCTestnetW       *BTCWallet
	BTCTestnetSegwitW *BTCWallet
)
func init(){
	BTCW = NewBTCWallet(NewWalletConfig(BTCMainnetConf))
	BTCSegwitW = NewBTCWallet(NewWalletConfig(BTCMainnetSegwitConf))
	BTCTestnetW = NewBTCWallet(NewWalletConfig(BTCTestnetConf))
	BTCTestnetSegwitW = NewBTCWallet(NewWalletConfig(BTCTestnetSegwitConf))
}

type BTCWallet struct{
	wallet *common.Wallet
	mnemonicLen int
}

func NewBTCWallet(wc *BTCWalletConfig) *BTCWallet{
	return &BTCWallet{wallet:common.NewWallet(wc.coinType,wc.isSegwit,false,wc.netParam),mnemonicLen:12}
}

func (BTCw *BTCWallet) GenerateAddressFromMnemonic(mnemonic,language string) (string, error) {
	return BTCw.wallet.GenerateAddressFromMnemonic(mnemonic,language)
}

func (BTCw *BTCWallet) GenerateAddressFromPrivateKey(privateKey string) (string,error){
	return BTCw.wallet.GenerateAddressFromPrivateKey(privateKey)
}

func (BTCw *BTCWallet) ExportPrivateKeyFromMnemonic(mnemonic,language string) (string,error) {
	return BTCw.wallet.ExportPrivateKeyFromMnemonic(mnemonic,language)
}

func (BTCw *BTCWallet) CheckAddress(address string) (bool, error) {
	return common.CheckAddress(address, BTCw.wallet.NetParam)
}

func (BTCw *BTCWallet) CheckPrivateKey(privateKey string) (bool, error) {
	return common.CheckPrivateKey(privateKey,BTCw.wallet.NetParam)
}

// get publickey from privatekey
func (BTCw *BTCWallet) GetPubKeyFromPrivateKey(privateKey string) (string, error) {

	isValid, err := BTCw.CheckPrivateKey(privateKey)
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