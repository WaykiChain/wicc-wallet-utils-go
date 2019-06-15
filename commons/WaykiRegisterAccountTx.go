package commons

import (
	"bytes"
	"encoding/hex"

	"github.com/btcsuite/btcutil"
	hash2 "github.com/WaykiChain/wicc-wallet-utils-go/commons/hash"
)

type WaykiRegisterAccountTx struct {
	WaykiBaseSignTx
	MinerId *UserIdWraper
	Fees    uint64
}

func (tx WaykiRegisterAccountTx) SignTx(wifKey *btcutil.WIF) string {

	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)
	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.Version)
	writer.WriteVarInt(tx.ValidHeight)
	writer.WriteUserId(tx.UserId)
	writer.WriteUserId(tx.MinerId)

	writer.WriteVarInt(int64(tx.Fees))
	signedBytes := tx.doSign(wifKey)
	writer.WriteBytes(signedBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return rawTx
}

func (tx WaykiRegisterAccountTx) doSign(wifKey *btcutil.WIF) []byte {
	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)
	writer.WriteVarInt(tx.Version)
	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.ValidHeight)
	writer.WriteUserId(tx.UserId)
	writer.WriteUserId(tx.MinerId)
	writer.WriteVarInt(int64(tx.Fees))

	hash := hash2.DoubleHash256(buf.Bytes())
	key := wifKey.PrivKey
	ss, _ := key.Sign(hash)
	return ss.Serialize()
}
