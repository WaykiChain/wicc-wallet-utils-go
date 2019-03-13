package wiccwallet

import (
	"encoding/hex"

	"github.com/WaykiChain/wicc-wallet-utils-go/commons"
	"github.com/btcsuite/btcutil"
)

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

//GetAddressFromPrivateKey get address from private key
// netType: WAYKI_TESTNET or WAYKI_MAINTNET
func GetAddressFromPrivateKey(words string, netType int) string {
	address := commons.ImportPrivateKey(words, commons.Network(netType))
	return address
}

//SignRegisterAccountTx sign for register account tx
// returns the signature hex string and nil error, or returns empty string and the error if it has error
func SignRegisterAccountTx(privateKey string, param *RegisterAccountTxParam) (string, error) {
	// check and convert params
	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return "", ERR_INVALID_PRIVATE_KEY
	}

	var tx commons.WaykiRegisterAccountTx
	if param.ValidHeight < 0 {
		return "", ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight

	if !checkMinTxFee(param.Fees) {
		return "", ERR_FEE_SMALLER_MIN
	}
	tx.Fees = uint64(param.Fees)

	tx.UserId = commons.NewPubKeyUid(wifKey.PrivKey.PubKey().SerializeCompressed())
	tx.TxType = commons.REG_ACCT_TX
	tx.Version = 1

	hash := tx.SignTx(wifKey)
	return hash, nil
}

//SignCommonTx sign for common tx
// returns the signature hex string and nil error, or returns empty string and the error if it has error
func SignCommonTx(privateKey string, param *CommonTxParam) (string, error) {

	// check and convert params
	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return "", ERR_INVALID_PRIVATE_KEY
	}

	var tx commons.WaykiCommonTx

	if param.ValidHeight < 0 {
		return "", ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight

	tx.UserId = parseRegId(param.SrcRegId)
	if tx.UserId == nil {
		return "", ERR_INVALID_SRC_REG_ID
	}

	tx.DestId = parseUserId(param.DestAddr)
	if tx.DestId == nil {
		return "", ERR_INVALID_DEST_ADDR
	}

	if !checkMoneyRange(param.Values) {
		return "", ERR_RANGE_VALUES
	}
	tx.Values = uint64(param.Values)

	if !checkMoneyRange(param.Fees) {
		return "", ERR_RANGE_FEE
	}

	if !checkMinTxFee(param.Fees) {
		return "", ERR_FEE_SMALLER_MIN
	}
	tx.Fees = uint64(param.Fees)

	tx.TxType = commons.COMMON_TX
	tx.Version = TX_VERSION
	hash := tx.SignTx(wifKey)
	return hash, nil
}

func checkPubKey(pubKey []byte) bool {
	return len(pubKey) == 33
}

//SignDelegateTx sign for delegate tx
// returns the signature hex string and nil error, or returns empty string and the error if it has error
func SignDelegateTx(privateKey string, param *DelegateTxParam) (string, error) {

	// check and convert params
	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return "", ERR_INVALID_PRIVATE_KEY
	}

	var tx commons.WaykiDelegateTx

	if param.ValidHeight < 0 {
		return "", ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight

	tx.UserId = parseRegId(param.SrcRegId)
	if tx.UserId == nil {
		return "", ERR_INVALID_SRC_REG_ID
	}

	if !checkMoneyRange(param.Fees) {
		return "", ERR_RANGE_FEE
	}

	if !checkMinTxFee(param.Fees) {
		return "", ERR_FEE_SMALLER_MIN
	}

	if len(param.Votes.voteArray) == 0 {
		return "", ERR_EMPTY_VOTES
	}
	var voteData []commons.OperVoteFund
	for i := 0; i < len(param.Votes.voteArray); i++ {
		inVote := param.Votes.voteArray[i]
		var v commons.OperVoteFund
		if !checkPubKey(inVote.PubKey) {
			return "", ERR_INVALID_VOTE_PUBKEY
		}

		v.VoteValue = abs(inVote.VoteValue)
		if !checkMoneyRange(v.VoteValue) {
			return "", ERR_RANGE_VOTE_VALUE
		}
		v.VoteType = commons.GetVoteTypeByValue(inVote.VoteValue)

		var pubKey commons.PubKeyId = inVote.PubKey
		v.PubKey = &pubKey
		voteData = append(voteData, v)
	}

	tx.TxType = commons.DELEGATE_TX
	tx.Version = 1
	tx.OperVoteFunds = voteData
	hash := tx.SignTx(wifKey)
	return hash, nil
}

var ()

//SignCallContractTx sign for call contract tx
// returns the signature hex string and nil error, or returns empty string and the error if it has error
func SignCallContractTx(privateKey string, param *CallContractTxParam) (string, error) {
	// check and convert params
	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return "", ERR_INVALID_PRIVATE_KEY
	}

	var tx commons.WaykiCallContractTx

	if param.ValidHeight < 0 {
		return "", ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight

	tx.UserId = parseRegId(param.SrcRegId)
	if tx.UserId == nil {
		return "", ERR_INVALID_SRC_REG_ID
	}

	tx.AppId = parseRegId(param.AppId)
	if tx.AppId == nil {
		return "", ERR_INVALID_APP_ID
	}

	if !checkMoneyRange(param.Values) {
		return "", ERR_RANGE_VALUES
	}
	tx.Values = uint64(param.Values)

	if !checkMoneyRange(param.Fees) {
		return "", ERR_RANGE_FEE
	}
	if !checkMinTxFee(param.Fees) {
		return "", ERR_FEE_SMALLER_MIN
	}
	tx.Fees = uint64(param.Fees)

	binary, errHex := hex.DecodeString(param.ContractHex)
	if errHex != nil {
		return "", ERR_INVALID_CONTRACT_HEX
	}
	tx.Contract = []byte(binary)

	tx.TxType = commons.CONTRACT_TX
	tx.Version = TX_VERSION
	hash := tx.SignTx(wifKey)
	return hash, nil
}

//SignRegisterContractTx sign for call register contract tx
// returns the signature hex string and nil error, or returns empty string and the error if it has error
func SignRegisterContractTx(privateKey string, param *RegisterContractTxParam) (string, error) {

	// check and convert params
	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return "", ERR_INVALID_PRIVATE_KEY
	}

	var tx commons.WaykiRegisterContractTx

	if param.ValidHeight < 0 {
		return "", ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight

	tx.UserId = parseRegId(param.SrcRegId)
	if tx.UserId == nil {
		return "", ERR_INVALID_SRC_REG_ID
	}

	if !checkMoneyRange(param.Fees) {
		return "", ERR_RANGE_FEE
	}
	if !checkMinTxFee(param.Fees) {
		return "", ERR_FEE_SMALLER_MIN
	}
	tx.Fees = uint64(param.Fees)

	tx.TxType = commons.REG_CONT_TX
	tx.Version = TX_VERSION

	if len(param.Script) == 0 || len(param.Script) > CONTRACT_SCRIPT_MAX_SIZE {
		return "", ERR_INVALID_SCRIPT

	}
	tx.Script = param.Script

	if len(param.Description) > CONTRACT_SCRIPT_DESC_MAX_SIZE {
		return "", ERR_INVALID_SCRIPT_DESC
	}
	tx.Description = param.Description

	hash := tx.SignTx(wifKey)
	return hash, nil
}
