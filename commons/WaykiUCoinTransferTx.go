package commons

import (
	"bytes"
	"encoding/hex"

	"github.com/btcsuite/btcutil"
	hash2 "github.com/WaykiChain/wicc-wallet-utils-go/commons/hash"
)

type Dest struct {
	CoinSymbol string   //From Coin Type
	CoinAmount uint64
	DestAddr    *UserIdWraper
}

type WaykiUCoinTransferTx struct {
	WaykiBaseSignTx
	Fees   uint64
	FeeSymbol string      //Fee Type (WICC/WUSD)
	Dests    []Dest
	Memo       string
}

func (tx WaykiUCoinTransferTx) SignTx(wifKey *btcutil.WIF) string {
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
	signedBytes := tx.doSignTx(wifKey)
	writer.WriteBytes(signedBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return rawTx
}

func (tx WaykiUCoinTransferTx) doSignTx(wifKey *btcutil.WIF) []byte {

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

	hash := hash2.DoubleHash256(buf.Bytes())
	key := wifKey.PrivKey
	ss, _ := key.Sign(hash)
	return ss.Serialize()
}
