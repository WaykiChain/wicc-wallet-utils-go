package wicc_wallet_utils_go

import (
	"github.com/WaykiChain/wicc-wallet-utils-go/waykichain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/base58"
	"regexp"
)
const (
	// const for transactions
	TX_VERSION                    = 1                             // transaction version
	INITIAL_COIN                  = 210000000                     // initial coin, unit: wicc
	MONEY_PER_COIN                = 100000000                     // money per coin, unit: sawi
	MAX_MONEY                     = INITIAL_COIN * MONEY_PER_COIN // the max money in WaykiChain
	MIN_TX_FEE                    = 10000                         // tx fee min value, unit: sawi
	CONTRACT_SCRIPT_MAX_SIZE      = 65536                         //64 KB max for contract script size, unit: bytes
	CONTRACT_SCRIPT_DESC_MAX_SIZE = 512                           //max for contract script description size, unit: bytes
	MIN_TX_FEE_CDP                = 100000                        // cdp tx fee min value, unit: sawi
)


func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func checkMoneyRange(value int64) bool {
	return value >= 0 && value <= MAX_MONEY
}


func checkCdpMinTxFee(fees int64) bool {

	return fees >= MIN_TX_FEE_CDP
}

func checkMinTxFee(fees int64) bool {
	return fees >= MIN_TX_FEE
}

func checkAssetSymbol(symbol string) bool {
	var symbolMatch="^[A-Z]{6,7}$"
	var match=regexp.MustCompile(symbolMatch)
	var ma=match.MatchString(symbol)
	return ma
}

func parseRegId(idStr string) *waykichain.UserIdWraper {
	if waykichain.IsRegIdStr(idStr) {
		return waykichain.NewRegUid(*waykichain.ParseRegId(idStr))
	}
	return nil
}

func parseAddressId(idStr string) *waykichain.UserIdWraper {
	addrBytes, _, err := base58.CheckDecode(idStr)
	if err != nil {
		return nil
	}
	return waykichain.NewAdressUid(waykichain.AddressId(addrBytes))
}

func parseUserId(idStr string) *waykichain.UserIdWraper {
	userId := parseRegId(idStr)
	if userId == nil {
		userId = parseAddressId(idStr)
	}
	return userId
}

func checkPubKey(pubKey []byte) bool {
	return len(pubKey) == 33
}
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
// OperVoteFund operation of vote fund
type OperVoteFund struct {
	PubKey    []byte //< public key, binary format
	VoteValue int64  //< add fund if >= 0, minus fund if < 0
}

// OperVoteFunds array of OperVoteFund
type OperVoteFunds struct {
	VoteArray []*OperVoteFund
}

//NewOperVoteFunds create new OperVoteFunds
func NewOperVoteFunds() *OperVoteFunds {
	return &OperVoteFunds{}
}

//Len get the length of OperVoteFunds
func (votes *OperVoteFunds) Len(index int) int {
	return len(votes.VoteArray)
}

//Get element of OperVoteFund by index
func (votes *OperVoteFunds) Get(index int) *OperVoteFund {
	return votes.VoteArray[index]
}

//Add element to OperVoteFund
//pubKey is binary bytes
//voteValue add fund if >= 0, minus fund if < 0
func (votes *OperVoteFunds) Add(fund *OperVoteFund) {
	votes.VoteArray = append(votes.VoteArray, fund)
}
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

type CreateRawTxResult struct{
	* waykichain.SignResult
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

type VerifyMsgSignResult struct{
	IsRight bool
	Address string
}