package wiccwallet

import (
	"encoding/hex"
	hash2 "github.com/WaykiChain/wicc-wallet-utils-go/commons/hash"
	"github.com/WaykiChain/wicc-wallet-utils-go/commons"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
	"strings"
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

func CheckWalletAddress(address string, netType int) (bool, error) {
	versionAndDataBytes, version, error := base58.CheckDecode(address)
	if (len(versionAndDataBytes) < 1 || error != nil) {
		return false, ERR_INVALID_DEST_ADDR
	}
	//version:=versionAndDataBytes[0] & 0xFF
	netParams, err := commons.NetworkToChainConfig(commons.Network(netType))
	if (err != nil) {
		return false, ERR_INVALID_NETWORK
	}
	if (netParams.PubKeyHashAddrID == version) {
		return true, nil
	} else {
		return false, nil
	}
}

func CheckPrivateKey(privateKey string, netType int) (bool, error) {
	_, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return false, ERR_INVALID_PRIVATE_KEY
	}
	versionAndDataBytes := base58.Decode(privateKey)
	if (len(versionAndDataBytes) < 1) {
		return false, ERR_INVALID_PRIVATE_KEY
	}
	version := versionAndDataBytes[0] & 0xFF
	netParams, err := commons.NetworkToChainConfig(commons.Network(netType))
	if (err != nil) {
		return false, ERR_INVALID_NETWORK
	}
	if (netParams.PrivateKeyID == version) {
		return true, nil
	} else {
		return false, nil
	}
}

//助记词转换地址
// netType: WAYKI_TESTNET or WAYKI_MAINTNET
func GetAddressFromMnemonic(words string, netType int) (string, error) {
	wordArr:=strings.Split(words," ")
	if(len(wordArr)!=12){
		return "", ERR_INVALID_MNEMONIC
	}
	address := commons.GetAddressFromMnemonic(words, commons.Network(netType))
	return address,nil
}

//助记词转私钥
// netType: WAYKI_TESTNET or WAYKI_MAINTNET
func GetPrivateKeyFromMnemonic(words string, netType int) (string, error) {
	wordArr:=strings.Split(words," ")
	if(len(wordArr)!=12){
		return "", ERR_INVALID_MNEMONIC
	}
	privateKey := commons.GetPrivateKeyFromMnemonic(words, commons.Network(netType))
	return privateKey,nil
}

// get publickey from privatekey
func GetPubKeyFromPrivateKey(privKey string) (string, error) {
	wifKey, err := btcutil.DecodeWIF(privKey)
	if err != nil {
		return "", ERR_INVALID_PRIVATE_KEY
	}
	pubHex := hex.EncodeToString(wifKey.SerializePubKey())
	return pubHex, nil
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
	tx.Memo = param.Memo
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

//SignUCoinCallContractTx sign for call contract tx
// returns the signature hex string and nil error, or returns empty string and the error if it has error
func SignUCoinCallContractTx(privateKey string, param *UCoinContractTxParam) (string, error) {
	// check and convert params
	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return "", ERR_INVALID_PRIVATE_KEY
	}

	var tx commons.WaykiUCoinCallContractTx

	if param.ValidHeight < 0 {
		return "", ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight
	tx.UserId = parseRegId(param.SrcRegId)
	tx.AppId = parseRegId(param.AppId)
	if tx.AppId == nil {
		return "", ERR_INVALID_APP_ID
	}
	if !checkMoneyRange(param.CoinAmount) {
		return "", ERR_RANGE_VALUES
	}
	tx.CoinAmount = int64(param.CoinAmount)
	if !checkMoneyRange(param.Fees) {
		return "", ERR_RANGE_FEE
	}
	if !checkMinTxFee(param.Fees) {
		return "", ERR_FEE_SMALLER_MIN
	}
	tx.Fees = int64(param.Fees)
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
	if param.CoinSymbol == "" {
		return "", ERR_COIN_TYPE
	}
	tx.CoinSymbol = param.CoinSymbol
	if (param.FeeSymbol == "") {
		tx.FeeSymbol = string(commons.WICC)
	} else {
		tx.FeeSymbol = string(param.FeeSymbol)
	}
	tx.TxType = commons.UCOIN_CONTRACT_INVOKE_TX
	tx.Version = TX_VERSION
	hash := tx.SignTx(wifKey)
	return hash, nil
}

//Sign message by private Key
func SignMessage(privateKey string, message string) (*SignMessageParam, error) {
	hash := hash2.Hash256(hash2.Hash160([]byte(message)))
	wifKey, errorDecode := btcutil.DecodeWIF(privateKey)
	if (errorDecode != nil) {
		return nil, ERR_INVALID_PRIVATE_KEY
	}
	key := wifKey.PrivKey
	signMsg, errorSign := key.Sign(hash)
	if (errorSign != nil) {
		return nil, ERR_SIGNATURE_ERROR
	}
	sign := SignMessageParam{hex.EncodeToString(wifKey.SerializePubKey()),
		hex.EncodeToString(signMsg.Serialize())}
	return &sign, nil
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
	if (tx.UserId == nil && tx.PubKey == nil) {
		return "", ERR_INVALID_SRC_REG_ID
	}
	hash := tx.SignTx(wifKey)
	return hash, nil
}

//SignUCoinRegisterContractTx sign for call register contract tx
// returns the signature hex string and nil error, or returns empty string and the error if it has error
func SignUCoinRegisterContractTx(privateKey string, param *UCoinRegisterContractTxParam) (string, error) {

	// check and convert params
	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return "", ERR_INVALID_PRIVATE_KEY
	}

	var tx commons.WaykiUCoinRegisterContractTx

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

	tx.TxType = commons.UCONTRACT_DEPLOY_TX
	tx.Version = TX_VERSION

	if len(param.Script) == 0 || len(param.Script) > CONTRACT_SCRIPT_MAX_SIZE {
		return "", ERR_INVALID_SCRIPT
	}
	tx.Script = param.Script

	if len(param.Description) > CONTRACT_SCRIPT_DESC_MAX_SIZE {
		return "", ERR_INVALID_SCRIPT_DESC
	}
	tx.Description = param.Description
	if (tx.UserId == nil && tx.PubKey == nil) {
		return "", ERR_INVALID_SRC_REG_ID
	}

	if (param.FeeSymbol == "") {
		tx.FeeSymbol = string(commons.WICC)
	} else {
		tx.FeeSymbol = string(param.FeeSymbol)
	}

	hash := tx.SignTx(wifKey)
	return hash, nil
}

//SignUCoinTransferTx sign for Multi-currency transfer
// returns the signature hex string and nil error, or returns empty string and the error if it has error
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

	var dests []commons.Dest
	for i:=0;i< len(param.Dests.destArray);i++{
		var dest commons.Dest
		dest.DestAddr = parseUserId(param.Dests.destArray[i].DestAddr)
		if param.Dests.destArray[i].CoinAmount < 0 {
			return "", ERR_RANGE_VALUES
		}
		dest.CoinAmount = uint64(param.Dests.destArray[i].CoinAmount)
		if param.Dests.destArray[i].CoinSymbol == "" {
			return "", ERR_COIN_TYPE
		}
		dest.CoinSymbol = string(param.Dests.destArray[i].CoinSymbol)
		dests=append(dests,dest)
	}
	tx.Dests=dests
	if (param.FeeSymbol == "") {
		tx.FeeSymbol = string(commons.WICC)
	} else {
		tx.FeeSymbol = string(param.FeeSymbol)
	}
	tx.Memo = param.Memo
	hash := tx.SignTx(wifKey)
	return hash, nil
}

//SignCdpStakeTx sign for create a cdp tx
// returns the signature hex string and nil error, or returns empty string and the error if it has error
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
	if  param.ScoinMint <= 0 {
		return "", ERR_CDP_STAKE_NUMBER
	}

	tx.ScoinValues = uint64(param.ScoinMint)

	if  param.ScoinSymbol == "" {
		return "", ERR_COIN_TYPE
	}
	if (param.FeeSymbol == "") {
		tx.FeeSymbol = string(commons.WICC)
	} else {
		tx.FeeSymbol = string(param.FeeSymbol)
	}
	var models []commons.AssetModel
	for i := 0; i < len(param.Assets.assetArray); i++ {
		asset := param.Assets.assetArray[i]
		var model commons.AssetModel
		model.AssetSymbol=asset.AssetSymbol
		if asset.AssetSymbol=="" {
			return "", ERR_COIN_TYPE
		}

		model.AssetAmount = abs(asset.AssetAmount)
		if !checkMoneyRange(model.AssetAmount) {
			return "", ERR_CDP_STAKE_NUMBER
		}
		models = append(models, model)
	}
	tx.Assets=models
	tx.ScoinSymbol = string(param.ScoinSymbol)
	if (param.CdpTxid == "") {
		param.CdpTxid = "0000000000000000000000000000000000000000000000000000000000000000"
	}
	txHash, err := hex.DecodeString(param.CdpTxid)
	if (err != nil) {
		return "", ERR_CDP_TX_HASH
	}
	tx.CdpTxHash = txHash

	hash := tx.SignTx(wifKey)
	return hash, nil
}

//SignCdpRedeemTx sign for redeem a cdp tx
// returns the signature hex string and nil error, or returns empty string and the error if it has error
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
	if param.ScoinsToRepay < 0 {
		return "", ERR_CDP_STAKE_NUMBER
	}
	tx.ScoinValues = uint64(param.ScoinsToRepay)

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

	var models []commons.AssetModel
	for i := 0; i < len(param.Assets.assetArray); i++ {
		asset := param.Assets.assetArray[i]
		var model commons.AssetModel
		model.AssetSymbol=asset.AssetSymbol
		if asset.AssetSymbol=="" {
			return "", ERR_COIN_TYPE
		}

		model.AssetAmount = abs(asset.AssetAmount)
		if !checkMoneyRange(model.AssetAmount) {
			return "", ERR_CDP_STAKE_NUMBER
		}
		models = append(models, model)
	}
	tx.Assets=models

	hash := tx.SignTx(wifKey)
	return hash, nil
}

//SignCdpLiquidateTx sign for liquidate a cdp tx
// returns the signature hex string and nil error, or returns empty string and the error if it has error
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
	tx.ScoinsLiquidate = uint64(param.ScoinsLiquidate)

	if (param.FeeSymbol == "") {
		tx.FeeSymbol = string(commons.WICC)
	} else {
		tx.FeeSymbol = string(param.FeeSymbol)
	}

	if param.AssetSymbol==""{
		return "",ERR_COIN_TYPE
	}
	tx.AseetSymbol=param.AssetSymbol

	txHash, err := hex.DecodeString(param.CdpTxid)
	if (err != nil) {
		return "", ERR_CDP_TX_HASH
	}
	tx.CdpTxHash = txHash

	hash := tx.SignTx(wifKey)
	return hash, nil
}

//SignDexSellLimitTx sign for dex sell limit price tx
// returns the signature hex string and nil error, or returns empty string and the error if it has error
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
	if !checkMinTxFee(param.Fees) {
		return "", ERR_FEE_SMALLER_MIN
	}
	tx.Fees = uint64(param.Fees)

	tx.TxType = commons.DEX_SELL_LIMIT_ORDER_TX
	tx.Version = TX_VERSION
	if param.Price <= 0 {
		return "", ERR_ASK_PRICE
	}
	tx.AskPrice = uint64(param.Price)
	if param.CoinSymbol == "" || param.AssetSymbol == "" {
		return "", ERR_COIN_TYPE
	}
	tx.AssetSymbol = param.AssetSymbol
	tx.CoinSymbol = param.CoinSymbol
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

//SignDexMarketSellTx sign for dex sell market price tx
// returns the signature hex string and nil error, or returns empty string and the error if it has error
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
	tx.AssetSymbol = param.AssetSymbol
	tx.CoinSymbol = param.CoinSymbol
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

//SignDexBuyLimitTx sign for dex buy limit Price tx
// returns the signature hex string and nil error, or returns empty string and the error if it has error
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
	if param.Price <= 0 {
		return "", ERR_ASK_PRICE
	}
	tx.AskPrice = uint64(param.Price)
	if param.CoinSymbol == "" || param.AssetSymbol == "" {
		return "", ERR_COIN_TYPE
	}
	tx.AssetSymbol = param.AssetSymbol
	tx.CoinSymbol = param.CoinSymbol
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

//SignDexMarketBuyTx sign for dex buy market price tx
// returns the signature hex string and nil error, or returns empty string and the error if it has error
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
	tx.AssetSymbol = param.AssetSymbol
	tx.CoinSymbol = param.CoinSymbol
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

//SignDexMarketBuyTx sign for cancel dex tx
// returns the signature hex string and nil error, or returns empty string and the error if it has error
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

//SignAssetIssueTx sign
// returns the signature hex string and nil error, or returns empty string and the error if it has error
func SignAssetCreateTx(privateKey string, param *AssetIssueTxParam) (string, error) {

	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return "", ERR_INVALID_PRIVATE_KEY
	}
	var tx commons.WaykiAssetIssueTx

	if param.ValidHeight < 0 {
		return "", ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight
	tx.UserId = parseRegId(param.SrcRegId)

	if (tx.UserId == nil) {
		return "", ERR_INVALID_SRC_REG_ID
	}

	tx.Fees = uint64(param.Fees)

	tx.TxType = commons.ASSET_ISSUE_TX
	tx.Version = TX_VERSION

	if !checkAssetSymbol(param.AssetSymbol) {
		return "", ERR_SYMBOL_ERROR
	}
	tx.AssetSymbol = param.AssetSymbol
	if (param.FeeSymbol == "") {
		tx.FeeSymbol = string(commons.WICC)
	} else {
		tx.FeeSymbol = string(param.FeeSymbol)
	}
	tx.MinTable = param.MinTable
	if (param.AssetName == "") {
		return "", ERR_ASSET_NAME_ERROR
	}
	tx.AssetName = param.AssetName
	if (param.AssetTotal < 100000000) {
		return "", ERR_TOTAl_SUPPLY_ERROR
	}
	tx.AssetTotal = uint64(param.AssetTotal)
	tx.AssetOwner = parseUserId(param.AssetOwner)
	if (tx.AssetOwner == nil) {
		return "", ERR_ASSET_UPDATE_OWNER_ERROR
	}
	hash := tx.SignTx(wifKey)
	return hash, nil
}

//SignAssetUpdateTx sign
// returns the signature hex string and nil error, or returns empty string and the error if it has error
func SignAssetUpdateTx(privateKey string, param *AssetUpdateTxParam) (string, error) {

	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return "", ERR_INVALID_PRIVATE_KEY
	}
	var tx commons.WaykiAssetUpdateTx

	if param.ValidHeight < 0 {
		return "", ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight
	tx.UserId = parseRegId(param.SrcRegId)
	if (tx.UserId == nil) {
		return "", ERR_INVALID_SRC_REG_ID
	}

	tx.Fees = uint64(param.Fees)

	tx.TxType = commons.ASSET_UPDATE_TX
	tx.Version = TX_VERSION

	if !checkAssetSymbol(param.AssetSymbol) {
		return "", ERR_SYMBOL_ERROR
	}
	tx.AssetSymbol = param.AssetSymbol
	if (param.FeeSymbol == "") {
		tx.FeeSymbol = string(commons.WICC)
	} else {
		tx.FeeSymbol = string(param.FeeSymbol)
	}
    tx.UpdateType=param.UpdateType
	switch param.UpdateType {
	case 1:
		tx.AssetOwner = parseUserId(param.AssetOwner)
		if (tx.AssetOwner == nil) {
			return "", ERR_ASSET_UPDATE_OWNER_ERROR
		}
		break
	case 2:
		if (param.AssetName == "") {
			return "", ERR_ASSET_NAME_ERROR
		}
		tx.AssetName = param.AssetName
		break
	case 3:
		if (param.AssetTotal < 100000000) {
			return "", ERR_TOTAl_SUPPLY_ERROR
		}
		tx.AssetTotal = uint64(param.AssetTotal)
		break
	default:
		return "", ERR_ASSET_UPDATE_TYPE_ERROR
	}

	hash := tx.SignTx(wifKey)
	return hash, nil
}
