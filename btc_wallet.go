package wicc_wallet_utils_go

import (
	"encoding/hex"
	"github.com/WaykiChain/wicc-wallet-utils-go/common"
	"github.com/btcsuite/btcutil"
)

var(
	BTCMainnet 		 ="BTCMainnet"
	BTCSegwitMainnet ="BTCSegwitMainnet"
	BTCTestnet       ="BTCTestnet"
	BTCSegwitTestnet ="BTCSegwitTestnet"

	BTCW 			  *BTCWallet
	BTCSegwitW  	  *BTCWallet
	BTCTestnetW       *BTCWallet
	BTCTestnetSegwitW *BTCWallet
)
func init(){
	BTCW = NewBTCWallet(BTCMainnet)
	BTCSegwitW = NewBTCWallet(BTCSegwitMainnet)
	BTCTestnetW = NewBTCWallet(BTCTestnet)
	BTCTestnetSegwitW = NewBTCWallet(BTCSegwitTestnet)
}

type BTCWallet struct{
	wallet *common.Wallet
	mnemonicLen int
}

func NewBTCWallet(wc string) *BTCWallet{
	newWallet := BTCWallet{}
	switch wc {
	case BTCMainnet:
		newWallet =  BTCWallet{wallet:common.NewWallet(common.BTC,false,false,&common.BTCParams),mnemonicLen:12}
	case BTCSegwitMainnet:
		newWallet =  BTCWallet{wallet:common.NewWallet(common.BTC,true,false,&common.BTCParams),mnemonicLen:12}
	case BTCTestnet:
		newWallet =  BTCWallet{wallet:common.NewWallet(common.BTC_TESTNET,false,false,&common.BTCTestnetParams),mnemonicLen:12}
	case BTCSegwitTestnet:
		newWallet =  BTCWallet{wallet:common.NewWallet(common.BTC_TESTNET,true,false,&common.BTCTestnetParams),mnemonicLen:12}
	default:
		newWallet =  BTCWallet{wallet:common.NewWallet(common.BTC,false,false,&common.BTCParams),mnemonicLen:12}
	}
	return  &newWallet
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