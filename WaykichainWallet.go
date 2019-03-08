package wiccwallet

import (
	"encoding/hex"

	"github.com/WaykiChain/wicc-wallet-utils-go/commons"
	"github.com/btcsuite/btcutil"
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

type RegisterAccountTxParam struct {
	ValidHeight int64
	Fees        int64
}

//注册账户交易签名
func SignRegisterAccountTx(privateKey string, params *RegisterAccountTxParam) string {
	var tx commons.WaykiRegisterAccountTx
	tx.PrivateKey = privateKey
	tx.ValidHeight = params.ValidHeight
	tx.Fees = uint64(params.Fees)
	tx.TxType = commons.REG_ACCT_TX
	tx.Version = 1
	wif, _ := btcutil.DecodeWIF(tx.PrivateKey)
	tx.UserId = commons.NewPubKeyUid(*commons.NewPubKeyIdByKey(wif.PrivKey))
	hash := tx.SignTx()
	return hash
}

type CommonTxParam struct {
	ValidHeight int64
	SrcRegId    string
	DestAddr    string
	Values      int64
	Fees        int64
}

//普通交易签名
func SignCommonTx(privateKey string, params *CommonTxParam) string {
	var tx commons.WaykiCommonTx
	tx.PrivateKey = privateKey
	tx.ValidHeight = params.ValidHeight
	tx.UserId = commons.NewRegUidByStr(params.SrcRegId)
	tx.DestId = commons.NewAdressUidByStr(params.DestAddr)
	tx.Values = uint64(params.Values)
	tx.Fees = uint64(params.Fees)
	tx.TxType = commons.COMMON_TX
	tx.Version = 1
	hash := tx.SignTx()
	return hash
}

type DelegateTxParam struct {
	ValidHeight int64
	SrcRegId    string
	Fees        int64
	Votes       *OperVoteFunds
}

//投票交易签名
func SignDelegateTx(privateKey string, params *DelegateTxParam) string {

	var voteData []commons.OperVoteFund
	for i := 0; i < len(params.Votes.voteArray); i++ {
		inVote := params.Votes.voteArray[i]
		var v commons.OperVoteFund
		var pubKey commons.PubKeyId = inVote.PubKey
		v.PubKey = &pubKey
		v.VoteType = commons.GetVoteTypeByValue(inVote.VoteValue)
		v.VoteValue = abs(inVote.VoteValue)
		voteData = append(voteData, v)
	}

	var tx commons.WaykiDelegateTx
	tx.PrivateKey = privateKey
	tx.UserId = commons.NewRegUidByStr(params.SrcRegId)
	tx.ValidHeight = params.ValidHeight
	tx.Fees = uint64(params.Fees)
	tx.TxType = commons.DELEGATE_TX
	tx.Version = 1
	tx.OperVoteFunds = voteData
	hash := tx.SignTx()
	return hash
}

type CallContractTxParam struct {
	ValidHeight int64
	SrcRegId    string
	AppId       string
	Fees        int64
	Values      int64
	ContractHex string
}

//智能合约交易签名
func SignCallContractTx(privateKey string, params *CallContractTxParam) string {
	var tx commons.WaykiCallContractTx
	tx.Values = uint64(params.Values)
	tx.PrivateKey = privateKey
	tx.UserId = commons.NewRegUidByStr(params.SrcRegId)
	tx.AppId = commons.NewRegUidByStr(params.AppId)
	tx.ValidHeight = params.ValidHeight
	tx.Fees = uint64(params.Fees)
	tx.TxType = commons.CONTRACT_TX
	tx.Version = 1
	binary, _ := hex.DecodeString(params.ContractHex)
	tx.Contract = []byte(binary)
	hash := tx.SignTx()
	return hash
}

type RegisterContractTxParam struct {
	ValidHeight int64
	SrcRegId    string
	Fees        int64
	Script      []byte
	Description string
}

func SignRegisterContractTx(privateKey string, params *RegisterContractTxParam) string {

	var tx commons.WaykiRegisterContractTx
	tx.PrivateKey = privateKey
	tx.TxType = commons.REG_CONT_TX
	tx.Version = 1
	tx.ValidHeight = params.ValidHeight
	tx.UserId = commons.NewRegUidByStr(params.SrcRegId)
	tx.Script = params.Script
	tx.Description = params.Description

	tx.Fees = 110000001
	hash := tx.SignTx()
	return hash
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
