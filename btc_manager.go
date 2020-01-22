package wicc_wallet_utils_go

import (
	"encoding/hex"
	"fmt"
	"github.com/WaykiChain/wicc-wallet-utils-go/bitcoin"
	"github.com/WaykiChain/wicc-wallet-utils-go/log"
	"github.com/imroc/req"
	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
	"math"
	"strings"

	//"strconv"
)

type BTCWalletManager struct {
	Wallet          *BTCWallet
	WalletClient    *bitcoin.Client                       // 节点客户端
	BitcoreClient   *bitcoin.Bitcorer                     // 浏览器API客户端
	Config          *bitcoin.Config                 //钱包管理配置
	Log             *log.OWLogger                 //日志工具
}

func NewBTCWalletManager() *BTCWalletManager {
	wm := BTCWalletManager{}
	wm.Config = bitcoin.NewConfig()
	wm.Wallet = NewBTCWallet(wm.Config.WalletConfig)
	wm.BitcoreClient = bitcoin.NewExplorer(wm.Config.ServerAPI,wm.Config.Debug)
	return &wm
}

//EstimateFee 预估手续费
func (wm *BTCWalletManager) EstimateFee(inputs, outputs int64, feeRate decimal.Decimal) (uint64, error) {

	var piece int64 = 1

	//UTXO如果大于设定限制，则分拆成多笔交易单发送
	if inputs > int64(wm.Config.MaxTxInputs) {
		piece = int64(math.Ceil(float64(inputs) / float64(wm.Config.MaxTxInputs)))
	}

	//size计算公式如下：148 * 输入数额 + 34 * 输出数额 + 10
	size := decimal.New(inputs*148+outputs*34+piece*10,0)
	fmt.Println("size=",size.String())
	fee := size.Div(decimal.New(1000, 0)).Mul(feeRate)  //unit :BTC ,可能结果小数点后大于8位
	trx_fee := fee.Round(wm.Config.Decimals)   //unit :BTC  ,精确到小数点后8位

	//是否低于最小手续费
	if trx_fee.LessThan(wm.Config.MinFees) {
		trx_fee = wm.Config.MinFees
	}
	trx_fee_satoshis := trx_fee.Mul(decimal.New(1, wm.Config.Decimals)).IntPart()

	return uint64(trx_fee_satoshis), nil
}

//EstimateFeeRate 预估的没KB手续费率
func (wm *BTCWalletManager) EstimateFeeRate() (decimal.Decimal, error) {

	if wm.Config.RPCServerType == bitcoin.RPCServerExplorer {
		return wm.BitcoreClient.EstimateFeeRateByBitcore()
	} else {
		return wm.WalletClient.EstimateFeeRateByNode()
	}
}


//ListUnspent 获取未花记录
func (wm *BTCWalletManager) ListUnspent(min uint64, addresses ...string) ([]*Unspent, error) {

	//:分页限制

	var (
		limit       = 100
		searchAddrs = make([]string, 0)
		max         = len(addresses)
		step        = max / limit
		utxo        = make([]*Unspent, 0)
		pice        []*Unspent
		err         error
	)

	for i := 0; i <= step; i++ {
		begin := i * limit
		end := (i + 1) * limit
		if end > max {
			end = max
		}

		searchAddrs = addresses[begin:end]

		if len(searchAddrs) == 0 {
			continue
		}

		if wm.Config.RPCServerType == bitcoin.RPCServerExplorer {
			pice, err = wm.listUnspentByBitcore(min, searchAddrs...)
			if err != nil {
				return nil, err
			}
		} else {
			pice, err = wm.listUnspentByCore(min, searchAddrs...)
			if err != nil {
				return nil, err
			}
		}
		utxo = append(utxo, pice...)
	}
	return utxo, nil
}


//listUnspentByBitcore 获取未花交易
func (wm *BTCWalletManager) listUnspentByBitcore(min uint64, address ...string) ([]*Unspent, error) {

	var (
		utxos = make([]*Unspent, 0)
	)

	addrs := strings.Join(address, ",")

	request := req.Param{
		"addrs": addrs,
	}

	path := "addrs/utxo"

	result, err := wm.BitcoreClient.Call(path, request, "POST")
	if err != nil {
		return nil, err
	}

	array := result.Array()
	for _, a := range array {
		u := NewUnspent(&a)
		if u.Confirmations >= min {
			utxos = append(utxos, NewUnspent(&a))
		}
	}

	return utxos, nil

}

//getTransactionByCore 获取交易单
func (wm *BTCWalletManager) listUnspentByCore(min uint64, addresses ...string) ([]*Unspent, error) {

	var (
		utxos = make([]*Unspent, 0)
	)

	request := []interface{}{
		min,
		9999999,
	}

	if len(addresses) > 0 {
		request = append(request, addresses)
	}

	result, err := wm.WalletClient.Call("listunspent", request)
	if err != nil {
		return nil, err
	}

	array := result.Array()
	for _, a := range array {
		utxos = append(utxos, NewUnspent(&a))
	}

	return utxos, nil
}

func (wm *BTCWalletManager) newTxByExplorer(json *gjson.Result) *Transaction {

	/*
			{
			"txid": "9f5eae5b95016825a437ceb9c9224d3e30d3b351f1100e4df5cc0cacac4e668c",
			"version": 1,
			"locktime": 1433760,
			"vin": [],
			"vout": [],
			"blockhash": "0000000000003ac968ee1ae321f35f76d4dcb685045968d60fc39edb20b0eed0",
			"blockheight": 1433761,
			"confirmations": 5,
			"time": 1539050096,
			"blocktime": 1539050096,
			"valueOut": 0.14652549,
			"size": 814,
			"valueIn": 0.14668889,
			"fees": 0.0001634
		}
	*/
	obj := Transaction{}
	//解析json
	obj.TxID = gjson.Get(json.Raw, "txid").String()
	obj.Version = gjson.Get(json.Raw, "version").Uint()
	obj.LockTime = gjson.Get(json.Raw, "locktime").Int()
	obj.BlockHash = gjson.Get(json.Raw, "blockhash").String()
	blockHeight := gjson.Get(json.Raw, "blockheight").Int()
	if blockHeight < 0 {
		obj.BlockHeight = 0
	} else {
		obj.BlockHeight = uint64(blockHeight)
	}

	obj.Confirmations = gjson.Get(json.Raw, "confirmations").Uint()
	obj.Blocktime = gjson.Get(json.Raw, "blocktime").Int()
	obj.Size = gjson.Get(json.Raw, "size").Uint()
	obj.Fees = gjson.Get(json.Raw, "fees").String()

	obj.Vins = make([]*Vin, 0)
	if vins := gjson.Get(json.Raw, "vin"); vins.IsArray() {
		for _, vin := range vins.Array() {
			input := newTxVinByExplorer(&vin)
			if input != nil {
				obj.Vins = append(obj.Vins, input)
			}
		}
	}

	obj.Vouts = make([]*Vout, 0)
	if vouts := gjson.Get(json.Raw, "vout"); vouts.IsArray() {
		for _, vout := range vouts.Array() {
			output := wm.newTxVoutByExplorer(&vout)
			if output != nil {
				obj.Vouts = append(obj.Vouts, output)
			}
		}
	}

	return &obj
}

func newTxVinByExplorer(json *gjson.Result) *Vin {

	/*
		{
			"txid": "b8c00fff9208cb02f694666084fe0d65c471e92e45cdc3fb2e43af3a772e702d",
			"vout": 0,
			"sequence": 4294967294,
			"n": 0,
			"scriptSig": {
				"hex": "47304402201f77d18435931a6cb51b6dd183decf067f933e92647562f71a33e80988fbc8f6022012abe6824ffa70e5ccb7326e0dbb66144ba71133c1d4a1215da0b17358d7ca660121024d7be1242bd44619779a976cd1cd2d9351fcf58df59929b30a0c69d852302fb5",
				"asm": "304402201f77d18435931a6cb51b6dd183decf067f933e92647562f71a33e80988fbc8f6022012abe6824ffa70e5ccb7326e0dbb66144ba71133c1d4a1215da0b17358d7ca66[ALL] 024d7be1242bd44619779a976cd1cd2d9351fcf58df59929b30a0c69d852302fb5"
			},
			"addr": "msYiUQquCtGucnk3ZaWeJenYmY8WxRoeuv",
			"valueSat": 990000,
			"value": 0.0099,
			"doubleSpentTxID": null
		}
	*/
	obj := Vin{}
	//解析json
	obj.TxID = gjson.Get(json.Raw, "txid").String()
	obj.Vout = gjson.Get(json.Raw, "vout").Uint()
	obj.N = gjson.Get(json.Raw, "n").Uint()
	obj.Addr = gjson.Get(json.Raw, "addr").String()
	obj.Value = gjson.Get(json.Raw, "value").String()
	obj.Coinbase = gjson.Get(json.Raw, "coinbase").String()

	return &obj
}

func (wm *BTCWalletManager) newTxVoutByExplorer(json *gjson.Result) *Vout {

	/*
		{
			"value": "0.01652549",
			"n": 0,
			"scriptPubKey": {
				"hex": "76a9142760a760e8d22b5facb380444920e1197f272ea888ac",
				"asm": "OP_DUP OP_HASH160 2760a760e8d22b5facb380444920e1197f272ea8 OP_EQUALVERIFY OP_CHECKSIG",
				"addresses": ["mj7ASAGw8ia2o7Hqvo2XS1d7jGWr5UgEU9"],
				"type": "pubkeyhash"
			},
			"spentTxId": null,
			"spentIndex": null,
			"spentHeight": null
		}
	*/
	obj := Vout{}
	//解析json
	obj.Value = gjson.Get(json.Raw, "value").String()
	obj.N = gjson.Get(json.Raw, "n").Uint()
	obj.ScriptPubKey = gjson.Get(json.Raw, "scriptPubKey.hex").String()
	asm := gjson.Get(json.Raw, "scriptPubKey.asm").String()

	if len(obj.ScriptPubKey) == 0 {
		scriptPubKey, err := DecodeScript(asm)
		if err == nil {
			obj.ScriptPubKey = hex.EncodeToString(scriptPubKey)
		}
	}

	//提取地址
	if addresses := gjson.Get(json.Raw, "scriptPubKey.addresses"); addresses.IsArray() {
		obj.Addr = addresses.Array()[0].String()
	}

	obj.Type = gjson.Get(json.Raw, "scriptPubKey.type").String()

	/*	if len(obj.Addr) == 0 {

			scriptBytes, _ := hex.DecodeString(obj.ScriptPubKey)
			obj.Addr, _ = wm.Decoder.ScriptPubKeyToBech32Address(scriptBytes)
		}
	*/ //wjq
	if strings.HasPrefix(asm, "OP_RETURN") {
		//OP_RETURN的脚本
		obj.Type = "OP_RETURN"
	}

	return &obj
}