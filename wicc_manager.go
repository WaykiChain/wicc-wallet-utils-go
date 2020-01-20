package wicc_wallet_utils_go

import (
	"github.com/WaykiChain/wicc-wallet-utils-go/log"
	"github.com/WaykiChain/wicc-wallet-utils-go/waykichain"
)

type WICCWalletManager struct {
	Wallet          *WICCWallet
	WalletClient    *waykichain.Client                 // 节点客户端
	BaaSClient      *waykichain.BaaSClient
	Config          *waykichain.Config                 //钱包管理配置
	Log             *log.OWLogger           //日志工具
}

func NewWICCWalletManager() *WICCWalletManager {
	wm := WICCWalletManager{}
	return &wm
}

//获取当前链上最新高度
func (wm *WICCWalletManager) GetSynBlockHeight() (int64, error) {

	if wm.Config.ServerType == waykichain.Node {
		return wm.WalletClient.GetSynBlockHeight()
	} else {
		return wm.BaaSClient.GetBaaSSynBlockHeight()
	}
}

//广播交易
func (wm *WICCWalletManager) SubmitTxRaw(rawtx string) (string, error) {

	if wm.Config.ServerType == waykichain.Node {
		return wm.WalletClient.SubmitTxRaw(rawtx)
	} else {
		return wm.BaaSClient.SubmitTxRaw(rawtx)
	}
}