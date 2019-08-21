package commons

import (
	"bytes"
	"encoding/hex"

	"github.com/btcsuite/btcutil"
	hash2 "github.com/WaykiChain/wicc-wallet-utils-go/commons/hash"
)

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

func (tx WaykiDexSellLimitTx) SignTx(wifKey *btcutil.WIF) string {
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
	signedBytes := tx.doSignTx(wifKey)
	writer.WriteBytes(signedBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return rawTx
}

func (tx WaykiDexSellLimitTx) doSignTx(wifKey *btcutil.WIF) []byte {

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
	writer.WriteString(tx.CoinSymbol)
	writer.WriteString(tx.AssetSymbol)
	writer.WriteVarInt(int64(tx.AssetAmount))
	writer.WriteVarInt(int64(tx.AskPrice))

	hash := hash2.DoubleHash256(buf.Bytes())
	key := wifKey.PrivKey
	ss, _ := key.Sign(hash)
	return ss.Serialize()
}
