package wicc_wallet_utils_go

import (
	"github.com/WaykiChain/wicc-wallet-utils-go/common"
	"github.com/WaykiChain/wicc-wallet-utils-go/waykichain"
	"testing"
)

var (
	WWM *WICCWalletManager
)

func init(){

	WWM = NewWICCWalletManager()
	WWM.Config = waykichain.NewConfig()
	WWM.Config.NodeServerAPI = "http://10.0.0.20:6968" //testnet
//	tw.Config.NodeServerAPI = "http://10.0.0.21:6968" //mainnet
	WWM.Config.BaaSServerAPI = "https://baas-test.wiccdev.org/v2/" //testnet
//	tw.Config.BaaSServerAPI = "https://baas.wiccdev.org/v2/" //mainnet
	WWM.Config.RpcUser = "wayki"
	WWM.Config.RpcPassword = "admin@123"
	WWM.Config.WalletConfig = waykichain.WICCTestnetConfig
	WWM.Config.Debug = true
	WWM.Config.ServerType = waykichain.BaaS
	WWM.Wallet = NewWICCWallet(WWM.Config.WalletConfig)
	token := waykichain.BasicAuth(WWM.Config.RpcUser, WWM.Config.RpcPassword)
	WWM.WalletClient = waykichain.NewClient(WWM.Config.NodeServerAPI, token, WWM.Config.Debug)
	WWM.BaaSClient = waykichain.NewBaaSClient(WWM.Config.BaaSServerAPI, WWM.Config.Debug)
}

/*
1、从钱包获取私钥
2、创建rawtx：从外部获取签名时需要的有效高度：ValidHeight
3、广播交易
*/
func TestSendUCoinTransferTx(t *testing.T){

	//钱包
	mnemonic := "wreck bullet carpet nerve belt border often trust exchange believe defense rebel"
	privateKey,err := WICCTestnetW.ExportPrivateKeyFromMnemonic(mnemonic,common.English)
	if err != nil {
		t.Error("ExportPrivateKeyFromMnemonic:",err)
	}

	//创建rawtx
	var txParam UCoinTransferTxParam
	txParam.FeeSymbol = string(common.WICC)
	txParam.Fees = 1000000
	txParam.ValidHeight,_ = WWM.GetSynBlockHeight()
	txParam.SrcRegId = ""
	txParam.Dests=NewDestArr()
	dest:=Dest{string(common.WICC),1000000, "wLYLCxsBDjbRiPVEzvbX2bgFftqnWuQxB7"}
	txParam.Dests.Add(&dest)
	txParam.PubKey = "031b27286c65b81ac13cfd4067b030398a19eb147f439c094fbb19a2f3ab9ec10b"
	txParam.Memo = ""
	result, err := txParam.CreateRawTx(privateKey)
	if err != nil {
		t.Error("CreateUCoinTransferRawTx err: ", err)
	}
	t.Log("rawTx=",result.RawTx)
	t.Log("txid=",result.Txid)

	//广播交易
	submitxid,err := WWM.SubmitTxRaw(result.RawTx)
	if err != nil {
		t.Error("SubmitTxRaw err: ", err)
	}
	t.Log("submitxid=",submitxid)
}