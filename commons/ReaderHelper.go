package commons

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/JKinGH/go-hdwallet"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
	"strconv"
)

type WaykiDecodeUCoinTransferTx struct {
	TxType      WaykiTxType
	Version     int64
	ValidHeight int64
	UserId      string // regid/publicKey
	Fees   		uint64
	FeeSymbol 	string      //Fee Type (WICC/WUSD)
	Dests    	[]DecodeDest
	Memo        string
	Signature   string
}

type DecodeDest struct {
	CoinSymbol string   //From Coin Type
	CoinAmount uint64
	DestAddr   string
}

func DecodeUCoinTransferTx() WaykiDecodeUCoinTransferTx{
	data := "0b01df390684f0c10c82480457494343bc834002141c758724cc60db35dd387bcf619a478ec3c065f20457494343bc8340142af03ec43eb893039b5dd5bab612d73034cf1b610457555344858c1f00473045022100d68782ebf4059ac26b169ae035ca2a8c1533c4f5639c9fd64445f205d86fbf2c022008b7ed1467ec9321382284ce9d762967a604602a26295f4d569f9a15b643e1db"
	dataBytes,_ := hex.DecodeString(data)
	buf := bytes.NewBuffer(dataBytes)

	//交易类型
	v1 := ReadVarInt(buf)
	fmt.Println("v1=",v1)
	fmt.Println("after txType=",buf.Bytes())

	//版本号
	v2 := ReadVarInt(buf)
	fmt.Println("v2=",v2)
	fmt.Println("after version=",buf.Bytes())

	//有效高度
	v3  := ReadVarInt(buf)
	fmt.Println("v3=",v3)
	fmt.Println("after vaildheight=",buf.Bytes())

	//regid/pubkey
	regid,pubkey := ReadUserId(buf)
	fmt.Println("regid=",regid,"pubkey=",pubkey)
	fmt.Println("after userId=",buf.Bytes())

	//feeSymbol
	feeSymbol := ReadString(buf)
	fmt.Println("feeSymbol=",feeSymbol)
	fmt.Println("after feeSymbol=",buf.Bytes())

	//fee
	fee := ReadVarInt(buf)
	fmt.Println("fee=",fee)
	fmt.Println("after fee=",buf.Bytes())

	//destaddr
	dests,_ := ReadUCoinDestAddr(buf,&hdwallet.WICCTestnetParams)
	for i ,dest := range dests.destArray{
		fmt.Printf("dest[%d]=%+v\n",i,dest)
	}
	fmt.Println("after destaddr=",buf.Bytes())

	//memo
	memo := commons.ReadString(buf)
	fmt.Println("memo=",memo)
	fmt.Println("after memo=",buf.Bytes())

	//signature
	signature := commons.ReadHex(buf)
	fmt.Println("signature=",signature)
	fmt.Println("after signature=",buf.Bytes())
}



//return value in buf and number of value bytes
//func ReadVarInt(data []byte) (uint64,int){
func ReadVarInt(buf *bytes.Buffer) (uint64){

	n := uint64(0)
	for i:= 0 ; true ;i++ {
		fmt.Println("i=",i)
		c := uint64(buf.Bytes()[i])
		fmt.Println("c=", c)
		n = (n << 7) | (c & 0x7F)
		if ((c & 0x80) != uint64(0)) {
			n ++
		} else {
			len := i+1
			fmt.Println("value=",n,"len=",len)
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

	fmt.Println("regid=",strconv.FormatInt(int64(height),10) + "-" + strconv.FormatInt(int64(index),10))

	return strconv.FormatInt(int64(height),10) + "-" + strconv.FormatInt(int64(index),10)
}

//return regid + publicKey
func ReadUserId(buf *bytes.Buffer) (string,string) {
	publicKey := ""
	regid := ""
	idLen := ReadVarInt(buf)
	fmt.Println("idLen=",idLen)
	fmt.Println("after=",buf.Bytes())
	if idLen == 33 {//公钥
		publicKey = ReadPubkey(buf)
	}else { //regid
		regid = ReadRegid(buf)
	}

	return regid,publicKey
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

func GetAddrFrom20BytePubKeyHash( pubKeyHash []byte, netParams *chaincfg.Params) (string,error){

	if len(pubKeyHash) != ripemd160.Size{
		return "",errors.New("The len of pubKeyHash error!")
	}

	return base58.CheckEncode(pubKeyHash[:ripemd160.Size], netParams.PubKeyHashAddrID),nil
}

func ReadUCoinDestAddr(buf *bytes.Buffer,netParams *chaincfg.Params) (*DestArr,error){
	//数组数
	size := ReadVarInt(buf)
	Dests := NewDestArr()

	fmt.Println("size=",size)

	for i:=0; i < int(size) ;i ++{
		keyid_len := ReadVarInt(buf)
		keyid := buf.Bytes()[:keyid_len]
		address,err := GetAddrFrom20BytePubKeyHash(keyid,netParams)
		if err != nil {
			return nil,err
		}
		buf.Next(int(keyid_len))
		coinSymbol := ReadString(buf)
		transferAmount := ReadVarInt(buf)

		fmt.Println("coinSymbol=",coinSymbol)
		fmt.Println("transferAmount=",transferAmount)
		fmt.Println("address=",address)
		//dest:=Dest{string(commons.WICC),1000000, "wLKf2NqwtHk3BfzK5wMDfbKYN1SC3weyR4"}
		Dests.Add(&Dest{coinSymbol,transferAmount,address})
	}

	return Dests,nil
}




