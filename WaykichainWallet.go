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
	address := commons.GetAddressFromMnemonic(words, commons.Network(netType))
	return address
}

//助记词转私钥
// netType: WAYKI_TESTNET or WAYKI_MAINTNET
func GetPrivateKeyFromMnemonic(words string, netType int) string {
	privateKey := commons.GetPrivateKeyFromMnemonic(words, commons.Network(netType))
	return privateKey
}

func GetPubKeyFromPrivateKey(privKey string) (string,error) {
	wifKey, err := btcutil.DecodeWIF(privKey)
	if err != nil {
		return "", ERR_INVALID_PRIVATE_KEY
	}
	pubHex:=hex.EncodeToString(wifKey.SerializePubKey())
	return pubHex,nil
}

func checkPubKey(pubKey []byte) bool {
	return len(pubKey) == 33
}

//GetAddressFromPrivateKey get address from private key
// netType: WAYKI_TESTNET or WAYKI_MAINTNET
func GetAddressFromPrivateKey(privateKey string, netType int) string {
	address := commons.GetAddressFromPrivateKey(privateKey, commons.Network(netType))
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

	pubKey, err := hex.DecodeString(param.PubKey)
	if (err != nil) {
		return "", ERR_USER_PUBLICKEY
	}
	tx.PubKey = pubKey
	if (tx.UserId == nil && tx.PubKey == nil) {
		return "", ERR_INVALID_SRC_REG_ID
	}

	tx.TxType = commons.COMMON_TX
	tx.Version = TX_VERSION
	hash := tx.SignTx(wifKey)
	return hash, nil
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

	pubKey, err := hex.DecodeString(param.PubKey)
	if (err != nil) {
		return "", ERR_USER_PUBLICKEY
	}
	tx.PubKey = pubKey
	if (tx.UserId == nil && tx.PubKey == nil) {
		return "", ERR_INVALID_SRC_REG_ID
	}
	tx.TxType = commons.DELEGATE_TX
	tx.Version = 1
	tx.OperVoteFunds = voteData
	hash := tx.SignTx(wifKey)
	return hash, nil
}

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
	pubKey, err := hex.DecodeString(param.PubKey)
	if (err != nil) {
		return "", ERR_USER_PUBLICKEY
	}
	tx.PubKey = pubKey
	if (tx.UserId == nil && tx.PubKey == nil) {
		return "", ERR_INVALID_SRC_REG_ID
	}
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
	pubKey, err := hex.DecodeString(param.PubKey)
	if (err != nil) {
		return "", ERR_USER_PUBLICKEY
	}
	tx.PubKey = pubKey
	if (tx.UserId == nil && tx.PubKey == nil) {
		return "", ERR_INVALID_SRC_REG_ID
	}
	hash := tx.SignTx(wifKey)
	return hash, nil
}

func SignUCoinTransferTx(privateKey string, param *UCoinTransferTxParam) (string, error) {

	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return "", ERR_INVALID_PRIVATE_KEY
	}
	var tx commons.WaykiUCoinTransferTx

	if param.ValidHeight < 0 {
		return "", ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight

	tx.UserId = parseRegId(param.SrcRegId)

	pubKey, err := hex.DecodeString(param.PubKey)
	if (err != nil) {
		return "", ERR_USER_PUBLICKEY
	}
	tx.PubKey = pubKey
	if (tx.UserId == nil && tx.PubKey == nil) {
		return "", ERR_INVALID_SRC_REG_ID
	}
	if !checkMoneyRange(param.Fees) {
		return "", ERR_RANGE_FEE
	}
	tx.Fees = uint64(param.Fees)

	tx.TxType = commons.UCOIN_TRANSFER_TX
	tx.Version = TX_VERSION
	if param.CoinAmount < 0 {
		return "", ERR_RANGE_VALUES
	}

	if param.CoinSymbol == "" {
		return "", ERR_COIN_TYPE
	}
	if (param.FeeSymbol == "") {
		tx.FeeSymbol = string(commons.WICC)
	} else {
		tx.FeeSymbol = string(param.FeeSymbol)
	}
	tx.CoinSymbol = string(param.CoinSymbol)
	tx.CoinAmount = param.CoinAmount
	tx.Memo = param.Memo
	hash := tx.SignTx(wifKey)
	return hash, nil
}

func SignCdpStakeTx(privateKey string, param *CdpStakeTxParam) (string, error) {

	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return "", ERR_INVALID_PRIVATE_KEY
	}
	var tx commons.WaykiCdpStakeTx

	if param.ValidHeight < 0 {
		return "", ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight

	tx.UserId = parseRegId(param.SrcRegId)

	pubKey, err := hex.DecodeString(param.PubKey)
	if (err != nil) {
		return "", ERR_USER_PUBLICKEY
	}
	tx.PubKey = pubKey
	if (tx.UserId == nil && tx.PubKey == nil) {
		return "", ERR_INVALID_SRC_REG_ID
	}
	if !checkCdpMinTxFee(param.Fees) {
		return "", ERR_FEE_SMALLER_MIN
	}
	tx.Fees = uint64(param.Fees)

	tx.TxType = commons.CDP_STAKE_TX
	tx.Version = TX_VERSION
	if param.BcoinStake < 0 || param.ScoinMint < 0 {
		return "", ERR_CDP_STAKE_NUMBER
	}
	tx.BcoinValues = param.BcoinStake
	tx.ScoinValues = param.ScoinMint

	if param.BcoinSymbol == "" || param.ScoinSymbol == "" {
		return "", ERR_COIN_TYPE
	}
	if (param.FeeSymbol == "") {
		tx.FeeSymbol = string(commons.WICC)
	} else {
		tx.FeeSymbol = string(param.FeeSymbol)
	}
	tx.BcoinSymbol = string(param.BcoinSymbol)
	tx.ScoinSymbol = string(param.ScoinSymbol)
	if(param.CdpTxid==""){
      param.CdpTxid="0000000000000000000000000000000000000000000000000000000000000000"
	}
	txHash, err := hex.DecodeString(param.CdpTxid)
	if (err != nil) {
		return "", ERR_CDP_TX_HASH
	}
	tx.CdpTxHash = txHash

	hash := tx.SignTx(wifKey)
	return hash, nil
}

func SignCdpRedeemTx(privateKey string, param *CdpRedeemTxParam) (string, error) {

	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return "", ERR_INVALID_PRIVATE_KEY
	}
	var tx commons.WaykiCdpRedeemTx

	if param.ValidHeight < 0 {
		return "", ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight

	tx.UserId = parseRegId(param.SrcRegId)

	pubKey, err := hex.DecodeString(param.PubKey)
	if (err != nil) {
		return "", ERR_USER_PUBLICKEY
	}
	tx.PubKey = pubKey
	if (tx.UserId == nil && tx.PubKey == nil) {
		return "", ERR_INVALID_SRC_REG_ID
	}
	if !checkCdpMinTxFee(param.Fees) {
		return "", ERR_FEE_SMALLER_MIN
	}
	tx.Fees = uint64(param.Fees)

	tx.TxType = commons.CDP_REDEEMP_TX
	tx.Version = TX_VERSION
	if param.BcoinsToRedeem < 0 || param.ScoinsToRepay < 0 {
		return "", ERR_CDP_STAKE_NUMBER
	}
	tx.ScoinValues = param.ScoinsToRepay
	tx.BcoinValues = param.BcoinsToRedeem

	if (param.FeeSymbol == "") {
		tx.FeeSymbol = string(commons.WICC)
	} else {
		tx.FeeSymbol = string(param.FeeSymbol)
	}

	txHash, err := hex.DecodeString(param.CdpTxid)
	if (err != nil) {
		return "", ERR_CDP_TX_HASH
	}
	tx.CdpTxHash = txHash

	hash := tx.SignTx(wifKey)
	return hash, nil
}

func SignCdpLiquidateTx(privateKey string, param *CdpLiquidateTxParam) (string, error) {

	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return "", ERR_INVALID_PRIVATE_KEY
	}
	var tx commons.WaykiCdpLiquidateTx

	if param.ValidHeight < 0 {
		return "", ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight

	tx.UserId = parseRegId(param.SrcRegId)

	pubKey, err := hex.DecodeString(param.PubKey)
	if (err != nil) {
		return "", ERR_USER_PUBLICKEY
	}
	tx.PubKey = pubKey
	if (tx.UserId == nil && tx.PubKey == nil) {
		return "", ERR_INVALID_SRC_REG_ID
	}
	if !checkCdpMinTxFee(param.Fees) {
		return "", ERR_FEE_SMALLER_MIN
	}
	tx.Fees = uint64(param.Fees)

	tx.TxType = commons.CDP_LIQUIDATE_TX
	tx.Version = TX_VERSION
	if param.ScoinsLiquidate < 0 {
		return "", ERR_CDP_STAKE_NUMBER
	}
	tx.ScoinsLiquidate = param.ScoinsLiquidate

	if (param.FeeSymbol == "") {
		tx.FeeSymbol = string(commons.WICC)
	} else {
		tx.FeeSymbol = string(param.FeeSymbol)
	}

	txHash, err := hex.DecodeString(param.CdpTxid)
	if (err != nil) {
		return "", ERR_CDP_TX_HASH
	}
	tx.CdpTxHash = txHash

	hash := tx.SignTx(wifKey)
	return hash, nil
}

func SignDexSellLimitTx(privateKey string, param *DexLimitTxParam) (string, error) {
	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return "", ERR_INVALID_PRIVATE_KEY
	}
	var tx commons.WaykiDexSellLimitTx

	if param.ValidHeight < 0 {
		return "", ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight

	tx.UserId = parseRegId(param.SrcRegId)

	pubKey, err := hex.DecodeString(param.PubKey)
	if (err != nil) {
		return "", ERR_USER_PUBLICKEY
	}
	tx.PubKey = pubKey
	if (tx.UserId == nil && tx.PubKey == nil) {
		return "", ERR_INVALID_SRC_REG_ID
	}
	if !checkCdpMinTxFee(param.Fees) {
		return "", ERR_FEE_SMALLER_MIN
	}
	tx.Fees = uint64(param.Fees)

	tx.TxType = commons.DEX_SELL_LIMIT_ORDER_TX
	tx.Version = TX_VERSION
	if param.AskPrice <= 0 {
		return "", ERR_ASK_PRICE
	}
	tx.AskPrice = uint64(param.AskPrice)
	if param.CoinSymbol == "" || param.AssetSymbol == "" {
		return "", ERR_COIN_TYPE
	}
	tx.AssetSymbol=param.AssetSymbol
	tx.CoinSymbol=param.CoinSymbol
	if (param.FeeSymbol == "") {
		tx.FeeSymbol = string(commons.WICC)
	} else {
		tx.FeeSymbol = string(param.FeeSymbol)
	}
	if !checkMoneyRange(param.AssetAmount) {
		return "", ERR_RANGE_VALUES
	}
	tx.AssetAmount = uint64(param.AssetAmount)
	hash := tx.SignTx(wifKey)
	return hash, nil
}

func SignDexMarketSellTx(privateKey string, param *DexMarketTxParam) (string, error) {
	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return "", ERR_INVALID_PRIVATE_KEY
	}
	var tx commons.WaykiDexMarketTx

	if param.ValidHeight < 0 {
		return "", ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight

	tx.UserId = parseRegId(param.SrcRegId)

	pubKey, err := hex.DecodeString(param.PubKey)
	if (err != nil) {
		return "", ERR_USER_PUBLICKEY
	}
	tx.PubKey = pubKey
	if (tx.UserId == nil && tx.PubKey == nil) {
		return "", ERR_INVALID_SRC_REG_ID
	}
	if !checkCdpMinTxFee(param.Fees) {
		return "", ERR_FEE_SMALLER_MIN
	}
	tx.Fees = uint64(param.Fees)

	tx.TxType = commons.DEX_SELL_MARKET_ORDER_TX
	tx.Version = TX_VERSION
	if param.CoinSymbol == "" || param.AssetSymbol == "" {
		return "", ERR_COIN_TYPE
	}
	tx.AssetSymbol=param.AssetSymbol
	tx.CoinSymbol=param.CoinSymbol
	if (param.FeeSymbol == "") {
		tx.FeeSymbol = string(commons.WICC)
	} else {
		tx.FeeSymbol = string(param.FeeSymbol)
	}
	if !checkMoneyRange(param.AssetAmount) {
		return "", ERR_RANGE_VALUES
	}
	tx.AssetAmount = uint64(param.AssetAmount)
	hash := tx.SignTx(wifKey)
	return hash, nil
}

func SignDexBuyLimitTx(privateKey string, param *DexLimitTxParam) (string, error) {
	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return "", ERR_INVALID_PRIVATE_KEY
	}
	var tx commons.WaykiDexSellLimitTx

	if param.ValidHeight < 0 {
		return "", ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight

	tx.UserId = parseRegId(param.SrcRegId)

	pubKey, err := hex.DecodeString(param.PubKey)
	if (err != nil) {
		return "", ERR_USER_PUBLICKEY
	}
	tx.PubKey = pubKey
	if (tx.UserId == nil && tx.PubKey == nil) {
		return "", ERR_INVALID_SRC_REG_ID
	}
	if !checkCdpMinTxFee(param.Fees) {
		return "", ERR_FEE_SMALLER_MIN
	}
	tx.Fees = uint64(param.Fees)

	tx.TxType = commons.DEX_BUY_LIMIT_ORDER_TX
	tx.Version = TX_VERSION
	if param.AskPrice <= 0 {
		return "", ERR_ASK_PRICE
	}
	tx.AskPrice = uint64(param.AskPrice)
	if param.CoinSymbol == "" || param.AssetSymbol == "" {
		return "", ERR_COIN_TYPE
	}
	tx.AssetSymbol=param.AssetSymbol
	tx.CoinSymbol=param.CoinSymbol
	if (param.FeeSymbol == "") {
		tx.FeeSymbol = string(commons.WICC)
	} else {
		tx.FeeSymbol = string(param.FeeSymbol)
	}
	if !checkMoneyRange(param.AssetAmount) {
		return "", ERR_RANGE_VALUES
	}
	tx.AssetAmount = uint64(param.AssetAmount)
	hash := tx.SignTx(wifKey)
	return hash, nil
}

func SignDexMarketBuyTx(privateKey string, param *DexMarketTxParam) (string, error) {
	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return "", ERR_INVALID_PRIVATE_KEY
	}
	var tx commons.WaykiDexMarketTx

	if param.ValidHeight < 0 {
		return "", ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight

	tx.UserId = parseRegId(param.SrcRegId)

	pubKey, err := hex.DecodeString(param.PubKey)
	if (err != nil) {
		return "", ERR_USER_PUBLICKEY
	}
	tx.PubKey = pubKey
	if (tx.UserId == nil && tx.PubKey == nil) {
		return "", ERR_INVALID_SRC_REG_ID
	}
	if !checkCdpMinTxFee(param.Fees) {
		return "", ERR_FEE_SMALLER_MIN
	}
	tx.Fees = uint64(param.Fees)

	tx.TxType = commons.DEX_BUY_MARKET_ORDER_TX
	tx.Version = TX_VERSION
	if param.CoinSymbol == "" || param.AssetSymbol == "" {
		return "", ERR_COIN_TYPE
	}
	tx.AssetSymbol=param.AssetSymbol
	tx.CoinSymbol=param.CoinSymbol
	if (param.FeeSymbol == "") {
		tx.FeeSymbol = string(commons.WICC)
	} else {
		tx.FeeSymbol = string(param.FeeSymbol)
	}
	if !checkMoneyRange(param.AssetAmount) {
		return "", ERR_RANGE_VALUES
	}
	tx.AssetAmount = uint64(param.AssetAmount)
	hash := tx.SignTx(wifKey)
	return hash, nil
}

func SignDexCancelTx(privateKey string, param *DexCancelTxParam) (string, error) {
	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return "", ERR_INVALID_PRIVATE_KEY
	}
	var tx commons.WaykiDexCancelTx

	if param.ValidHeight < 0 {
		return "", ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight

	tx.UserId = parseRegId(param.SrcRegId)

	pubKey, err := hex.DecodeString(param.PubKey)
	if (err != nil) {
		return "", ERR_USER_PUBLICKEY
	}
	tx.PubKey = pubKey
	if (tx.UserId == nil && tx.PubKey == nil) {
		return "", ERR_INVALID_SRC_REG_ID
	}
	if !checkCdpMinTxFee(param.Fees) {
		return "", ERR_FEE_SMALLER_MIN
	}
	tx.Fees = uint64(param.Fees)

	tx.TxType = commons.DEX_CANCEL_ORDER_TX
	tx.Version = TX_VERSION
	if (param.FeeSymbol == "") {
		tx.FeeSymbol = string(commons.WICC)
	} else {
		tx.FeeSymbol = string(param.FeeSymbol)
	}
	txHash, err := hex.DecodeString(param.DexTxid)
	if (err != nil) {
		return "", ERR_CDP_TX_HASH
	}
	tx.DexHash = txHash
	hash := tx.SignTx(wifKey)
	return hash, nil
}