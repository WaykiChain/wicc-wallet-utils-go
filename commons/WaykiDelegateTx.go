package commons

import (
	"bytes"
	"encoding/hex"

	"github.com/btcsuite/btcutil"
)

type OperVoteFund struct {
	VoteType  WaykiVoteType
	PubKey    []byte
	VoteValue int64
}

type WaykiDelegateTxParams struct {
	BaseSignTxParams
	OperVoteFunds []OperVoteFund
	Fees          uint64
}

func (params WaykiDelegateTxParams) SignTX() string {
	srcId := ParseRegId(params.UserId)

	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)

	writer.WriteByte(byte(params.TxType))
	writer.WriteVarInt(params.Version)
	writer.WriteVarInt(params.ValidHeight)
	writer.WriteRegId(*srcId)
	writer.WriteVarInt(int64(len(params.OperVoteFunds)))
	for _, fund := range params.OperVoteFunds {
		writer.WriteByte(byte(fund.VoteType))
		writer.WriteBytes(fund.PubKey)
		writer.WriteVarInt(fund.VoteValue)
	}
	writer.WriteVarInt(int64(params.Fees))
	signedBytes := params.doSignTX()
	writer.WriteBytes(signedBytes)

	rawTx := hex.EncodeToString(buf.Bytes())
	return rawTx
}

func (params WaykiDelegateTxParams) doSignTX() []byte {
	uid := ParseRegId(params.UserId)

	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)

	writer.WriteVarInt(params.Version)
	writer.WriteByte(byte(params.TxType))
	writer.WriteVarInt(params.ValidHeight)
	writer.WriteRegId(*uid)
	writer.WriteVarInt(int64(len(params.OperVoteFunds)))
	for _, fund := range params.OperVoteFunds {
		writer.WriteByte(byte(fund.VoteType))
		writer.WriteBytes(fund.PubKey)
		writer.WriteVarInt(fund.VoteValue)
	}
	writer.WriteVarInt(int64(params.Fees))
	hash, _ := HashDoubleSha256(buf.Bytes())
	wif, _ := btcutil.DecodeWIF(params.PrivateKey)
	key := wif.PrivKey
	ss, _ := key.Sign(hash)
	return ss.Serialize()
}
