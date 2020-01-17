package ethereum

import "C"
import (
	"github.com/WaykiChain/wicc-wallet-utils-go/common"
)

const (
	GAS_LIMIT         = 21000
	GAS_PRICE         = 500000000000
)

const (
	ETH_GET_TOKEN_BALANCE_METHOD      = "0x70a08231"
	ETH_TRANSFER_TOKEN_METHOD 		  = "0xa9059cbb"
	ETH_TRANSFER_EVENT_ID             = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
)

const (
	SOLIDITY_TYPE_ADDRESS = "address"
	SOLIDITY_TYPE_UINT256 = "uint256"
	SOLIDITY_TYPE_UINT160 = "uint160"
)

const(
	LEATEST  = "latest"
	EARLIEST = "earliest"
	PENDING  = "pending"
)

const(
	Mainnet = "mainnet"
	Ropsten = "ropsten"
	Rinkeby = "rinkeby"
	Local	= "local"
)

//ETHWallet config
var(
	ETHWalletConf 		 = "ETHWalletConf"
	ETHWalletLedgerConf =  "ETHWalletLedgerConf"
)
type ETHWalletConfig struct {
	coinType uint32
	isLedger bool
}
func NewWalletConfig(wc string) *ETHWalletConfig{
	newConf := ETHWalletConfig{}
	switch wc {
	case ETHWalletConf:
		newConf = ETHWalletConfig{common.ETH,false}
	case ETHWalletLedgerConf:
		newConf = ETHWalletConfig{common.ETH,true}
	default:
		newConf = ETHWalletConfig{common.ETH,false}
	}
	return &newConf
}

type Config struct {
	//钱包配置
	WalletConfig string
	//钱包服务API
	ServerAPI  string
	//RPC认证账户名
	RpcUser string
	//RPC认证账户密码
	RpcPassword string
	//小数位精度
	Decimals int32
	//是否开启Debug
	Debug bool
	//ChainID ， Mainnet：0 ，Ropsten：3 ， Rinkeby：4
	ChainID int64
	//Keystore 存储路径
	KeystoreDir string
}

func NewConfig() *Config {
	c := Config{}
	c.KeystoreDir = "./"
	return &c
}
