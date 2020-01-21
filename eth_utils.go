package wicc_wallet_utils_go

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/WaykiChain/wicc-wallet-utils-go/common"
	"github.com/WaykiChain/wicc-wallet-utils-go/log"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"strconv"
	"strings"
	tool "github.com/blocktree/openwallet/common"
)

type SolidityParam struct {
	ParamType  string
	ParamValue interface{}
}

func makeRepeatString(c string, count uint) string {
	cs := make([]string, 0)
	for i := 0; i < int(count); i++ {
		cs = append(cs, c)
	}
	return strings.Join(cs, "")
}

func makeTransactionData(methodId string, params []SolidityParam) (string, error) {

	data := methodId
	for i, _ := range params {
		var param string
		if params[i].ParamType == SOLIDITY_TYPE_ADDRESS {
			param = strings.ToLower(params[i].ParamValue.(string))
			if strings.Index(param, "0x") != -1 {
				param = tool.Substr(param, 2, len(param))
			}

			if len(param) != 40 {
				return "", errors.New("length of address error.")
			}
			param = makeRepeatString("0", 24) + param
		} else if params[i].ParamType == SOLIDITY_TYPE_UINT256 {
			intParam := params[i].ParamValue.(*big.Int)
			param = intParam.Text(16)
			l := len(param)
			if l > 64 {
				return "", errors.New("integer overflow.")
			}
			param = makeRepeatString("0", uint(64-l)) + param
			//fmt.Println("makeTransactionData intParam:", intParam.String(), " param:", param)
		} else {
			return "", errors.New("not support solidity type")
		}

		data += param
	}
	return data, nil
}

func AppendOxToAddress(addr string) string {
	if strings.Index(addr, "0x") == -1 {
		return "0x" + addr
	}
	return addr
}


func makeERC20TokenTransferData(contractAddr string, toAddr string, amount *big.Int) (string, error) {
	var funcParams []SolidityParam
	funcParams = append(funcParams, SolidityParam{
		ParamType:  SOLIDITY_TYPE_ADDRESS,
		ParamValue: toAddr,
	})

	funcParams = append(funcParams, SolidityParam{
		ParamType:  SOLIDITY_TYPE_UINT256,
		ParamValue: amount,
	})

	//fmt.Println("make token transfer data, amount:", amount.String())
	data, err := makeTransactionData(ETH_TRANSFER_TOKEN_METHOD, funcParams)
	if err != nil {
		log.Errorf("make transaction data failed, err = %v", err)
		return "", err
	}
	log.Debugf("data:%v", data)
	return data, nil
}

func makeGasEstimatePara(fromAddr string, toAddr string, value *big.Int, data string) map[string]interface{} {
	paraMap := make(map[string]interface{})
	paraMap["from"] = AppendOxToAddress(fromAddr)
	paraMap["to"] = AppendOxToAddress(toAddr)
	if data != "" {
		paraMap["data"] = data
	}

	if value != nil {
		paraMap["value"] = "0x" + value.Text(16)
	}
	return paraMap
}

func StringToPrivateKey(privateKeyStr string) (*ecdsa.PrivateKey, error) {
	privateKeyByte, err := hex.DecodeString(common.RemoveOxFromHex(privateKeyStr))
	if err != nil {
		return nil, err
	}
	privateKey, err := crypto.ToECDSA(privateKeyByte)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func ConvertToUint64(value string, base int) (uint64, error) {
	v := value
	if base == 16 {
		v = common.RemoveOxFromHex(v)
	}

	rst, err := strconv.ParseUint(v, base, 64)
	if err != nil {
		log.Errorf("convert string[%v] to int failed, err = %v", value, err)
		return 0, err
	}
	return rst, nil
}

func ConvertToInt64(value string, base int) (int64, error) {
	v := value
	if base == 16 {
		v = common.RemoveOxFromHex(v)
	}

	rst, err := strconv.ParseInt(v, base, 64)
	if err != nil {
		log.Errorf("convert string[%v] to int failed, err = %v", value, err)
		return 0, err
	}
	return rst, nil
}

func ConvertToBigInt(value string, base int) (*big.Int, error) {
	bigvalue := new(big.Int)
	var success bool
	if base == 16 {
		value = common.RemoveOxFromHex(value)
	}

	if value == "" {
		value = "0"
	}

	_, success = bigvalue.SetString(value, base)
	if !success {
		errInfo := fmt.Sprintf("convert value [%v] to bigint failed, check the value and base passed through\n", value)
		log.Errorf(errInfo)
		return big.NewInt(0), errors.New(errInfo)
	}
	return bigvalue, nil
}

func convertToETHTranscation(tx ETHTransactionStr) (*ETHTransaction,error) {

	amountBig ,err := ConvertToBigInt(tx.Amount,10)
	if err != nil{
		return nil, err
	}

	nonceBig ,err := strconv.ParseUint(tx.Nonce,10,64)
	if err != nil{
		return nil, err
	}

	gasBig,err := ConvertToBigInt(tx.GasLimit,10)
	if err != nil{
		return nil, err
	}

	gasPriceBig,err := ConvertToBigInt(tx.GasPrice,10)
	if err != nil{
		return nil, err
	}

	dataBytes ,err := hex.DecodeString(common.RemoveOxFromHex(tx.Data))
	if err != nil{
		return nil ,errors.New("The ERC20 Transfer data failed!")
	}

	return &ETHTransaction{
		tx.FromAddress,
		tx.ToAddress,
		amountBig,
		nonceBig,
		gasBig,
		gasPriceBig,
		dataBytes,

	},nil
}

const (
	GAS_LIMIT         = 21000
	GAS_PRICE         = 500000000000
)

const (
	ETH_GET_TOKEN_BALANCE_METHOD      = "0x70a08231"
	ETH_TRANSFER_TOKEN_METHOD 		  = "0xa9059cbb"
	ETH_TRANSFER_EVENT_ID             = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
)

const (
	SOLIDITY_TYPE_ADDRESS = "address"
	SOLIDITY_TYPE_UINT256 = "uint256"
	SOLIDITY_TYPE_UINT160 = "uint160"
)

const(
	LEATEST  = "latest"
	EARLIEST = "earliest"
	PENDING  = "pending"
)