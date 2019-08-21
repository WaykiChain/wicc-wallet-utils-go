package commons

import (
	"bytes"
	"encoding/hex"

	"github.com/btcsuite/btcutil"
	hash2 "github.com/WaykiChain/wicc-wallet-utils-go/commons/hash"
)

type WaykiRegisterContractTx struct {
	WaykiBaseSignTx
	Script      []byte
	Description string
	Fees        uint64
}

// sign transaction
func (tx WaykiRegisterContractTx) SignTx(wifKey *btcutil.WIF) string {

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
	WriteContractScript(writer, tx.Script, tx.Description)

	writer.WriteVarInt(int64(tx.Fees))
	signedBytes := tx.doSign(wifKey)
	writer.WriteBytes(signedBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return rawTx
}

func (tx WaykiRegisterContractTx) doSign(wifKey *btcutil.WIF) []byte {
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
	WriteContractScript(writer, tx.Script, tx.Description)
	writer.WriteVarInt(int64(tx.Fees))

	hash := hash2.DoubleHash256(buf.Bytes())
	key := wifKey.PrivKey
	ss, _ := key.Sign(hash)
	return ss.Serialize()
}

func WriteContractScript(writer *WriterHelper, script []byte, description string) {

	scriptWriter := NewWriterHelper(bytes.NewBuffer([]byte{}))
	scriptWriter.WriteBytes(script)
	scriptWriter.WriteString(description)
	writer.WriteBytes(scriptWriter.GetBuf().Bytes())
}
