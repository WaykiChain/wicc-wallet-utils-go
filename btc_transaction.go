package wicc_wallet_utils_go

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/WaykiChain/wicc-wallet-utils-go/common"
	"github.com/WaykiChain/wicc-wallet-utils-go/common/hash"
	"github.com/WaykiChain/wicc-wallet-utils-go/log"
	"github.com/blocktree/go-owcdrivers/btcTransaction"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
)

type VOutPuts struct{
	VOuts []VOutPut
}
//用于签名输出的结构体
type VOutPut struct{
	Vout *btcTransaction.Vout
}

type VInPuts struct{
	VIns []VInPut
}
//用于签名输入的结构体
type VInPut struct{
	AddrInfo      *FromInfo  //
	PrevTxid 	  string     //输入对应的上一笔输出的交易哈希
	VoutIndex 	  uint64     //输入对应的上一笔交易哈希的索引
	Amount 		  uint64     //输入金额
}
type FromInfo struct{
	WIFPrivateKey string	 //WIF格式私钥，用于签名
	BTCWallet  *BTCWallet    //私钥是否为隔离见证地址标识; 主网、测试网
	Address string			 //输入(From)地址,用于校验是否与私钥对应
}

//用于生成脚本
type Sricpt struct{
	lockScript string
	redeemScript string
}

//return Lockscript 、Redeemscript
func GetSricpt(txin VInPut) (*Sricpt, error){

	address :=  txin.AddrInfo.Address
	pubKeyHash,_,_ :=  base58.CheckDecode(address)

	//如果地址是隔离见证,需要获取Lockscript和Redeemscript
	//Lockscript :
	// 				隔离见证: a914xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx87  23字节  P2WPKH = P2SH
	//				普通地址：76a914xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxac  25字节  P2PKH
	//Redeemscript: (只有隔离见证时需要)
	//						 0014xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx  22字节
	if txin.AddrInfo.BTCWallet.wallet.IsSegwit == true {
		//Get Lockscript
		lockScriptBytes,_ := PayToScriptHashScript(pubKeyHash)
		lockScript := hex.EncodeToString(lockScriptBytes)
		//Get Redeemscript
		wif, err := btcutil.DecodeWIF(txin.AddrInfo.WIFPrivateKey)
		if err != nil {
			return nil, err
		}
		pubKey := wif.PrivKey.PubKey().SerializeCompressed() //33 bytes 公钥
		pubHash := hash.Hash160(pubKey)
		redeemScriptBytes, err := PayToWitnessPubKeyHashScript(pubHash)
		if err != nil {
			return nil, err
		}
		redeemScript := hex.EncodeToString(redeemScriptBytes)

		return &Sricpt{lockScript ,redeemScript},nil

	} else if txin.AddrInfo.BTCWallet.wallet.IsSegwit == false{//如果地址是普通地址,只需要获取Lockscript
		lockScriptBytes,_ := PayToPubKeyHashScript(pubKeyHash)
		lockScript := hex.EncodeToString(lockScriptBytes)
		return &Sricpt{lockScript ,""},nil
	}else{
		return nil,errors.New("txin error!")
	}
}

//将WIF私钥转换成hex
func ConvertWIFToHex(WIF string) ([]byte,error){
	wif, err := btcutil.DecodeWIF(WIF)
	if err != nil {
		return nil,err
	}
	return wif.PrivKey.Serialize(), nil
}


type AddrBalance struct{
	Address string
	ConfirmedBalance uint64
	UnconfirmBalance uint64
	Balance uint64
}

//calculateAddrBalance 通过未花计算余额
func calculateAddrBalance(utxos []*Unspent) map[string]*AddrBalance {

	addrBalanceMap := make(map[string]*AddrBalance)

	for _, utxo := range utxos {

		obj, exist := addrBalanceMap[utxo.Address]
		if !exist {
			obj = &AddrBalance{}
		}

		tu := obj.UnconfirmBalance
		tb := obj.ConfirmedBalance

		if utxo.Spendable {
			if utxo.Confirmations > 0 {
				tb = tb + utxo.Satoshis
			} else {
				tu = tu + utxo.Satoshis
			}
		}

		obj.Address = utxo.Address
		obj.ConfirmedBalance = tb
		obj.UnconfirmBalance = tu
		obj.Balance = tb + tu

		addrBalanceMap[utxo.Address] = obj
	}

	return addrBalanceMap
}


func newFromInfoMap(from FromInfo) (map[string]*FromInfo,error){

	addrMap := make(map[string]*FromInfo)
	//根据私钥获取地址
	address,err := from.BTCWallet.wallet.GenerateAddressFromPrivateKey(from.WIFPrivateKey)
	if err != nil{
		return nil, err
	}
	//判断输入地址存在时是否与私钥匹配
	if (len(from.Address) != 0) && (address != from.Address) {
		return nil, errors.New("Private key and address don't match！")
	}
	//记录所有键值对
	addrMap[address] = &FromInfo{from.WIFPrivateKey,from.BTCWallet,address}

	return addrMap, nil
}


//创建BTC转账rawtx,依赖外部服务查询UTXO、估算手续费等
func (wm *BTCWalletManager) CreateBTCTransferRawTx(ins *VInPuts, outs *VOutPuts) (string,error){

	var(
		addressInfo  = make(map[string]*FromInfo) //地址 - 私钥/是否隔离见证/链网络类型/地址 的映射
		toAmount     = uint64(0) //输出总金额（真正的转账金额）
		fees    	 = uint64(0) //转账手续费, 从链上估算
		usedUTXO 	 = make([]*Unspent, 0) //真正作为输入的utxo
		balance      = uint64(0)//usedUTXO的总金额
		actualOut    = uint64(0) //toAmount + fees(转账金额 + 手续费)
		finalTxIns   = make([]VInPut, 0) //最终用于签名的输入
		finalTxOuts  = make([]VOutPut, 0) //最终用于签名的输出
	)
	//判断输入地址是否存在
	if len(ins.VIns) <= 0 {
		return "", errors.New("[]InAddr have not addresses")
	}

	//遍历读取传入的的输入 地址数组
	searchAddrs := make([]string, 0)
	for _, in := range ins.VIns {
		//获取From地址的映射，目的用于获取私钥等信息,组装到最终签名的输入结构体中
		addressInfo, err := newFromInfoMap(*in.AddrInfo)
		if err != nil{
			return "", err
		}

		searchAddrs = append(searchAddrs, addressInfo[in.AddrInfo.Address].Address)
	}

	//查找账户的utxo，包含已确认和未确认的
	unspents, err := wm.ListUnspent(0, searchAddrs...)
	if err != nil {
		return "", err
	}

	//判断输入地址的utxo是否存在，不存在则说明余额不足
	if len(unspents) <= 0 {
		return "", errors.New("error, balance is enough: UTXO don't exist")
	}

	if len(outs.VOuts) <= 0 {
		return "", errors.New("Receiver addresses is empty!")
	}

	//计算总输出金额+ 构建最终用于签名的输出
	for _, out := range outs.VOuts {
		toAmount = toAmount + uint64(out.Vout.Amount)
		finalTxOuts = append(finalTxOuts,out)
	}

	//暂不排序
/*	sort.Sort(UnspentSort{unspents, func(a, b *Unspent) int {
		if a.Amount > b.Amount {
			return 1
		} else {
			return -1
		}
	}})*/

	//估算费率 unit: BTC/Kb
	feesRate, err := wm.EstimateFeeRate()
	if err != nil {
		return "", err
	}

	log.Info("Calculating wallet unspent record to build transaction...")

	//计算一个可用于支付的余额
	for _, unspent := range unspents {
		//将查询且用到的uxto成员赋值给用于最终签名的输入
		finalTxin := VInPut{addressInfo[unspent.Address],unspent.TxID,unspent.Vout,unspent.Satoshis}
		finalTxIns = append(finalTxIns, finalTxin)

		balance = balance + unspent.Satoshis
		usedUTXO = append(usedUTXO, unspent)
		//计算手续费，输出数量需要在真实的输出数量上 + 1个找零地址
		fees, err = wm.EstimateFee(int64(len(usedUTXO)), int64(len(outs.VOuts)+1), feesRate)
		if err != nil {
			return "", err
		}

		//真正需要花出去的钱 = 输出的金额 + 手续费
		actualOut = toAmount + fees
		//如果选中的utxo总余额 >= 真正需要花出去的钱，即退出
		if balance >= actualOut {
			break
		}
	}

	//如果所有utxo的余额总和小于真正需要花出去的钱
	if balance < actualOut {
		return "", errors.New(fmt.Sprintf("The balance: %v is not enough! ",balance))
	}

	//UTXO如果大于设定限制，则分拆成多笔交易单发送
	if len(usedUTXO) > wm.Config.MaxTxInputs {
		errStr := fmt.Sprintf("The transaction is use max inputs over: %d", wm.Config.MaxTxInputs)
		return "", errors.New(errStr)
	}

	//取账户第一个地址作为找零地址
	changeAddress := usedUTXO[0].Address
	//找零金额 = 选中的utxo总额 - (转账金额 + 手续费)
	changeAmount := balance - actualOut

	wm.Log.Std.Notice("-----------------------------------------------")
	wm.Log.Std.Notice("To Address: %v", outs)
	wm.Log.Std.Notice("Use UTXO balance: %v", balance)
	wm.Log.Std.Notice("Fees: %v", fees)
	wm.Log.Std.Notice("toAmount: %v", toAmount)
	wm.Log.Std.Notice("Change Amount: %v", changeAmount)
	wm.Log.Std.Notice("Change Address: %v", changeAddress)
	wm.Log.Std.Notice("-----------------------------------------------")

	//如果找零金额 > 0,把找零地址装进最终签名的输出
	if changeAmount > 0 {
		finalTxOuts = append(finalTxOuts, VOutPut{&btcTransaction.Vout{changeAddress,changeAmount}})
	}

	//签名
	rawTxHex, err := CreateBTCTransferRawTx(&VInPuts{finalTxIns},&VOutPuts{finalTxOuts})
	if err != nil{
		return "", err
	}

	return rawTxHex,nil
}

//创建交易广播之前的rawtx
func CreateBTCTransferRawTx( txins *VInPuts,  txouts *VOutPuts) (string, error) {

	var (
		vins = make([]btcTransaction.Vin,0)
		vouts = make([]btcTransaction.Vout,0)
		txUnlocks = make([]btcTransaction.TxUnlock,0)
		addressPrefix btcTransaction.AddressPrefix
		segwit = false
	)

	for _, txin := range txins.VIns {
		vin := btcTransaction.Vin{txin.PrevTxid,uint32(txin.VoutIndex)}
		vins = append(vins,vin)

		script,_ := GetSricpt(txin)
		txunlock := btcTransaction.TxUnlock{script.lockScript,script.redeemScript,uint64(txin.Amount),btcTransaction.SigHashAll}
		txUnlocks = append(txUnlocks,txunlock)
	}

	//check segwit or not
	for _, txin := range txins.VIns {
		if txin.AddrInfo.BTCWallet.wallet.IsSegwit == true{
			segwit = true
			break
		}
	}

	//create vouts
	for _, txout := range txouts.VOuts{
		vout := btcTransaction.Vout{txout.Vout.Address,txout.Vout.Amount}
		vouts = append(vouts,vout)
	}

	//锁定时间:可以设置某个高度之后交易才入块, 一般设置0
	lockTime := uint32(0)

	//追加手续费支持
	replaceable := false

	if txins.VIns[0].AddrInfo.BTCWallet.wallet.NetParam == &common.BTCTestnetParams{
		addressPrefix = btcTransaction.BTCTestnetAddressPrefix
	}else if txins.VIns[0].AddrInfo.BTCWallet.wallet.NetParam == &common.BTCParams {
		addressPrefix = btcTransaction.BTCMainnetAddressPrefix
	}else{
		return "", errors.New("Net err! Neither the Mainnet nor the Testnet")
	}

	///////构建空交易单
	emptyTrans, err := btcTransaction.CreateEmptyRawTransaction(vins, vouts, lockTime, replaceable, addressPrefix)
	if err != nil {
		return "", err
	}

	/////////计算待签名交易单哈希
	transHash, err := btcTransaction.CreateRawTransactionHashForSig(emptyTrans, txUnlocks, segwit, addressPrefix)
	if err != nil {
		return "", err
	}

	for i, txin := range txins.VIns {
		privateKey,_ := ConvertWIFToHex(txin.AddrInfo.WIFPrivateKey)
		sigPub, err := btcTransaction.SignRawTransactionHash(transHash[i].Hash, privateKey)
		if err != nil {
			return "", err
		}

		transHash[i].Normal.SigPub = *sigPub
	}

	//交易单合并
	signedTrans, err := btcTransaction.InsertSignatureIntoEmptyTransaction(emptyTrans, transHash,txUnlocks , segwit)
	if err != nil {
		return "", err
	}

	// 验证交易单
	pass := btcTransaction.VerifyRawTransaction(signedTrans, txUnlocks, segwit, addressPrefix)
	if pass == false {
		return "", errors.New("VerifyRawTransaction failed!")
	}

	return signedTrans ,nil
}


// payToPubKeyHashScript creates a new script to pay a transaction
// output to a 20-byte pubkey hash. It is expected that the input is a valid
// hash.

/* P2PKH
	address = mxeBxFWLFAY3G1RKijr91B3kzsX2mTvnYv :普通地址
	pubKeyHash =  base58.CheckDecode(address)
*/
func PayToPubKeyHashScript(pubKeyHash []byte) ([]byte, error) {
	return txscript.NewScriptBuilder().AddOp(txscript.OP_DUP).AddOp(txscript.OP_HASH160).
		AddData(pubKeyHash).AddOp(txscript.OP_EQUALVERIFY).AddOp(txscript.OP_CHECKSIG).
		Script()
}

// payToWitnessPubKeyHashScript creates a new script to pay to a version 0
// pubkey hash witness program. The passed hash is expected to be valid.

/* P2WPKH
	address =
	pubKeyHash =  base58.CheckDecode(address)
*/
func PayToWitnessPubKeyHashScript(pubKeyHash []byte) ([]byte, error) {
	return txscript.NewScriptBuilder().AddOp(txscript.OP_0).AddData(pubKeyHash).Script()
}

// payToScriptHashScript creates a new script to pay a transaction output to a
// script hash. It is expected that the input is a valid hash.

/* P2SH
	address =
	pubKeyHash =  base58.CheckDecode(address)
*/
func PayToScriptHashScript(scriptHash []byte) ([]byte, error) {
	return txscript.NewScriptBuilder().AddOp(txscript.OP_HASH160).AddData(scriptHash).
		AddOp(txscript.OP_EQUAL).Script()
}

// payToWitnessPubKeyHashScript creates a new script to pay to a version 0
// script hash witness program. The passed hash is expected to be valid.

/*P2WSH
*/
func PayToWitnessScriptHashScript(scriptHash []byte) ([]byte, error) {
	return txscript.NewScriptBuilder().AddOp(txscript.OP_0).AddData(scriptHash).Script()
}

// payToPubkeyScript creates a new script to pay a transaction output to a
// public key. It is expected that the input is a valid pubkey.
/* P2PK
	address =
	pubKeyHash =  base58.CheckDecode(address)
*/
func PayToPubKeyScript(serializedPubKey []byte) ([]byte, error) {
	return txscript.NewScriptBuilder().AddData(serializedPubKey).
		AddOp(txscript.OP_CHECKSIG).Script()
}
