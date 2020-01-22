package wicc_wallet_utils_go

type Wallet interface {
	//wallet
	GenerateAddressFromMnemonic(mnemonic,language string) 	(string, error)
	GenerateAddressFromPrivateKey(privateKey string)	(string, error)
	ExportPrivateKeyFromMnemonic(mnemonic,language string) 	(string, error)
	CheckAddress(address string) (bool, error)
	CheckPrivateKey(privateKey string) 	(bool, error)
	GetPubKeyFromPrivateKey(privateKey string) (string, error)
}

type ETHTransactions interface{
	//ETH、ERC20 transfer
	CreateRawTx(privateKeyStr string,chainId int64) (string, error)
}

type WICCTransactions interface {
	//All Tx type
	CreateRawTx(privateKey string) (* CreateRawTxResult, error)
}
