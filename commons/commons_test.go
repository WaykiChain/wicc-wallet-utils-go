package commons

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
)

func TestSignContractTx(t *testing.T) {
	var waykiContract WaykiContractTxParams
	waykiContract.Values = 10000
	waykiContract.PrivateKey = "YAHcraeGRDpvwBWVccV7NLGAU6uK39nNUTip8srbJSu6HKSTfDcC"
	waykiContract.UserId = "25312-1"
	waykiContract.AppId = "470867-1"
	waykiContract.ValidHeight = 670532
	waykiContract.Fees = 100000
	waykiContract.TxType = TX_CONTRACT
	waykiContract.Version = 1
	binary, _ := hex.DecodeString("f0140000151d000000000000151d000000000000")
	waykiContract.Contract = binary
	hash := waykiContract.SignTX()
	println(hash)
}

func TestSignDelegateTx(t *testing.T) {
	var waykiDelegate WaykiDelegateTxParams
	waykiDelegate.PrivateKey = "YAHcraeGRDpvwBWVccV7NLGAU6uK39nNUTip8srbJSu6HKSTfDcC"
	waykiDelegate.ValidHeight = 663956
	waykiDelegate.Fees = 10000
	waykiDelegate.UserId = "25312-1"
	waykiDelegate.TxType = TX_DELEGATE
	waykiDelegate.Version = 1

	wif, _ := btcutil.DecodeWIF("YAa1wFCfFnZ5bt4hg9MDeDevTMd1Nu874Mn83hEXwtfAL2vkQE9f")
	key := wif.PrivKey
	delegateList := []OperVoteFund{OperVoteFund{
		MINUS_FUND, key.PubKey().SerializeCompressed(), 10000,
	}}
	waykiDelegate.OperVoteFunds = delegateList

	hash := waykiDelegate.SignTX()
	println(hash)
}

func TestSignRegisterTx(t *testing.T) {
	var waykiRegister WaykiRegisterTxParams
	waykiRegister.PrivateKey = "Y9sx4Y8sBAbWDAqAWytYuUnJige3ZPwKDZp1SCDqqRby1YMgRG9c"
	waykiRegister.ValidHeight = 7783
	waykiRegister.Fees = 10001
	waykiRegister.TxType = TX_REGISTERACCOUNT
	waykiRegister.Version = 1
	hash := waykiRegister.SignTX()
	println(hash)
}

func TestSignCommonTx(t *testing.T) {
	var waykicommon WaykiCommonTxParams
	waykicommon.Values = 10000
	waykicommon.DestAddress = "wSSbTePArv6BkDsQW9gpGCTX55AXVxVKbd"
	waykicommon.PrivateKey = "Y7V1jwCRr8D3tyPTkcsjgBTHwZN45b1U3ueZfJ5oWVJqwcKpArou"
	waykicommon.UserId = "158-1"
	waykicommon.ValidHeight = 8107
	waykicommon.Fees = 10000
	waykicommon.TxType = TX_COMMON
	waykicommon.Version = 1
	hash := waykicommon.SignTX()
	println(hash)
}

func TestSignTx(t *testing.T) {

	privateKey1 := "Y9XMqNzseQFSK32SvMDNF9J7xz1CQmHRsmY1hMYiqZyTck8pYae3" //GeneratePrivateKey(mnemonic1, WAYKI_TESTNET)
	srcAddress := ImportPrivateKey(privateKey1, TESTNET)
	fmt.Println("私钥1:  " + privateKey1)
	fmt.Println("地址1:  " + srcAddress)

	mnemonic2 := "empty regular curve turtle student prize toy accuse develop spike scatter ginger"
	privateKey2 := GeneratePrivateKey(mnemonic2, TESTNET)
	destAddress := "wZujmSBQ7sNhxA7WfEuN46HAyZpw1B8NBA" //ImportPrivateKey(privateKey2,WAYKI_TESTNET)
	fmt.Println("私钥2:  " + privateKey2)
	fmt.Println("地址2:  " + destAddress)

	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)
	ss1 := int64(TX_COMMON) //txtype
	ss2 := int64(1)         //tx version
	ss3 := int64(662251)    //tx height
	ss4 := int64(4)         //
	ss5 := int64(30947)     //regheight
	ss6 := int64(1)         //regindex
	writer.WriteByte(byte(ss1))
	writer.WriteVarInt(ss2)
	writer.WriteVarInt(ss3)
	writer.WriteVarInt(ss4)
	writer.WriteVarInt(ss5)
	writer.WriteVarInt(ss6)

	ss7, _, _ := base58.CheckDecode(destAddress) //base58.Decode(destAddress)
	//ss81,_:= btcutil.Address()//btcutil.NewAddressScriptHash(ss7,&WaykiTestParams)
	//ss81,_:=btcutil.NewAddressWitnessScriptHash(ss7,&WaykiTestParams)
	//ss8:=ss81.ScriptAddress()
	//println(ss81.IsForNet(&WaykiTestParams))
	writer.WriteBytes(ss7)

	writer.WriteVarInt(100000) // Fees

	writer.WriteVarInt(0) // Contract, empty bytes, only write the length=0

	ss9 := signhash(destAddress, privateKey1)
	writer.WriteBytes(ss9)
	println(hex.EncodeToString(buf.Bytes()))
}

func signhash(destAddress string, priv string) []byte {

	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)
	ss1 := int64(1)             //txtype
	ss2 := int64(3)             //tx version
	ss3 := int64(662251)        //tx height
	ss4 := int64(4)             //
	ss5 := int64(30947)         //regheight
	ss6 := int64(1)             //regindex
	writer.WriteVarInt(ss2)     // Version
	writer.WriteByte(byte(ss1)) // txtype
	writer.WriteVarInt(ss3)
	writer.WriteVarInt(ss4)
	writer.WriteVarInt(ss5)
	writer.WriteVarInt(ss6)

	//ss7,_:=btcutil.DecodeAddress(destAddress,&WaykiTestParams) //toaddress
	//ss8, _ := txscript.PayToAddrScript(ss7)
	ss7, _, _ := base58.CheckDecode(destAddress) //base58.Decode(destAddress)
	//ss81,_:=btcutil2.NewAddressScriptHash(ss7,&WaykiTestParams)
	//ss8 := ss81.Hash160()//txscript.PayToAddrScript(ss7)
	writer.WriteBytes(ss7)
	writer.WriteVarInt(100000) // Fees
	writer.WriteVarInt(0)      // Contract, empty bytes, only write the length=0

	data1, _ := HashDoubleSha256(buf.Bytes())
	wif, _ := btcutil.DecodeWIF(priv)
	key := wif.PrivKey
	ss, _ := key.Sign(data1) //secp256k1.Sign(data1,wif.PrivKey.Serialize())
	return ss.Serialize()
}
