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
	ADD_FUND   //投票
	MINUS_FUND //撤销投票
)

type WalletAddress struct {
	key     ecdsa.PrivateKey
	privKey string
	address string
}

type BaseSignTxParams struct {
	PrivateKey  string
	RegId       string
	UserPubKey  []byte
	MinerPubKey []byte
	ValidHeight int64
	Fees        int64
	TxType      WaykiTxType
	Version     int64
}

func parseRegId(regId string) []int64 {
	regidStr := strings.Split(regId, "-")
	regHeight, _ := strconv.ParseInt(regidStr[0], 10, 64)
	regIndex, _ := strconv.ParseInt(regidStr[1], 10, 64)
	regIdArray := []int64{regHeight, regIndex}
	return regIdArray
}
