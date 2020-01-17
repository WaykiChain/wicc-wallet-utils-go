package waykichain

import (
	"github.com/WaykiChain/wicc-wallet-utils-go/log"
)

type WalletManager struct {
	Wallet          *WICCWallet
	WalletClient    *Client                 // 节点客户端
	BaaSClient      *BaaSClient
	Config          *Config                 //钱包管理配置
	Log             *log.OWLogger           //日志工具
}

func NewWalletManager() *WalletManager {
	wm := WalletManager{}
	return &wm
}

//获取当前链上最新高度
func (wm *WalletManager) GetSynBlockHeight() (int64, error) {

	if wm.Config.ServerType == Node {
		return wm.WalletClient.GetSynBlockHeight()
	} else {
		return wm.BaaSClient.GetBaaSSynBlockHeight()
	}
}

//广播交易
func (wm *WalletManager) SubmitTxRaw(rawtx string) (string, error) {

	if wm.Config.ServerType == Node {
		return wm.WalletClient.SubmitTxRaw(rawtx)
	} else {
		return wm.BaaSClient.SubmitTxRaw(rawtx)
	}
}