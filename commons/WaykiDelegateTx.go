package commons

import (
	"bytes"
	"encoding/hex"
	"github.com/btcsuite/btcutil"
)

type OperVoteFund struct {
	VoteType WaykiVoteType
	PubKey []byte
	VoteValue int64
}

type WaykiDelegateTxParams struct {
	BaseSignTxParams
	OperVoteFunds []OperVoteFund
}

func (waykidelegate WaykiDelegateTxParams)SignTX()string{
	regId:=parseRegId(waykidelegate.RegId)
	bytesBuffer := bytes.NewBuffer([]byte{})
	bytesBuffer.WriteByte(byte(waykidelegate.TxType))
	bytesBuffer.Write(EncodeInOldWay(waykidelegate.Version))
	bytesBuffer.Write(EncodeInOldWay(waykidelegate.ValidHeight))
	bytesBuffer.Write(EncodeInOldWay(4))
	bytesBuffer.Write(EncodeInOldWay(regId[0]))
	bytesBuffer.Write(EncodeInOldWay(regId[1]))
	bytesBuffer.Write(EncodeInOldWay(int64(len(waykidelegate.OperVoteFunds))))
	for _, fund := range waykidelegate.OperVoteFunds {
		bytesBuffer.Write(EncodeInOldWay(int64(fund.VoteType)))
		bytesBuffer.Write(EncodeInOldWay(33))
		bytesBuffer.Write(fund.PubKey)
		bytesBuffer.Write(EncodeInOldWay(int64(fund.VoteValue)))
	}
	bytesBuffer.Write(EncodeInOldWay(waykidelegate.Fees))
	ss9:=signDelegateTX(waykidelegate)
	bytesBuffer.Write(EncodeInOldWay(int64(len(ss9))))
	bytesBuffer.Write(ss9)
	//println(hex.EncodeToString(bytesBuffer.Bytes()))
	signHex:=hex.EncodeToString(bytesBuffer.Bytes())
	return signHex
}

func signDelegateTX(waykidelegate WaykiDelegateTxParams) []byte{
	regId:=parseRegId(waykidelegate.RegId)
	bytesBuffer := bytes.NewBuffer([]byte{})
	bytesBuffer.Write(EncodeInOldWay(waykidelegate.Version))
	bytesBuffer.WriteByte(byte(waykidelegate.TxType))
	bytesBuffer.Write(EncodeInOldWay(waykidelegate.ValidHeight))
	bytesBuffer.Write(EncodeInOldWay(4))
	bytesBuffer.Write(EncodeInOldWay(regId[0]))
	bytesBuffer.Write(EncodeInOldWay(regId[1]))
	bytesBuffer.Write(EncodeInOldWay(int64(len(waykidelegate.OperVoteFunds))))
	for _, fund := range waykidelegate.OperVoteFunds {
		bytesBuffer.Write(EncodeInOldWay(int64(fund.VoteType)))
		bytesBuffer.Write(EncodeInOldWay(33))
		bytesBuffer.Write(fund.PubKey)
		bytesBuffer.Write(EncodeInOldWay(int64(fund.VoteValue)))
	}
	bytesBuffer.Write(EncodeInOldWay(waykidelegate.Fees))

	data1,_:=HashDoubleSha256(bytesBuffer.Bytes())
	wif,_ := btcutil.DecodeWIF(waykidelegate.PrivateKey)
	key:=wif.PrivKey
	ss,_:=key.Sign(data1)
	return ss.Serialize()
}
