
package bitcoin

import (
	"errors"
	"fmt"
	"github.com/WaykiChain/wicc-wallet-utils-go/log"
	"github.com/imroc/req"
	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
	"net/http"
)

// Bitcorer是由bitpay提供的区块数据查询接口
// 具体接口说明查看https://github.com/bitpay/bitcore/blob/master/packages/bitcore-node/docs/api-documentation.md
type Bitcorer struct {
	BaseURL     string
	AccessToken string
	Debug       bool
	Client      *req.Req
}

func NewExplorer(url string, debug bool) *Bitcorer {
	c := Bitcorer{
		BaseURL: url,
		//AccessToken: token,
		Debug: debug,
	}

	api := req.New()
	c.Client = api

	return &c
}

// Call calls a remote procedure on another node, specified by the path.
func (b *Bitcorer) Call(path string, request interface{}, method string) (*gjson.Result, error) {

	if b.Client == nil {
		return nil, errors.New("API url is not setup. ")
	}

	if b.Debug {
		log.Std.Debug("Start Request API...")
	}

	url := b.BaseURL + path

	fmt.Println("url=",url)

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

	return &resp, nil
}

//isError 是否报错
func (b *Bitcorer) isError(resp *req.Resp) error {

	if resp == nil || resp.Response() == nil {
		return errors.New("Response is empty! ")
	}

	if resp.Response().StatusCode != http.StatusOK {
		return fmt.Errorf("%s", resp.String())
	}

	return nil
}

/*//getBlockByExplorer 获取区块数据
func (wm *WalletManager) getBlockByExplorer(hash string) (*Block, error) {

	path := fmt.Sprintf("block/%s", hash)

	result, err := wm.BitcoreClient.Call(path, nil, "GET")
	if err != nil {
		return nil, err
	}

	return newBlockByExplorer(result), nil
}*/

/*//getBlockHashByExplorer 获取区块hash
func (wm *WalletManager) getBlockHashByExplorer(height uint64) (string, error) {

	path := fmt.Sprintf("block-index/%d", height)

	result, err := wm.BitcoreClient.Call(path, nil, "GET")
	if err != nil {
		return "", err
	}

	return result.Get("blockHash").String(), nil
}*/

//getBlockHeightByBitcore 获取区块链高度
func (b *Bitcorer) getBlockHeightByBitcore() (uint64, error) {

	path := "block/tip"

	result, err := b.Call(path, nil, "GET")
	if err != nil {
		return 0, err
	}

	height := result.Get("height").Uint()

	return height, nil
}

//getTxIDsInMemPoolByBitcore 获取待处理的交易池中的交易单IDs
func (b *Bitcorer) getTxIDsInMemPoolByBitcore() ([]string, error) {

	return nil, nil
}

/*//getTransactionByBitcire 获取交易单
func (wm *WalletManager) getTransactionByBitcire(txid string) (*Transaction, error) {

	path := fmt.Sprintf("tx/%s", txid)

	result, err := wm.BitcoreClient.Call(path, nil, "GET")
	if err != nil {
		return nil, err
	}

	tx := wm.newTxByExplorer(result)

	return tx, nil

}*/


/*
func newBlockByExplorer(json *gjson.Result) *Block {

	/*
		{
			"hash": "0000000000002bd2475d1baea1de4067ebb528523a8046d5f9d8ef1cb60460d3",
			"size": 549,
			"height": 1434016,
			"version": 536870912,
			"merkleroot": "ae4310c991ec16cfc7404aaad9fe5fbd533d0b6617c03eb1ac644c89d58b3e18",
			"tx": ["6767a8acc1a63c7978186c582fdea26c47da5e04b0b2b34740a1728bfd959a05", "226dee96373aedd8a3dd00021684b190b7f23f5e16bb186cee11d0560406c19d"],
			"time": 1539066282,
			"nonce": 4089837546,
			"bits": "1a3fffc0",
			"difficulty": 262144,
			"chainwork": "0000000000000000000000000000000000000000000000c6fce84fddeb57e5fb",
			"confirmations": 279,
			"previousblockhash": "0000000000001fdabb5efc93d15ccaf6980642918cd898df6b3ff5fbf26c19c4",
			"nextblockhash": "00000000000024f2bd323157e595613291f83485ddfbbf311323ed0c0dc46545",
			"reward": 0.78125,
			"isMainChain": true,
			"poolInfo": {}
		}

	obj := &Block{}
	//解析json
	obj.Hash = gjson.Get(json.Raw, "hash").String()
	obj.Confirmations = gjson.Get(json.Raw, "confirmations").Uint()
	obj.Merkleroot = gjson.Get(json.Raw, "merkleroot").String()

	txs := make([]string, 0)
	for _, tx := range gjson.Get(json.Raw, "tx").Array() {
		txs = append(txs, tx.String())
	}

	obj.tx = txs
	obj.Previousblockhash = gjson.Get(json.Raw, "previousblockhash").String()
	obj.Height = gjson.Get(json.Raw, "height").Uint()
	//obj.Version = gjson.Get(json.Raw, "version").String()
	obj.Time = gjson.Get(json.Raw, "time").Uint()

	return obj
}
*/



/*//getBalanceByExplorer 获取地址余额
func (wm *WalletManager) getBalanceByExplorer(address string) (*openwallet.Balance, error) {

	path := fmt.Sprintf("addr/%s?noTxList=1", address)

	result, err := wm.BitcoreClient.Call(path, nil, "GET")
	if err != nil {
		return nil, err
	}

	return newBalanceByExplorer(result), nil
}*/
/*
func newBalanceByExplorer(json *gjson.Result) *openwallet.Balance {

	/*

		{
			"addrStr": "mnMSQs3HZ5zhJrCEKbqGvcDLjAAxvDJDCd",
			"balance": 3136.82244887,
			"balanceSat": 313682244887,
			"totalReceived": 3136.82244887,
			"totalReceivedSat": 313682244887,
			"totalSent": 0,
			"totalSentSat": 0,
			"unconfirmedBalance": 0,
			"unconfirmedBalanceSat": 0,
			"unconfirmedTxApperances": 0,
			"txApperances": 3909
		}

	*//*
	//log.Debug(json.Raw)
	obj := openwallet.Balance{}
	//解析json
	obj.Address = gjson.Get(json.Raw, "addrStr").String()
	obj.ConfirmBalance = gjson.Get(json.Raw, "balance").String()
	obj.UnconfirmBalance = gjson.Get(json.Raw, "unconfirmedBalance").String()
	u, _ := decimal.NewFromString(obj.ConfirmBalance)
	b, _ := decimal.NewFromString(obj.UnconfirmBalance)
	obj.Balance = u.Add(b).String()

	return &obj
}
*/
/*//getMultiAddrTransactionsByExplorer 获取多个地址的交易单数组
func (wm *WalletManager) getMultiAddrTransactionsByExplorer(offset, limit int, address ...string) ([]*Transaction, error) {

	var (
		trxs = make([]*Transaction, 0)
	)

	addrs := strings.Join(address, ",")

	request := req.Param{
		"addrs": addrs,
		"from":  offset,
		"to":    offset + limit,
	}

	path := fmt.Sprintf("addrs/txs")

	result, err := wm.BitcoreClient.Call(path, request, "POST")
	if err != nil {
		return nil, err
	}

	if items := result.Get("items"); items.IsArray() {
		for _, obj := range items.Array() {
			tx := wm.newTxByExplorer(&obj)
			trxs = append(trxs, tx)
		}
	}

	return trxs, nil
}*/

//estimateFeeRateByExplorer 通过浏览器获取费率
func (e *Bitcorer) EstimateFeeRateByBitcore() (decimal.Decimal, error) {

	defaultRate, _ := decimal.NewFromString("0.00001")

	path := fmt.Sprintf("fee/%d", 2)

	result, err := e.Call(path, nil, "GET")
	if err != nil {
		return decimal.New(0, 0), err
	}

	feeRate, _ := decimal.NewFromString(result.Get("2").String())

	if feeRate.LessThan(defaultRate) {
		feeRate = defaultRate
	}

	return feeRate, nil
}

//sendRawTransactionByExplorer 广播交易
func (b *Bitcorer) sendRawTransactionByBitcore(txHex string) (string, error) {

	request := req.Param{
		"rawtx": txHex,
	}

	path := fmt.Sprintf("tx/send")

	result, err := b.Call(path, request, "POST")
	if err != nil {
		return "", err
	}

	return result.Get("txid").String(), nil
}
