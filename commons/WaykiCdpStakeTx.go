package commons

import (
"bytes"
"encoding/hex"

"github.com/btcsuite/btcutil"
hash2 "github.com/WaykiChain/wicc-wallet-utils-go/commons/hash"
)

type AssetModel struct {
	AssetAmount    int64 //
	AssetSymbol string  //
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
    writer.WriteCdpAsset(tx.Assets)
	ss:=[]byte(tx.ScoinSymbol)
	writer.WriteVarInt(int64(len(ss)))
	writer.Write(ss)
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

	hash := hash2.DoubleHash256(buf.Bytes())
	key := wifKey.PrivKey
	ss, _ := key.Sign(hash)
	return ss.Serialize()
}
