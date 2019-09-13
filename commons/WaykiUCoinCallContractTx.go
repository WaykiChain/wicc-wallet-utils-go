package commons

import (
	"bytes"
	"encoding/hex"

	"github.com/btcsuite/btcutil"
	hash2 "github.com/WaykiChain/wicc-wallet-utils-go/commons/hash"
)

type WaykiUCoinCallContractTx struct {
	WaykiBaseSignTx
	AppId    *UserIdWraper //user regid or user key id or app regid
	Fees     int64
	CoinAmount   int64 //transfer amount
	FeeSymbol string      //Fee Type (WICC/WUSD)
	CoinSymbol string   //From Coin Type
	Contract []byte
}

func (tx WaykiUCoinCallContractTx) SignTx(wifKey *btcutil.WIF) string {

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


	signedBytes := tx.doSignTx(wifKey)
	writer.WriteBytes(signedBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return rawTx
}

func (tx WaykiUCoinCallContractTx) doSignTx(wifKey *btcutil.WIF) []byte {

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
	hash := hash2.DoubleHash256(buf.Bytes())

	key := wifKey.PrivKey
	ss, _ := key.Sign(hash)
	return ss.Serialize()
}
