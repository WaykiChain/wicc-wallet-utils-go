package commons

import (
	"bytes"
	"encoding/hex"

	"github.com/btcsuite/btcutil"
	hash2 "github.com/WaykiChain/wicc-wallet-utils-go/commons/hash"
)

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

func (tx WaykiAssetIssueTx) SignTx(wifKey *btcutil.WIF) string {
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
	writer.WriteString(tx.AssetSymbol)
	writer.WriteUserId(tx.UserId)
	writer.WriteString(tx.AssetName)
	if(tx.MinTable){
		writer.WriteByte(1)
	}else {
		writer.WriteByte(0)
	}
	writer.WriteVarInt(int64(tx.AssetTotal))

	signedBytes := tx.doSignTx(wifKey)
	writer.WriteBytes(signedBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return rawTx
}

func (tx WaykiAssetIssueTx) doSignTx(wifKey *btcutil.WIF) []byte {

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
	writer.WriteString(tx.AssetSymbol)
	writer.WriteUserId(tx.UserId)
	writer.WriteString(tx.AssetName)
	if(tx.MinTable){
		writer.WriteByte(1)
	}else {
		writer.WriteByte(0)
	}
	writer.WriteVarInt(int64(tx.AssetTotal))

	hash := hash2.DoubleHash256(buf.Bytes())
	key := wifKey.PrivKey
	ss, _ := key.Sign(hash)
	return ss.Serialize()
}
