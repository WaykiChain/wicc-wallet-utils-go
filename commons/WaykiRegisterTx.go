package commons

import (
	"bytes"
	"encoding/hex"

	"github.com/btcsuite/btcutil"
)

type WaykiRegisterTxParams struct {
	BaseSignTxParams
	Fees uint64
}

func (waykiRegister WaykiRegisterTxParams) SignTX() string {
	wif, _ := btcutil.DecodeWIF(waykiRegister.PrivateKey)
	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)
	writer.WriteByte(byte(waykiRegister.TxType))
	writer.WriteVarInt(waykiRegister.Version)
	writer.WriteVarInt(waykiRegister.ValidHeight)
	writer.WriteBytes(wif.PrivKey.PubKey().SerializeCompressed()) // 33
	writer.WriteBytes([]byte{})                                   // minerid

	writer.WriteVarInt(int64(waykiRegister.Fees))
	signedBytes := waykiRegister.doSign()
	writer.WriteBytes(signedBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return rawTx
}

func (waykiRegister WaykiRegisterTxParams) doSign() []byte {
	wif, _ := btcutil.DecodeWIF(waykiRegister.PrivateKey)
	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)
	writer.WriteVarInt(waykiRegister.Version)
	writer.WriteByte(byte(waykiRegister.TxType))
	writer.WriteVarInt(waykiRegister.ValidHeight)
	writer.WriteBytes(wif.PrivKey.PubKey().SerializeCompressed())
	writer.WriteBytes([]byte{}) // miner id
	writer.WriteVarInt(int64(waykiRegister.Fees))

	hash, _ := HashDoubleSha256(buf.Bytes())
	key := wif.PrivKey
	ss, _ := key.Sign(hash)
	return ss.Serialize()
}
