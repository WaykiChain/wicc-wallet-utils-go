package wicc_wallet_utils_go

import (
	"fmt"
	"math/big"
	"testing"
)


func TestEthCreateSimpleRawTransaction(t *testing.T) {

	privateKeyStr := "6B93D965D9981F9066CCC44B9DBF32B50F411C0DCEDF4A41CA4E7424ABDB617F"
	from := "0x232D23C22543144B988F738C701Df6dfd6eAcA4c"
	to := "81FD1F7aE91041aAc5fCF7d8Ed3e1dd88Cc1359a"
	amount := big.NewInt(100000000000000000)

	txcount ,err:= tw.WalletClient.ethGetTransactionCount(from,LEATEST)
	if err != nil {
		t.Errorf("Failed to ethGetTransactionCount: %v",err)
	}
	t.Log("txcount=",txcount)

	txFeeInfo, err := tw.WalletClient.GetTransactionFeeEstimated(from,to,amount,"")
	if err != nil {
		t.Errorf("Failed to GetTransactionFeeEstimated: %v",err)
	}
	t.Logf("txFeeInfo.GasLimit=%v, txFeeInfo.GasPrice=%v\n",txFeeInfo.GasLimit,txFeeInfo.GasPrice)


	tx := NewETHSimpleTransaction(to,amount,txcount,txFeeInfo.GasLimit,txFeeInfo.GasPrice)
	rawtx ,err :=tx.CreateRawTx(privateKeyStr,tw.Config.ChainID)
	if err != nil{
		fmt.Println("err=",err)
	}
	fmt.Println("rawtx=",rawtx)

}


func TestERC20CreateTransferRawTransaction(t *testing.T) {

	privateKeyStr := "6B93D965D9981F9066CCC44B9DBF32B50F411C0DCEDF4A41CA4E7424ABDB617F"
	from := "0x232D23C22543144B988F738C701Df6dfd6eAcA4c"
	contractAddr := "0x8E1dA42EbC22F91d528ceB9865f241167Ebb8A0f"  //WICC合约
	to := "81FD1F7aE91041aAc5fCF7d8Ed3e1dd88Cc1359a"
	amount := big.NewInt(100000000)

	txcount ,err:= tw.WalletClient.ethGetTransactionCount(from,LEATEST)
	if err != nil {
		t.Errorf("Failed to ethGetTransactionCount: %v",err)
	}
	t.Log("txcount=",txcount)

	data ,err := makeERC20TokenTransferData(contractAddr,to,amount)
	if err != nil {
		t.Errorf("makeERC20TokenTransData, err=%v", err)
	}

	txFeeInfo,err := tw.WalletClient.GetTransactionFeeEstimated(from,contractAddr,big.NewInt(0),data)
	if err != nil {
		t.Errorf("Failed to GetTransactionFeeEstimated: %v",err)
	}
	t.Logf("txFeeInfo.GasLimit=%v, txFeeInfo.GasPrice=%v\n",txFeeInfo.GasLimit,txFeeInfo.GasPrice)

	tx,err := NewERC20TransferTransaction(contractAddr,big.NewInt(0),txcount,txFeeInfo.GasLimit,txFeeInfo.GasPrice,data)
	if err != nil {
		t.Errorf("Failed to NewERC20TransferTransaction: %v",err)
	}
	rawtx ,err :=tx.CreateRawTx(privateKeyStr,tw.Config.ChainID)
	if err != nil{
		fmt.Println("err=",err)
	}
	fmt.Println("rawtx=",rawtx)
}

func TestEthSendRawTransaction(t *testing.T) {

	//rawtx := "f86b06843b9aca008252089481fd1f7ae91041aac5fcf7d8ed3e1dd88cc1359a88016345785d8a00008029a018adef9b4ec654de5ecb7976b52faf6e5ddfc22902dffa8ab5f5eabd07f4c1eca020d35fc7f7b58ee5e53a7e3ea4849eba1b26a7748db373fbb4f48b7f30e5e4ea"
	//rawtx := "f86308843b9aca0082d0ed948e1da42ebc22f91d528ceb9865f241167ebb8a0f808029a034fce598f3abe285c12d53f1bae5b790b598e168b33b7f1ca0e55afb0652837fa06c4a00544fbe49be8bcc321a7bd271c775113d79a1550f583d0887df899b20bd"
	//rawtx := "f8a809843b9aca0082d0ed948ne1da42ebc22f91d528ceb9865f241167ebb8a0f80b844a9059cbb00000000000000000000000081fd1f7ae91041aac5fcf7d8ed3e1dd88cc1359a0000000000000000000000000000000000000000000000000000000005f5e10029a01bef8dd99141dd4cb6461f901e4c03b1df28c6422dc2ee5478ae2b73cb2610d2a032ed60e4137b9619d41930fef0193205fa4d57a174bf0014bb1a246e36cb4480"
	rawtx := "f86b0b843b9aca008252089481fd1f7ae91041aac5fcf7d8ed3e1dd88cc1359a88016345785d8a00008029a0fc16ac404d694561f44e12e9f1ab4916244c352add8836ecc05ae80ef9f7e3e9a05a2c6e6b6fc65373700c6fbf003d055a38c49d626b370b71eec1809978f027fe"

	txid ,err  := tw.WalletClient.EthSendRawTransaction(rawtx)
	if err != nil{
		fmt.Println("err=",err)
	}
	fmt.Println("txid=",txid)
}



