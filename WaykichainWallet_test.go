package wiccwallet

import (
	"testing"
	"fmt"
	"wiccwallet/commons"
	"bytes"
	"encoding/hex"
	"btcutil"
	"btcutil/base58"
)

func TestMnemonic(t *testing.T) {
	mnemonic := "empty regular curve turtle student prize toy accuse develop spike scatter ginger"
	//seed := bip.NewSeed(mnemonic, "")
	////fmt.Println(hex.EncodeToString(seed))
	address := commons.GenerateAddress(mnemonic, WAYKI_MAINTNET)
	fmt.Println("地址"+address)
}

func TestMnemonicWIF(t *testing.T) {
	mnemonic := "empty regular curve turtle student prize toy accuse develop spike scatter ginger"
	privateKey := commons.GeneratePrivateKey(mnemonic, WAYKI_MAINTNET)
	fmt.Println("私钥"+privateKey)
	fmt.Println("地址"+commons.ImportPrivateKey(privateKey,WAYKI_MAINTNET))
}

func TestSignContractTx(t *testing.T) {
	var waykiContract commons.WaykiContractTxParams
	waykiContract.Value=10000
	waykiContract.BaseSignTxParams.PrivateKey="YAHcraeGRDpvwBWVccV7NLGAU6uK39nNUTip8srbJSu6HKSTfDcC"
	waykiContract.BaseSignTxParams.RegId="25312-1"
	waykiContract.Appid="470867-1"
	waykiContract.BaseSignTxParams.ValidHeight=670532
	waykiContract.BaseSignTxParams.Fees=100000
	waykiContract.BaseSignTxParams.TxType=commons.TX_CONTRACT
	waykiContract.BaseSignTxParams.Version=1
	binary,_:=hex.DecodeString("f0140000151d000000000000151d000000000000")
	waykiContract.ContractBytes= []byte(binary)
	hash:=waykiContract.SignTX()
	println(hash)
}


func TestSignDelegateTx(t *testing.T) {
	var waykiDelegate commons.WaykiDelegateTxParams
	waykiDelegate.BaseSignTxParams.PrivateKey="YAHcraeGRDpvwBWVccV7NLGAU6uK39nNUTip8srbJSu6HKSTfDcC"
	waykiDelegate.BaseSignTxParams.ValidHeight=663956
	waykiDelegate.BaseSignTxParams.Fees=10000
	waykiDelegate.BaseSignTxParams.RegId="25312-1"
	waykiDelegate.BaseSignTxParams.TxType=commons.TX_DELEGATE
	waykiDelegate.BaseSignTxParams.Version=1

	wif,_ := btcutil.DecodeWIF("YAa1wFCfFnZ5bt4hg9MDeDevTMd1Nu874Mn83hEXwtfAL2vkQE9f")
	key:=wif.PrivKey
	delegateList:=[]commons.OperVoteFund{commons.OperVoteFund{
		commons.MINUS_FUND,key.PubKey().SerializeCompressed(),10000,
	}}
	waykiDelegate.OperVoteFunds=delegateList

	hash:=waykiDelegate.SignTX()
	println(hash)
}

func TestSignRegisterTx(t *testing.T) {
	var waykiRegister commons.WaykiRegisterTxParams
	waykiRegister.BaseSignTxParams.PrivateKey="Y7W4t1wtXmdojGPeHt23HSipZAykpTzwbng9gghT3ePzMtSx1g6y"
	waykiRegister.BaseSignTxParams.ValidHeight=663168
	waykiRegister.BaseSignTxParams.Fees=10000
	waykiRegister.BaseSignTxParams.TxType=commons.TX_REGISTERACCOUNT
	waykiRegister.BaseSignTxParams.Version=1
	hash:=waykiRegister.SignTX()
	println(hash)
}

func TestSignCommonTx(t *testing.T) {
	var waykicommon commons.WaykiCommonTxParams
	waykicommon.Value=10000
	waykicommon.DestAddress="wZujmSBQ7sNhxA7WfEuN46HAyZpw1B8NBA"
	waykicommon.BaseSignTxParams.PrivateKey="Y9XMqNzseQFSK32SvMDNF9J7xz1CQmHRsmY1hMYiqZyTck8pYae3"
	waykicommon.BaseSignTxParams.RegId="30947-1"
	waykicommon.BaseSignTxParams.ValidHeight=662788
	waykicommon.BaseSignTxParams.Fees=10000
	waykicommon.BaseSignTxParams.TxType=commons.TX_COMMON
	waykicommon.BaseSignTxParams.Version=1
	hash:=waykicommon.SignTX()
	println(hash)
}

func TestSignTx(t *testing.T) {

	//mnemonic1 :="fragile chalk speed absorb enter weasel hurdle eternal tooth acoustic cost boss"
	privateKey1 := "Y9XMqNzseQFSK32SvMDNF9J7xz1CQmHRsmY1hMYiqZyTck8pYae3"//commons.GeneratePrivateKey(mnemonic1, WAYKI_TESTNET)
	srcAddress  := commons.ImportPrivateKey(privateKey1,WAYKI_TESTNET)
	fmt.Println("私钥1:  "+privateKey1)
	fmt.Println("地址1:  "+srcAddress)


	mnemonic2 := "empty regular curve turtle student prize toy accuse develop spike scatter ginger"
	privateKey2 := commons.GeneratePrivateKey(mnemonic2, WAYKI_TESTNET)
	destAddress  := "wZujmSBQ7sNhxA7WfEuN46HAyZpw1B8NBA"//commons.ImportPrivateKey(privateKey2,WAYKI_TESTNET)
	fmt.Println("私钥2:  "+privateKey2)
	fmt.Println("地址2:  "+destAddress)

	bytesBuffer := bytes.NewBuffer([]byte{})
	ss1 :=int64(commons.TX_COMMON)  //txtype
    ss2 :=int64(1) //tx version
    ss3:=int64(662251) //tx height
    ss4:=int64(4) //
    ss5:=int64(30947)  //regheight
	ss6:=int64(1)  //regindex
	bytesBuffer.Write(commons.EncodeInOldWay(ss1))
	bytesBuffer.Write(commons.EncodeInOldWay(ss2))
	bytesBuffer.Write(commons.EncodeInOldWay(ss3))
	bytesBuffer.Write(commons.EncodeInOldWay(ss4))
	bytesBuffer.Write(commons.EncodeInOldWay(ss5))
	bytesBuffer.Write(commons.EncodeInOldWay(ss6))




	 ss7,_,_:= base58.CheckDecode(destAddress)//base58.Decode(destAddress)
	//ss81,_:= btcutil.Address()//btcutil.NewAddressScriptHash(ss7,&commons.WaykiTestParams)
	//ss81,_:=btcutil.NewAddressWitnessScriptHash(ss7,&commons.WaykiTestParams)
	//ss8:=ss81.ScriptAddress()
	//println(ss81.IsForNet(&commons.WaykiTestParams))
	 bytesBuffer.Write(commons.EncodeInOldWay(int64(len(ss7))))
	bytesBuffer.Write(ss7)
	bytesBuffer.Write(commons.EncodeInOldWay(int64(10000)))
	bytesBuffer.Write(commons.EncodeInOldWay(int64(10000)))
	bytesBuffer.Write(commons.EncodeInOldWay(int64(0)))

	ss9:=signhash(destAddress,privateKey1)
	bytesBuffer.Write(commons.EncodeInOldWay(int64(len(ss9))))
	bytesBuffer.Write(ss9)
	println(hex.EncodeToString(bytesBuffer.Bytes()))
}
func signhash(destAddress string,priv string)  []byte{

	bytesBuffer := bytes.NewBuffer([]byte{})
	ss1 :=int64(1)  //txtype
	ss2 :=3 //tx version
	ss3:=int64(662251) //tx height
	ss4:=int64(4) //
	ss5:=int64(30947)  //regheight
	ss6:=int64(1)  //regindex
	bytesBuffer.Write(commons.EncodeInOldWay(ss1))
	bytesBuffer.WriteByte(byte(ss2))
	bytesBuffer.Write(commons.EncodeInOldWay(ss3))
	bytesBuffer.Write(commons.EncodeInOldWay(ss4))
	bytesBuffer.Write(commons.EncodeInOldWay(ss5))
	bytesBuffer.Write(commons.EncodeInOldWay(ss6))

	//ss7,_:=btcutil.DecodeAddress(destAddress,&commons.WaykiTestParams) //toaddress
	//ss8, _ := txscript.PayToAddrScript(ss7)
	ss7,_,_:= base58.CheckDecode(destAddress)//base58.Decode(destAddress)
	//ss81,_:=btcutil2.NewAddressScriptHash(ss7,&commons.WaykiTestParams)
	//ss8 := ss81.Hash160()//txscript.PayToAddrScript(ss7)
	bytesBuffer.Write(commons.EncodeInOldWay(int64(len(ss7))))
	bytesBuffer.Write(ss7)

	bytesBuffer.Write(commons.EncodeInOldWay(10000))
	bytesBuffer.Write(commons.EncodeInOldWay(10000))
	bytesBuffer.Write(commons.EncodeInOldWay(0))

	data1,_:=commons.HashDoubleSha256(bytesBuffer.Bytes())
	wif, _ := btcutil.DecodeWIF(priv)
	key:=wif.PrivKey
	ss,_:=key.Sign(data1)//secp256k1.Sign(data1,wif.PrivKey.Serialize())
	return ss.Serialize()
}



