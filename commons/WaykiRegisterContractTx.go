package commons

import (
	"bytes"
	"encoding/hex"

	"github.com/btcsuite/btcutil"
)

type WaykiRegisterContractTx struct {
	WaykiBaseSignTx
	Script      []byte
	Description string
	Fees        uint64
}

// sign transaction
func (tx WaykiRegisterContractTx) SignTx() string {

	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)
	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.Version)
	writer.WriteVarInt(tx.ValidHeight)
	writer.WriteUserId(tx.UserId)
	WriteContractScript(writer, tx.Script, tx.Description)

	writer.WriteVarInt(int64(tx.Fees))
	signedBytes := tx.doSign()
	writer.WriteBytes(signedBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return rawTx
}

func (tx WaykiRegisterContractTx) doSign() []byte {
	wif, _ := btcutil.DecodeWIF(tx.PrivateKey)
	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)
	writer.WriteVarInt(tx.Version)
	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.ValidHeight)
	writer.WriteUserId(tx.UserId)
	WriteContractScript(writer, tx.Script, tx.Description)
	writer.WriteVarInt(int64(tx.Fees))

	hash, _ := HashDoubleSha256(buf.Bytes())
	key := wif.PrivKey
	ss, _ := key.Sign(hash)
	return ss.Serialize()
}

func WriteContractScript(writer *WriterHelper, script []byte, description string) {

	scriptWriter := NewWriterHelper(bytes.NewBuffer([]byte{}))
	scriptWriter.WriteBytes(script)
	scriptWriter.WriteString(description)
	writer.WriteBytes(scriptWriter.GetBuf().Bytes())
}
