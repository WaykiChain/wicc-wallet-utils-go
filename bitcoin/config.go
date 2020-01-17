package bitcoin

import (
	"github.com/WaykiChain/wicc-wallet-utils-go/common"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/shopspring/decimal"
)

var(
	BTCMainnetConf 		 = "BTCMainnetConf"
	BTCMainnetSegwitConf = "BTCMainnetSegwitConf"
	BTCTestnetConf 		 = "BTCTestnetConf"
	BTCTestnetSegwitConf = "BTCTestnetSegwitConf"
)

type BTCWalletConfig struct {
	coinType uint32
	isSegwit bool
	netParam *chaincfg.Params
}
func NewWalletConfig(wc string) *BTCWalletConfig{
	newConf := BTCWalletConfig{}
	switch wc {
	case BTCMainnetConf:
		newConf = BTCWalletConfig{common.BTC,false, &common.BTCParams}
	case BTCMainnetSegwitConf:
		newConf = BTCWalletConfig{common.BTC,true, &common.BTCParams}
	case BTCTestnetConf:
		newConf = BTCWalletConfig{common.BTC_TESTNET,false, &common.BTCTestnetParams}
	case BTCTestnetSegwitConf:
		newConf = BTCWalletConfig{common.BTC_TESTNET,true, &common.BTCTestnetParams}
	default:
		newConf = BTCWalletConfig{common.BTC,false, &common.BTCParams}
	}
	return &newConf
}

const (
	RPCServerCore     = 0 //RPC服务，bitcoin核心钱包
	RPCServerExplorer = 1 //RPC服务，insight-API
)

type Config struct {
	//钱包配置
	WalletConfig string
	//钱包服务API
	ServerAPI  string
	//RPC认证账户名
	RpcUser string
	//RPC认证账户密码
	RpcPassword string
	//后台数据源类型
	RPCServerType int
	//小数位精度
	Decimals int32
	//最大的输入数量
	MaxTxInputs int
	//最低手续费
	MinFees decimal.Decimal
	//是否开启Debug
	Debug bool
}


func NewConfig() *Config {

	c := Config{}

	//c.ServerAPI = "https://insight.bitpay.com/api/"
	c.ServerAPI = "http://10.0.0.11:3000/api/BTC/testnet/"
	//c.ServerAPI = "http://10.0.0.11:18332"
	c.RpcUser = ""
	c.RpcPassword = ""
	c.RPCServerType = RPCServerExplorer
	c.Decimals = 8
	c.MaxTxInputs = 150
	c.MinFees = decimal.Zero
	c.Debug = true

	return &c
}