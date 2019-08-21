package commons

import (
	"bytes"
	"encoding/hex"

	"github.com/btcsuite/btcutil"
	hash2 "github.com/WaykiChain/wicc-wallet-utils-go/commons/hash"
)

type WaykiCommonTx struct {
	WaykiBaseSignTx
	Fees   uint64
	Values uint64
	DestId *UserIdWraper //< the dest id(reg id or address or public key) received the wicc values
}

func (tx WaykiCommonTx) SignTx(wifKey *btcutil.WIF) string {
	//uid := ParseRegId(tx.UserId)
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
	writer.WriteUserId(tx.DestId)
	writer.WriteVarInt(int64(tx.Fees))
	writer.WriteVarInt(int64(tx.Values))
	writer.WriteVarInt(0) // write the empty contract script data

	signedBytes := tx.doSignTx(wifKey)
	writer.WriteBytes(signedBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return rawTx
}

func (tx WaykiCommonTx) doSignTx(wifKey *btcutil.WIF) []byte {

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
	writer.WriteUserId(tx.DestId)
	writer.WriteVarInt(int64(tx.Fees))
	writer.WriteVarInt(int64(tx.Values))
	writer.WriteVarInt(0) // write the empty contract script data

	hash := hash2.DoubleHash256(buf.Bytes())
	key := wifKey.PrivKey
	ss, _ := key.Sign(hash)
	return ss.Serialize()
}
