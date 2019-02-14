package commons

import (
	"bytes"
	"encoding/hex"

	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
)

type WaykiCommonTxParams struct {
	BaseSignTxParams
	Fees        uint64
	Values      uint64
	DestAddress string
}

func (params WaykiCommonTxParams) SignTX() string {
	uid := ParseRegId(params.UserId)
	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)

	writer.WriteByte(byte(params.TxType))
	writer.WriteVarInt(params.Version)
	writer.WriteVarInt(params.ValidHeight)

	// WriteRegId
	writer.WriteRegId(*uid)

	destAddress, _, _ := base58.CheckDecode(params.DestAddress)
	writer.WriteBytes(destAddress)

	writer.WriteVarInt(int64(params.Fees))
	writer.WriteVarInt(int64(params.Values))
	writer.WriteVarInt(0) // write the empty contract script data

	signedBytes := params.doSignTx()
	writer.WriteBytes(signedBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return rawTx
}

func (params WaykiCommonTxParams) doSignTx() []byte {
	uid := ParseRegId(params.UserId)

	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)

	writer.WriteVarInt(params.Version)
	writer.WriteByte(byte(params.TxType))
	writer.WriteVarInt(params.ValidHeight)
	writer.WriteRegId(*uid)
	destAddress, _, _ := base58.CheckDecode(params.DestAddress)
	writer.WriteBytes(destAddress)
	writer.WriteVarInt(int64(params.Fees))
	writer.WriteVarInt(int64(params.Values))
	writer.WriteVarInt(0) // write the empty contract script data

	hash, _ := HashDoubleSha256(buf.Bytes())
	wif, _ := btcutil.DecodeWIF(params.PrivateKey)
	key := wif.PrivKey
	ss, _ := key.Sign(hash)
	return ss.Serialize()
}
