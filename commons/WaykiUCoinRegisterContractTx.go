package commons

import (
	"bytes"
	"encoding/hex"

	"github.com/btcsuite/btcutil"
	hash2 "github.com/WaykiChain/wicc-wallet-utils-go/commons/hash"
)

type WaykiUCoinRegisterContractTx struct {
	WaykiBaseSignTx
	Script      []byte
	Description string
	Fees        uint64
	FeeSymbol   string
}

// sign transaction
func (tx WaykiUCoinRegisterContractTx) SignTx(wifKey *btcutil.WIF) string {

	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)
	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.Version)
	writer.WriteVarInt(tx.ValidHeight)
	writer.WriteUserId(tx.UserId)
	writer.WriteVarInt(int64(tx.Fees))
	writer.WriteString(tx.FeeSymbol)
	WriteContractScript(writer, tx.Script, tx.Description)
	signedBytes := tx.doSign(wifKey)
	writer.WriteBytes(signedBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return rawTx
}

func (tx WaykiUCoinRegisterContractTx) doSign(wifKey *btcutil.WIF) []byte {
	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)
	writer.WriteVarInt(tx.Version)
	writer.WriteByte(byte(tx.TxType))
	writer.WriteVarInt(tx.ValidHeight)
	writer.WriteUserId(tx.UserId)
	writer.WriteVarInt(int64(tx.Fees))
	writer.WriteString(tx.FeeSymbol)
	WriteContractScript(writer, tx.Script, tx.Description)


	hash := hash2.DoubleHash256(buf.Bytes())
	key := wifKey.PrivKey
	ss, _ := key.Sign(hash)
	return ss.Serialize()
}

