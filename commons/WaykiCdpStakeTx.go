package commons

import (
"bytes"
"encoding/hex"

"github.com/btcsuite/btcutil"
hash2 "github.com/WaykiChain/wicc-wallet-utils-go/commons/hash"
)

type WaykiCdpStakeTx struct {
	WaykiBaseSignTx
	Fees   uint64
	ScoinValues uint64   //Stake Coin
	BcoinValues uint64   // get Coin
	FeeSymbol string      //Fee Type (WICC/WUSD)
    ScoinSymbol string   //From Coin Type
    BcoinSymbol string  //Get Coin Type
    CdpTxHash []byte
}

func (tx WaykiCdpStakeTx) SignTx(wifKey *btcutil.WIF) string {
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
	bs:=[]byte(tx.BcoinSymbol)
	writer.WriteVarInt(int64(len(bs)))
	writer.Write(bs)
	ss:=[]byte(tx.ScoinSymbol)
	writer.WriteVarInt(int64(len(ss)))
	writer.Write(ss)
	writer.WriteVarInt(int64(tx.BcoinValues))
	writer.WriteVarInt(int64(tx.ScoinValues))
	signedBytes := tx.doSignTx(wifKey)
	writer.WriteBytes(signedBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return rawTx
}

func (tx WaykiCdpStakeTx) doSignTx(wifKey *btcutil.WIF) []byte {

	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)

	writer.WriteVarInt(tx.Version)
	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.ValidHeight)
	if(tx.UserId!=nil){
		writer.WriteUserId(tx.UserId)
	}else if(tx.PubKey!=nil){
		writer.WriteReverse(tx.PubKey)
	}
	writer.WriteString(tx.FeeSymbol)
	writer.WriteVarInt(int64(tx.Fees))
	writer.WriteReverse(tx.CdpTxHash)
	bscoin:=[]byte(tx.BcoinSymbol)
	writer.WriteVarInt(int64(len(bscoin)))
	writer.Write(bscoin)
	ssoin:=[]byte(tx.ScoinSymbol)
	writer.WriteVarInt(int64(len(ssoin)))
	writer.Write(ssoin)
	writer.WriteVarInt(int64(tx.BcoinValues))
	writer.WriteVarInt(int64(tx.ScoinValues))

	hash := hash2.DoubleHash256(buf.Bytes())
	key := wifKey.PrivKey
	ss, _ := key.Sign(hash)
	return ss.Serialize()
}
