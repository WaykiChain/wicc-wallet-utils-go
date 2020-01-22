package wicc_wallet_utils_go

import (
	"encoding/hex"
	"github.com/WaykiChain/wicc-wallet-utils-go/common"
	"github.com/WaykiChain/wicc-wallet-utils-go/waykichain"
	"github.com/btcsuite/btcutil"
)


//CreateUCoinTransferRawTx sign for Multi-currency transfer
// returns the signature hex string and nil error, or returns empty string and the error if it has error
func (param *UCoinTransferTxParam) CreateRawTx(privateKey string) (*CreateRawTxResult,error) {

	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return nil, common.ERR_INVALID_PRIVATEKEY
	}
	var tx waykichain.WaykiUCoinTransferTx

	if param.ValidHeight < 0 {
		return nil, common.ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight

	tx.UserId = parseRegId(param.SrcRegId)

	pubKey, err := hex.DecodeString(param.PubKey)
	if (err != nil) {
		return nil, common.ERR_USER_PUBLICKEY
	}
	tx.PubKey = pubKey
	if (tx.UserId == nil && tx.PubKey == nil) {
		return nil, common.ERR_INVALID_SRC_REG_ID
	}
	if !checkMoneyRange(param.Fees) {
		return nil, common.ERR_RANGE_FEE
	}
	tx.Fees = uint64(param.Fees)

	tx.TxType = waykichain.UCOIN_TRANSFER_TX
	tx.Version = TX_VERSION

	var dests []waykichain.UCoinTransferDest
	for i:=0;i< len(param.Dests.DestArray);i++{
		var dest waykichain.UCoinTransferDest
		dest.DestAddr = parseUserId(param.Dests.DestArray[i].DestAddr)
		if param.Dests.DestArray[i].CoinAmount < 0 {
			return nil, common.ERR_RANGE_VALUES
		}
		dest.CoinAmount = uint64(param.Dests.DestArray[i].CoinAmount)
		if param.Dests.DestArray[i].CoinSymbol == "" {
			return nil, common.ERR_COIN_TYPE
		}
		dest.CoinSymbol = string(param.Dests.DestArray[i].CoinSymbol)
		dests=append(dests,dest)
	}
	tx.Dests=dests
	if (param.FeeSymbol == "") {
		tx.FeeSymbol = string(common.WICC)
	} else {
		tx.FeeSymbol = string(param.FeeSymbol)
	}
	tx.Memo = param.Memo

	result,err := tx.CreateRawTxid(wifKey)
	if err !=nil{
		return nil,err
	}

	return &CreateRawTxResult{result},nil
}



// returns the signature hex string and nil error, or returns empty string and the error if it has error
func (param * UCoinContractTxParam) CreateRawTx(privateKey string) (* CreateRawTxResult, error) {
	// check and convert params
	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return nil, common.ERR_INVALID_PRIVATEKEY
	}

	var tx waykichain.WaykiUCoinCallContractTx

	if param.ValidHeight < 0 {
		return nil, common.ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight
	tx.UserId = parseRegId(param.SrcRegId)
	tx.AppId = parseRegId(param.AppId)
	if tx.AppId == nil {
		return nil, common.ERR_INVALID_APP_ID
	}
	if !checkMoneyRange(param.CoinAmount) {
		return nil, common.ERR_RANGE_VALUES
	}
	tx.CoinAmount = int64(param.CoinAmount)
	if !checkMoneyRange(param.Fees) {
		return nil, common.ERR_RANGE_FEE
	}
	if !checkMinTxFee(param.Fees) {
		return nil, common.ERR_FEE_SMALLER_MIN
	}
	tx.Fees = int64(param.Fees)
	binary, errHex := hex.DecodeString(param.ContractHex)
	if errHex != nil {
		return nil, common.ERR_INVALID_CONTRACT_HEX
	}
	tx.Contract = []byte(binary)
	pubKey, err := hex.DecodeString(param.PubKey)
	if (err != nil) {
		return nil, common.ERR_USER_PUBLICKEY
	}
	tx.PubKey = pubKey
	if (tx.UserId == nil && tx.PubKey == nil) {
		return nil, common.ERR_INVALID_SRC_REG_ID
	}
	if param.CoinSymbol == "" {
		return nil, common.ERR_COIN_TYPE
	}
	tx.CoinSymbol = param.CoinSymbol
	if (param.FeeSymbol == "") {
		tx.FeeSymbol = string(common.WICC)
	} else {
		tx.FeeSymbol = string(param.FeeSymbol)
	}
	tx.TxType = waykichain.UCOIN_CONTRACT_INVOKE_TX
	tx.Version = TX_VERSION

	result,err := tx.CreateRawTxid(wifKey)
	if err !=nil{
		return nil,err
	}

	return &CreateRawTxResult{result},nil
}




//CreateUCoinDeployContractTx sign for call register contract tx
// returns the signature hex string and nil error, or returns empty string and the error if it has error
func (param *UCoinRegisterContractTxParam) CreateRawTx(privateKey string) (* CreateRawTxResult,error) {

	// check and convert params
	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return nil, common.ERR_INVALID_PRIVATEKEY
	}

	var tx waykichain.WaykiUCoinRegisterContractTx

	if param.ValidHeight < 0 {
		return nil, common.ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight

	tx.UserId = parseRegId(param.SrcRegId)

	if !checkMoneyRange(param.Fees) {
		return nil, common.ERR_RANGE_FEE
	}
	if !checkMinTxFee(param.Fees) {
		return nil ,common.ERR_FEE_SMALLER_MIN
	}
	tx.Fees = uint64(param.Fees)

	tx.TxType = waykichain.UCONTRACT_DEPLOY_TX
	tx.Version = TX_VERSION

	if len(param.Script) == 0 || len(param.Script) > CONTRACT_SCRIPT_MAX_SIZE {
		return nil,common.ERR_INVALID_SCRIPT
	}
	tx.Script = param.Script

	if len(param.Description) > CONTRACT_SCRIPT_DESC_MAX_SIZE {
		return nil, common.ERR_INVALID_SCRIPT_DESC
	}
	tx.Description = param.Description
	if (tx.UserId == nil && tx.PubKey == nil) {
		return nil, common.ERR_INVALID_SRC_REG_ID
	}

	if (param.FeeSymbol == "") {
		tx.FeeSymbol = string(common.WICC)
	} else {
		tx.FeeSymbol = string(param.FeeSymbol)
	}

	result,err := tx.CreateRawTxid(wifKey)
	if err !=nil{
		return nil,err
	}

	return &CreateRawTxResult{result},nil
}

// returns the signature hex string and nil error, or returns empty string and the error if it has error
func (param *DelegateTxParam) CreateRawTx(privateKey string) (* CreateRawTxResult, error) {

	// check and convert params
	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return nil, common.ERR_INVALID_PRIVATEKEY
	}

	var tx waykichain.WaykiDelegateTx

	if param.ValidHeight < 0 {
		return nil, common.ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight

	tx.UserId = parseRegId(param.SrcRegId)

	if !checkMoneyRange(param.Fees) {
		return nil, common.ERR_RANGE_FEE
	}

	if !checkMinTxFee(param.Fees) {
		return nil, common.ERR_FEE_SMALLER_MIN
	}

	if len(param.Votes.VoteArray) == 0 {
		return nil, common.ERR_EMPTY_VOTES
	}
	var voteData []waykichain.OperVoteFundTx
	for i := 0; i < len(param.Votes.VoteArray); i++ {
		inVote := param.Votes.VoteArray[i]
		var v waykichain.OperVoteFundTx
		if !checkPubKey(inVote.PubKey) {
			return nil, common.ERR_INVALID_VOTE_PUBKEY
		}

		v.VoteValue = abs(inVote.VoteValue)
		if !checkMoneyRange(v.VoteValue) {
			return nil, common.ERR_RANGE_VOTE_VALUE
		}
		v.VoteType = waykichain.GetVoteTypeByValue(inVote.VoteValue)

		var pubKey waykichain.PubKeyId = inVote.PubKey
		v.PubKey = &pubKey
		voteData = append(voteData, v)
	}

	pubKey, err := hex.DecodeString(param.PubKey)
	if (err != nil) {
		return nil, common.ERR_USER_PUBLICKEY
	}
	tx.PubKey = pubKey
	if (tx.UserId == nil && tx.PubKey == nil) {
		return nil, common.ERR_INVALID_SRC_REG_ID
	}
	tx.TxType = waykichain.DELEGATE_TX
	tx.Version = 1
	tx.OperVoteFunds = voteData

	result,err := tx.CreateRawTxid(wifKey)
	if err !=nil{
		return nil,err
	}

	return &CreateRawTxResult{result},nil
}


// returns the signature hex string and nil error, or returns empty string and the error if it has error
func (param *CdpStakeTxParam) CreateRawTx(privateKey string) (* CreateRawTxResult, error) {

	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return nil, common.ERR_INVALID_PRIVATEKEY
	}
	var tx waykichain.WaykiCdpStakeTx

	if param.ValidHeight < 0 {
		return nil,common.ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight

	tx.UserId = parseRegId(param.SrcRegId)

	pubKey, err := hex.DecodeString(param.PubKey)
	if (err != nil) {
		return nil, common.ERR_USER_PUBLICKEY
	}
	tx.PubKey = pubKey
	if (tx.UserId == nil && tx.PubKey == nil) {
		return nil, common.ERR_INVALID_SRC_REG_ID
	}
	if !checkCdpMinTxFee(param.Fees) {
		return nil, common.ERR_FEE_SMALLER_MIN
	}
	tx.Fees = uint64(param.Fees)

	tx.TxType = waykichain.CDP_STAKE_TX
	tx.Version = TX_VERSION
	if  param.ScoinMint < 0 {
		return nil, common.ERR_CDP_STAKE_NUMBER
	}

	tx.ScoinValues = uint64(param.ScoinMint)

	if  param.ScoinSymbol == "" {
		return nil, common.ERR_COIN_TYPE
	}
	if (param.FeeSymbol == "") {
		tx.FeeSymbol = string(common.WICC)
	} else {
		tx.FeeSymbol = string(param.FeeSymbol)
	}
	var models []waykichain.AssetModel
	for i := 0; i < len(param.Assets.AssetArray); i++ {
		asset := param.Assets.AssetArray[i]
		var model waykichain.AssetModel
		model.AssetSymbol=asset.AssetSymbol
		if asset.AssetSymbol=="" {
			return nil, common.ERR_COIN_TYPE
		}

		model.AssetAmount = abs(asset.AssetAmount)
		if !checkMoneyRange(model.AssetAmount) {
			return nil, common.ERR_CDP_STAKE_NUMBER
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
		return nil, common.ERR_CDP_TX_HASH
	}
	tx.CdpTxHash = txHash

	result,err := tx.CreateRawTxid(wifKey)
	if err !=nil{
		return nil,err
	}

	return &CreateRawTxResult{result},nil
}



// returns the signature hex string and nil error, or returns empty string and the error if it has error
func (param *CdpRedeemTxParam) CreateRawTx(privateKey string) (* CreateRawTxResult ,error) {

	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return nil,common.ERR_INVALID_PRIVATEKEY
	}
	var tx waykichain.WaykiCdpRedeemTx

	if param.ValidHeight < 0 {
		return nil,common.ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight

	tx.UserId = parseRegId(param.SrcRegId)

	pubKey, err := hex.DecodeString(param.PubKey)
	if (err != nil) {
		return nil, common.ERR_USER_PUBLICKEY
	}
	tx.PubKey = pubKey
	if (tx.UserId == nil && tx.PubKey == nil) {
		return nil, common.ERR_INVALID_SRC_REG_ID
	}
	if !checkCdpMinTxFee(param.Fees) {
		return nil, common.ERR_FEE_SMALLER_MIN
	}
	tx.Fees = uint64(param.Fees)

	tx.TxType = waykichain.CDP_REDEEMP_TX
	tx.Version = TX_VERSION
	if param.ScoinsToRepay < 0 {
		return nil, common.ERR_CDP_STAKE_NUMBER
	}
	tx.ScoinValues = uint64(param.ScoinsToRepay)

	if (param.FeeSymbol == "") {
		tx.FeeSymbol = string(common.WICC)
	} else {
		tx.FeeSymbol = string(param.FeeSymbol)
	}

	txHash, err := hex.DecodeString(param.CdpTxid)
	if (err != nil) {
		return nil, common.ERR_CDP_TX_HASH
	}
	tx.CdpTxHash = txHash

	var models []waykichain.AssetModel
	for i := 0; i < len(param.Assets.AssetArray); i++ {
		asset := param.Assets.AssetArray[i]
		var model waykichain.AssetModel
		model.AssetSymbol=asset.AssetSymbol
		if asset.AssetSymbol=="" {
			return nil, common.ERR_COIN_TYPE
		}

		model.AssetAmount = abs(asset.AssetAmount)
		if !checkMoneyRange(model.AssetAmount) {
			return nil, common.ERR_CDP_STAKE_NUMBER
		}
		models = append(models, model)
	}
	tx.Assets=models

	result,err := tx.CreateRawTxid(wifKey)
	if err !=nil{
		return nil,err
	}

	return &CreateRawTxResult{result},nil
}


// returns the signature hex string and nil error, or returns empty string and the error if it has error
func (param *CdpLiquidateTxParam) CreateRawTx(privateKey string) (* CreateRawTxResult, error) {

	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return nil, common.ERR_INVALID_PRIVATEKEY
	}
	var tx waykichain.WaykiCdpLiquidateTx

	if param.ValidHeight < 0 {
		return nil, common.ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight

	tx.UserId = parseRegId(param.SrcRegId)

	pubKey, err := hex.DecodeString(param.PubKey)
	if (err != nil) {
		return nil, common.ERR_USER_PUBLICKEY
	}
	tx.PubKey = pubKey
	if (tx.UserId == nil && tx.PubKey == nil) {
		return nil, common.ERR_INVALID_SRC_REG_ID
	}
	if !checkCdpMinTxFee(param.Fees) {
		return nil, common.ERR_FEE_SMALLER_MIN
	}
	tx.Fees = uint64(param.Fees)

	tx.TxType = waykichain.CDP_LIQUIDATE_TX
	tx.Version = TX_VERSION
	if param.ScoinsLiquidate < 0 {
		return nil, common.ERR_CDP_STAKE_NUMBER
	}
	tx.ScoinsLiquidate = uint64(param.ScoinsLiquidate)

	if (param.FeeSymbol == "") {
		tx.FeeSymbol = string(common.WICC)
	} else {
		tx.FeeSymbol = string(param.FeeSymbol)
	}

	if param.AssetSymbol==""{
		return nil,common.ERR_COIN_TYPE
	}
	tx.AseetSymbol=param.AssetSymbol

	txHash, err := hex.DecodeString(param.CdpTxid)
	if (err != nil) {
		return nil, common.ERR_CDP_TX_HASH
	}
	tx.CdpTxHash = txHash

	result,err := tx.CreateRawTxid(wifKey)
	if err !=nil{
		return nil,err
	}

	return &CreateRawTxResult{result},nil
}


// returns the signature hex string and nil error, or returns empty string and the error if it has error
func (param *DexLimitSellTxParam) CreateRawTx(privateKey string) (* CreateRawTxResult,error) {
	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return nil, common.ERR_INVALID_PRIVATEKEY
	}
	var tx waykichain.WaykiDexSellLimitTx

	if param.ValidHeight < 0 {
		return nil, common.ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight

	tx.UserId = parseRegId(param.SrcRegId)

	pubKey, err := hex.DecodeString(param.PubKey)
	if (err != nil) {
		return nil, common.ERR_USER_PUBLICKEY
	}
	tx.PubKey = pubKey
	if (tx.UserId == nil && tx.PubKey == nil) {
		return nil, common.ERR_INVALID_SRC_REG_ID
	}
	if !checkMinTxFee(param.Fees) {
		return nil, common.ERR_FEE_SMALLER_MIN
	}
	tx.Fees = uint64(param.Fees)

	tx.TxType = waykichain.DEX_SELL_LIMIT_ORDER_TX
	tx.Version = TX_VERSION
	if param.Price <= 0 {
		return nil, common.ERR_ASK_PRICE
	}
	tx.AskPrice = uint64(param.Price)
	if param.CoinSymbol == "" || param.AssetSymbol == "" {
		return nil, common.ERR_COIN_TYPE
	}
	tx.AssetSymbol = param.AssetSymbol
	tx.CoinSymbol = param.CoinSymbol
	if (param.FeeSymbol == "") {
		tx.FeeSymbol = string(common.WICC)
	} else {
		tx.FeeSymbol = string(param.FeeSymbol)
	}
	if !checkMoneyRange(param.AssetAmount) {
		return nil, common.ERR_RANGE_VALUES
	}
	tx.AssetAmount = uint64(param.AssetAmount)

	result,err := tx.CreateRawTxid(wifKey)
	if err !=nil{
		return nil,err
	}

	return &CreateRawTxResult{result},nil
}

// returns the signature hex string and nil error, or returns empty string and the error if it has error
func (param *DexLimitBuyTxParam) CreateRawTx(privateKey string) (* CreateRawTxResult, error) {
	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return nil, common.ERR_INVALID_PRIVATEKEY
	}
	var tx waykichain.WaykiDexSellLimitTx

	if param.ValidHeight < 0 {
		return nil, common.ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight

	tx.UserId = parseRegId(param.SrcRegId)

	pubKey, err := hex.DecodeString(param.PubKey)
	if (err != nil) {
		return nil, common.ERR_USER_PUBLICKEY
	}
	tx.PubKey = pubKey
	if (tx.UserId == nil && tx.PubKey == nil) {
		return nil, common.ERR_INVALID_SRC_REG_ID
	}
	if !checkCdpMinTxFee(param.Fees) {
		return nil, common.ERR_FEE_SMALLER_MIN
	}
	tx.Fees = uint64(param.Fees)

	tx.TxType = waykichain.DEX_BUY_LIMIT_ORDER_TX
	tx.Version = TX_VERSION
	if param.Price <= 0 {
		return nil,common.ERR_ASK_PRICE
	}
	tx.AskPrice = uint64(param.Price)
	if param.CoinSymbol == "" || param.AssetSymbol == "" {
		return nil, common.ERR_COIN_TYPE
	}
	tx.AssetSymbol = param.AssetSymbol
	tx.CoinSymbol = param.CoinSymbol
	if (param.FeeSymbol == "") {
		tx.FeeSymbol = string(common.WICC)
	} else {
		tx.FeeSymbol = string(param.FeeSymbol)
	}
	if !checkMoneyRange(param.AssetAmount) {
		return nil, common.ERR_RANGE_VALUES
	}
	tx.AssetAmount = uint64(param.AssetAmount)

	result,err := tx.CreateRawTxid(wifKey)
	if err !=nil{
		return nil,err
	}

	return &CreateRawTxResult{result},nil
}


// returns the signature hex string and nil error, or returns empty string and the error if it has error
func (param *DexMarketSellTxParam) CreateRawTx(privateKey string) (* CreateRawTxResult,error) {
	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return nil, common.ERR_INVALID_PRIVATEKEY
	}
	var tx waykichain.WaykiDexMarketTx

	if param.ValidHeight < 0 {
		return nil, common.ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight

	tx.UserId = parseRegId(param.SrcRegId)

	pubKey, err := hex.DecodeString(param.PubKey)
	if (err != nil) {
		return nil, common.ERR_USER_PUBLICKEY
	}
	tx.PubKey = pubKey
	if (tx.UserId == nil && tx.PubKey == nil) {
		return nil, common.ERR_INVALID_SRC_REG_ID
	}
	if !checkCdpMinTxFee(param.Fees) {
		return nil, common.ERR_FEE_SMALLER_MIN
	}
	tx.Fees = uint64(param.Fees)

	tx.TxType = waykichain.DEX_SELL_MARKET_ORDER_TX
	tx.Version = TX_VERSION
	if param.CoinSymbol == "" || param.AssetSymbol == "" {
		return nil,common.ERR_COIN_TYPE
	}
	tx.AssetSymbol = param.AssetSymbol
	tx.CoinSymbol = param.CoinSymbol
	if (param.FeeSymbol == "") {
		tx.FeeSymbol = string(common.WICC)
	} else {
		tx.FeeSymbol = string(param.FeeSymbol)
	}
	if !checkMoneyRange(param.AssetAmount) {
		return nil, common.ERR_RANGE_VALUES
	}
	tx.AssetAmount = uint64(param.AssetAmount)

	result,err := tx.CreateRawTxid(wifKey)
	if err !=nil{
		return nil,err
	}

	return &CreateRawTxResult{result},nil
}

// returns the signature hex string and nil error, or returns empty string and the error if it has error
func (param *DexMarketBuyTxParam) CreateRawTx(privateKey string) (* CreateRawTxResult, error) {
	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return nil, common.ERR_INVALID_PRIVATEKEY
	}
	var tx waykichain.WaykiDexMarketTx

	if param.ValidHeight < 0 {
		return nil, common.ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight

	tx.UserId = parseRegId(param.SrcRegId)

	pubKey, err := hex.DecodeString(param.PubKey)
	if (err != nil) {
		return nil, common.ERR_USER_PUBLICKEY
	}
	tx.PubKey = pubKey
	if (tx.UserId == nil && tx.PubKey == nil) {
		return nil, common.ERR_INVALID_SRC_REG_ID
	}
	if !checkCdpMinTxFee(param.Fees) {
		return nil, common.ERR_FEE_SMALLER_MIN
	}
	tx.Fees = uint64(param.Fees)

	tx.TxType = waykichain.DEX_BUY_MARKET_ORDER_TX
	tx.Version = TX_VERSION
	if param.CoinSymbol == "" || param.AssetSymbol == "" {
		return nil, common.ERR_COIN_TYPE
	}
	tx.AssetSymbol = param.AssetSymbol
	tx.CoinSymbol = param.CoinSymbol
	if (param.FeeSymbol == "") {
		tx.FeeSymbol = string(common.WICC)
	} else {
		tx.FeeSymbol = string(param.FeeSymbol)
	}
	if !checkMoneyRange(param.CoinAmount) {
		return nil, common.ERR_RANGE_VALUES
	}
	tx.AssetAmount = uint64(param.CoinAmount)

	result,err := tx.CreateRawTxid(wifKey)
	if err !=nil{
		return nil,err
	}

	return &CreateRawTxResult{result},nil
}


// returns the signature hex string and nil error, or returns empty string and the error if it has error
func (param *DexCancelTxParam) CreateRawTx(privateKey string) (* CreateRawTxResult,error) {
	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return nil, common.ERR_INVALID_PRIVATEKEY
	}
	var tx waykichain.WaykiDexCancelTx

	if param.ValidHeight < 0 {
		return nil, common.ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight

	tx.UserId = parseRegId(param.SrcRegId)

	pubKey, err := hex.DecodeString(param.PubKey)
	if (err != nil) {
		return nil, common.ERR_USER_PUBLICKEY
	}
	tx.PubKey = pubKey
	if (tx.UserId == nil && tx.PubKey == nil) {
		return nil, common.ERR_INVALID_SRC_REG_ID
	}
	if !checkCdpMinTxFee(param.Fees) {
		return nil, common.ERR_FEE_SMALLER_MIN
	}
	tx.Fees = uint64(param.Fees)

	tx.TxType = waykichain.DEX_CANCEL_ORDER_TX
	tx.Version = TX_VERSION
	if (param.FeeSymbol == "") {
		tx.FeeSymbol = string(common.WICC)
	} else {
		tx.FeeSymbol = string(param.FeeSymbol)
	}
	txHash, err := hex.DecodeString(param.DexTxid)
	if (err != nil) {
		return nil, common.ERR_CDP_TX_HASH
	}
	tx.DexHash = txHash

	result,err := tx.CreateRawTxid(wifKey)
	if err !=nil{
		return nil,err
	}

	return &CreateRawTxResult{result},nil
}



// returns the signature hex string and nil error, or returns empty string and the error if it has error
func (param *AssetIssueTxParam) CreateRawTx(privateKey string) (* CreateRawTxResult, error) {

	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return nil, common.ERR_INVALID_PRIVATEKEY
	}
	var tx waykichain.WaykiAssetIssueTx

	if param.ValidHeight < 0 {
		return nil, common.ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight
	tx.UserId = parseRegId(param.SrcRegId)

	if (tx.UserId == nil) {
		return nil, common.ERR_INVALID_SRC_REG_ID
	}

	tx.Fees = uint64(param.Fees)

	tx.TxType = waykichain.ASSET_ISSUE_TX
	tx.Version = TX_VERSION

	if !checkAssetSymbol(param.AssetSymbol) {
		return nil, common.ERR_SYMBOL_ERROR
	}
	tx.AssetSymbol = param.AssetSymbol
	if (param.FeeSymbol == "") {
		tx.FeeSymbol = string(common.WICC)
	} else {
		tx.FeeSymbol = string(param.FeeSymbol)
	}
	tx.MinTable = param.MinTable
	if (param.AssetName == "") {
		return nil, common.ERR_ASSET_NAME_ERROR
	}
	tx.AssetName = param.AssetName
	if (param.AssetTotal < 100000000) {
		return nil, common.ERR_TOTAl_SUPPLY_ERROR
	}
	tx.AssetTotal = uint64(param.AssetTotal)
	tx.AssetOwner = parseUserId(param.AssetOwner)
	if (tx.AssetOwner == nil) {
		return nil, common.ERR_ASSET_UPDATE_OWNER_ERROR
	}

	result,err := tx.CreateRawTxid(wifKey)
	if err !=nil{
		return nil,err
	}

	return &CreateRawTxResult{result},nil
}


// returns the signature hex string and nil error, or returns empty string and the error if it has error
func (param *AssetUpdateTxParam) CreateRawTx(privateKey string) (* CreateRawTxResult,error) {

	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return nil, common.ERR_INVALID_PRIVATEKEY
	}
	var tx waykichain.WaykiAssetUpdateTx

	if param.ValidHeight < 0 {
		return nil, common.ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight
	tx.UserId = parseRegId(param.SrcRegId)
	if (tx.UserId == nil) {
		return nil, common.ERR_INVALID_SRC_REG_ID
	}

	tx.Fees = uint64(param.Fees)

	tx.TxType = waykichain.ASSET_UPDATE_TX
	tx.Version = TX_VERSION

	if !checkAssetSymbol(param.AssetSymbol) {
		return nil, common.ERR_SYMBOL_ERROR
	}
	tx.AssetSymbol = param.AssetSymbol
	if (param.FeeSymbol == "") {
		tx.FeeSymbol = string(common.WICC)
	} else {
		tx.FeeSymbol = string(param.FeeSymbol)
	}
	tx.UpdateType=param.UpdateType
	switch param.UpdateType {
	case 1:
		tx.AssetOwner = parseUserId(param.AssetOwner)
		if (tx.AssetOwner == nil) {
			return nil, common.ERR_ASSET_UPDATE_OWNER_ERROR
		}
		break
	case 2:
		if (param.AssetName == "") {
			return nil, common.ERR_ASSET_NAME_ERROR
		}
		tx.AssetName = param.AssetName
		break
	case 3:
		if (param.AssetTotal < 100000000) {
			return nil, common.ERR_TOTAl_SUPPLY_ERROR
		}
		tx.AssetTotal = uint64(param.AssetTotal)
		break
	default:
		return nil, common.ERR_ASSET_UPDATE_TYPE_ERROR
	}

	result,err := tx.CreateRawTxid(wifKey)
	if err !=nil{
		return nil,err
	}

	return &CreateRawTxResult{result},nil
}


// returns the signature hex string and nil error, or returns empty string and the error if it has error
func (param *RegisterAccountTxParam) CreateRawTx(privateKey string) (* CreateRawTxResult, error) {
	// check and convert params
	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return nil, common.ERR_INVALID_PRIVATEKEY
	}

	var tx waykichain.WaykiRegisterAccountTx
	if param.ValidHeight < 0 {
		return nil, common.ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight

	if !checkMinTxFee(param.Fees) {
		return nil, common.ERR_FEE_SMALLER_MIN
	}
	tx.Fees = uint64(param.Fees)

	tx.UserId = waykichain.NewPubKeyUid(wifKey.PrivKey.PubKey().SerializeCompressed())
	tx.TxType = waykichain.REG_ACCT_TX
	tx.Version = 1

	result,err := tx.CreateRawTxid(wifKey)
	if err !=nil{
		return nil,err
	}

	return &CreateRawTxResult{result},nil
}


// returns the signature hex string and nil error, or returns empty string and the error if it has error
func (param *CommonTxParam) CreateRawTx(privateKey string) (* CreateRawTxResult, error) {

	// check and convert params
	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return nil, common.ERR_INVALID_PRIVATEKEY
	}

	var tx waykichain.WaykiCommonTx

	if param.ValidHeight < 0 {
		return nil, common.ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight

	tx.UserId = parseRegId(param.SrcRegId)

	tx.DestId = parseUserId(param.DestAddr)
	if tx.DestId == nil {
		return nil, common.ERR_INVALID_DEST_ADDR
	}

	if !checkMoneyRange(param.Values) {
		return nil, common.ERR_RANGE_VALUES
	}
	tx.Values = uint64(param.Values)

	if !checkMoneyRange(param.Fees) {
		return nil, common.ERR_RANGE_FEE
	}

	if !checkMinTxFee(param.Fees) {
		return nil, common.ERR_FEE_SMALLER_MIN
	}
	tx.Fees = uint64(param.Fees)

	pubKey, err := hex.DecodeString(param.PubKey)
	if (err != nil) {
		return nil, common.ERR_USER_PUBLICKEY
	}
	tx.PubKey = pubKey
	if (tx.UserId == nil && tx.PubKey == nil) {
		return nil, common.ERR_INVALID_SRC_REG_ID
	}
	tx.TxType = waykichain.COMMON_TX
	tx.Version = TX_VERSION
	tx.Memo = param.Memo

	result,err := tx.CreateRawTxid(wifKey)
	if err !=nil{
		return nil,err
	}

	return &CreateRawTxResult{result},nil
}


// returns the signature hex string and nil error, or returns empty string and the error if it has error
func (param *CallContractTxParam) CreateRawTx(privateKey string) (* CreateRawTxResult,error) {
	// check and convert params
	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return nil, common.ERR_INVALID_PRIVATEKEY
	}

	var tx  waykichain.WaykiCallContractTx

	if param.ValidHeight < 0 {
		return nil, common.ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight

	tx.UserId = parseRegId(param.SrcRegId)

	tx.AppId = parseRegId(param.AppId)
	if tx.AppId == nil {
		return nil, common.ERR_INVALID_APP_ID
	}

	if !checkMoneyRange(param.Values) {
		return nil, common.ERR_RANGE_VALUES
	}
	tx.Values = uint64(param.Values)

	if !checkMoneyRange(param.Fees) {
		return nil, common.ERR_RANGE_FEE
	}
	if !checkMinTxFee(param.Fees) {
		return nil, common.ERR_FEE_SMALLER_MIN
	}
	tx.Fees = uint64(param.Fees)

	binary, errHex := hex.DecodeString(param.ContractHex)
	if errHex != nil {
		return nil, common.ERR_INVALID_CONTRACT_HEX
	}
	tx.Contract = []byte(binary)
	pubKey, err := hex.DecodeString(param.PubKey)
	if (err != nil) {
		return nil, common.ERR_USER_PUBLICKEY
	}
	tx.PubKey = pubKey
	if (tx.UserId == nil && tx.PubKey == nil) {
		return nil, common.ERR_INVALID_SRC_REG_ID
	}
	tx.TxType = waykichain.CONTRACT_TX
	tx.Version = TX_VERSION

	result,err := tx.CreateRawTxid(wifKey)
	if err !=nil{
		return nil,err
	}

	return &CreateRawTxResult{result},nil
}



// returns the signature hex string and nil error, or returns empty string and the error if it has error
func (param *RegisterContractTxParam) CreateRawTx(privateKey string) (* CreateRawTxResult, error) {

	// check and convert params
	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return nil, common.ERR_INVALID_PRIVATEKEY
	}

	var tx waykichain.WaykiRegisterContractTx

	if param.ValidHeight < 0 {
		return nil, common.ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight

	tx.UserId = parseRegId(param.SrcRegId)

	if !checkMoneyRange(param.Fees) {
		return nil, common.ERR_RANGE_FEE
	}
	if !checkMinTxFee(param.Fees) {
		return nil, common.ERR_FEE_SMALLER_MIN
	}
	tx.Fees = uint64(param.Fees)

	tx.TxType = waykichain.REG_CONT_TX
	tx.Version = TX_VERSION

	if len(param.Script) == 0 || len(param.Script) > CONTRACT_SCRIPT_MAX_SIZE {
		return nil, common.ERR_INVALID_SCRIPT
	}
	tx.Script = param.Script

	if len(param.Description) > CONTRACT_SCRIPT_DESC_MAX_SIZE {
		return nil, common.ERR_INVALID_SCRIPT_DESC
	}
	tx.Description = param.Description
	if (tx.UserId == nil && tx.PubKey == nil) {
		return nil, common.ERR_INVALID_SRC_REG_ID
	}

	result,err := tx.CreateRawTxid(wifKey)
	if err !=nil{
		return nil,err
	}

	return &CreateRawTxResult{result},nil
}

