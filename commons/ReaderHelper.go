package commons

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
	"strconv"
)

type DecodeBaseParams struct{
	TxType      string
	Version     uint64
	ValidHeight uint64
	UserId      string      // regid/publicKey
	FeeSymbol 	string      //Fee Type (WICC/WUSD)
	Fees   		uint64
	Signature   string
}

type DecodeUCoinTransferParams struct {
	BaseParams  DecodeBaseParams
	Memo        string
	Dests    	[]DecodeUCoinTransferDest
}

type DecodeUContractInvokeParams struct {
	BaseParams    DecodeBaseParams
	ContractRegId string
	CoinSymbol    string
	CoinAmount    uint64
	ContractParam string
}

type DecodeUCoinTransferDest struct {
	CoinSymbol string   //From Coin Type
	CoinAmount uint64
	DestAddr   string
}

// Support UCOIN_TRANSFER_TX and UCOIN_CONTRACT_INVOKE_TX only
func DecodeRawTx(rawTx string,netType int) (interface{},error){
	rawTxBytes, err := hex.DecodeString(rawTx)
	if  err != nil{
		return nil,err
	}
	buf := bytes.NewBuffer(rawTxBytes)
	//交易类型
	txType := ReadVarInt(buf)
	switch txType {
	case UCOIN_TRANSFER_TX:
		result,err := DecodeUCoinTransferTx(buf,netType)
		if err != nil{
			return nil,err
		}
		return result, nil
	case UCOIN_CONTRACT_INVOKE_TX:
		result,err := DecodeUContractInvokeTx(buf)
		if err != nil{
			return nil,err
		}
		return result, nil
	default:
		return nil ,errors.New("This txType have't benn supported ")
	}

}

//Decode UCoinTransfer Rawtx by buf after read txType
func DecodeUCoinTransferTx(buf *bytes.Buffer,netType int) (DecodeUCoinTransferParams,error){
	//version
	version := ReadVarInt(buf)
	//vaildHeight
	vaildHeight  := ReadVarInt(buf)
	//regid/pubkey
	userId := ReadUserId(buf)
	//feeSymbol
	feeSymbol := ReadString(buf)
	//fee
	fee := ReadVarInt(buf)
	//dests
	dests,err := ReadUCoinDestAddr(buf,netType)
	if err != nil{
		return DecodeUCoinTransferParams{},err
	}
	//memo
	memo := ReadString(buf)
	//signature
	signature := ReadHex(buf)

	return DecodeUCoinTransferParams{
		BaseParams  : DecodeBaseParams{"UCOIN_TRANSFER_TX",version,
			vaildHeight,userId,feeSymbol,fee,signature},
		Memo   	    : memo,
		Dests  		: dests,
	},nil
}

//Decode UContractInvoke RawTx by buf after read txType
func DecodeUContractInvokeTx(buf *bytes.Buffer) (DecodeUContractInvokeParams,error){
	//version
	version := ReadVarInt(buf)
	//vaildHeight
	vaildHeight  := ReadVarInt(buf)
	//regid/pubkey
	userId := ReadUserId(buf)
	//The contract regid which be Invoke
	contractRegid := ReadContractRegid(buf)
	//The param which to invoke contract
	contractParam := ReadHex(buf)
	//fee
	fee := ReadVarInt(buf)
	//feeSymbol
	feeSymbol := ReadString(buf)
	//coinSymbol
	coinSymbol := ReadString(buf)
	//coinAmount
	coinAmount := ReadVarInt(buf)
	//signature
	signature := ReadHex(buf)

	return DecodeUContractInvokeParams{
		BaseParams  : DecodeBaseParams{"UCONTRACT_INVOKE_TX",version,
			vaildHeight,userId,feeSymbol,fee,signature},
		ContractRegId: contractRegid,
		CoinSymbol   : coinSymbol,
		CoinAmount    :coinAmount,
		ContractParam :contractParam,
	},nil
}

//return value in buf and number of value bytes
//func ReadVarInt(data []byte) (uint64,int){
func ReadVarInt(buf *bytes.Buffer) (uint64){

	n := uint64(0)
	for i:= 0 ; true ;i++ {
		c := uint64(buf.Bytes()[i])
		n = (n << 7) | (c & 0x7F)
		if ((c & 0x80) != uint64(0)) {
			n ++
		} else {
			len := i+1
			buf.Next(len)
			return n
		}
	}
	return uint64(0)
}

//return publicKey hex : 33bytes
func ReadPubkey(buf *bytes.Buffer) string{
	pubkeyHex := hex.EncodeToString(buf.Bytes()[:33])
	buf.Next(33)
	return pubkeyHex
}

//return regid
func ReadRegid(buf *bytes.Buffer) string{
	height := ReadVarInt(buf)
	index := ReadVarInt(buf)
	return strconv.FormatInt(int64(height),10) + "-" + strconv.FormatInt(int64(index),10)
}

//return regid + publicKey
func ReadUserId(buf *bytes.Buffer) (string) {

	idLen := ReadVarInt(buf)
	if idLen == 33 {//公钥
		return ReadPubkey(buf)
	}else { //regid
		return ReadRegid(buf)
	}
}

//return contract regid
func ReadContractRegid(buf *bytes.Buffer) string{
	ReadVarInt(buf) //contractRegidLen
	return ReadRegid(buf)
}

func ReadString(buf *bytes.Buffer) string{
	stringLen := ReadVarInt(buf)
	data := string(buf.Bytes()[:stringLen])
	buf.Next(int(stringLen))
	return data
}

func ReadHex(buf *bytes.Buffer) string{
	hexLen := ReadVarInt(buf)
	hexString := hex.EncodeToString(buf.Bytes()[:hexLen])
	buf.Next(int(hexLen))
	return hexString
}

func GetAddrFrom20BytePubKeyHash( pubKeyHash []byte, netType int) (string,error){

	if len(pubKeyHash) != ripemd160.Size{
		return "",errors.New("The len of pubKeyHash error!")
	}

	netParams, err := NetworkToChainConfig(Network(netType))
	if (err != nil) {
		fmt.Errorf("invalid network")
		return "",err
	}

	return base58.CheckEncode(pubKeyHash[:ripemd160.Size], netParams.PubKeyHashAddrID),nil
}

func ReadUCoinDestAddr(buf *bytes.Buffer,netType int) ([]DecodeUCoinTransferDest,error){
	//数组数
	size := ReadVarInt(buf)
	dests := make([]DecodeUCoinTransferDest,0)

	for i:=0; i < int(size) ;i ++{
		keyid_len := ReadVarInt(buf)
		keyid := buf.Bytes()[:keyid_len]
		address,err := GetAddrFrom20BytePubKeyHash(keyid,netType)
		if err != nil {
			return nil,err
		}
		buf.Next(int(keyid_len))
		coinSymbol := ReadString(buf)
		transferAmount := ReadVarInt(buf)
		dest := DecodeUCoinTransferDest{coinSymbol, transferAmount, address}
		//dest:=Dest{string(commons.WICC),1000000, "wLKf2NqwtHk3BfzK5wMDfbKYN1SC3weyR4"}
		dests  = append(dests,dest)
	}

	return dests,nil
}


