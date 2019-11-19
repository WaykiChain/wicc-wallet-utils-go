package commons

import (
	"bytes"
	"encoding/hex"

	"github.com/btcsuite/btcutil"
	hash2 "github.com/WaykiChain/wicc-wallet-utils-go/commons/hash"
)

type WaykiCdpRedeemTx struct {
	WaykiBaseSignTx
	Fees   uint64
	ScoinValues uint64   //Stake Coin
	Assets   []AssetModel
	FeeSymbol string      //Fee Type (WICC/WUSD)
	CdpTxHash []byte
}

func (tx WaykiCdpRedeemTx) SignTx(wifKey *btcutil.WIF) string {
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
	signedBytes := tx.doSignTx(wifKey)
	writer.WriteBytes(signedBytes)
	rawTx := hex.EncodeToString(buf.Bytes())
	return rawTx
}

func (tx WaykiCdpRedeemTx) doSignTx(wifKey *btcutil.WIF) []byte {

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

	hash := hash2.DoubleHash256(buf.Bytes())
	key := wifKey.PrivKey
	ss, _ := key.Sign(hash)
	return ss.Serialize()
}

