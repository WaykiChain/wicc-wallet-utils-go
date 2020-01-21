package wicc_wallet_utils_go

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
)

type txFeeInfo struct {
	GasLimit *big.Int
	GasPrice *big.Int
	Fee      *big.Int
}

type ETHTransactionStr struct{
	FromAddress string
	ToAddress string
	Amount string
	Nonce string
	GasLimit string
	GasPrice string
	Data string
}

type ETHTransaction struct{
	FromAddress string
	ToAddress string
	Amount *big.Int
	Nonce uint64
	GasLimit *big.Int
	GasPrice *big.Int
	Data []byte
}

func (this *txFeeInfo) CalcFee() error {
	fee := new(big.Int)
	fee.Mul(this.GasLimit, this.GasPrice)
	this.Fee = fee
	return nil
}


func NewETHSimpleTransaction(to,amount,nonce ,gas, gasPrice string) (*ETHTransaction,error){

	TxStr := ETHTransactionStr{
		ToAddress:to,
		Amount : amount,
		Nonce : nonce,
		GasLimit : gas,
		GasPrice : gasPrice,
	}

	return convertToETHTranscation(TxStr)
}

func NewERC20TransferTransaction(to, ETHamount, nonce, gas, gasPrice, data string) (*ETHTransaction, error){

	TxStr := ETHTransactionStr{
		ToAddress:to,
		Amount : ETHamount,
		Nonce : nonce,
		GasLimit : gas,
		GasPrice : gasPrice,
		Data: data,
	}

	return convertToETHTranscation(TxStr)
}

func (transaction *ETHTransaction) CreateRawTx(privateKeyStr string, chainId int64) (string, error) {
	privateKey, err := StringToPrivateKey(privateKeyStr)
	if err != nil {
		return "", err
	}

	tx := types.NewTransaction(transaction.Nonce, common.HexToAddress(transaction.ToAddress), transaction.Amount, transaction.GasLimit.Uint64(), transaction.GasPrice, transaction.Data)
	signTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(chainId)), privateKey)
	if err != nil {
		return "", nil
	}

	b, err := rlp.EncodeToBytes(signTx)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}



