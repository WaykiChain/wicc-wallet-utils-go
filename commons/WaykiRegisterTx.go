package commons

import (
	"bytes"
	"encoding/hex"
	"github.com/btcsuite/btcutil"
)

type WaykiRegisterTxParams struct {
	BaseSignTxParams
}

func (waykiRegister WaykiRegisterTxParams)SignTX()string{
	wif,_ := btcutil.DecodeWIF(waykiRegister.PrivateKey)
	bytesBuffer := bytes.NewBuffer([]byte{})
	bytesBuffer.WriteByte(byte(waykiRegister.TxType))
	bytesBuffer.Write(EncodeInOldWay(waykiRegister.Version))
	bytesBuffer.Write(EncodeInOldWay(waykiRegister.ValidHeight))
	bytesBuffer.Write(EncodeInOldWay(33))
	bytesBuffer.Write(wif.PrivKey.PubKey().SerializeCompressed())
	bytesBuffer.Write(EncodeInOldWay(0))
	bytesBuffer.Write(EncodeInOldWay(waykiRegister.Fees))
	ss9:=signRegisterTX(waykiRegister)
	bytesBuffer.Write(EncodeInOldWay(int64(len(ss9))))
	bytesBuffer.Write(ss9)
	//println(hex.EncodeToString(bytesBuffer.Bytes()))
	signHex:=hex.EncodeToString(bytesBuffer.Bytes())
	return signHex
}

func signRegisterTX(waykiRegister WaykiRegisterTxParams) []byte{
	wif,_ := btcutil.DecodeWIF(waykiRegister.PrivateKey)
	bytesBuffer := bytes.NewBuffer([]byte{})
	bytesBuffer.Write(EncodeInOldWay(waykiRegister.Version))
	bytesBuffer.WriteByte(byte(waykiRegister.TxType))
	bytesBuffer.Write(EncodeInOldWay(waykiRegister.ValidHeight))
	bytesBuffer.Write(EncodeInOldWay(33))
	bytesBuffer.Write(wif.PrivKey.PubKey().SerializeCompressed())
	bytesBuffer.Write(EncodeInOldWay(0))
	bytesBuffer.Write(EncodeInOldWay(waykiRegister.Fees))

	data1,_:=HashDoubleSha256(bytesBuffer.Bytes())
	key:=wif.PrivKey
	ss,_:=key.Sign(data1)
	return ss.Serialize()
}
