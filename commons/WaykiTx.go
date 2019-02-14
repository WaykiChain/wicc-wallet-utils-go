package commons

import (
	"crypto/ecdsa"
	"strconv"
	"strings"
)

type WaykiTxType uint32

const (
	TX_REGISTERACCOUNT WaykiTxType = 2 + iota
	TX_COMMON
	TX_CONTRACT
	REG_APP_TX
	TX_DELEGATE
)

type WaykiVoteType uint32

const (
	NULL_OPER  WaykiVoteType = iota
	ADD_FUND                 //投票
	MINUS_FUND               //撤销投票
)

type WalletAddress struct {
	key     ecdsa.PrivateKey
	privKey string
	address string
}

type BaseSignTxParams struct {
	PrivateKey  string
	TxType      WaykiTxType
	Version     int64
	ValidHeight int64
	UserId      string // current operating user id, which one will do something
}

func parseRegId(regId string) []int64 {
	regidStr := strings.Split(regId, "-")
	regHeight, _ := strconv.ParseInt(regidStr[0], 10, 64)
	regIndex, _ := strconv.ParseInt(regidStr[1], 10, 64)
	regIdArray := []int64{regHeight, regIndex}
	return regIdArray
}

type UserIdType int

const (
	UID_REG UserIdType = iota
	UID_PUB_KEY
	UID_ADDRESS
)

type UserId struct {
	idType UserIdType
	id     interface{}
}

func ParseUserId(str string) *UserId {
	// TODO:...
	return nil
}

func NewPubKeyId(pubKey []byte) *UserId {
	return &UserId{UID_PUB_KEY, pubKey}
}

func (uid UserId) GetType() UserIdType {
	return uid.idType
}

func (uid UserId) GetId() interface{} {
	return uid.id
}

type RegId struct {
	Height uint64
	Index  uint64
}

func ParseRegId(regId string) *RegId {
	regidStr := strings.Split(regId, "-")
	regHeight, _ := strconv.ParseInt(regidStr[0], 10, 64)
	regIndex, _ := strconv.ParseInt(regidStr[1], 10, 64)
	return &RegId{uint64(regHeight), uint64(regIndex)}
}
