package waykichain

import (
	"bytes"
	"encoding/hex"
	"github.com/WaykiChain/wicc-wallet-utils-go/common"
	"github.com/btcsuite/btcutil"
)

//Base Params for Tx
type WaykiBaseSignTx struct {
	TxType      WaykiTxType
	Version     int64
	ValidHeight int64
	PubKey      []byte
	UserId      *UserIdWraper // current operating user id, the one want to do something
}
/***************************************UCOIN_TRANSFER_TX*************************/
type UCoinTransferDest struct {
	CoinSymbol string   //From Coin Type
	CoinAmount uint64
	DestAddr    *UserIdWraper
}

type WaykiUCoinTransferTx struct {
	WaykiBaseSignTx
	Fees   uint64
	FeeSymbol string      //Fee Type (WICC/WUSD)
	Dests    []UCoinTransferDest
	Memo       string
}
//return rawtx + txid
func (tx WaykiUCoinTransferTx) createRawTx(wifKey *btcutil.WIF) (string,string,error){
	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)

	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.Version)
	writer.WriteVarInt(tx.ValidHeight)
	if(tx.UserId!=nil){
		writer.WriteUserId(tx.UserId)
	}else if(tx.PubKey!=nil){
		writer.WritePubKeyId(tx.PubKey)
	}
	writer.WriteString(tx.FeeSymbol)
	writer.WriteVarInt(int64(tx.Fees))
	writer.WriteCompactSize(uint64(len(tx.Dests)))
	for i:=0;i<len(tx.Dests);i++  {
		writer.WriteUserId(tx.Dests[i].DestAddr)
		writer.WriteString(tx.Dests[i].CoinSymbol)
		writer.WriteVarInt(int64(tx.Dests[i].CoinAmount))
	}
	writer.WriteString(tx.Memo)
	signatureBytes,txid ,err:= tx.CalSignatureTxid(wifKey)
	if err != nil{
		return "","",err
	}
	writer.WriteBytes(signatureBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return rawTx, txid,nil
}

//return signature and txid
func (tx WaykiUCoinTransferTx) CalSignatureTxid(wifKey *btcutil.WIF) ([]byte,string,error){

	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)

	writer.WriteVarInt(tx.Version)
	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.ValidHeight)
	if(tx.UserId!=nil){
		writer.WriteUserId(tx.UserId)
	}else if(tx.PubKey!=nil){
		writer.WritePubKeyId(tx.PubKey)
	}
	writer.WriteString(tx.FeeSymbol)
	writer.WriteVarInt(int64(tx.Fees))
	writer.WriteCompactSize(uint64(len(tx.Dests)))
	for i:=0;i<len(tx.Dests);i++  {
		writer.WriteUserId(tx.Dests[i].DestAddr)
		writer.WriteString(tx.Dests[i].CoinSymbol)
		writer.WriteVarInt(int64(tx.Dests[i].CoinAmount))
	}
	writer.WriteString(tx.Memo)
	return calSignatureTxid(buf.Bytes(),wifKey)
}

//CreateUCoinTransferRawTx sign for Multi-currency transfer
// returns the signature hex string and nil error, or returns empty string and the error if it has error
func (param *UCoinTransferTxParam) CreateRawTx(privateKey string) (string, string,error) {

	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return "","", common.ERR_INVALID_PRIVATEKEY
	}
	var tx WaykiUCoinTransferTx

	if param.ValidHeight < 0 {
		return "","", common.ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight

	tx.UserId = parseRegId(param.SrcRegId)

	pubKey, err := hex.DecodeString(param.PubKey)
	if (err != nil) {
		return "","", common.ERR_USER_PUBLICKEY
	}
	tx.PubKey = pubKey
	if (tx.UserId == nil && tx.PubKey == nil) {
		return "","", common.ERR_INVALID_SRC_REG_ID
	}
	if !checkMoneyRange(param.Fees) {
		return "","", common.ERR_RANGE_FEE
	}
	tx.Fees = uint64(param.Fees)

	tx.TxType = UCOIN_TRANSFER_TX
	tx.Version = TX_VERSION

	var dests []UCoinTransferDest
	for i:=0;i< len(param.Dests.DestArray);i++{
		var dest UCoinTransferDest
		dest.DestAddr = parseUserId(param.Dests.DestArray[i].DestAddr)
		if param.Dests.DestArray[i].CoinAmount < 0 {
			return "","", common.ERR_RANGE_VALUES
		}
		dest.CoinAmount = uint64(param.Dests.DestArray[i].CoinAmount)
		if param.Dests.DestArray[i].CoinSymbol == "" {
			return "","", common.ERR_COIN_TYPE
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

	return tx.createRawTx(wifKey)
}

/***************************************UCOIN_CONTRACT_INVOKE_TX*************************/
type WaykiUCoinCallContractTx struct {
	WaykiBaseSignTx
	AppId    *UserIdWraper //user regid or user key id or app regid
	Fees     int64
	CoinAmount   int64 //transfer amount
	FeeSymbol string      //Fee Type (WICC/WUSD)
	CoinSymbol string   //From Coin Type
	Contract []byte
}

func (tx WaykiUCoinCallContractTx) createRawTx(wifKey *btcutil.WIF) (string,string,error) {

	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)

	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.Version)
	writer.WriteVarInt(tx.ValidHeight)
	if(tx.UserId!=nil){
		writer.WriteUserId(tx.UserId)
	}else if(tx.PubKey!=nil){
		writer.WritePubKeyId(tx.PubKey)
	}
	writer.WriteUserId(tx.AppId)
	writer.WriteBytes(tx.Contract)
	writer.WriteVarInt(int64(tx.Fees))
	writer.WriteString(tx.FeeSymbol)
	writer.WriteString(tx.CoinSymbol)
	writer.WriteVarInt(int64(tx.CoinAmount))
	signatureBytes,txid,err := tx.CalSignatureTxid(wifKey)
	if err != nil{
		return "","",err
	}
	writer.WriteBytes(signatureBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return rawTx, txid,nil
}

//return signature and txid
func (tx WaykiUCoinCallContractTx) CalSignatureTxid(wifKey *btcutil.WIF) ([]byte,string,error) {
	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)

	writer.WriteVarInt(tx.Version)
	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.ValidHeight)
	if(tx.UserId!=nil){
		writer.WriteUserId(tx.UserId)
	}else if(tx.PubKey!=nil){
		writer.WritePubKeyId(tx.PubKey)
	}
	writer.WriteUserId(tx.AppId)
	writer.WriteBytes(tx.Contract)
	writer.WriteVarInt(int64(tx.Fees))
	writer.WriteString(tx.FeeSymbol)
	writer.WriteString(tx.CoinSymbol)
	writer.WriteVarInt(int64(tx.CoinAmount))

	return calSignatureTxid(buf.Bytes(),wifKey)
}
// returns the signature hex string and nil error, or returns empty string and the error if it has error
func (param *UCoinContractTxParam) CreateRawTx(privateKey string) (string,string, error) {
	// check and convert params
	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return "","", common.ERR_INVALID_PRIVATEKEY
	}

	var tx WaykiUCoinCallContractTx

	if param.ValidHeight < 0 {
		return "","", common.ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight
	tx.UserId = parseRegId(param.SrcRegId)
	tx.AppId = parseRegId(param.AppId)
	if tx.AppId == nil {
		return "","", common.ERR_INVALID_APP_ID
	}
	if !checkMoneyRange(param.CoinAmount) {
		return "","", common.ERR_RANGE_VALUES
	}
	tx.CoinAmount = int64(param.CoinAmount)
	if !checkMoneyRange(param.Fees) {
		return "","", common.ERR_RANGE_FEE
	}
	if !checkMinTxFee(param.Fees) {
		return "","", common.ERR_FEE_SMALLER_MIN
	}
	tx.Fees = int64(param.Fees)
	binary, errHex := hex.DecodeString(param.ContractHex)
	if errHex != nil {
		return "", "",common.ERR_INVALID_CONTRACT_HEX
	}
	tx.Contract = []byte(binary)
	pubKey, err := hex.DecodeString(param.PubKey)
	if (err != nil) {
		return "","", common.ERR_USER_PUBLICKEY
	}
	tx.PubKey = pubKey
	if (tx.UserId == nil && tx.PubKey == nil) {
		return "","", common.ERR_INVALID_SRC_REG_ID
	}
	if param.CoinSymbol == "" {
		return "","", common.ERR_COIN_TYPE
	}
	tx.CoinSymbol = param.CoinSymbol
	if (param.FeeSymbol == "") {
		tx.FeeSymbol = string(common.WICC)
	} else {
		tx.FeeSymbol = string(param.FeeSymbol)
	}
	tx.TxType = UCOIN_CONTRACT_INVOKE_TX
	tx.Version = TX_VERSION

	return tx.createRawTx(wifKey)
}


type WaykiUCoinRegisterContractTx struct {
	WaykiBaseSignTx
	Script      []byte
	Description string
	Fees        uint64
	FeeSymbol   string
}

func WriteContractScript(writer *WriterHelper, script []byte, description string) {

	scriptWriter := NewWriterHelper(bytes.NewBuffer([]byte{}))
	scriptWriter.WriteBytes(script)
	scriptWriter.WriteString(description)
	writer.WriteBytes(scriptWriter.GetBuf().Bytes())
}
// sign transaction
func (tx WaykiUCoinRegisterContractTx) createRawTx(wifKey *btcutil.WIF) (string,string,error) {

	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)
	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.Version)
	writer.WriteVarInt(tx.ValidHeight)
	writer.WriteUserId(tx.UserId)
	writer.WriteVarInt(int64(tx.Fees))
	writer.WriteString(tx.FeeSymbol)
	WriteContractScript(writer, tx.Script, tx.Description)
	signatureBytes,txid,err := tx.CalSignatureTxid(wifKey)
	if err != nil{
		return "","",err
	}
	writer.WriteBytes(signatureBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return rawTx, txid,nil
}

func (tx WaykiUCoinRegisterContractTx) CalSignatureTxid(wifKey *btcutil.WIF) ([]byte,string,error) {
	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)
	writer.WriteVarInt(tx.Version)
	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.ValidHeight)
	writer.WriteUserId(tx.UserId)
	writer.WriteVarInt(int64(tx.Fees))
	writer.WriteString(tx.FeeSymbol)
	WriteContractScript(writer, tx.Script, tx.Description)

	return calSignatureTxid(buf.Bytes(),wifKey)
}

//CreateUCoinDeployContractTx sign for call register contract tx
// returns the signature hex string and nil error, or returns empty string and the error if it has error
func (param *UCoinRegisterContractTxParam) CreateRawTx(privateKey string) (string, string,error) {

	// check and convert params
	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return "","", common.ERR_INVALID_PRIVATEKEY
	}

	var tx WaykiUCoinRegisterContractTx

	if param.ValidHeight < 0 {
		return "","", common.ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight

	tx.UserId = parseRegId(param.SrcRegId)

	if !checkMoneyRange(param.Fees) {
		return "","", common.ERR_RANGE_FEE
	}
	if !checkMinTxFee(param.Fees) {
		return "", "",common.ERR_FEE_SMALLER_MIN
	}
	tx.Fees = uint64(param.Fees)

	tx.TxType = UCONTRACT_DEPLOY_TX
	tx.Version = TX_VERSION

	if len(param.Script) == 0 || len(param.Script) > CONTRACT_SCRIPT_MAX_SIZE {
		return "", "",common.ERR_INVALID_SCRIPT
	}
	tx.Script = param.Script

	if len(param.Description) > CONTRACT_SCRIPT_DESC_MAX_SIZE {
		return "","", common.ERR_INVALID_SCRIPT_DESC
	}
	tx.Description = param.Description
	if (tx.UserId == nil && tx.PubKey == nil) {
		return "","", common.ERR_INVALID_SRC_REG_ID
	}

	if (param.FeeSymbol == "") {
		tx.FeeSymbol = string(common.WICC)
	} else {
		tx.FeeSymbol = string(param.FeeSymbol)
	}

	return tx.createRawTx(wifKey)
}


type OperVoteFundTx struct {
	VoteType  WaykiVoteType
	PubKey    *PubKeyId
	VoteValue int64
}

type WaykiDelegateTx struct {
	WaykiBaseSignTx
	OperVoteFunds []OperVoteFundTx
	Fees          uint64
}

//return rawtx + txid
func (tx WaykiDelegateTx) createRawTx(wifKey *btcutil.WIF) (string,string,error) {

	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)

	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.Version)
	writer.WriteVarInt(tx.ValidHeight)
	if(tx.UserId!=nil){
		writer.WriteUserId(tx.UserId)
	}else if(tx.PubKey!=nil){
		writer.WritePubKeyId(tx.PubKey)
	}
	writer.WriteVarInt(int64(len(tx.OperVoteFunds)))
	for _, fund := range tx.OperVoteFunds {
		writer.WriteByte(byte(fund.VoteType))
		writer.WritePubKeyId(*fund.PubKey)
		writer.WriteVarInt(fund.VoteValue)
	}
	writer.WriteVarInt(int64(tx.Fees))
	signatureBytes,txid,err := tx.CalSignatureTxid(wifKey)
	if err != nil{
		return "","",err
	}
	writer.WriteBytes(signatureBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return rawTx, txid,nil
}

//return signature and txid
func (tx WaykiDelegateTx) CalSignatureTxid(wifKey *btcutil.WIF) ([]byte,string,error){

	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)

	writer.WriteVarInt(tx.Version)
	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.ValidHeight)
	if(tx.UserId!=nil){
		writer.WriteUserId(tx.UserId)
	}else if(tx.PubKey!=nil){
		writer.WritePubKeyId(tx.PubKey)
	}
	writer.WriteVarInt(int64(len(tx.OperVoteFunds)))
	for _, fund := range tx.OperVoteFunds {
		writer.WriteByte(byte(fund.VoteType))
		writer.WriteBytes(*fund.PubKey)
		writer.WriteVarInt(fund.VoteValue)
	}
	writer.WriteVarInt(int64(tx.Fees))

	return calSignatureTxid(buf.Bytes(),wifKey)
}
// returns the signature hex string and nil error, or returns empty string and the error if it has error
func (param *DelegateTxParam) CreateRawTx(privateKey string) (string,string, error) {

	// check and convert params
	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return "","", common.ERR_INVALID_PRIVATEKEY
	}

	var tx WaykiDelegateTx

	if param.ValidHeight < 0 {
		return "","", common.ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight

	tx.UserId = parseRegId(param.SrcRegId)

	if !checkMoneyRange(param.Fees) {
		return "","", common.ERR_RANGE_FEE
	}

	if !checkMinTxFee(param.Fees) {
		return "","", common.ERR_FEE_SMALLER_MIN
	}

	if len(param.Votes.VoteArray) == 0 {
		return "","", common.ERR_EMPTY_VOTES
	}
	var voteData []OperVoteFundTx
	for i := 0; i < len(param.Votes.VoteArray); i++ {
		inVote := param.Votes.VoteArray[i]
		var v OperVoteFundTx
		if !checkPubKey(inVote.PubKey) {
			return "","", common.ERR_INVALID_VOTE_PUBKEY
		}

		v.VoteValue = abs(inVote.VoteValue)
		if !checkMoneyRange(v.VoteValue) {
			return "","", common.ERR_RANGE_VOTE_VALUE
		}
		v.VoteType = GetVoteTypeByValue(inVote.VoteValue)

		var pubKey PubKeyId = inVote.PubKey
		v.PubKey = &pubKey
		voteData = append(voteData, v)
	}

	pubKey, err := hex.DecodeString(param.PubKey)
	if (err != nil) {
		return "","", common.ERR_USER_PUBLICKEY
	}
	tx.PubKey = pubKey
	if (tx.UserId == nil && tx.PubKey == nil) {
		return "","", common.ERR_INVALID_SRC_REG_ID
	}
	tx.TxType = DELEGATE_TX
	tx.Version = 1
	tx.OperVoteFunds = voteData

	return tx.createRawTx(wifKey)
}


type WaykiCdpStakeTx struct {
	WaykiBaseSignTx
	Fees   uint64
	ScoinValues uint64   //get Coin amount
	FeeSymbol string      //Fee Type (WICC/WUSD)
	ScoinSymbol string   //get Coin Type
	Assets   []AssetModel
	CdpTxHash []byte
}

//return rawtx + txid
func (tx WaykiCdpStakeTx) createRawTx(wifKey *btcutil.WIF) (string,string,error){
	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)

	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.Version)
	writer.WriteVarInt(tx.ValidHeight)
	if(tx.UserId!=nil){
		writer.WriteUserId(tx.UserId)
	}else if(tx.PubKey!=nil){
		writer.WritePubKeyId(tx.PubKey)
	}
	writer.WriteString(tx.FeeSymbol)
	writer.WriteVarInt(int64(tx.Fees))
	writer.WriteReverse(tx.CdpTxHash)
	writer.WriteCdpAsset(tx.Assets)
	ss:=[]byte(tx.ScoinSymbol)
	writer.WriteVarInt(int64(len(ss)))
	writer.Write(ss)
	writer.WriteVarInt(int64(tx.ScoinValues))
	signatureBytes,txid,err := tx.CalSignatureTxid(wifKey)
	if err != nil{
		return "","",err
	}
	writer.WriteBytes(signatureBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return rawTx, txid,nil
}

//return signature and txid
func (tx WaykiCdpStakeTx) CalSignatureTxid(wifKey *btcutil.WIF) ([]byte,string,error){

	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)

	writer.WriteVarInt(tx.Version)
	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.ValidHeight)
	if(tx.UserId!=nil){
		writer.WriteUserId(tx.UserId)
	}else if(tx.PubKey!=nil){
		writer.WritePubKeyId(tx.PubKey)
	}
	writer.WriteString(tx.FeeSymbol)
	writer.WriteVarInt(int64(tx.Fees))
	writer.WriteReverse(tx.CdpTxHash)
	writer.WriteCdpAsset(tx.Assets)
	ssoin:=[]byte(tx.ScoinSymbol)
	writer.WriteVarInt(int64(len(ssoin)))
	writer.Write(ssoin)
	writer.WriteVarInt(int64(tx.ScoinValues))

	return calSignatureTxid(buf.Bytes(),wifKey)
}
// returns the signature hex string and nil error, or returns empty string and the error if it has error
func (param *CdpStakeTxParam) CreateRawTx(privateKey string) (string,string, error) {

	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return "","", common.ERR_INVALID_PRIVATEKEY
	}
	var tx WaykiCdpStakeTx

	if param.ValidHeight < 0 {
		return "", "",common.ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight

	tx.UserId = parseRegId(param.SrcRegId)

	pubKey, err := hex.DecodeString(param.PubKey)
	if (err != nil) {
		return "","", common.ERR_USER_PUBLICKEY
	}
	tx.PubKey = pubKey
	if (tx.UserId == nil && tx.PubKey == nil) {
		return "","", common.ERR_INVALID_SRC_REG_ID
	}
	if !checkCdpMinTxFee(param.Fees) {
		return "","", common.ERR_FEE_SMALLER_MIN
	}
	tx.Fees = uint64(param.Fees)

	tx.TxType = CDP_STAKE_TX
	tx.Version = TX_VERSION
	if  param.ScoinMint < 0 {
		return "","", common.ERR_CDP_STAKE_NUMBER
	}

	tx.ScoinValues = uint64(param.ScoinMint)

	if  param.ScoinSymbol == "" {
		return "","", common.ERR_COIN_TYPE
	}
	if (param.FeeSymbol == "") {
		tx.FeeSymbol = string(common.WICC)
	} else {
		tx.FeeSymbol = string(param.FeeSymbol)
	}
	var models []AssetModel
	for i := 0; i < len(param.Assets.AssetArray); i++ {
		asset := param.Assets.AssetArray[i]
		var model AssetModel
		model.AssetSymbol=asset.AssetSymbol
		if asset.AssetSymbol=="" {
			return "","", common.ERR_COIN_TYPE
		}

		model.AssetAmount = abs(asset.AssetAmount)
		if !checkMoneyRange(model.AssetAmount) {
			return "","", common.ERR_CDP_STAKE_NUMBER
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
		return "","", common.ERR_CDP_TX_HASH
	}
	tx.CdpTxHash = txHash

	return tx.createRawTx(wifKey)
}


type WaykiCdpRedeemTx struct {
	WaykiBaseSignTx
	Fees   uint64
	ScoinValues uint64   //Stake Coin
	Assets   []AssetModel
	FeeSymbol string      //Fee Type (WICC/WUSD)
	CdpTxHash []byte
}

//return rawtx + txid
func (tx WaykiCdpRedeemTx) createRawTx(wifKey *btcutil.WIF) (string,string,error){
	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)

	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.Version)
	writer.WriteVarInt(tx.ValidHeight)
	if(tx.UserId!=nil){
		writer.WriteUserId(tx.UserId)
	}else if(tx.PubKey!=nil){
		writer.WritePubKeyId(tx.PubKey)
	}
	writer.WriteString(tx.FeeSymbol)
	writer.WriteVarInt(int64(tx.Fees))
	writer.WriteReverse(tx.CdpTxHash)
	writer.WriteVarInt(int64(tx.ScoinValues))
	writer.WriteCdpAsset(tx.Assets)
	signatureBytes,txid,err := tx.CalSignatureTxid(wifKey)
	if err != nil{
		return "","",err
	}
	writer.WriteBytes(signatureBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return rawTx, txid,nil
}

//return signature and txid
func (tx WaykiCdpRedeemTx) CalSignatureTxid(wifKey *btcutil.WIF) ([]byte,string,error){

	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)

	writer.WriteVarInt(tx.Version)
	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.ValidHeight)
	if(tx.UserId!=nil){
		writer.WriteUserId(tx.UserId)
	}else if(tx.PubKey!=nil){
		writer.WritePubKeyId(tx.PubKey)
	}
	writer.WriteString(tx.FeeSymbol)
	writer.WriteVarInt(int64(tx.Fees))
	writer.WriteReverse(tx.CdpTxHash)
	writer.WriteVarInt(int64(tx.ScoinValues))
	writer.WriteCdpAsset(tx.Assets)

	return calSignatureTxid(buf.Bytes(),wifKey)
}
// returns the signature hex string and nil error, or returns empty string and the error if it has error
func (param *CdpRedeemTxParam) CreateRawTx(privateKey string) (string,string ,error) {

	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return "", "",common.ERR_INVALID_PRIVATEKEY
	}
	var tx WaykiCdpRedeemTx

	if param.ValidHeight < 0 {
		return "", "",common.ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight

	tx.UserId = parseRegId(param.SrcRegId)

	pubKey, err := hex.DecodeString(param.PubKey)
	if (err != nil) {
		return "","", common.ERR_USER_PUBLICKEY
	}
	tx.PubKey = pubKey
	if (tx.UserId == nil && tx.PubKey == nil) {
		return "","", common.ERR_INVALID_SRC_REG_ID
	}
	if !checkCdpMinTxFee(param.Fees) {
		return "","", common.ERR_FEE_SMALLER_MIN
	}
	tx.Fees = uint64(param.Fees)

	tx.TxType = CDP_REDEEMP_TX
	tx.Version = TX_VERSION
	if param.ScoinsToRepay < 0 {
		return "","", common.ERR_CDP_STAKE_NUMBER
	}
	tx.ScoinValues = uint64(param.ScoinsToRepay)

	if (param.FeeSymbol == "") {
		tx.FeeSymbol = string(common.WICC)
	} else {
		tx.FeeSymbol = string(param.FeeSymbol)
	}

	txHash, err := hex.DecodeString(param.CdpTxid)
	if (err != nil) {
		return "","", common.ERR_CDP_TX_HASH
	}
	tx.CdpTxHash = txHash

	var models []AssetModel
	for i := 0; i < len(param.Assets.AssetArray); i++ {
		asset := param.Assets.AssetArray[i]
		var model AssetModel
		model.AssetSymbol=asset.AssetSymbol
		if asset.AssetSymbol=="" {
			return "","", common.ERR_COIN_TYPE
		}

		model.AssetAmount = abs(asset.AssetAmount)
		if !checkMoneyRange(model.AssetAmount) {
			return "","", common.ERR_CDP_STAKE_NUMBER
		}
		models = append(models, model)
	}
	tx.Assets=models

	return tx.createRawTx(wifKey)
}

type WaykiCdpLiquidateTx struct {
	WaykiBaseSignTx
	Fees   uint64
	ScoinsLiquidate uint64   //Scoin  Liquidate
	FeeSymbol string      //Fee Type (WICC/WUSD)
	CdpTxHash []byte
	AseetSymbol string
}
//return rawtx + txid
func (tx WaykiCdpLiquidateTx) createRawTx(wifKey *btcutil.WIF) (string,string,error){
	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)

	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.Version)
	writer.WriteVarInt(tx.ValidHeight)
	if(tx.UserId!=nil){
		writer.WriteUserId(tx.UserId)
	}else if(tx.PubKey!=nil){
		writer.WritePubKeyId(tx.PubKey)
	}
	writer.WriteString(tx.FeeSymbol)
	writer.WriteVarInt(int64(tx.Fees))
	writer.WriteReverse(tx.CdpTxHash)
	writer.WriteString(tx.AseetSymbol)
	writer.WriteVarInt(int64(tx.ScoinsLiquidate))
	signatureBytes,txid,err := tx.CalSignatureTxid(wifKey)
	if err != nil{
		return "","",err
	}
	writer.WriteBytes(signatureBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return rawTx, txid,nil
}

//return signature and txid
func (tx WaykiCdpLiquidateTx) CalSignatureTxid(wifKey *btcutil.WIF) ([]byte,string,error){

	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)

	writer.WriteVarInt(tx.Version)
	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.ValidHeight)
	if(tx.UserId!=nil){
		writer.WriteUserId(tx.UserId)
	}else if(tx.PubKey!=nil){
		writer.WritePubKeyId(tx.PubKey)
	}
	writer.WriteString(tx.FeeSymbol)
	writer.WriteVarInt(int64(tx.Fees))
	writer.WriteReverse(tx.CdpTxHash)
	writer.WriteString(tx.AseetSymbol)
	writer.WriteVarInt(int64(tx.ScoinsLiquidate))
	return calSignatureTxid(buf.Bytes(),wifKey)
}
// returns the signature hex string and nil error, or returns empty string and the error if it has error
func (param *CdpLiquidateTxParam) CreateRawTx(privateKey string) (string,string, error) {

	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return "","", common.ERR_INVALID_PRIVATEKEY
	}
	var tx WaykiCdpLiquidateTx

	if param.ValidHeight < 0 {
		return "","", common.ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight

	tx.UserId = parseRegId(param.SrcRegId)

	pubKey, err := hex.DecodeString(param.PubKey)
	if (err != nil) {
		return "","", common.ERR_USER_PUBLICKEY
	}
	tx.PubKey = pubKey
	if (tx.UserId == nil && tx.PubKey == nil) {
		return "","", common.ERR_INVALID_SRC_REG_ID
	}
	if !checkCdpMinTxFee(param.Fees) {
		return "","", common.ERR_FEE_SMALLER_MIN
	}
	tx.Fees = uint64(param.Fees)

	tx.TxType = CDP_LIQUIDATE_TX
	tx.Version = TX_VERSION
	if param.ScoinsLiquidate < 0 {
		return "","", common.ERR_CDP_STAKE_NUMBER
	}
	tx.ScoinsLiquidate = uint64(param.ScoinsLiquidate)

	if (param.FeeSymbol == "") {
		tx.FeeSymbol = string(common.WICC)
	} else {
		tx.FeeSymbol = string(param.FeeSymbol)
	}

	if param.AssetSymbol==""{
		return "","",common.ERR_COIN_TYPE
	}
	tx.AseetSymbol=param.AssetSymbol

	txHash, err := hex.DecodeString(param.CdpTxid)
	if (err != nil) {
		return "","", common.ERR_CDP_TX_HASH
	}
	tx.CdpTxHash = txHash

	return tx.createRawTx(wifKey)
}

type WaykiDexSellLimitTx struct {
	WaykiBaseSignTx
	Fees   uint64
	FeeSymbol string      //Fee Type (WICC/WUSD)
	CoinSymbol string   //From Coin Type
	AssetSymbol string
	AssetAmount uint64
	AskPrice uint64
	DestId *UserIdWraper //< the dest id(reg id or address or public key) received the wicc values
}

//return rawtx + txid
func (tx WaykiDexSellLimitTx) createRawTx(wifKey *btcutil.WIF) (string,string,error){
	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)

	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.Version)
	writer.WriteVarInt(tx.ValidHeight)
	if(tx.UserId!=nil){
		writer.WriteUserId(tx.UserId)
	}else if(tx.PubKey!=nil){
		writer.WritePubKeyId(tx.PubKey)
	}
	writer.WriteString(tx.FeeSymbol)
	writer.WriteVarInt(int64(tx.Fees))
	writer.WriteString(tx.CoinSymbol)
	writer.WriteString(tx.AssetSymbol)
	writer.WriteVarInt(int64(tx.AssetAmount))
	writer.WriteVarInt(int64(tx.AskPrice))
	signatureBytes,txid,err := tx.CalSignatureTxid(wifKey)
	if err != nil{
		return "","",err
	}
	writer.WriteBytes(signatureBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return rawTx, txid,nil
}

//return signature and txid
func (tx WaykiDexSellLimitTx) CalSignatureTxid(wifKey *btcutil.WIF) ([]byte,string,error) {

	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)

	writer.WriteVarInt(tx.Version)
	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.ValidHeight)
	if (tx.UserId != nil) {
		writer.WriteUserId(tx.UserId)
	} else if (tx.PubKey != nil) {
		writer.WritePubKeyId(tx.PubKey)
	}
	writer.WriteString(tx.FeeSymbol)
	writer.WriteVarInt(int64(tx.Fees))
	writer.WriteString(tx.CoinSymbol)
	writer.WriteString(tx.AssetSymbol)
	writer.WriteVarInt(int64(tx.AssetAmount))
	writer.WriteVarInt(int64(tx.AskPrice))

	return calSignatureTxid(buf.Bytes(), wifKey)
}
// returns the signature hex string and nil error, or returns empty string and the error if it has error
func (param *DexLimitSellTxParam) CreateRawTx(privateKey string) (string, string,error) {
	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return "","", common.ERR_INVALID_PRIVATEKEY
	}
	var tx WaykiDexSellLimitTx

	if param.ValidHeight < 0 {
		return "","", common.ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight

	tx.UserId = parseRegId(param.SrcRegId)

	pubKey, err := hex.DecodeString(param.PubKey)
	if (err != nil) {
		return "","", common.ERR_USER_PUBLICKEY
	}
	tx.PubKey = pubKey
	if (tx.UserId == nil && tx.PubKey == nil) {
		return "","", common.ERR_INVALID_SRC_REG_ID
	}
	if !checkMinTxFee(param.Fees) {
		return "","", common.ERR_FEE_SMALLER_MIN
	}
	tx.Fees = uint64(param.Fees)

	tx.TxType = DEX_SELL_LIMIT_ORDER_TX
	tx.Version = TX_VERSION
	if param.Price <= 0 {
		return "","", common.ERR_ASK_PRICE
	}
	tx.AskPrice = uint64(param.Price)
	if param.CoinSymbol == "" || param.AssetSymbol == "" {
		return "","", common.ERR_COIN_TYPE
	}
	tx.AssetSymbol = param.AssetSymbol
	tx.CoinSymbol = param.CoinSymbol
	if (param.FeeSymbol == "") {
		tx.FeeSymbol = string(common.WICC)
	} else {
		tx.FeeSymbol = string(param.FeeSymbol)
	}
	if !checkMoneyRange(param.AssetAmount) {
		return "","", common.ERR_RANGE_VALUES
	}
	tx.AssetAmount = uint64(param.AssetAmount)

	return tx.createRawTx(wifKey)
}

// returns the signature hex string and nil error, or returns empty string and the error if it has error
func (param *DexLimitBuyTxParam) CreateRawTx(privateKey string) (string,string, error) {
	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return "","", common.ERR_INVALID_PRIVATEKEY
	}
	var tx WaykiDexSellLimitTx

	if param.ValidHeight < 0 {
		return "","", common.ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight

	tx.UserId = parseRegId(param.SrcRegId)

	pubKey, err := hex.DecodeString(param.PubKey)
	if (err != nil) {
		return "","", common.ERR_USER_PUBLICKEY
	}
	tx.PubKey = pubKey
	if (tx.UserId == nil && tx.PubKey == nil) {
		return "","", common.ERR_INVALID_SRC_REG_ID
	}
	if !checkCdpMinTxFee(param.Fees) {
		return "","", common.ERR_FEE_SMALLER_MIN
	}
	tx.Fees = uint64(param.Fees)

	tx.TxType = DEX_BUY_LIMIT_ORDER_TX
	tx.Version = TX_VERSION
	if param.Price <= 0 {
		return "", "",common.ERR_ASK_PRICE
	}
	tx.AskPrice = uint64(param.Price)
	if param.CoinSymbol == "" || param.AssetSymbol == "" {
		return "","", common.ERR_COIN_TYPE
	}
	tx.AssetSymbol = param.AssetSymbol
	tx.CoinSymbol = param.CoinSymbol
	if (param.FeeSymbol == "") {
		tx.FeeSymbol = string(common.WICC)
	} else {
		tx.FeeSymbol = string(param.FeeSymbol)
	}
	if !checkMoneyRange(param.AssetAmount) {
		return "","", common.ERR_RANGE_VALUES
	}
	tx.AssetAmount = uint64(param.AssetAmount)

	return tx.createRawTx(wifKey)
}

type WaykiDexMarketTx struct {
	WaykiBaseSignTx
	Fees   uint64
	FeeSymbol string      //Fee Type (WICC/WUSD)
	CoinSymbol string   //From Coin Type
	AssetSymbol string
	AssetAmount uint64
	DestId *UserIdWraper //< the dest id(reg id or address or public key) received the wicc values
}

//return rawtx + txid
func (tx WaykiDexMarketTx) createRawTx(wifKey *btcutil.WIF) (string,string,error){
	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)

	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.Version)
	writer.WriteVarInt(tx.ValidHeight)
	if(tx.UserId!=nil){
		writer.WriteUserId(tx.UserId)
	}else if(tx.PubKey!=nil){
		writer.WritePubKeyId(tx.PubKey)
	}
	writer.WriteString(tx.FeeSymbol)
	writer.WriteVarInt(int64(tx.Fees))
	writer.WriteString(tx.CoinSymbol)
	writer.WriteString(tx.AssetSymbol)
	writer.WriteVarInt(int64(tx.AssetAmount))
	signatureBytes,txid,err := tx.CalSignatureTxid(wifKey)
	if err != nil{
		return "","",err
	}
	writer.WriteBytes(signatureBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return rawTx, txid,nil
}

//return signature and txid
func (tx WaykiDexMarketTx) CalSignatureTxid(wifKey *btcutil.WIF) ([]byte,string,error){

	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)

	writer.WriteVarInt(tx.Version)
	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.ValidHeight)
	if(tx.UserId!=nil){
		writer.WriteUserId(tx.UserId)
	}else if(tx.PubKey!=nil){
		writer.WritePubKeyId(tx.PubKey)
	}
	writer.WriteString(tx.FeeSymbol)
	writer.WriteVarInt(int64(tx.Fees))
	writer.WriteString(tx.CoinSymbol)
	writer.WriteString(tx.AssetSymbol)
	writer.WriteVarInt(int64(tx.AssetAmount))

	return calSignatureTxid(buf.Bytes(),wifKey)
}
// returns the signature hex string and nil error, or returns empty string and the error if it has error
func (param *DexMarketSellTxParam) CreateRawTx(privateKey string) (string, string,error) {
	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return "","", common.ERR_INVALID_PRIVATEKEY
	}
	var tx WaykiDexMarketTx

	if param.ValidHeight < 0 {
		return "","", common.ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight

	tx.UserId = parseRegId(param.SrcRegId)

	pubKey, err := hex.DecodeString(param.PubKey)
	if (err != nil) {
		return "","", common.ERR_USER_PUBLICKEY
	}
	tx.PubKey = pubKey
	if (tx.UserId == nil && tx.PubKey == nil) {
		return "","", common.ERR_INVALID_SRC_REG_ID
	}
	if !checkCdpMinTxFee(param.Fees) {
		return "","", common.ERR_FEE_SMALLER_MIN
	}
	tx.Fees = uint64(param.Fees)

	tx.TxType = DEX_SELL_MARKET_ORDER_TX
	tx.Version = TX_VERSION
	if param.CoinSymbol == "" || param.AssetSymbol == "" {
		return "", "",common.ERR_COIN_TYPE
	}
	tx.AssetSymbol = param.AssetSymbol
	tx.CoinSymbol = param.CoinSymbol
	if (param.FeeSymbol == "") {
		tx.FeeSymbol = string(common.WICC)
	} else {
		tx.FeeSymbol = string(param.FeeSymbol)
	}
	if !checkMoneyRange(param.AssetAmount) {
		return "","", common.ERR_RANGE_VALUES
	}
	tx.AssetAmount = uint64(param.AssetAmount)

	return tx.createRawTx(wifKey)
}

// returns the signature hex string and nil error, or returns empty string and the error if it has error
func (param *DexMarketBuyTxParam) CreateRawTx(privateKey string) (string,string, error) {
	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return "","", common.ERR_INVALID_PRIVATEKEY
	}
	var tx WaykiDexMarketTx

	if param.ValidHeight < 0 {
		return "","", common.ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight

	tx.UserId = parseRegId(param.SrcRegId)

	pubKey, err := hex.DecodeString(param.PubKey)
	if (err != nil) {
		return "","", common.ERR_USER_PUBLICKEY
	}
	tx.PubKey = pubKey
	if (tx.UserId == nil && tx.PubKey == nil) {
		return "","", common.ERR_INVALID_SRC_REG_ID
	}
	if !checkCdpMinTxFee(param.Fees) {
		return "","", common.ERR_FEE_SMALLER_MIN
	}
	tx.Fees = uint64(param.Fees)

	tx.TxType = DEX_BUY_MARKET_ORDER_TX
	tx.Version = TX_VERSION
	if param.CoinSymbol == "" || param.AssetSymbol == "" {
		return "","", common.ERR_COIN_TYPE
	}
	tx.AssetSymbol = param.AssetSymbol
	tx.CoinSymbol = param.CoinSymbol
	if (param.FeeSymbol == "") {
		tx.FeeSymbol = string(common.WICC)
	} else {
		tx.FeeSymbol = string(param.FeeSymbol)
	}
	if !checkMoneyRange(param.CoinAmount) {
		return "","", common.ERR_RANGE_VALUES
	}
	tx.AssetAmount = uint64(param.CoinAmount)

	return tx.createRawTx(wifKey)
}

type WaykiDexCancelTx struct {
	WaykiBaseSignTx
	Fees   uint64
	FeeSymbol string      //Fee Type (WICC/WUSD)
	DexHash []byte   //From Coin Type
}

//return rawtx + txid
func (tx WaykiDexCancelTx) createRawTx(wifKey *btcutil.WIF) (string,string,error){
	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)

	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.Version)
	writer.WriteVarInt(tx.ValidHeight)
	if(tx.UserId!=nil){
		writer.WriteUserId(tx.UserId)
	}else if(tx.PubKey!=nil){
		writer.WritePubKeyId(tx.PubKey)
	}
	writer.WriteString(tx.FeeSymbol)
	writer.WriteVarInt(int64(tx.Fees))
	writer.WriteReverse(tx.DexHash)
	signatureBytes,txid,err := tx.CalSignatureTxid(wifKey)
	if err != nil{
		return "","",err
	}
	writer.WriteBytes(signatureBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return rawTx, txid,nil
}

//return signature and txid
func (tx WaykiDexCancelTx) CalSignatureTxid(wifKey *btcutil.WIF) ([]byte,string,error){

	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)

	writer.WriteVarInt(tx.Version)
	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.ValidHeight)
	if(tx.UserId!=nil){
		writer.WriteUserId(tx.UserId)
	}else if(tx.PubKey!=nil){
		writer.WritePubKeyId(tx.PubKey)
	}
	writer.WriteString(tx.FeeSymbol)
	writer.WriteVarInt(int64(tx.Fees))
	writer.WriteReverse(tx.DexHash)
	return calSignatureTxid(buf.Bytes(),wifKey)
}
// returns the signature hex string and nil error, or returns empty string and the error if it has error
func (param *DexCancelTxParam) CreateRawTx(privateKey string) (string, string,error) {
	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return "","", common.ERR_INVALID_PRIVATEKEY
	}
	var tx WaykiDexCancelTx

	if param.ValidHeight < 0 {
		return "","", common.ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight

	tx.UserId = parseRegId(param.SrcRegId)

	pubKey, err := hex.DecodeString(param.PubKey)
	if (err != nil) {
		return "","", common.ERR_USER_PUBLICKEY
	}
	tx.PubKey = pubKey
	if (tx.UserId == nil && tx.PubKey == nil) {
		return "","", common.ERR_INVALID_SRC_REG_ID
	}
	if !checkCdpMinTxFee(param.Fees) {
		return "","", common.ERR_FEE_SMALLER_MIN
	}
	tx.Fees = uint64(param.Fees)

	tx.TxType = DEX_CANCEL_ORDER_TX
	tx.Version = TX_VERSION
	if (param.FeeSymbol == "") {
		tx.FeeSymbol = string(common.WICC)
	} else {
		tx.FeeSymbol = string(param.FeeSymbol)
	}
	txHash, err := hex.DecodeString(param.DexTxid)
	if (err != nil) {
		return "","", common.ERR_CDP_TX_HASH
	}
	tx.DexHash = txHash

	return tx.createRawTx(wifKey)
}

type WaykiAssetIssueTx struct {
	WaykiBaseSignTx
	Fees   uint64
	FeeSymbol string      //Fee Type (WICC/WUSD)
	AssetSymbol string   //From Coin Type
	AssetName   string
	AssetTotal   uint64
	AssetOwner   *UserIdWraper
	MinTable     bool
}

//return rawtx + txid
func (tx WaykiAssetIssueTx) createRawTx(wifKey *btcutil.WIF) (string,string,error){
	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)

	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.Version)
	writer.WriteVarInt(tx.ValidHeight)
	if(tx.UserId!=nil){
		writer.WriteUserId(tx.UserId)
	}
	writer.WriteString(tx.FeeSymbol)
	writer.WriteVarInt(int64(tx.Fees))
	writer.WriteString(tx.AssetSymbol)
	writer.WriteUserId(tx.AssetOwner)
	writer.WriteString(tx.AssetName)
	if(tx.MinTable){
		writer.WriteByte(1)
	}else {
		writer.WriteByte(0)
	}
	writer.WriteVarInt(int64(tx.AssetTotal))

	signatureBytes,txid,err := tx.CalSignatureTxid(wifKey)
	if err != nil{
		return "","",err
	}
	writer.WriteBytes(signatureBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return rawTx, txid,nil
}

//return signature and txid
func (tx WaykiAssetIssueTx) CalSignatureTxid(wifKey *btcutil.WIF) ([]byte,string,error){
	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)

	writer.WriteVarInt(tx.Version)
	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.ValidHeight)
	if(tx.UserId!=nil){
		writer.WriteUserId(tx.UserId)
	}
	writer.WriteString(tx.FeeSymbol)
	writer.WriteVarInt(int64(tx.Fees))
	writer.WriteString(tx.AssetSymbol)
	writer.WriteUserId(tx.AssetOwner)
	writer.WriteString(tx.AssetName)
	if(tx.MinTable){
		writer.WriteByte(1)
	}else {
		writer.WriteByte(0)
	}
	writer.WriteVarInt(int64(tx.AssetTotal))

	return calSignatureTxid(buf.Bytes(),wifKey)
}

// returns the signature hex string and nil error, or returns empty string and the error if it has error
func (param *AssetIssueTxParam) CreateRawTx(privateKey string) (string,string, error) {

	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return "","", common.ERR_INVALID_PRIVATEKEY
	}
	var tx WaykiAssetIssueTx

	if param.ValidHeight < 0 {
		return "","", common.ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight
	tx.UserId = parseRegId(param.SrcRegId)

	if (tx.UserId == nil) {
		return "","", common.ERR_INVALID_SRC_REG_ID
	}

	tx.Fees = uint64(param.Fees)

	tx.TxType = ASSET_ISSUE_TX
	tx.Version = TX_VERSION

	if !checkAssetSymbol(param.AssetSymbol) {
		return "","", common.ERR_SYMBOL_ERROR
	}
	tx.AssetSymbol = param.AssetSymbol
	if (param.FeeSymbol == "") {
		tx.FeeSymbol = string(common.WICC)
	} else {
		tx.FeeSymbol = string(param.FeeSymbol)
	}
	tx.MinTable = param.MinTable
	if (param.AssetName == "") {
		return "","", common.ERR_ASSET_NAME_ERROR
	}
	tx.AssetName = param.AssetName
	if (param.AssetTotal < 100000000) {
		return "","", common.ERR_TOTAl_SUPPLY_ERROR
	}
	tx.AssetTotal = uint64(param.AssetTotal)
	tx.AssetOwner = parseUserId(param.AssetOwner)
	if (tx.AssetOwner == nil) {
		return "","", common.ERR_ASSET_UPDATE_OWNER_ERROR
	}

	return tx.createRawTx(wifKey)
}

type WaykiAssetUpdateTx struct {
	WaykiBaseSignTx
	Fees   uint64
	UpdateType int
	FeeSymbol string      //Fee Type (WICC/WUSD)
	AssetSymbol string   //From Coin Type
	AssetName   string
	AssetTotal   uint64
	AssetOwner   *UserIdWraper
}

//return rawtx + txid
func (tx WaykiAssetUpdateTx) createRawTx(wifKey *btcutil.WIF) (string,string,error){
	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)

	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.Version)
	writer.WriteVarInt(tx.ValidHeight)
	if(tx.UserId!=nil){
		writer.WriteUserId(tx.UserId)
	}
	writer.WriteString(tx.FeeSymbol)
	writer.WriteVarInt(int64(tx.Fees))
	writer.WriteString(tx.AssetSymbol)
	writer.WriteByte(byte(tx.UpdateType))
	switch tx.UpdateType {
	case 1:
		writer.WriteUserId(tx.AssetOwner)
		break
	case 2:
		writer.WriteString(tx.AssetName)
		break
	case 3:
		writer.WriteVarInt(int64(tx.AssetTotal))
		break
	}

	signatureBytes,txid,err := tx.CalSignatureTxid(wifKey)
	if err != nil{
		return "","",err
	}
	writer.WriteBytes(signatureBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return rawTx, txid,nil
}

//return signature and txid
func (tx WaykiAssetUpdateTx) CalSignatureTxid(wifKey *btcutil.WIF) ([]byte,string,error){

	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)

	writer.WriteVarInt(tx.Version)
	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.ValidHeight)
	if(tx.UserId!=nil){
		writer.WriteUserId(tx.UserId)
	}
	writer.WriteString(tx.FeeSymbol)
	writer.WriteVarInt(int64(tx.Fees))
	writer.WriteString(tx.AssetSymbol)
	writer.WriteByte(byte(tx.UpdateType))
	switch tx.UpdateType {
	case 1:
		writer.WriteUserId(tx.AssetOwner)
		break
	case 2:
		writer.WriteString(tx.AssetName)
		break
	case 3:
		writer.WriteVarInt(int64(tx.AssetTotal))
		break
	}

	return calSignatureTxid(buf.Bytes(),wifKey)
}
// returns the signature hex string and nil error, or returns empty string and the error if it has error
func (param *AssetUpdateTxParam) CreateRawTx(privateKey string) (string, string,error) {

	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return "","", common.ERR_INVALID_PRIVATEKEY
	}
	var tx WaykiAssetUpdateTx

	if param.ValidHeight < 0 {
		return "","", common.ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight
	tx.UserId = parseRegId(param.SrcRegId)
	if (tx.UserId == nil) {
		return "","", common.ERR_INVALID_SRC_REG_ID
	}

	tx.Fees = uint64(param.Fees)

	tx.TxType = ASSET_UPDATE_TX
	tx.Version = TX_VERSION

	if !checkAssetSymbol(param.AssetSymbol) {
		return "","", common.ERR_SYMBOL_ERROR
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
			return "","", common.ERR_ASSET_UPDATE_OWNER_ERROR
		}
		break
	case 2:
		if (param.AssetName == "") {
			return "","", common.ERR_ASSET_NAME_ERROR
		}
		tx.AssetName = param.AssetName
		break
	case 3:
		if (param.AssetTotal < 100000000) {
			return "","", common.ERR_TOTAl_SUPPLY_ERROR
		}
		tx.AssetTotal = uint64(param.AssetTotal)
		break
	default:
		return "","", common.ERR_ASSET_UPDATE_TYPE_ERROR
	}

	return tx.createRawTx(wifKey)
}

type WaykiRegisterAccountTx struct {
	WaykiBaseSignTx
	MinerId *UserIdWraper
	Fees    uint64
}

//return rawtx + txid
func (tx WaykiRegisterAccountTx) createRawTx(wifKey *btcutil.WIF) (string,string,error){
	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)
	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.Version)
	writer.WriteVarInt(tx.ValidHeight)
	writer.WriteUserId(tx.UserId)
	writer.WriteUserId(tx.MinerId)

	writer.WriteVarInt(int64(tx.Fees))
	signatureBytes,txid,err := tx.CalSignatureTxid(wifKey)
	if err != nil{
		return "","",err
	}
	writer.WriteBytes(signatureBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return rawTx, txid,nil
}

//return signature and txid
func (tx WaykiRegisterAccountTx) CalSignatureTxid(wifKey *btcutil.WIF) ([]byte,string,error){
	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)
	writer.WriteVarInt(tx.Version)
	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.ValidHeight)
	writer.WriteUserId(tx.UserId)
	writer.WriteUserId(tx.MinerId)
	writer.WriteVarInt(int64(tx.Fees))

	return calSignatureTxid(buf.Bytes(),wifKey)
}

// returns the signature hex string and nil error, or returns empty string and the error if it has error
func (param *RegisterAccountTxParam) CreateRawTx(privateKey string) (string,string, error) {
	// check and convert params
	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return "","", common.ERR_INVALID_PRIVATEKEY
	}

	var tx WaykiRegisterAccountTx
	if param.ValidHeight < 0 {
		return "","", common.ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight

	if !checkMinTxFee(param.Fees) {
		return "","", common.ERR_FEE_SMALLER_MIN
	}
	tx.Fees = uint64(param.Fees)

	tx.UserId = NewPubKeyUid(wifKey.PrivKey.PubKey().SerializeCompressed())
	tx.TxType = REG_ACCT_TX
	tx.Version = 1

	return tx.createRawTx(wifKey)
}

type WaykiCommonTx struct {
	WaykiBaseSignTx
	Fees   uint64
	Values uint64
	Memo   string
	DestId *UserIdWraper //< the dest id(reg id or address or public key) received the wicc values
}

//return rawtx + txid
func (tx WaykiCommonTx) createRawTx(wifKey *btcutil.WIF) (string,string,error){
	//uid := ParseRegId(tx.UserId)
	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)

	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.Version)
	writer.WriteVarInt(tx.ValidHeight)
	if(tx.UserId!=nil){
		writer.WriteUserId(tx.UserId)
	}else if(tx.PubKey!=nil){
		writer.WritePubKeyId(tx.PubKey)
	}
	writer.WriteUserId(tx.DestId)
	writer.WriteVarInt(int64(tx.Fees))
	writer.WriteVarInt(int64(tx.Values))
	writer.WriteString(tx.Memo)

	signatureBytes,txid,err := tx.CalSignatureTxid(wifKey)
	if err != nil{
		return "","",err
	}
	writer.WriteBytes(signatureBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return rawTx, txid,nil
}

//return signature and txid
func (tx WaykiCommonTx) CalSignatureTxid(wifKey *btcutil.WIF) ([]byte,string,error){

	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)

	writer.WriteVarInt(tx.Version)
	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.ValidHeight)
	if(tx.UserId!=nil){
		writer.WriteUserId(tx.UserId)
	}else if(tx.PubKey!=nil){
		writer.WritePubKeyId(tx.PubKey)
	}
	writer.WriteUserId(tx.DestId)
	writer.WriteVarInt(int64(tx.Fees))
	writer.WriteVarInt(int64(tx.Values))
	writer.WriteString(tx.Memo)

	return calSignatureTxid(buf.Bytes(),wifKey)
}
// returns the signature hex string and nil error, or returns empty string and the error if it has error
func (param *CommonTxParam) CreateRawTx(privateKey string) (string,string, error) {

	// check and convert params
	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return "","", common.ERR_INVALID_PRIVATEKEY
	}

	var tx WaykiCommonTx

	if param.ValidHeight < 0 {
		return "","", common.ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight

	tx.UserId = parseRegId(param.SrcRegId)

	tx.DestId = parseUserId(param.DestAddr)
	if tx.DestId == nil {
		return "","", common.ERR_INVALID_DEST_ADDR
	}

	if !checkMoneyRange(param.Values) {
		return "","", common.ERR_RANGE_VALUES
	}
	tx.Values = uint64(param.Values)

	if !checkMoneyRange(param.Fees) {
		return "","", common.ERR_RANGE_FEE
	}

	if !checkMinTxFee(param.Fees) {
		return "","", common.ERR_FEE_SMALLER_MIN
	}
	tx.Fees = uint64(param.Fees)

	pubKey, err := hex.DecodeString(param.PubKey)
	if (err != nil) {
		return "","", common.ERR_USER_PUBLICKEY
	}
	tx.PubKey = pubKey
	if (tx.UserId == nil && tx.PubKey == nil) {
		return "","", common.ERR_INVALID_SRC_REG_ID
	}
	tx.TxType = COMMON_TX
	tx.Version = TX_VERSION
	tx.Memo = param.Memo

	return tx.createRawTx(wifKey)
}

type WaykiCallContractTx struct {
	WaykiBaseSignTx
	AppId    *UserIdWraper //user regid or user key id or app regid
	Fees     uint64
	Values   uint64 //transfer amount
	Contract []byte
}

//return rawtx + txid
func (tx WaykiCallContractTx) createRawTx(wifKey *btcutil.WIF) (string,string,error){
	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)

	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.Version)
	writer.WriteVarInt(tx.ValidHeight)
	if(tx.UserId!=nil){
		writer.WriteUserId(tx.UserId)
	}else if(tx.PubKey!=nil){
		writer.WritePubKeyId(tx.PubKey)
	}
	writer.WriteUserId(tx.AppId)
	writer.WriteVarInt(int64(tx.Fees))
	writer.WriteVarInt(int64(tx.Values))
	writer.WriteBytes(tx.Contract)

	signatureBytes,txid,err := tx.CalSignatureTxid(wifKey)
	if err != nil{
		return "","",err
	}
	writer.WriteBytes(signatureBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return rawTx, txid,nil
}

//return signature and txid
func (tx WaykiCallContractTx) CalSignatureTxid(wifKey *btcutil.WIF) ([]byte,string,error){

	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)

	writer.WriteVarInt(tx.Version)
	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.ValidHeight)
	if(tx.UserId!=nil){
		writer.WriteUserId(tx.UserId)
	}else if(tx.PubKey!=nil){
		writer.WritePubKeyId(tx.PubKey)
	}
	writer.WriteUserId(tx.AppId)
	writer.WriteVarInt(int64(tx.Fees))
	writer.WriteVarInt(int64(tx.Values))
	writer.WriteBytes(tx.Contract)
	return calSignatureTxid(buf.Bytes(),wifKey)
}
// returns the signature hex string and nil error, or returns empty string and the error if it has error
func (param *CallContractTxParam) CreateRawTx(privateKey string) (string,string,error) {
	// check and convert params
	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return "","", common.ERR_INVALID_PRIVATEKEY
	}

	var tx  WaykiCallContractTx

	if param.ValidHeight < 0 {
		return "","", common.ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight

	tx.UserId = parseRegId(param.SrcRegId)

	tx.AppId = parseRegId(param.AppId)
	if tx.AppId == nil {
		return "","", common.ERR_INVALID_APP_ID
	}

	if !checkMoneyRange(param.Values) {
		return "","", common.ERR_RANGE_VALUES
	}
	tx.Values = uint64(param.Values)

	if !checkMoneyRange(param.Fees) {
		return "","", common.ERR_RANGE_FEE
	}
	if !checkMinTxFee(param.Fees) {
		return "","", common.ERR_FEE_SMALLER_MIN
	}
	tx.Fees = uint64(param.Fees)

	binary, errHex := hex.DecodeString(param.ContractHex)
	if errHex != nil {
		return "","", common.ERR_INVALID_CONTRACT_HEX
	}
	tx.Contract = []byte(binary)
	pubKey, err := hex.DecodeString(param.PubKey)
	if (err != nil) {
		return "","", common.ERR_USER_PUBLICKEY
	}
	tx.PubKey = pubKey
	if (tx.UserId == nil && tx.PubKey == nil) {
		return "","", common.ERR_INVALID_SRC_REG_ID
	}
	tx.TxType = CONTRACT_TX
	tx.Version = TX_VERSION

	return tx.createRawTx(wifKey)
}

type WaykiRegisterContractTx struct {
	WaykiBaseSignTx
	Script      []byte
	Description string
	Fees        uint64
}

//return rawtx + txid
func (tx WaykiRegisterContractTx) createRawTx(wifKey *btcutil.WIF) (string,string,error){

	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)
	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.Version)
	writer.WriteVarInt(tx.ValidHeight)
	writer.WriteUserId(tx.UserId)
	WriteContractScript(writer, tx.Script, tx.Description)

	writer.WriteVarInt(int64(tx.Fees))
	signatureBytes,txid,err := tx.CalSignatureTxid(wifKey)
	if err != nil{
		return "","",err
	}
	writer.WriteBytes(signatureBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return rawTx, txid,nil
}

//return signature and txid
func (tx WaykiRegisterContractTx) CalSignatureTxid(wifKey *btcutil.WIF) ([]byte,string,error){
	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)
	writer.WriteVarInt(tx.Version)
	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.ValidHeight)
	writer.WriteUserId(tx.UserId)
	WriteContractScript(writer, tx.Script, tx.Description)
	writer.WriteVarInt(int64(tx.Fees))

	return calSignatureTxid(buf.Bytes(),wifKey)
}

// returns the signature hex string and nil error, or returns empty string and the error if it has error
func (param *RegisterContractTxParam) CreateRawTx(privateKey string) (string,string, error) {

	// check and convert params
	wifKey, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return "","", common.ERR_INVALID_PRIVATEKEY
	}

	var tx WaykiRegisterContractTx

	if param.ValidHeight < 0 {
		return "","", common.ERR_NEGATIVE_VALID_HEIGHT
	}
	tx.ValidHeight = param.ValidHeight

	tx.UserId = parseRegId(param.SrcRegId)

	if !checkMoneyRange(param.Fees) {
		return "","", common.ERR_RANGE_FEE
	}
	if !checkMinTxFee(param.Fees) {
		return "","", common.ERR_FEE_SMALLER_MIN
	}
	tx.Fees = uint64(param.Fees)

	tx.TxType = REG_CONT_TX
	tx.Version = TX_VERSION

	if len(param.Script) == 0 || len(param.Script) > CONTRACT_SCRIPT_MAX_SIZE {
		return "","", common.ERR_INVALID_SCRIPT
	}
	tx.Script = param.Script

	if len(param.Description) > CONTRACT_SCRIPT_DESC_MAX_SIZE {
		return "","", common.ERR_INVALID_SCRIPT_DESC
	}
	tx.Description = param.Description
	if (tx.UserId == nil && tx.PubKey == nil) {
		return "","", common.ERR_INVALID_SRC_REG_ID
	}

	return tx.createRawTx(wifKey)
}

type WaykiRewardTx struct {
	WaykiBaseSignTx
	Values uint64 // reward values
}

//return rawtx + txid
func (tx WaykiRewardTx) createRawTx(wifKey *btcutil.WIF) (string,string,error) {

	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)

	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.Version)
	writer.WriteUserId(tx.UserId)
	writer.WriteVarInt(int64(tx.Values))
	writer.WriteVarInt(tx.ValidHeight)

	signatureBytes,txid,err := tx.CalSignatureTxid(wifKey)
	if err != nil{
		return "","",err
	}
	writer.WriteBytes(signatureBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return rawTx, txid,nil
}

//return signature and txid
func (tx WaykiRewardTx) CalSignatureTxid(wifKey *btcutil.WIF)  ([]byte,string,error) {

	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)

	writer.WriteVarInt(tx.Version)
	writer.WriteByte(byte(tx.TxType))
	writer.WriteUserId(tx.UserId)
	writer.WriteVarInt(int64(tx.Values))
	writer.WriteVarInt(tx.ValidHeight)

	return calSignatureTxid(buf.Bytes(),wifKey)
}