package waykichain

import (
	"encoding/hex"
	"github.com/WaykiChain/wicc-wallet-utils-go/common"
	"github.com/WaykiChain/wicc-wallet-utils-go/common/hash"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
	"regexp"
	"strconv"
	"strings"
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

func parseRegId(idStr string) *UserIdWraper {
	if IsRegIdStr(idStr) {
		return NewRegUid(*ParseRegId(idStr))
	}
	return nil
}

func parseAddressId(idStr string) *UserIdWraper {
	addrBytes, _, err := base58.CheckDecode(idStr)
	if err != nil {
		return nil
	}
	return NewAdressUid(AddressId(addrBytes))
}

func parseUserId(idStr string) *UserIdWraper {
	userId := parseRegId(idStr)
	if userId == nil {
		userId = parseAddressId(idStr)
	}
	return userId
}

func checkPubKey(pubKey []byte) bool {
	return len(pubKey) == 33
}



type WaykiTxType uint32

const (
	REWARD_TX   WaykiTxType = 1 + iota //!< reward tx
	REG_ACCT_TX                        //!< tx that used to register account
	COMMON_TX                          //!< transfer coin from one account to another
	CONTRACT_TX                        //!< contract tx
	REG_CONT_TX                        //!< register contract
	DELEGATE_TX                        //!< delegate tx

	FCOIN_STAKE_TX           = 8
	ASSET_ISSUE_TX=9  //!< a user issues onchain asset
	ASSET_UPDATE_TX=10   //!< a user update onchain asset
	UCOIN_TRANSFER_TX        = 11
	UCONTRACT_DEPLOY_TX         = 14   //!< universal VM contract deployment
	UCOIN_CONTRACT_INVOKE_TX = 15
	PRICE_FEED_TX            = 16

	CDP_STAKE_TX     = 21
	CDP_REDEEMP_TX   = 22
	CDP_LIQUIDATE_TX = 23

	DEX_SETTLE_TX            = 89 //!< dex settle Tx
	DEX_CANCEL_ORDER_TX      = 88 //!< dex cancel order Tx
	DEX_BUY_LIMIT_ORDER_TX   = 84 //!< dex buy limit price order Tx
	DEX_SELL_LIMIT_ORDER_TX  = 85 //!< dex sell limit price order Tx
	DEX_BUY_MARKET_ORDER_TX  = 86 //!< dex buy market price order Tx
	DEX_SELL_MARKET_ORDER_TX = 87 //!< dex sell market price order Tx
)

type WaykiVoteType uint32

const (
	NULL_OPER  WaykiVoteType = iota
	ADD_FUND                 //投票
	MINUS_FUND               //撤销投票
)

func GetVoteTypeByValue(value int64) WaykiVoteType {
	ret := ADD_FUND
	if value < 0 {
		ret = MINUS_FUND
	}
	return ret
}

type RegId struct {
	Height uint64
	Index  uint64
}

func IsRegIdStr(regId string) bool {
	re := regexp.MustCompile(`^\s*(\d+)\-(\d+)\s*$`)
	return re.MatchString(regId)
}

func ParseRegId(regId string) *RegId {
	regidStr := strings.Split(regId, "-")
	regHeight, _ := strconv.ParseInt(regidStr[0], 10, 64)
	regIndex, _ := strconv.ParseInt(regidStr[1], 10, 64)
	return &RegId{uint64(regHeight), uint64(regIndex)}
}

type PubKeyId []byte

func NewPubKeyIdByKey(privKey *btcec.PrivateKey) *PubKeyId {
	var myid PubKeyId = privKey.PubKey().SerializeCompressed()
	return &myid
}

func NewPubKeyIdByStr(str string) *PubKeyId {
	myid, _ := hex.DecodeString(str)
	return (*PubKeyId)(&myid)
}

type AddressId []byte

type UserIdType int

const (
	UID_REG     UserIdType = iota //< reg id
	UID_PUB_KEY                   //< public key
	UID_ADDRESS                   //< wicc address
)

type UserIdWraper struct {
	idType UserIdType
	id     interface{}
}

func NewPubKeyUid(pubKey PubKeyId) *UserIdWraper {
	return &UserIdWraper{UID_PUB_KEY, pubKey}
}

func NewPubKeyUidByStr(pubKey string) *UserIdWraper {
	id, _ := hex.DecodeString(pubKey)
	return NewPubKeyUid(id)
}

func NewAdressUid(address AddressId) *UserIdWraper {
	return &UserIdWraper{UID_ADDRESS, address}
}

func NewAdressUidByStr(address string) *UserIdWraper {
	addrBytes, _, _ := base58.CheckDecode(address)
	return &UserIdWraper{UID_ADDRESS, AddressId(addrBytes)}
}

func NewRegUid(regId RegId) *UserIdWraper {
	return &UserIdWraper{UID_REG, regId}
}

func NewRegUidByStr(regId string) *UserIdWraper {
	return NewRegUid(*ParseRegId(regId))
}

func (uid UserIdWraper) GetType() UserIdType {
	return uid.idType
}

func (uid UserIdWraper) GetId() interface{} {
	return uid.id
}


func calSignatureTxid(paramHash []byte,wifKey *btcutil.WIF)([]byte,string,error){
	txidRevrse := hash.DoubleHash256(paramHash)
	txidHex := common.Reverse(txidRevrse)
	signature, err := wifKey.PrivKey.Sign(txidRevrse)
	if err != nil{
		return nil,"",err
	}
	return signature.Serialize(),hex.EncodeToString(txidHex),nil
}