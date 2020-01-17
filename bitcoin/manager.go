package bitcoin

import (
	"fmt"
	"github.com/WaykiChain/wicc-wallet-utils-go/log"
	"github.com/shopspring/decimal"
	"math"
	//"strconv"
)

type WalletManager struct {
	Wallet          *BTCWallet
	WalletClient    *Client                       // 节点客户端
	BitcoreClient   *Bitcorer                     // 浏览器API客户端
	Config          *Config                 //钱包管理配置
	Log             *log.OWLogger                 //日志工具
}

func NewWalletManager() *WalletManager {
	wm := WalletManager{}
	wm.Config = NewConfig()
	wm.Wallet = NewBTCWallet(NewWalletConfig(wm.Config.WalletConfig))
	wm.BitcoreClient = NewExplorer(wm.Config.ServerAPI,wm.Config.Debug)
	return &wm
}

//EstimateFee 预估手续费
func (wm *WalletManager) EstimateFee(inputs, outputs int64, feeRate decimal.Decimal) (uint64, error) {

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
func (wm *WalletManager) EstimateFeeRate() (decimal.Decimal, error) {

	if wm.Config.RPCServerType == RPCServerExplorer {
		return wm.BitcoreClient.estimateFeeRateByBitcore()
	} else {
		return wm.WalletClient.estimateFeeRateByNode()
	}
}


//ListUnspent 获取未花记录
func (wm *WalletManager) ListUnspent(min uint64, addresses ...string) ([]*Unspent, error) {

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

		if wm.Config.RPCServerType == RPCServerExplorer {
			pice, err = wm.BitcoreClient.listUnspentByBitcore(min, searchAddrs...)
			if err != nil {
				return nil, err
			}
		} else {
			pice, err = wm.WalletClient.getListUnspentByCore(min, searchAddrs...)
			if err != nil {
				return nil, err
			}
		}
		utxo = append(utxo, pice...)
	}
	return utxo, nil
}

