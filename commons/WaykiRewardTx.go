package commons

import (
	"bytes"
	"encoding/hex"

	"github.com/btcsuite/btcutil"
)

type WaykiRewardTx struct {
	WaykiBaseSignTx
	Values uint64 // reward values
}

func (tx WaykiRewardTx) SignTx() string {

	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)

	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.Version)
	writer.WriteUserId(tx.UserId)
	writer.WriteVarInt(int64(tx.Values))
	writer.WriteVarInt(tx.ValidHeight)

	signedBytes := tx.doSignTx()
	writer.WriteBytes(signedBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return rawTx
}

func (tx WaykiRewardTx) doSignTx() []byte {

	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)

	writer.WriteVarInt(tx.Version)
	writer.WriteByte(byte(tx.TxType))
	writer.WriteUserId(tx.UserId)
	writer.WriteVarInt(int64(tx.Values))
	writer.WriteVarInt(tx.ValidHeight)

	hash, _ := HashDoubleSha256(buf.Bytes())
	wif, _ := btcutil.DecodeWIF(tx.PrivateKey)
	key := wif.PrivKey
	ss, _ := key.Sign(hash)
	return ss.Serialize()
}
