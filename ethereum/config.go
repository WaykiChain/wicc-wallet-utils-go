package ethereum

const(
	Mainnet = "mainnet"
	Ropsten = "ropsten"
	Rinkeby = "rinkeby"
	Local	= "local"
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
