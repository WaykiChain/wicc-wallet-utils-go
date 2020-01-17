package ethereum

import (
	"encoding/hex"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	wicc_common "github.com/WaykiChain/wicc-wallet-utils-go/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
)

type txFeeInfo struct {
	GasLimit *big.Int
	GasPrice *big.Int
	Fee      *big.Int
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


func NewETHSimpleTransaction(to string,amount *big.Int,nonce uint64,gas, gasPrice *big.Int) (*ETHTransaction){

	return &ETHTransaction{
		ToAddress: to,
		Amount : amount,
		Nonce : nonce,
		GasLimit : gas,
		GasPrice : gasPrice,
		Data :nil,
	}
}

func NewERC20TransferTransaction(to string,ETHamount *big.Int,nonce uint64,gas, gasPrice *big.Int, data string) (*ETHTransaction, error){

	dataBytes ,err := hex.DecodeString(wicc_common.RemoveOxFromHex(data))
	if err != nil{
		return nil ,errors.New("The ERC20 Transfer data failed!")
	}
	return &ETHTransaction{
		ToAddress: to,
		Amount : ETHamount,
		Nonce : nonce,
		GasLimit : gas,
		GasPrice : gasPrice,
		Data : dataBytes,
	},nil
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



