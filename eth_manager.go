package wicc_wallet_utils_go

import (
	"github.com/WaykiChain/wicc-wallet-utils-go/ethereum"
	"github.com/WaykiChain/wicc-wallet-utils-go/log"
)

type ETHWalletManager struct {
	Wallet          *ETHWallet
	WalletClient    *Client                  // 节点客户端
	Config          *ethereum.Config                 //钱包管理配置
	Log             *log.OWLogger                 //日志工具
}

func NewETHWalletManager() *ETHWalletManager {
	wm := ETHWalletManager{}
	return &wm
}