package commons

import (
	"bytes"
	"encoding/hex"

	"github.com/btcsuite/btcutil"
	hash2 "github.com/WaykiChain/wicc-wallet-utils-go/commons/hash"
)

type WaykiCallContractTx struct {
	WaykiBaseSignTx
	AppId    *UserIdWraper //user regid or user key id or app regid
	Fees     uint64
	Values   uint64 //transfer amount
	Contract []byte
}

func (tx WaykiCallContractTx) SignTx(wifKey *btcutil.WIF) string {

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
	writer.WriteUserId(tx.AppId)
	writer.WriteVarInt(int64(tx.Fees))
	writer.WriteVarInt(int64(tx.Values))
	writer.WriteBytes(tx.Contract)

	signedBytes := tx.doSignTx(wifKey)
	writer.WriteBytes(signedBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return rawTx
}

func (tx WaykiCallContractTx) doSignTx(wifKey *btcutil.WIF) []byte {

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
	writer.WriteUserId(tx.AppId)
	writer.WriteVarInt(int64(tx.Fees))
	writer.WriteVarInt(int64(tx.Values))
	writer.WriteBytes(tx.Contract)
	hash := hash2.DoubleHash256(buf.Bytes())

	key := wifKey.PrivKey
	ss, _ := key.Sign(hash)
	return ss.Serialize()
}
