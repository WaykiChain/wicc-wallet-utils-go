package wiccwallet

import (
	"errors"

	"github.com/WaykiChain/wicc-wallet-utils-go/commons"
	"github.com/btcsuite/btcutil/base58"
	"regexp"
)

const (
	// network type
	WAYKI_MAINTNET int = 1 // mainnet
	WAYKI_TESTNET  int = 2 // testnet

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

// OperVoteFund operation of vote fund
type OperVoteFund struct {
	PubKey    []byte //< public key, binary format
	VoteValue int64  //< add fund if >= 0, minus fund if < 0
}

// OperVoteFunds array of OperVoteFund
type OperVoteFunds struct {
	voteArray []*OperVoteFund
}

//NewOperVoteFunds create new OperVoteFunds
func NewOperVoteFunds() *OperVoteFunds {
	return &OperVoteFunds{}
}

//Len get the length of OperVoteFunds
func (votes *OperVoteFunds) Len(index int) int {
	return len(votes.voteArray)
}

//Get element of OperVoteFund by index
func (votes *OperVoteFunds) Get(index int) *OperVoteFund {
	return votes.voteArray[index]
}

//Add element to OperVoteFund
//pubKey is binary bytes
//voteValue add fund if >= 0, minus fund if < 0
func (votes *OperVoteFunds) Add(fund *OperVoteFund) {
	votes.voteArray = append(votes.voteArray, fund)
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

//RegisterUCoinContractTxParam param of the register contract tx
type UCoinRegisterContractTxParam struct {
	ValidHeight int64  // valid height Within the height of the latest block
	SrcRegId    string // the reg id of the register
	Fees        int64  // fees for mining
	FeeSymbol   string //WICC/WUSD
	Script      []byte // the contract script, binary format
	Description string // description of contract
}

type SignMessageParam struct {
	PublicKey string
	SignMessage string
}

// cdp stake Asset
type AssetModel struct {
	AssetAmount    int64 //stake asset amount
	AssetSymbol string  //stake asset symbol
}
//cdp stake asset list
type AssetModels struct {
	assetArray []*AssetModel
}

func (assetModels *AssetModels) Add(model *AssetModel)  {
	assetModels.assetArray= append(assetModels.assetArray, model)
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
	destArray []*Dest
}
func NewDestArr() *DestArr {
	return &DestArr{}
}

func (dests *DestArr) Add(dest *Dest)  {
	dests.destArray= append(dests.destArray, dest)
}
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

//Dex Sell Limit param of the tx
type DexLimitTxParam struct {
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
type DexMarketTxParam struct {
	ValidHeight int64  // valid height Within the height of the latest block
	SrcRegId    string // the reg id of the register
	PubKey      string
	FeeSymbol   string
	Fees        int64 // fees for mining
	AssetSymbol string
	CoinSymbol  string
	AssetAmount  int64
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

// errors
var (
	ERR_INVALID_MNEMONIC       = errors.New("Invalid Mnemonic")
	ERR_INVALID_NETWORK   = errors.New("Invalid Network type")
	ERR_INVALID_PRIVATE_KEY   = errors.New("privateKey invalid")
	ERR_NEGATIVE_VALID_HEIGHT = errors.New("ValidHeight can not be negative")
	ERR_INVALID_SRC_REG_ID    = errors.New("SrcRegId must be a valid RegID")
	ERR_INVALID_DEST_ADDR     = errors.New("DestAddr must be a valid RegID or Address")
	ERR_RANGE_VALUES          = errors.New("Values out of range")
	ERR_RANGE_FEE             = errors.New("Fees out of range")
	ERR_FEE_SMALLER_MIN       = errors.New("Fees smaller than MinTxFee")
	ERR_EMPTY_VOTES           = errors.New("Votes can be not empty")
	ERR_INVALID_VOTE_PUBKEY   = errors.New("Vote PubKey invalid, PubKey len must equal 33")
	ERR_RANGE_VOTE_VALUE      = errors.New("VoteValue out of range")
	ERR_INVALID_APP_ID        = errors.New("AppId must be a valid RegID")
	ERR_INVALID_CONTRACT_HEX  = errors.New("ContractHex must be valid hex format")
	ERR_INVALID_SCRIPT        = errors.New("Script can not be empty or is too large")
	ERR_INVALID_SCRIPT_DESC   = errors.New("Description of script is too large")

	ERR_CDP_TX_HASH      = errors.New("CDP tx hash error")
	ERR_CDP_STAKE_NUMBER = errors.New("CDP stake number error")
	ERR_COIN_TYPE        = errors.New("Coin type error")
	ERR_USER_PUBLICKEY   = errors.New("PublicKey invalid")

	ERR_ASK_PRICE   = errors.New("Ask Price invalid")
	ERR_SIGNATURE_ERROR       = errors.New("Signature error")
	ERR_SYMBOL_ERROR       = errors.New("Symbol Capital letter A-Z 6-7 digits [A_Z] error")
	ERR_ASSET_NAME_ERROR       = errors.New("Asset Name error")
	ERR_TOTAl_SUPPLY_ERROR       = errors.New("Asset Total Supply error")
	ERR_ASSET_UPDATE_TYPE_ERROR       = errors.New("Asset Update Type error")
	ERR_ASSET_UPDATE_OWNER_ERROR       = errors.New("Asset Update Owner error")
)

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func checkMoneyRange(value int64) bool {
	return value > 0 && value <= MAX_MONEY
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

func parseRegId(idStr string) *commons.UserIdWraper {
	if commons.IsRegIdStr(idStr) {
		return commons.NewRegUid(*commons.ParseRegId(idStr))
	}
	return nil
}

func parseAddressId(idStr string) *commons.UserIdWraper {
	addrBytes, _, err := base58.CheckDecode(idStr)
	if err != nil {
		return nil
	}
	return commons.NewAdressUid(commons.AddressId(addrBytes))
}

func parseUserId(idStr string) *commons.UserIdWraper {
	userId := parseRegId(idStr)
	if userId == nil {
		userId = parseAddressId(idStr)
	}
	return userId
}

func checkPubKey(pubKey []byte) bool {
	return len(pubKey) == 33
}
