package commons

import (
	"crypto/ecdsa"
	"encoding/hex"
	"regexp"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil/base58"
)

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

type WalletAddress struct {
	key     ecdsa.PrivateKey
	privKey string
	address string
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

type WaykiBaseSignTx struct {
	TxType      WaykiTxType
	Version     int64
	ValidHeight int64
	PubKey      []byte
	UserId      *UserIdWraper // current operating user id, the one want to do something
}
