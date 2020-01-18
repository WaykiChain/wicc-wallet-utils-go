package waykichain

import "github.com/btcsuite/btcd/chaincfg"

/***************************************UCOIN_TRANSFER_TX*************************/
//UCoin Transfer param of the tx
type UCoinTransferTxParam struct {
	ValidHeight int64  // valid height Within the height of the latest block
	SrcRegId    string // the reg id of the register
	PubKey      string
	FeeSymbol   string
	Fees        int64 // fees for mining
	Dests       *DestArr
	Memo        string
}

/***************************************UCOIN_CONTRACT_INVOKE_TX*************************/
//UCoin Contract param of the tx
type UCoinContractTxParam struct {
	ValidHeight int64  // valid height Within the height of the latest block
	SrcRegId    string // the reg id of the caller
	AppId       string // the reg id of the contract app
	Fees        int64  // fees for mining
	CoinAmount      int64  // the values send to the contract app
	ContractHex string // the command of contract, hex format
	PubKey      string
	FeeSymbol string      //Fee Type (WICC/WUSD)
	CoinSymbol string   //From Coin Type
}

/***************************************UCONTRACT_DEPLOY_TX*************************/
//RegisterUCoinContractTxParam param of the register contract tx
type UCoinRegisterContractTxParam struct {
	ValidHeight int64  // valid height Within the height of the latest block
	SrcRegId    string // the reg id of the register
	Fees        int64  // fees for mining
	FeeSymbol   string //WICC/WUSD
	Script      []byte // the contract script, binary format
	Description string // description of contract
}


//DelegateTxParam param of the delegate tx
type DelegateTxParam struct {
	ValidHeight int64          // valid height Within the height of the latest block
	SrcRegId    string         // the reg id of the voter
	Fees        int64          // fees for mining
	Votes       *OperVoteFunds // vote list
	PubKey      string
}

//CallContractTxParam param of the call contract tx
type CallContractTxParam struct {
	ValidHeight int64  // valid height Within the height of the latest block
	SrcRegId    string // the reg id of the caller
	AppId       string // the reg id of the contract app
	Fees        int64  // fees for mining
	Values      int64  // the values send to the contract app
	ContractHex string // the command of contract, hex format
	PubKey      string
}

//RegisterContractTxParam param of the register contract tx
type RegisterContractTxParam struct {
	ValidHeight int64  // valid height Within the height of the latest block
	SrcRegId    string // the reg id of the register
	Fees        int64  // fees for mining
	Script      []byte // the contract script, binary format
	Description string // description of contract
}


type SignMsgResult struct {
	PublicKey string
	Signature string
}

type SignMsgInput struct {
	PrivateKey string
	Data string
}
type VerifySignInput struct {
	Signature string
	PublicKey string
	Data string
	NetParams chaincfg.Params
}



// cdp stake Asset
type AssetModel struct {
	AssetAmount    int64 //stake asset amount
	AssetSymbol string  //stake asset symbol
}
//cdp stake asset list
type AssetModels struct {
	AssetArray []*AssetModel
}

func (assetModels *AssetModels) Add(model *AssetModel)  {
	assetModels.AssetArray= append(assetModels.AssetArray, model)
}

func NewCdpAssets() *AssetModels {
	return &AssetModels{}
}

//Cdp Stake param of the tx
type CdpStakeTxParam struct {
	ValidHeight int64  // valid height Within the height of the latest block
	SrcRegId    string // the reg id of the register
	PubKey      string
	FeeSymbol   string
	Fees        int64 // fees for mining
	CdpTxid     string
	ScoinSymbol string  //get coin symbol
	ScoinMint   int64 // get coin amount
	Assets   *AssetModels
}

//Cdp Redeem param of the tx
type CdpRedeemTxParam struct {
	ValidHeight int64  // valid height Within the height of the latest block
	SrcRegId    string // the reg id of the register
	PubKey      string
	FeeSymbol   string
	Fees        int64 // fees for mining
	CdpTxid     string
	ScoinsToRepay  int64 //repay coin number
	Assets   *AssetModels
}

//Cdp Redeem param of the tx
type CdpLiquidateTxParam struct {
	ValidHeight int64  // valid height Within the height of the latest block
	SrcRegId    string // the reg id of the register
	PubKey      string
	FeeSymbol   string
	Fees        int64 // fees for mining
	CdpTxid     string // target CDP to liquidate
	ScoinsLiquidate  int64   // partial liquidation is allowed, must include penalty fees in
	AssetSymbol  string  //stake asset symbol
}

type Dest struct {
	CoinSymbol string   //From Coin Type
	CoinAmount uint64
	DestAddr    string
}

type DestArr struct {
	DestArray []*Dest
}
func NewDestArr() *DestArr {
	return &DestArr{}
}

func (dests *DestArr) Add(dest *Dest)  {
	dests.DestArray= append(dests.DestArray, dest)
}


//Dex Sell Limit param of the tx
type DexLimitSellTxParam struct {
	ValidHeight int64  // valid height Within the height of the latest block
	SrcRegId    string // the reg id of the register
	PubKey      string
	FeeSymbol   string
	Fees        int64 // fees for mining
	AssetSymbol string
	CoinSymbol  string
	AssetAmount  int64
	Price     int64
}

//Dex Buy Limit param of the tx
type DexLimitBuyTxParam struct {
	ValidHeight int64  // valid height Within the height of the latest block
	SrcRegId    string // the reg id of the register
	PubKey      string
	FeeSymbol   string
	Fees        int64 // fees for mining
	AssetSymbol string
	CoinSymbol  string
	AssetAmount  int64
	Price     int64
}

//Dex market Sell param of the tx
type DexMarketSellTxParam struct {
	ValidHeight int64  // valid height Within the height of the latest block
	SrcRegId    string // the reg id of the register
	PubKey      string
	FeeSymbol   string
	Fees        int64 // fees for mining
	AssetSymbol string
	CoinSymbol  string
	AssetAmount  int64
}

//Dex market Buy param of the tx
type DexMarketBuyTxParam struct {
	ValidHeight int64  // valid height Within the height of the latest block
	SrcRegId    string // the reg id of the register
	PubKey      string
	FeeSymbol   string
	Fees        int64 // fees for mining
	AssetSymbol string
	CoinSymbol  string
	CoinAmount  int64
}

type DexCancelTxParam struct {
	ValidHeight int64  // valid height Within the height of the latest block
	SrcRegId    string // the reg id of the register
	PubKey      string
	FeeSymbol   string
	Fees        int64 // fees for mining
	DexTxid     string
}

type AssetIssueTxParam struct {
	ValidHeight int64  // valid height Within the height of the latest block
	SrcRegId    string // the reg id of the register
	Fees       int64
	FeeSymbol string      //Fee Type (WICC/WUSD)
	AssetSymbol string   //From Coin Type
	AssetName   string
	AssetTotal   int64
	AssetOwner   string //owner regid
	MinTable     bool
}

type AssetUpdateTxParam struct {
	ValidHeight int64  // valid height Within the height of the latest block
	SrcRegId    string // the reg id of the register
	Fees       int64
	UpdateType int
	FeeSymbol string      //Fee Type (WICC/WUSD)
	AssetSymbol string   //From Coin Type
	AssetName   string
	AssetTotal   int64
	AssetOwner   string //owner regid
}

//RegisterAccountTxParam register account tx param
type RegisterAccountTxParam struct {
	ValidHeight int64 // valid height Within the height of the latest block
	Fees        int64 // fees for mining
}

type CommonTxParam struct {
	ValidHeight int64  // valid height Within the height of the latest block
	SrcRegId    string // the source reg id that the transaction send from
	DestAddr    string // the dest address that the transaction send to
	Values      int64  // transfer values
	Fees        int64  // fees for mining
	PubKey      string
	Memo        string
}

const (
	ASSET_OWNER_UID   int = 1
	ASSET_NAME                        = 2
	ASSET_MINT_AMOUNT                 = 3
)

type SignResult struct{
	RawTx string
	Txid  string
}
