package commons

import (
	"bytes"
	"encoding/hex"
	"github.com/btcsuite/btcutil"
)

type WaykiContractTxParams struct {
	BaseSignTxParams
	Value int64
	Appid string
	ContractBytes []byte
}

func (waykiContract WaykiContractTxParams)SignTX()string{
	srcId:=parseRegId(waykiContract.RegId)
	destId:=parseRegId(waykiContract.Appid)
	bytesBuffer := bytes.NewBuffer([]byte{})
	bytesBuffer.WriteByte(byte(waykiContract.TxType))
	bytesBuffer.Write(EncodeInOldWay(waykiContract.Version))
	bytesBuffer.Write(EncodeInOldWay(waykiContract.ValidHeight))
	bytesBuffer.Write(EncodeInOldWay(4))
	bytesBuffer.Write(EncodeInOldWay(srcId[0]))
	bytesBuffer.Write(EncodeInOldWay(srcId[1]))
	bytesBuffer.Write(EncodeInOldWay(4))
	bytesBuffer.Write(EncodeInOldWay(destId[0]))
	bytesBuffer.Write(EncodeInOldWay(destId[1]))
	bytesBuffer.Write(EncodeInOldWay(waykiContract.Fees))
	bytesBuffer.Write(EncodeInOldWay(waykiContract.Value))
	bytesBuffer.Write(EncodeInOldWay(int64(len(waykiContract.ContractBytes))))
	bytesBuffer.Write(waykiContract.ContractBytes)
	ss9:=signContractTX(waykiContract)
	bytesBuffer.Write(EncodeInOldWay(int64(len(ss9))))
	bytesBuffer.Write(ss9)
	signHex:=hex.EncodeToString(bytesBuffer.Bytes())
	return signHex
}

func signContractTX(waykiContract WaykiContractTxParams) []byte{
	srcId:=parseRegId(waykiContract.RegId)
	destId:=parseRegId(waykiContract.Appid)
	bytesBuffer := bytes.NewBuffer([]byte{})
	bytesBuffer.Write(EncodeInOldWay(waykiContract.Version))
	bytesBuffer.WriteByte(byte(waykiContract.TxType))
	bytesBuffer.Write(EncodeInOldWay(waykiContract.ValidHeight))
	bytesBuffer.Write(EncodeInOldWay(4))
	bytesBuffer.Write(EncodeInOldWay(srcId[0]))
	bytesBuffer.Write(EncodeInOldWay(srcId[1]))
	bytesBuffer.Write(EncodeInOldWay(4))
	bytesBuffer.Write(EncodeInOldWay(destId[0]))
	bytesBuffer.Write(EncodeInOldWay(destId[1]))
	bytesBuffer.Write(EncodeInOldWay(waykiContract.Fees))
	bytesBuffer.Write(EncodeInOldWay(waykiContract.Value))
	bytesBuffer.Write(EncodeInOldWay(int64(len(waykiContract.ContractBytes))))
	bytesBuffer.Write(waykiContract.ContractBytes)
	data1,_:=HashDoubleSha256(bytesBuffer.Bytes())
	wif,_ := btcutil.DecodeWIF(waykiContract.PrivateKey)
	key:=wif.PrivKey
	ss,_:=key.Sign(data1)
	return ss.Serialize()
}
