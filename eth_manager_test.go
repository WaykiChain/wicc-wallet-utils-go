package wicc_wallet_utils_go

import (
	"github.com/WaykiChain/wicc-wallet-utils-go/ethereum"
	"math/big"
	"testing"
)

var (
	tw 		*ETHWalletManager
)

func init(){

	tw = NewETHWalletManager()
	tw.Config = ethereum.NewConfig()
	tw.Config.ServerAPI = "https://ropsten.infura.io/v3/5c93a142071540709b5d953478797194"
	tw.Config.ChainID = 3
	tw.Config.RpcPassword = ""
	tw.Config.RpcPassword = ""
	tw.Config.WalletConfig = ETH
	tw.Config.KeystoreDir = keystoreDir
	tw.Wallet = NewETHWallet(tw.Config.WalletConfig)
	//token := BasicAuth(tw.Config.RpcUser, tw.Config.RpcPassword)
	tw.WalletClient = NewClient(tw.Config.ServerAPI, "", true)

}

func TestSendETHTransaction(t *testing.T){

	privateKeyStr := "6B93D965D9981F9066CCC44B9DBF32B50F411C0DCEDF4A41CA4E7424ABDB617F"
	from := "0x232D23C22543144B988F738C701Df6dfd6eAcA4c"  //WJQ GenerateAddressFromPrivateKey
	to := "81FD1F7aE91041aAc5fCF7d8Ed3e1dd88Cc1359a"
	amount := big.NewInt(100000000000000000)

	//读取链上from地址的nonce,此做法有风险，可以采用本地记录的方式
	nonce ,err:= tw.WalletClient.ethGetTransactionCount(from,LEATEST)
	if err != nil {
		t.Errorf("Failed to ethGetTransactionCount: %v",err)
	}
	t.Log("nonce=",nonce)

	//从链上估算手续费
	txFeeInfo, err := tw.WalletClient.GetTransactionFeeEstimated(from,to,amount,"")
	if err != nil {
		t.Errorf("Failed to GetTransactionFeeEstimated: %v",err)
	}
	t.Logf("txFeeInfo.GasLimit=%v, txFeeInfo.GasPrice=%v\n",txFeeInfo.GasLimit,txFeeInfo.GasPrice)

	//离线签名
	tx := NewETHSimpleTransaction(to,amount,nonce,txFeeInfo.GasLimit,txFeeInfo.GasPrice)
	rawtx ,err :=tx.CreateRawTx(privateKeyStr,tw.Config.ChainID)
	if err != nil{
		t.Error("CreateRawTx=",err)
	}
	t.Log("rawtx=",rawtx)

	//提交广播交易
	txid ,err  := tw.WalletClient.EthSendRawTransaction(rawtx)
	if err != nil{
		t.Error("EthSendRawTransaction=",err)
	}
	t.Log("txid=",txid)//0xab5a428263314427eb01cffa5ac63ecf4ad4d29601b6c9b4476b5e5c4df840cf
}

