package hdwallet

func init() {
	coins[WICC_TESTNET] = newWICCTestnet
}

type wicc_testnet struct {
	name   string
	symbol string
	key    *Key
}

func newWICCTestnet(key *Key) Wallet {

	key.opt.Params = &WICCTestnetParams
	return &wicc_testnet{
		name:   "WaykiChain",
		symbol: "WICC_TEST",
		key:    key,
	}
}

func (c *wicc_testnet) GetType() uint32 {
	return c.key.opt.CoinType
}

func (c *wicc_testnet) GetName() string {
	return c.name
}

func (c *wicc_testnet) GetSymbol() string {
	return c.symbol
}

func (c *wicc_testnet) GetKey() *Key {
	return c.key
}

func (c *wicc_testnet) GetAddress() (string, error) {
	return c.key.AddressBTC()
}
