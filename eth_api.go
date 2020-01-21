/*
 * Copyright 2018 The openwallet Authors
 * This file is part of the openwallet library.
 *
 * The openwallet library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The openwallet library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 */
package wicc_wallet_utils_go

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/WaykiChain/wicc-wallet-utils-go/common"
	"github.com/WaykiChain/wicc-wallet-utils-go/log"
	"github.com/imroc/req"
	"github.com/tidwall/gjson"
	//"log"
	"math/big"
	"strconv"
	"strings"
)

type Client struct {
	BaseURL string
	AccessToken string
	Debug   bool
	client      *req.Req
}

type Response struct {
	Id      int         `json:"id"`
	Version string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
}


func NewClient(url, token string, debug bool) *Client {
	c := Client{
		BaseURL:     url,
		AccessToken: token,
		Debug:       debug,
	}
	c.client = req.New()

	return &c
}


/*func NewClient(baseUrl, token string ,debug bool ) *Client{
	return &Client{BaseURL: baseUrl, AccessToken:token ,Debug:debug}
}*/

// See 2 (end of page 4) http://www.ietf.org/rfc/rfc2617.txt
// "To receive authorization, the client sends the userid and password,
// separated by a single colon (":") character, within a base64
// encoded string in the credentials."
// It is not meant to be urlencoded.
func BasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

/*
1. eth block example
   "result": {
        "difficulty": "0x1a4f1f",
        "extraData": "0xd98301080d846765746888676f312e31302e338664617277696e",
        "gasLimit": "0x47e7c4",
        "gasUsed": "0x5b61",
        "hash": "0x85319757555e1cf069684dde286e3c34331dc27d2e54bed24e7291f1b84a0cc5",
        "logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
        "miner": "0x50068fd632c1a6e6c5bd407b4ccf8861a589e776",
        "mixHash": "0xb0cb0abb00c3fc77014abb2a520e3d2a14047cfa30a3b954f18fbeefd1a92f7b",
        "nonce": "0x4df323f58b7a7fd0",
        "number": "0x169cf",
        "parentHash": "0x3df7035473ec98c8c18d2785d5a345193a32b95fcf1ac2d3f09a93109feed3bc",
        "receiptsRoot": "0x441a5be885777bfdf0e985a8ef5046316b3384dd49db7ef95b2c546611c1e2fc",
        "sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
        "size": "0x2aa",
        "stateRoot": "0xb0d76a848be723c72c9639b2de591320f4456b665354995be08a8fa83897efbb",
        "timestamp": "0x5b7babbe",
        "totalDifficulty": "0x2a844e200a",
        "transactions": [
            {
                "blockHash": "0x85319757555e1cf069684dde286e3c34331dc27d2e54bed24e7291f1b84a0cc5",
                "blockNumber": "0x169cf",
                "from": "0x50068fd632c1a6e6c5bd407b4ccf8861a589e776",
                "gas": "0x15f90",
                "gasPrice": "0x430e23400",
                "hash": "0x925e33ac3ebaf40bb44a843860b6589ea2df78c955a27f9df16edcf789519671",
                "input": "0x70a082310000000000000000000000002a63b2203955b84fefe52baca3881b3614991b34",
                "nonce": "0x45",
                "to": "0x8847e5f841458ace82dbb0692c97115799fe28d3",
                "transactionIndex": "0x0",
                "value": "0x0",
                "v": "0x3c",
                "r": "0x8d2ffbe7cb7ac1159a999dfa4352fa27f5cce0df8755254393838aab229ecd33",
                "s": "0xe8ed1f7f8de902ccb008824fe39b2903b94f89e3ea0d5b9f9b880c302bae6cf"
            }
        ],
        "transactionsRoot": "0xa8cb62696679bc3d72762bd2aa5842fdd8aed9c9691fe82064c13e854c13d5cb",
        "uncles": []
    }
*/



type BlockTransaction struct {
	Hash             string `json:"hash" storm:"id"`
	BlockNumber      string `json:"blockNumber" storm:"index"`
	BlockHash        string `json:"blockHash" storm:"index"`
	From             string `json:"from"`
	To               string `json:"to"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	Value            string `json:"value"`
	Data             string `json:"input"`
	TransactionIndex string `json:"transactionIndex"`
	Timestamp        string `json:"timestamp"`
	BlockHeight      uint64 //transaction scanning 的时候对其进行赋值
	Status           uint64
}


func (c *Client) ethGetTransactionCount(addr ,state string) (uint64, error) {
	params := []interface{}{
		AppendOxToAddress(addr),
		state,
	}

	result, err := c.Call("eth_getTransactionCount", params)
	if err != nil {
		//errInfo := fmt.Sprintf("get block[%v] failed, err = %v \n", blockNumStr,  err)
		log.Errorf("get transaction count failed, err = %v \n", err)
		return 0, err
	}

	if result.Type != gjson.String {
		log.Errorf("result type failed. ")
		return 0, errors.New("result type failed. ")
	}

	//blockNum, err := ConvertToBigInt(result.String(), 16)
	nonceStr := result.String()
	nonceStr = strings.ToLower(nonceStr)
	nonceStr = common.RemoveOxFromHex(nonceStr)
	nonce, err := strconv.ParseUint(nonceStr, 16, 64)
	if err != nil {
		log.Errorf("parse nounce failed, err=%v", err)
		return 0, err
	}
	return nonce, nil
}

func (c *Client) EthGetTransactionByHash(txid string) (*BlockTransaction, error) {
	params := []interface{}{
		AppendOxToAddress(txid),
	}

	var tx BlockTransaction

	result, err := c.Call("eth_getTransactionByHash", params)
	if err != nil {
		//errInfo := fmt.Sprintf("get block[%v] failed, err = %v \n", blockNumStr,  err)
		log.Errorf("get transaction[%v] failed, err = %v \n", AppendOxToAddress(txid), err)
		return nil, err
	}

	if result.Type != gjson.JSON {
		errInfo := fmt.Sprintf("get transaction[%v] result type failed, result type is %v", AppendOxToAddress(txid), result.Type)
		log.Errorf(errInfo)
		return nil, errors.New(errInfo)
	}

	err = json.Unmarshal([]byte(result.Raw), &tx)
	if err != nil {
		log.Errorf("decode json [%v] failed, err=%v", result.Raw, err)
		return nil, err
	}

	return &tx, nil
}

func (c *Client) ERC20GetAddressBalance(address string, contractAddr string) (*big.Int, error) {

	contractAddr = "0x" + strings.TrimPrefix(contractAddr, "0x")
	var funcParams []SolidityParam
	funcParams = append(funcParams, SolidityParam{
		ParamType:  SOLIDITY_TYPE_ADDRESS,
		ParamValue: address,
	})
	trans := make(map[string]interface{})
	data, err := makeTransactionData(ETH_GET_TOKEN_BALANCE_METHOD, funcParams)
	if err != nil {
		log.Errorf("make transaction data failed, err = %v", err)
		return nil, err
	}

	trans["to"] = contractAddr
	trans["data"] = data
	params := []interface{}{
		trans,
		"latest",
	}
	result, err := c.Call("eth_call", params)
	if err != nil {
		log.Errorf(fmt.Sprintf("get addr[%v] erc20 balance failed, err=%v\n", address, err))
		return big.NewInt(0), err
	}
	if result.Type != gjson.String {
		errInfo := fmt.Sprintf("get addr[%v] erc20 balance result type error, result type is %v\n", address, result.Type)
		log.Errorf(errInfo)
		return big.NewInt(0), errors.New(errInfo)
	}

	balance, err := ConvertToBigInt(result.String(), 16)
	if err != nil {
		errInfo := fmt.Sprintf("convert addr[%v] erc20 balance format to bigint failed, response is %v, and err = %v\n", address, result.String(), err)
		log.Errorf(errInfo)
		return big.NewInt(0), errors.New(errInfo)
	}
	return balance, nil

}

func (c *Client) GetAddrBalance(address string, blockParameter string) (*big.Int, error) {
	if blockParameter != "latest" && blockParameter != "pending" {
		return nil, errors.New("unknown sign was put through.")
	}

	params := []interface{}{
		AppendOxToAddress(address),
		blockParameter,
	}
	result, err := c.Call("eth_getBalance", params)
	if err != nil {
		//log.Errorf(fmt.Sprintf("get addr[%v] balance failed, err=%v\n", address, err))
		return big.NewInt(0), err
	}
	if result.Type != gjson.String {
		errInfo := fmt.Sprintf("get addr[%v] balance result type error, result type is %v\n", address, result.Type)
		log.Errorf(errInfo)
		return big.NewInt(0), errors.New(errInfo)
	}

	balance, err := ConvertToBigInt(result.String(), 16)
	if err != nil {
		errInfo := fmt.Sprintf("convert addr[%v] balance format to bigint failed, response is %v, and err = %v\n", address, result.String(), err)
		log.Errorf(errInfo)
		return big.NewInt(0), errors.New(errInfo)
	}
	return balance, nil
}

func (c *Client) ethGetGasEstimated(paraMap map[string]interface{}) (*big.Int, error) {
	trans := make(map[string]interface{})
	var temp interface{}
	var exist bool
	var fromAddr string
	var toAddr string

	if temp, exist = paraMap["from"]; !exist {
		log.Errorf("from not found")
		return big.NewInt(0), errors.New("from not found")
	} else {
		fromAddr = temp.(string)
		trans["from"] = fromAddr
	}

	if temp, exist = paraMap["to"]; !exist {
		log.Errorf("to not found")
		return big.NewInt(0), errors.New("to not found")
	} else {
		toAddr = temp.(string)
		trans["to"] = toAddr
	}

	if temp, exist = paraMap["value"]; exist {
		amount := temp.(string)
		trans["value"] = amount
	}

	if temp, exist = paraMap["data"]; exist {
		data := temp.(string)
		trans["data"] = data
	}

	params := []interface{}{
		trans,
	}

	result, err := c.Call("eth_estimateGas", params)
	if err != nil {
		log.Errorf(fmt.Sprintf("get estimated gas limit from [%v] to [%v] faield, err = %v \n", fromAddr, toAddr, err))
		return big.NewInt(0), err
	}

	if result.Type != gjson.String {
		errInfo := fmt.Sprintf("get estimated gas from [%v] to [%v] result type error, result type is %v\n", fromAddr, toAddr, result.Type)
		log.Errorf(errInfo)
		return big.NewInt(0), errors.New(errInfo)
	}

	gasLimit, err := ConvertToBigInt(result.String(), 16)
	if err != nil {
		errInfo := fmt.Sprintf("convert estimated gas[%v] format to bigint failed, err = %v\n", result.String(), err)
		log.Errorf(errInfo)
		return big.NewInt(0), errors.New(errInfo)
	}
	return gasLimit, nil
}

func (c *Client) ethGetGasPrice() (*big.Int, error) {
	params := []interface{}{}
	result, err := c.Call("eth_gasPrice", params)
	if err != nil {
		log.Errorf(fmt.Sprintf("get gas price failed, err = %v \n", err))
		return big.NewInt(0), err
	}

	if result.Type != gjson.String {
		log.Errorf(fmt.Sprintf("get gas price failed, response is %v\n", err))
		return big.NewInt(0), err
	}

	gasLimit, err := ConvertToBigInt(result.String(), 16)
	if err != nil {
		errInfo := fmt.Sprintf("convert estimated gas[%v] format to bigint failed, err = %v\n", result.String(), err)
		log.Errorf(errInfo)
		return big.NewInt(0), errors.New(errInfo)
	}
	return gasLimit, nil
}


func (c *Client) EthSendRawTransaction(signedTx string) (string, error) {

	params := []interface{}{
		AppendOxToAddress(signedTx),
	}

	result, err := c.Call("eth_sendRawTransaction", params)
	if err != nil {
		log.Errorf(fmt.Sprintf("send raw transaction faield, err = %v \n", err))
		return "", err
	}

	if result.Type != gjson.String {
		log.Errorf("eth_sendRawTransaction result type error")
		return "", errors.New("eth_sendRawTransaction result type error")
	}
	return result.String(), nil
}


func (c * Client) GetTransactionFeeEstimated(from string, to string, value *big.Int, data string) (*txFeeInfo, error) {

	var (
		gas *big.Int
		err      error
	)
	/*if this.Config.FixGasLimit.Cmp(big.NewInt(0)) > 0 {
		//配置设置固定gasLimit
		gasLimit = this.Config.FixGasLimit
	} else {
		//动态计算gas消耗
		gasLimit, err = this.WalletClient.ethGetGasEstimated(makeGasEstimatePara(from, to, value, data))
		if err != nil {
			this.Log.Errorf(fmt.Sprintf("get gas limit failed, err = %v\n", err))
			return nil, err
		}
	}*/
	//动态计算gas消耗
	gas, err = c.ethGetGasEstimated(makeGasEstimatePara(from, to, value, data))
	if err != nil {
		log.Errorf(fmt.Sprintf("get gas limit failed, err = %v\n", err))
		return nil, err
	}

	gasPrice, err := c.ethGetGasPrice()
	if err != nil {
		log.Errorf(fmt.Sprintf("get gas price failed, err = %v\n", err))
		return nil, err
	}

	feeInfo := &txFeeInfo{
		GasLimit: gas,
		GasPrice: gasPrice,
	}

	feeInfo.CalcFee()
	return feeInfo, nil
}

func (c *Client) EthGetBlockNumber() (int64, error) {
	param := make([]interface{}, 0)
	result, err := c.Call("eth_blockNumber", param)
	if err != nil {
		log.Errorf("get block number faield, err = %v \n", err)
		return 0, err
	}

	if result.Type != gjson.String {
		log.Errorf("result of block number type error")
		return 0, errors.New("result of block number type error")
	}

	blockNum, err := ConvertToInt64(result.String(), 16)
	if err != nil {
		log.Errorf("parse block number to big.Int failed, err=%v", err)
		return 0, err
	}

	return blockNum, nil
}

//only form Infura
func (c *Client) GetInfuraChainId() (int64, error) {
	param := make([]interface{}, 0)
	result, err := c.Call("eth_chainId", param)
	if err != nil {
		log.Errorf("get chainId faield, err = %v \n", err)
		return -1, err
	}

	if result.Type != gjson.String {
		log.Errorf("result of chainId type error")
		return -1, errors.New("result of chainId type error")
	}

	chainId, err := ConvertToUint64(result.String(), 16)
	if err != nil {
		log.Errorf("parse chainId to big.Int failed, err=%v", err)
		return -1, err
	}

	return int64(chainId), nil
}

func (c *Client) Call(method string, params []interface{}) (*gjson.Result, error) {

	if c.client == nil {
		return nil, errors.New("API url is not setup. ")
	}
	authHeader := req.Header{
		"Accept":       "application/json",
		"Authorization": "Basic " + c.AccessToken,
		"Content-Type": "application/json",
	}
	body := make(map[string]interface{}, 0)
	body["jsonrpc"] = "2.0"
	body["id"] = 1
	body["method"] = method
	body["params"] = params

	if c.Debug {
		log.Debug("Start Request API...")
	}

	r, err := c.client.Post(c.BaseURL, req.BodyJSON(&body), authHeader)

	if c.Debug {
		log.Debug("Request API Completed")
	}

	if c.Debug {
		log.Debugf("%+v\n", r)
	}

	if err != nil {
		return nil, err
	}

	resp := gjson.ParseBytes(r.Bytes())
	err = isError(&resp)
	if err != nil {
		return nil, err
	}

	result := resp.Get("result")

	return &result, nil
}

//isError 是否报错
func isError(result *gjson.Result) error {
	var (
		err error
	)

	if !result.Get("error").IsObject() {

		if !result.Get("result").Exists() {
			return errors.New("Response is empty! ")
		}

		return nil
	}

	errInfo := fmt.Sprintf("[%d]%s",
		result.Get("error.code").Int(),
		result.Get("error.message").String())
	err = errors.New(errInfo)

	return err
}
