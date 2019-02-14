package commons

import (
	"bytes"
	"encoding/hex"

	"github.com/btcsuite/btcutil"
)

type WaykiContractTxParams struct {
	BaseSignTxParams
	AppId    string //user regid or user key id or app regid
	Fees     uint64
	Values   uint64 //transfer amount
	Contract []byte
}

func (params WaykiContractTxParams) SignTX() string {

	uid := ParseRegId(params.UserId)
	appId := ParseRegId(params.AppId)

	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)

	writer.WriteByte(byte(params.TxType))
	writer.WriteVarInt(params.Version)
	writer.WriteVarInt(params.ValidHeight)
	writer.WriteRegId(*uid)
	writer.WriteRegId(*appId)
	writer.WriteVarInt(int64(params.Fees))
	writer.WriteVarInt(int64(params.Values))
	writer.WriteBytes(params.Contract)

	signedBytes := params.doSignTx()
	writer.WriteBytes(signedBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return rawTx
}

func (params WaykiContractTxParams) doSignTx() []byte {

	uid := ParseRegId(params.UserId)
	appId := ParseRegId(params.AppId)

	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)

	writer.WriteVarInt(params.Version)
	writer.WriteByte(byte(params.TxType))
	writer.WriteVarInt(params.ValidHeight)
	writer.WriteRegId(*uid)
	writer.WriteRegId(*appId)
	writer.WriteVarInt(int64(params.Fees))
	writer.WriteVarInt(int64(params.Values))
	writer.WriteBytes(params.Contract)

	hash, _ := HashDoubleSha256(buf.Bytes())
	wif, _ := btcutil.DecodeWIF(params.PrivateKey)
	key := wif.PrivKey
	ss, _ := key.Sign(hash)
	return ss.Serialize()
}
