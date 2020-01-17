package ethereum

import "github.com/WaykiChain/wicc-wallet-utils-go/log"

type WalletManager struct {
	Wallet          *ETHWallet
	WalletClient    *Client                  // 节点客户端
	Config          *Config                 //钱包管理配置
	Log             *log.OWLogger                 //日志工具
}

func NewWalletManager() *WalletManager {
	wm := WalletManager{}
	return &wm
}