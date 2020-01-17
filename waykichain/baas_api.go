package waykichain

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/WaykiChain/wicc-wallet-utils-go/log"
	"github.com/imroc/req"
	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
	"net/http"
)

type BaaSClient struct {
	BaseURL     string
	AccessToken string
	Debug       bool
	Client      *req.Req
}

func NewBaaSClient(url string, debug bool) *BaaSClient {
	c := BaaSClient{
		BaseURL: url,
		//AccessToken: token,
		Debug: debug,
	}

	api := req.New()
	c.Client = api

	return &c
}

type BaaSChainInfo struct {
	Version string `json:"version"`
	ProtocolVersion int `json:"protocolversion"`
	NetType string `json:"nettype"`
	Proxy string `json:"proxy"`
	ConfDir string `json:"confdir"`
	DataDir string `json:"datadirectory"`
	BlockInterval int `json:"blockinterval"`
	GenBlock int `json:"genblock"`
	TimeOffset int `json:"timeoffset"`
	WalletBalance decimal.Decimal `json:"balance"`
	RelayFee decimal.Decimal `json:"relayfee"`
	TipBlockFuelRate int `json:"tipblockfuelrate"`
	TipBlockFuel int `json:"tipblockfuel"`
	TipBlockTime int `json:"tipblocktime"`
	TipBlockHash string `json:"tipblockhash"`
	TipBlockHeight uint64 `json:"blocks"`
	SynBlockHeight int64 `json:"syncheight"`
	Connections int `json:"connections"`
	Errors string `json:"errors"`
}

type SendTxRawBody struct{
	AllFields bool `json:"allFields"`
	RawTx string `json:"rawtx"`
}

//BaaS getinfo
func (c *BaaSClient) GetInfo() (*BaaSChainInfo, error) {

	var info BaaSChainInfo
	path := "api/block/getinfo"
	params := []interface{}{}

	result, err := c.Call(path, params, "POST")
	if err != nil {
		return nil, err
	}

	if result.Type != gjson.JSON {
		errInfo := fmt.Sprintf("getinfo result type failed, result type is %v",  result.Type)
		log.Errorf(errInfo)
		return nil, errors.New(errInfo)
	}

	err = json.Unmarshal([]byte(result.Raw), &info)
	if err != nil {
		log.Errorf("decode json [%v] failed, err=%v", result.Raw, err)
		fmt.Println("1236")
		return nil, err
	}

	return &info, nil
}

//获取链上最新高度,可用于签名时作为有效高度的输入
func (c *BaaSClient) GetBaaSSynBlockHeight() (int64, error) {
	info ,err := c.GetInfo()
	if err != nil{
		return 0,err
	}
	return info.SynBlockHeight,nil
}

//广播交易,返回txid
func (c *BaaSClient) SubmitTxRaw(rawtx string) (string, error) {

	path := "api/transaction/sendrawtx"
	body := make(map[string]interface{}, 0)
	body["allFields"] = false
	body["rawtx"] = rawtx

	result, err := c.Call(path, req.BodyJSON(&body), "POST")
	if err != nil {
		return "", err
	}

	if result.Type != gjson.JSON {
		errInfo := fmt.Sprintf("SubmitTxRaw result type failed, result type is %v",  result.Type)
		log.Errorf(errInfo)
		return "", errors.New(errInfo)
	}

	txid := result.Get("txid").Str
	return txid,nil
}

// Call calls a remote procedure on another node, specified by the path.
func (b *BaaSClient) Call(path string, request interface{}, method string) (*gjson.Result, error) {

	if b.Client == nil {
		return nil, errors.New("API url is not setup. ")
	}

	if b.Debug {
		log.Std.Debug("Start Request API...")
	}

	url := b.BaseURL + path

	r, err := b.Client.Do(method, url, request)

	if b.Debug {
		log.Std.Debug("Request API Completed")
	}

	if b.Debug {
		log.Std.Debug("%+v", r)
	}

	err = b.isError(r)
	if err != nil {
		return nil, err
	}

	resp := gjson.ParseBytes(r.Bytes())

	errCode := resp.Get("code").Num
	if errCode != 0{//报错
		errInfo := fmt.Sprintf("code: %f ; msg: %s", errCode,resp.Get("msg").Str)
		log.Errorf(errInfo)
		return nil, errors.New(errInfo)
	}

	data := resp.Get("data")

	return &data, nil
}

//isError 是否报错
func (b *BaaSClient) isError(resp *req.Resp) error {

	if resp == nil || resp.Response() == nil {
		return errors.New("Response is empty! ")
	}

	if resp.Response().StatusCode != http.StatusOK {
		return fmt.Errorf("%s", resp.String())
	}

	return nil
}
