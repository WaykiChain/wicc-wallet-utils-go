package bitcoin

import (
	"github.com/blocktree/go-owcdrivers/btcTransaction"
	"testing"
)


//依赖Chain服务，自动查询From的Unspent记录、余额、手续费率等信息，只创建签名生成rawtx,手动广播
func TestCreatetBTCRawTxRelyChain(t *testing.T){

	wm := NewWalletManager()
	//输入
	from1 := FromInfo{"cT7214EqFAbtpfuMfg36EDHBmbrkErXLb27ERSKTLX1RUr92LFSE",BTCTestnetSegwitW,"2MsNRcbHbMgwbbkfzx86Z4FdHkRp29NPjmD"}
	from2 := FromInfo{"cUQyhR3BbeFMwtqFRrMppFtPPcx6DMQNFEXm8C1yNSzMkEoRoGYM",BTCTestnetW,"mxeBxFWLFAY3G1RKijr91B3kzsX2mTvnYv"}
	fromInfos := []FromInfo{from1,from2}
	//输出
	out1 := VOut{&btcTransaction.Vout{"2N2f4HUsH1hFp36zwSUSoeGxVag5BUrJQwr",55000}}
	out2 := VOut{&btcTransaction.Vout{"mrA3J5RR2etH2FvUodfzBLXnKm5ozQtL7d",2000}}
	outs:= []VOut{out1,out2}

	rawtx ,err := wm.CreatetRawTxRelyChain(fromInfos,outs)
	if err != nil {
		t.Errorf("Failed to CreatetBTCRawTransaction: %v",err)
	}

	t.Log("rawtx=",rawtx)
}

//依赖Chain服务，自动查询From的Unspent记录、余额、手续费率等信息，进行转账 (包括生成rawtx + 广播交易)
func TestSendBTCTransactionRelyChain(t *testing.T){

/*	rawtx ,err := wm.CreatetBTCRawTxRelyChain(fromInfos,outs)
	if err != nil {
		t.Errorf("Failed to CreatetBTCRawTransaction: %v",err)
	}

	t.Log("rawtx=",rawtx)*/

	//SendTx(rawtx)
}


//手动查询Unspent记录，只创建签名生成rawtx,手动广播
func TestCreateBTCTransferRawTx(t *testing.T){
	wm := NewWalletManager()

	/*Case1: 隔离见证 -> 普通地址 -> 找零到from地址*/
/*	fromInfos := make([]FinalTxIn,0)
	tos := make([]btcTransaction.Vout,0)
	from1 := FromInfo{"cT7214EqFAbtpfuMfg36EDHBmbrkErXLb27ERSKTLX1RUr92LFSE",true,&common.BTCTestnetParams,"2MsNRcbHbMgwbbkfzx86Z4FdHkRp29NPjmD"}
	fromInfo1 := FinalTxIn{
		&from1,
		"730a2233013445cf970e489c2ceff52f04c590a4650536ba0e21b359756b215e", //unspent:2385082
		0,
		2385082,
	}
	to1 := btcTransaction.Vout{"mrA3J5RR2etH2FvUodfzBLXnKm5ozQtL7d",99999} //转账金额
	to2 := btcTransaction.Vout{"2MsNRcbHbMgwbbkfzx86Z4FdHkRp29NPjmD",1833083} //找零给自己 = 2385082 - 99999 - 452000
	//fee = 2000 * (1*148 + 2*34 + 10) = 452000
	fromInfos = append(fromInfos,fromInfo1)
	tos = append(tos,to1)
	tos = append(tos,to2)
	rawtx ,err := createtBTCRawTransaction(fromInfos,tos)
	if err != nil {
		t.Errorf("Failed to CreatetBTCRawTransaction: %v",err)
	}
	t.Log("rawtx=",rawtx)*/
	/*隔离见证 -> 普通地址 -> 找零到from地址 End
	txId:a1f80de6442aa80abd20bba7fd2cb1f84dccbfc3d9ccb8ed0fbc1dafcce22f92
	*/

	/*Case2: 隔离见证 -> 普通地址 */
	/*fromInfos := make([]FinalTxIn,0)
	tos := make([]btcTransaction.Vout,0)
	from1 := FromInfo{"cT7214EqFAbtpfuMfg36EDHBmbrkErXLb27ERSKTLX1RUr92LFSE",true,&common.BTCTestnetParams,"2MsNRcbHbMgwbbkfzx86Z4FdHkRp29NPjmD"}
	fromInfo1 := FinalTxIn{
		&from1,
		"a1f80de6442aa80abd20bba7fd2cb1f84dccbfc3d9ccb8ed0fbc1dafcce22f92", //unspent:1833083
		1,
		1833083,
	}
	to1 := btcTransaction.Vout{"mxeBxFWLFAY3G1RKijr91B3kzsX2mTvnYv",1449083} //转账金额
	//fee = 2000 * (1*148 + 1*34 + 10) = 384000
	fromInfos = append(fromInfos,fromInfo1)
	tos = append(tos,to1)

	rawtx ,err := CreatetBTCTransferRawTx(fromInfos,tos)
	if err != nil {
		t.Errorf("Failed to CreatetBTCRawTransaction: %v",err)
	}
	t.Log("rawtx=",rawtx)*/
	/*隔离见证 -> 普通地址
	txId:ee409ae5031af9b27b5dfd177dd5924407c79fca87e3f13cf3536e0ce2fd596b
	*/

	/*Case3: 普通地址 -> 隔离见证地址 */
	fromInfos := make([]FinalTxIn,0)
	from1 := FromInfo{"cUQyhR3BbeFMwtqFRrMppFtPPcx6DMQNFEXm8C1yNSzMkEoRoGYM",BTCTestnetW,"mxeBxFWLFAY3G1RKijr91B3kzsX2mTvnYv"}
	tos := make([]VOut,0)
	fromInfo1 := FinalTxIn{
		&from1,
		"ee409ae5031af9b27b5dfd177dd5924407c79fca87e3f13cf3536e0ce2fd596b", //unspent:1449083
		0,
		1449083,
	}
	to1 := VOut{&btcTransaction.Vout{"2N2f4HUsH1hFp36zwSUSoeGxVag5BUrJQwr",1065802}}//转账金额
	//fee = 2000 * (1*148 + 1*34 + 10) = 383281
	fromInfos = append(fromInfos,fromInfo1)
	tos = append(tos,to1)

	rawtx ,err := wm.CreateTransferRawTx(fromInfos,tos)
	if err != nil {
		t.Errorf("Failed to CreatetBTCRawTransaction: %v",err)
	}
	t.Log("rawtx=",rawtx)
	/*普通地址 -> 隔离见证地址
	txId:eb3180b0dd598b3d2ea4ee44d0a9a84847f630d9cd13d468439f8853f768200e
	*/

}




