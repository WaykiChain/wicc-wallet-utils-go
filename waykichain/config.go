package waykichain

import (
	"github.com/WaykiChain/wicc-wallet-utils-go/common"
	"github.com/btcsuite/btcd/chaincfg"
)

var(
	WICCMainnetConf 		 = "WICCMainnetConf"
	WICCTestnetConf 		 = "WICCTestnetConf"

	Node    = 0 //节点RPC服务
	BaaS    = 1 //BaaS服务
)

type WICCWalletConfig struct {
	coinType uint32
	netParam *chaincfg.Params
}
func NewWalletConfig(wc string) *WICCWalletConfig{
	newConf := WICCWalletConfig{}
	switch wc {
	case WICCMainnetConf:
		newConf = WICCWalletConfig{common.WICC, &common.WICCParams}
	case WICCTestnetConf:
		newConf = WICCWalletConfig{common.WICC_TESTNET,&common.WICCTestnetParams}
	default:
		newConf = WICCWalletConfig{common.WICC, &common.WICCParams}
	}
	return &newConf
}


type Config struct {
	//钱包配置
	WalletConfig string
	//RPC服务API
	NodeServerAPI  string
	//BaaS服务API
	BaaSServerAPI  string
	//RPC认证账户名
	RpcUser string
	//RPC认证账户密码
	RpcPassword string
	//小数位精度
	Decimals int32
	//是否开启Debug
	Debug bool
	//后台数据源类型
	ServerType int
}

func NewConfig() *Config {
	c := Config{}
	return &c
}