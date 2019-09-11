package commons

import (
	"bytes"
	"encoding/hex"

	"github.com/btcsuite/btcutil"
	hash2 "github.com/WaykiChain/wicc-wallet-utils-go/commons/hash"
)

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

func (tx WaykiAssetUpdateTx) SignTx(wifKey *btcutil.WIF) string {
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
	writer.WriteByte(byte(tx.UpdateType))
	switch tx.UpdateType {
	case 1:
		writer.WriteUserId(tx.UserId)
		break
	case 2:
		writer.WriteString(tx.AssetName)
		break
	case 3:
		writer.WriteVarInt(int64(tx.AssetTotal))
		break
	}

	signedBytes := tx.doSignTx(wifKey)
	writer.WriteBytes(signedBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return rawTx
}

func (tx WaykiAssetUpdateTx) doSignTx(wifKey *btcutil.WIF) []byte {

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
	writer.WriteByte(byte(tx.UpdateType))
	switch tx.UpdateType {
	case 1:
		writer.WriteUserId(tx.UserId)
		break
	case 2:
		writer.WriteString(tx.AssetName)
		break
	case 3:
		writer.WriteVarInt(int64(tx.AssetTotal))
		break
	}

	hash := hash2.DoubleHash256(buf.Bytes())
	key := wifKey.PrivKey
	ss, _ := key.Sign(hash)
	return ss.Serialize()
}
