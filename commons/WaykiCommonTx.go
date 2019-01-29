package commons

import (
	"bytes"
	"github.com/btcsuite/btcutil/base58"
	"github.com/btcsuite/btcutil"
	"encoding/hex"
)

type WaykiCommonTxParams struct {
	BaseSignTxParams
	Value int64
	DestAddress string
}

func (waykicommon WaykiCommonTxParams)SignTX()string{
	regId:=parseRegId(waykicommon.RegId)
	bytesBuffer := bytes.NewBuffer([]byte{})
	bytesBuffer.WriteByte(byte(waykicommon.TxType))
	bytesBuffer.Write(EncodeInOldWay(waykicommon.Version))
	bytesBuffer.Write(EncodeInOldWay(waykicommon.ValidHeight))
	bytesBuffer.Write(EncodeInOldWay(4))
	bytesBuffer.Write(EncodeInOldWay(regId[0]))
	bytesBuffer.Write(EncodeInOldWay(regId[1]))
	ss7,_,_:= base58.CheckDecode(waykicommon.DestAddress)
	bytesBuffer.Write(EncodeInOldWay(int64(len(ss7))))
	bytesBuffer.Write(ss7)
	bytesBuffer.Write(EncodeInOldWay(waykicommon.Fees))
	bytesBuffer.Write(EncodeInOldWay(waykicommon.Value))
	bytesBuffer.Write(EncodeInOldWay(0))

	ss9:=signCommonTX(waykicommon)
	bytesBuffer.Write(EncodeInOldWay(int64(len(ss9))))
	bytesBuffer.Write(ss9)
	//println(hex.EncodeToString(bytesBuffer.Bytes()))
	signHex:=hex.EncodeToString(bytesBuffer.Bytes())
	return signHex
}

func signCommonTX(waykicommon WaykiCommonTxParams) []byte{
	regId:=parseRegId(waykicommon.RegId)
	bytesBuffer := bytes.NewBuffer([]byte{})
	bytesBuffer.Write(EncodeInOldWay(waykicommon.Version))
	bytesBuffer.WriteByte(byte(waykicommon.TxType))
	bytesBuffer.Write(EncodeInOldWay(waykicommon.ValidHeight))
	bytesBuffer.Write(EncodeInOldWay(4))
	bytesBuffer.Write(EncodeInOldWay(regId[0]))
	bytesBuffer.Write(EncodeInOldWay(regId[1]))
	addressHash,_,_:= base58.CheckDecode(waykicommon.DestAddress)
	bytesBuffer.Write(EncodeInOldWay(int64(len(addressHash))))
	bytesBuffer.Write(addressHash)
	bytesBuffer.Write(EncodeInOldWay(waykicommon.Fees))
	bytesBuffer.Write(EncodeInOldWay(waykicommon.Value))
	bytesBuffer.Write(EncodeInOldWay(0))

	data1,_:=HashDoubleSha256(bytesBuffer.Bytes())
	wif,_ := btcutil.DecodeWIF(waykicommon.PrivateKey)
	key:=wif.PrivKey
	ss,_:=key.Sign(data1)
	return ss.Serialize()
}