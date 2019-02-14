package wiccwallet

import (
	"encoding/hex"

	"wicc-wallet-utils-go/commons"
)

const WAYKI_TESTNET commons.Network = 1
const WAYKI_MAINTNET commons.Network = 2

//Generate Mnemonics string, saprated by space, language is EN(english)
func GenerateMnemonics() string {
	mn := NewMnemonicWithLanguage(ENGLISH)
	words, err := mn.GenerateMnemonic()
	if err != nil {
		return ""
	}
	return words
}

//助记词转换地址
func Mnemonic2Address(words string, netType commons.Network) string {
	address := commons.GenerateAddress(words, netType)
	return address
}

//助记词转私钥
func Mnemonic2PrivateKey(words string, netType commons.Network) string {
	privateKey := commons.GeneratePrivateKey(words, netType)
	return privateKey
}

//私钥转地址
func PrivateKey2Address(words string, netType commons.Network) string {
	address := commons.ImportPrivateKey(words, netType)
	return address
}

//注册账户交易签名
func SignRegisterTx(height int64, fees uint64, privateKey string) string {
	var waykiRegister commons.WaykiRegisterTxParams
	waykiRegister.PrivateKey = privateKey
	waykiRegister.ValidHeight = height
	waykiRegister.Fees = fees
	waykiRegister.TxType = commons.TX_REGISTERACCOUNT
	waykiRegister.Version = 1
	hash := waykiRegister.SignTX()
	return hash
}

//普通交易签名
func SignCommonTx(values uint64, regid string, toAddr string, height int64, fees uint64, privateKey string) string {
	var waykicommon commons.WaykiCommonTxParams
	waykicommon.Values = values
	waykicommon.DestAddress = toAddr
	waykicommon.PrivateKey = privateKey
	waykicommon.UserId = regid
	waykicommon.ValidHeight = height
	waykicommon.Fees = fees
	waykicommon.TxType = commons.TX_COMMON
	waykicommon.Version = 1
	hash := waykicommon.SignTX()
	return hash
}

//投票交易签名
func SignDelegateTx(regid string, height int64, fees uint64, privateKey string, votes []commons.OperVoteFund) string {
	var waykiDelegate commons.WaykiDelegateTxParams
	waykiDelegate.PrivateKey = privateKey
	waykiDelegate.UserId = regid
	waykiDelegate.ValidHeight = height
	waykiDelegate.Fees = fees
	waykiDelegate.TxType = commons.TX_DELEGATE
	waykiDelegate.Version = 1
	waykiDelegate.OperVoteFunds = votes
	hash := waykiDelegate.SignTX()
	return hash
}

//智能合约交易签名
func SignContractTx(values uint64, height int64, fees uint64, privateKey string, regId string, appId string, contractStr string) string {
	var waykiContract commons.WaykiContractTxParams
	waykiContract.Values = values
	waykiContract.PrivateKey = privateKey
	waykiContract.UserId = regId
	waykiContract.AppId = appId
	waykiContract.ValidHeight = height
	waykiContract.Fees = fees
	waykiContract.TxType = commons.TX_CONTRACT
	waykiContract.Version = 1
	binary, _ := hex.DecodeString(contractStr)
	waykiContract.Contract = []byte(binary)
	hash := waykiContract.SignTX()
	return hash
}
