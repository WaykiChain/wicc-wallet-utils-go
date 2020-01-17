package wicc_wallet_utils_go

import (
	"github.com/WaykiChain/wicc-wallet-utils-go/bitcoin"
	"github.com/WaykiChain/wicc-wallet-utils-go/waykichain"
)

type Mnemonic interface{
	GenerateMnemonic(length int) (string, error)
	ChangeMnemonicLanguage(englishMnemonic string, language string) (string, error)
}

type Wallet interface {
	//wallet
	GenerateAddressFromMnemonic(mnemonic,language string) 	(string, error)
	GenerateAddressFromPrivateKey(privateKey string)	(string, error)
	ExportPrivateKeyFromMnemonic(mnemonic,language string) 	(string, error)
	CheckAddress(address string) (bool, error)
	CheckPrivateKey(privateKey string) 	(bool, error)
	GetPubKeyFromPrivateKey(privateKey string) (string, error)
}

type BTCTransaction interface{
	//BTC transfer
	CreatetRawTxRelyChain(ins []bitcoin.FromInfo, outs []bitcoin.VOut) (string,error)
	CreateTransferRawTx( txins []bitcoin.FinalTxIn,  txouts []bitcoin.VOut) (string, error)
}

type ETHTransaction interface{
	//ETH、ERC20 transfer
	CreateRawTx(privateKeyStr string,chainId int64) (string, error)
}

type WICCTransaction interface {
	//All Tx type
	CreateRawTx(privateKey string) (string, string, error)
}

type WICCSignMessage interface {
	SignMessage() (waykichain.SignMsgResult, error)
}

type WICCVerifyMsgSignature interface {
	VerifyMsgSignature() (isRight bool, addr string)
}