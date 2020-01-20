package waykichain

import (
	"bytes"
	"encoding/hex"
	"github.com/WaykiChain/wicc-wallet-utils-go/common"
	"github.com/WaykiChain/wicc-wallet-utils-go/common/hash"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
	"regexp"
	"strconv"
	"strings"
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
	UCONTRACT_DEPLOY_TX         = 14   //!< universal VM contract deployment
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

func calSignatureTxid(paramHash []byte,wifKey *btcutil.WIF)([]byte,string,error){
	txidRevrse := hash.DoubleHash256(paramHash)
	txidHex := common.Reverse(txidRevrse)
	signature, err := wifKey.PrivKey.Sign(txidRevrse)
	if err != nil{
		return nil,"",err
	}
	return signature.Serialize(),hex.EncodeToString(txidHex),nil
}

//Base Params for Tx
type WaykiBaseSignTx struct {
	TxType      WaykiTxType
	Version     int64
	ValidHeight int64
	PubKey      []byte
	UserId      *UserIdWraper // current operating user id, the one want to do something
}

type SignResult struct{
	RawTx string
	Txid  string
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
func (tx WaykiUCoinTransferTx) CreateRawTxid(wifKey *btcutil.WIF) (* SignResult,error){

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
		return nil,err
	}
	writer.WriteBytes(signatureBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return &SignResult{rawTx,txid},nil
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

func (tx WaykiUCoinCallContractTx) CreateRawTxid(wifKey *btcutil.WIF) (* SignResult,error) {

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
		return nil,err
	}
	writer.WriteBytes(signatureBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return &SignResult{rawTx,txid},nil
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
func (tx WaykiUCoinRegisterContractTx) CreateRawTxid(wifKey *btcutil.WIF) (* SignResult,error) {

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
		return nil,err
	}
	writer.WriteBytes(signatureBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return &SignResult{rawTx,txid},nil
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
func (tx WaykiDelegateTx) CreateRawTxid(wifKey *btcutil.WIF) (* SignResult,error) {

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
		return nil,err
	}
	writer.WriteBytes(signatureBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return &SignResult{rawTx,txid},nil
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

// cdp stake Asset
type AssetModel struct {
	AssetAmount    int64 //stake asset amount
	AssetSymbol string  //stake asset symbol
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
func (tx WaykiCdpStakeTx) CreateRawTxid(wifKey *btcutil.WIF) (* SignResult,error){
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
		return nil,err
	}
	writer.WriteBytes(signatureBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return &SignResult{rawTx,txid},nil
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

type WaykiCdpRedeemTx struct {
	WaykiBaseSignTx
	Fees   uint64
	ScoinValues uint64   //Stake Coin
	Assets   []AssetModel
	FeeSymbol string      //Fee Type (WICC/WUSD)
	CdpTxHash []byte
}

//return rawtx + txid
func (tx WaykiCdpRedeemTx) CreateRawTxid(wifKey *btcutil.WIF) (* SignResult,error){
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
		return nil,err
	}
	writer.WriteBytes(signatureBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return &SignResult{rawTx,txid},nil
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

type WaykiCdpLiquidateTx struct {
	WaykiBaseSignTx
	Fees   uint64
	ScoinsLiquidate uint64   //Scoin  Liquidate
	FeeSymbol string      //Fee Type (WICC/WUSD)
	CdpTxHash []byte
	AseetSymbol string
}
//return rawtx + txid
func (tx WaykiCdpLiquidateTx) CreateRawTxid(wifKey *btcutil.WIF) (* SignResult,error){
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
		return nil,err
	}
	writer.WriteBytes(signatureBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return &SignResult{rawTx,txid},nil
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
func (tx WaykiDexSellLimitTx) CreateRawTxid(wifKey *btcutil.WIF) (* SignResult,error){
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
		return nil,err
	}
	writer.WriteBytes(signatureBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return &SignResult{rawTx,txid},nil
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
func (tx WaykiDexMarketTx) CreateRawTxid(wifKey *btcutil.WIF) (* SignResult,error){
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
		return nil,err
	}
	writer.WriteBytes(signatureBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return &SignResult{rawTx,txid},nil
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

type WaykiDexCancelTx struct {
	WaykiBaseSignTx
	Fees   uint64
	FeeSymbol string      //Fee Type (WICC/WUSD)
	DexHash []byte   //From Coin Type
}

//return rawtx + txid
func (tx WaykiDexCancelTx) CreateRawTxid(wifKey *btcutil.WIF) (* SignResult,error){
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
		return nil,err
	}
	writer.WriteBytes(signatureBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return &SignResult{rawTx,txid},nil
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
func (tx WaykiAssetIssueTx) CreateRawTxid(wifKey *btcutil.WIF) (* SignResult,error){
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
		return nil,err
	}
	writer.WriteBytes(signatureBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return &SignResult{rawTx,txid},nil
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
func (tx WaykiAssetUpdateTx) CreateRawTxid(wifKey *btcutil.WIF) (* SignResult,error){
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
		return nil,err
	}
	writer.WriteBytes(signatureBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return &SignResult{rawTx,txid},nil
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

type WaykiRegisterAccountTx struct {
	WaykiBaseSignTx
	MinerId *UserIdWraper
	Fees    uint64
}

//return rawtx + txid
func (tx WaykiRegisterAccountTx) CreateRawTxid(wifKey *btcutil.WIF) (* SignResult,error){
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
		return nil,err
	}
	writer.WriteBytes(signatureBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return &SignResult{rawTx,txid},nil
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

type WaykiCommonTx struct {
	WaykiBaseSignTx
	Fees   uint64
	Values uint64
	Memo   string
	DestId *UserIdWraper //< the dest id(reg id or address or public key) received the wicc values
}

//return rawtx + txid
func (tx WaykiCommonTx) CreateRawTxid(wifKey *btcutil.WIF) (* SignResult,error){
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
		return nil,err
	}
	writer.WriteBytes(signatureBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return &SignResult{rawTx,txid},nil
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

type WaykiCallContractTx struct {
	WaykiBaseSignTx
	AppId    *UserIdWraper //user regid or user key id or app regid
	Fees     uint64
	Values   uint64 //transfer amount
	Contract []byte
}

//return rawtx + txid
func (tx WaykiCallContractTx) CreateRawTxid(wifKey *btcutil.WIF) (* SignResult,error){
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
		return nil,err
	}
	writer.WriteBytes(signatureBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return &SignResult{rawTx,txid},nil
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

type WaykiRegisterContractTx struct {
	WaykiBaseSignTx
	Script      []byte
	Description string
	Fees        uint64
}

//return rawtx + txid
func (tx WaykiRegisterContractTx) CreateRawTxid(wifKey *btcutil.WIF) (* SignResult,error){

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
		return nil,err
	}
	writer.WriteBytes(signatureBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return &SignResult{rawTx,txid},nil
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

type WaykiRewardTx struct {
	WaykiBaseSignTx
	Values uint64 // reward values
}

//return rawtx + txid
func (tx WaykiRewardTx) CreateRawTxid(wifKey *btcutil.WIF) (* SignResult,error) {

	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)

	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.Version)
	writer.WriteUserId(tx.UserId)
	writer.WriteVarInt(int64(tx.Values))
	writer.WriteVarInt(tx.ValidHeight)

	signatureBytes,txid,err := tx.CalSignatureTxid(wifKey)
	if err != nil{
		return nil,err
	}
	writer.WriteBytes(signatureBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return &SignResult{rawTx,txid},nil
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