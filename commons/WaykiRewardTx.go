package commons

import (
	"bytes"
	"encoding/hex"

	"github.com/btcsuite/btcutil"
	hash2 "github.com/WaykiChain/wicc-wallet-utils-go/commons/hash"
)

type WaykiRewardTx struct {
	WaykiBaseSignTx
	Values uint64 // reward values
}

func (tx WaykiRewardTx) SignTx(wifKey *btcutil.WIF) string {

	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)

	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.Version)
	writer.WriteUserId(tx.UserId)
	writer.WriteVarInt(int64(tx.Values))
	writer.WriteVarInt(tx.ValidHeight)

	signedBytes := tx.doSignTx(wifKey)
	writer.WriteBytes(signedBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return rawTx
}

func (tx WaykiRewardTx) doSignTx(wifKey *btcutil.WIF) []byte {

	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)

	writer.WriteVarInt(tx.Version)
	writer.WriteByte(byte(tx.TxType))
	writer.WriteUserId(tx.UserId)
	writer.WriteVarInt(int64(tx.Values))
	writer.WriteVarInt(tx.ValidHeight)

	hash := hash2.DoubleHash256(buf.Bytes())
	key := wifKey.PrivKey
	ss, _ := key.Sign(hash)
	return ss.Serialize()
}
