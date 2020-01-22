package wicc_wallet_utils_go

import "github.com/WaykiChain/wicc-wallet-utils-go/bitcoin"

var (
	BWM 	*BTCWalletManager
)

func init(){

	BWM = NewBTCWalletManager()
	BWM.Config.RpcUser = "wayki"
	BWM.Config.RpcPassword = "admin@123"
	token := bitcoin.BasicAuth(BWM.Config.RpcUser, BWM.Config.RpcPassword)
	BWM.WalletClient =  bitcoin.NewClient(BWM.Config.ServerAPI, token, false)
	BWM.BitcoreClient = bitcoin.NewExplorer(BWM.Config.ServerAPI,true)
}


