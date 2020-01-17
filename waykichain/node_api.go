package waykichain

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/WaykiChain/wicc-wallet-utils-go/log"
	"github.com/shopspring/decimal"

	"github.com/imroc/req"
	"github.com/tidwall/gjson"
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


type ChainInfo struct {
	Version string `json:"version"`
	ProtocolVersion int `json:"protocol_version"`
	NetType string `json:"net_type"`
	Proxy string `json:"proxy"`
	ExtIp string `json:"ext_ip"`
	ConfDir string `json:"conf_dir"`
	DataDir string `json:"data_dir"`
	BlockInterval int `json:"block_interval"`
	GenBlock int `json:"genblock"`
	TimeOffset int `json:"time_offset"`
	WalletBalance decimal.Decimal `json:"wallet_balance"`
	RelayFeePerk decimal.Decimal `json:"relay_fee_perkb"`
	TipBlockFuelRate int `json:"tipblock_fuel_rate"`
	TipBlockFuel int `json:"tipblock_fuel"`
	TipBlockTime int `json:"tipblock_time"`
	TipBlockHash string `json:"tipblock_hash"`
	TipBlockHeight uint64 `json:"tipblock_height"`
	SynBlockHeight int64 `json:"synblock_height"`
	Connections int `json:"connections"`
	Errors string `json:"errors"`
}

//RPC getinfo
func (c *Client) GetInfo() (*ChainInfo, error) {

	var info ChainInfo
	params := []interface{}{}

	result, err := c.Call("getinfo", params)
	if err != nil {
		log.Errorf("get getinfo failed, err = %v \n", err)
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
		return nil, err
	}

	return &info, nil
}

//获取链上最新高度,可用于签名时作为有效高度的输入
func (c *Client) GetSynBlockHeight() (int64, error) {
	info ,err := c.GetInfo()
	if err != nil{
		return 0,err
	}
	return info.SynBlockHeight,nil
}

//广播交易,返回txid
func (c *Client) SubmitTxRaw(rawtx string) (string, error) {

	params := []interface{}{
		rawtx,
	}
	result, err := c.Call("submittxraw", params)
	if err != nil {
		log.Errorf("submittxraw failed, err = %v \n", err)
		return "", err
	}

	if result.Type != gjson.JSON {
		errInfo := fmt.Sprintf("SubmitTxRaw result type failed, result type is %v",  result.Type)
		log.Errorf(errInfo)
		return "", errors.New(errInfo)
	}

	return result.Str,nil
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
