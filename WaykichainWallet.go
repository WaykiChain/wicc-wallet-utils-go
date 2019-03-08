package wiccwallet

import (
	"encoding/hex"

	"github.com/WaykiChain/wicc-wallet-utils-go/commons"
)

const (
	WAYKI_TESTNET  int = 1
	WAYKI_MAINTNET int = 2
)

type OperVoteFund struct {
	PubKey    []byte //< public key, binary format
	VoteValue int64  //< add fund if >= 0, minus fund if < 0
}

type OperVoteFunds struct {
	voteArray []*OperVoteFund
}

func NewOperVoteFunds() *OperVoteFunds {
	return &OperVoteFunds{}
}

func (votes *OperVoteFunds) Len(index int) int {
	return len(votes.voteArray)
}

func (votes *OperVoteFunds) Get(index int) *OperVoteFund {
	return votes.voteArray[index]
}

func (votes *OperVoteFunds) Add(pubKey []byte, voteValue int64) *OperVoteFund {
	vote := OperVoteFund{pubKey, voteValue}
	votes.voteArray = append(votes.voteArray, &vote)
	return &vote
}

//Generate Mnemonics string, saprated by space, default language is EN(english)
func GenerateMnemonics() string {
	mn := NewMnemonicWithLanguage(ENGLISH)
	words, err := mn.GenerateMnemonic()
	if err != nil {
		return ""
	}
	return words
}

//助记词转换地址
// netType: WAYKI_TESTNET or WAYKI_MAINTNET
func GetAddressFromMnemonic(words string, netType int) string {
	address := commons.GenerateAddress(words, commons.Network(netType))
	return address
}

//助记词转私钥
// netType: WAYKI_TESTNET or WAYKI_MAINTNET
func GetPrivateKeyFromMnemonic(words string, netType int) string {
	privateKey := commons.GeneratePrivateKey(words, commons.Network(netType))
	return privateKey
}

//私钥转地址
// netType: WAYKI_TESTNET or WAYKI_MAINTNET
func GetAddressFromPrivateKey(words string, netType int) string {
	address := commons.ImportPrivateKey(words, commons.Network(netType))
	return address
}

//注册账户交易签名
func SignRegisterTx(height int64, fees int64, privateKey string) string {
	var tx commons.WaykiRegisterAccountTx
	tx.PrivateKey = privateKey
	tx.ValidHeight = height
	tx.Fees = uint64(fees)
	tx.TxType = commons.REG_ACCT_TX
	tx.Version = 1
	hash := tx.SignTx()
	return hash
}

//普通交易签名
func SignCommonTx(values int64, regId string, toAddr string, height int64, fees int64, privateKey string) string {
	var tx commons.WaykiCommonTx
	tx.Values = uint64(values)
	tx.DestId = commons.NewAdressUidByStr(toAddr)
	tx.PrivateKey = privateKey
	tx.UserId = commons.NewRegUidByStr(regId)
	tx.ValidHeight = height
	tx.Fees = uint64(fees)
	tx.TxType = commons.COMMON_TX
	tx.Version = 1
	hash := tx.SignTx()
	return hash
}

//投票交易签名
func SignDelegateTx(regId string, height int64, fees int64, privateKey string, votes *OperVoteFunds) string {

	var voteData []commons.OperVoteFund
	for i := 0; i < len(votes.voteArray); i++ {
		inVote := votes.voteArray[i]
		var v commons.OperVoteFund
		var pubKey commons.PubKeyId = inVote.PubKey
		v.PubKey = &pubKey
		v.VoteType = commons.GetVoteTypeByValue(inVote.VoteValue)
		v.VoteValue = abs(inVote.VoteValue)
		voteData = append(voteData, v)
	}

	var tx commons.WaykiDelegateTx
	tx.PrivateKey = privateKey
	tx.UserId = commons.NewRegUidByStr(regId)
	tx.ValidHeight = height
	tx.Fees = uint64(fees)
	tx.TxType = commons.DELEGATE_TX
	tx.Version = 1
	tx.OperVoteFunds = voteData
	hash := tx.SignTx()
	return hash
}

//智能合约交易签名
func SignContractTx(values int64, height int64, fees int64, privateKey string, regId string, appId string, contractStr string) string {
	var tx commons.WaykiCallContractTx
	tx.Values = uint64(values)
	tx.PrivateKey = privateKey
	tx.UserId = commons.NewRegUidByStr(regId)
	tx.AppId = commons.NewRegUidByStr(appId)
	tx.ValidHeight = height
	tx.Fees = uint64(fees)
	tx.TxType = commons.CONTRACT_TX
	tx.Version = 1
	binary, _ := hex.DecodeString(contractStr)
	tx.Contract = []byte(binary)
	hash := tx.SignTx()
	return hash
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
