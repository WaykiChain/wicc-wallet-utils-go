package bitcoin

var (
	tw 				  *WalletManager
)

func init(){

	tw = NewWalletManager()
	tw.Config.RpcUser = "wayki"
	tw.Config.RpcPassword = "admin@123"
	token := BasicAuth(tw.Config.RpcUser, tw.Config.RpcPassword)
	tw.WalletClient = NewClient(tw.Config.ServerAPI, token, false)

	tw.BitcoreClient = NewExplorer(tw.Config.ServerAPI,true)
}


