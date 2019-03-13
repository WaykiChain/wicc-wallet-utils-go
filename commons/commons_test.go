package commons

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
)

func TestSignCallContractTx(t *testing.T) {

	privateKey := "Y9sx4Y8sBAbWDAqAWytYuUnJige3ZPwKDZp1SCDqqRby1YMgRG9c"
	wifKey, _ := btcutil.DecodeWIF(privateKey)

	var tx WaykiCallContractTx
	tx.Values = 10000
	tx.UserId = NewRegUidByStr("7849-1")
	tx.AppId = NewRegUidByStr("20988-1")
	tx.ValidHeight = 22365
	tx.Fees = 100000
	tx.TxType = CONTRACT_TX
	tx.Version = 1
	binary, _ := hex.DecodeString("f017")
	tx.Contract = binary
	hash := tx.SignTx(wifKey)
	println(hash)
}

func TestSignDelegateTx(t *testing.T) {
	privateKey := "Y9sx4Y8sBAbWDAqAWytYuUnJige3ZPwKDZp1SCDqqRby1YMgRG9c"
	wifKey, _ := btcutil.DecodeWIF(privateKey)

	var tx WaykiDelegateTx
	tx.ValidHeight = 40935
	tx.Fees = 1000
	tx.UserId = NewRegUidByStr("0-1")
	tx.TxType = DELEGATE_TX
	tx.Version = 1

	miner1Key := "Y5F2GraTdQqMbYrV6MG78Kbg4QE8p4B2DyxMdLMH7HmDNtiNmcbM"
	miner1Wif, _ := btcutil.DecodeWIF(miner1Key)
	miner1PubKey := NewPubKeyIdByKey(miner1Wif.PrivKey)
	miner2PubKey := NewPubKeyIdByStr("025a37cb6ec9f63bb17e562865e006f0bafa9afbd8a846bd87fc8ff9e35db1252e")
	delegateList := []OperVoteFund{
		OperVoteFund{ADD_FUND, miner1PubKey, 111},
		OperVoteFund{ADD_FUND, miner2PubKey, 222},
	}
	tx.OperVoteFunds = delegateList

	hash := tx.SignTx(wifKey)
	println(hash)
}

func TestSignRegisterAccountTx(t *testing.T) {

	privateKey := "Y9sx4Y8sBAbWDAqAWytYuUnJige3ZPwKDZp1SCDqqRby1YMgRG9c"
	wifKey, _ := btcutil.DecodeWIF(privateKey)

	var tx WaykiRegisterAccountTx

	tx.TxType = REG_ACCT_TX
	tx.Version = 1
	tx.ValidHeight = 7783
	tx.UserId = NewPubKeyUid(*NewPubKeyIdByKey(wifKey.PrivKey))
	tx.Fees = 10001
	hash := tx.SignTx(wifKey)
	println(hash)
}

func TestSignCommonTx(t *testing.T) {
	privateKey := "Y9sx4Y8sBAbWDAqAWytYuUnJige3ZPwKDZp1SCDqqRby1YMgRG9c"
	wifKey, _ := btcutil.DecodeWIF(privateKey)

	var tx WaykiCommonTx
	tx.TxType = COMMON_TX
	tx.ValidHeight = 14897
	tx.Version = 1
	tx.UserId = NewRegUidByStr("158-1")
	tx.DestId = NewAdressUidByStr("wSSbTePArv6BkDsQW9gpGCTX55AXVxVKbd")
	tx.Fees = 10000
	tx.Values = 10000
	hash := tx.SignTx(wifKey)
	println(hash)
}

func TestSignRegisterContractTx(t *testing.T) {
	privateKey := "Y9sx4Y8sBAbWDAqAWytYuUnJige3ZPwKDZp1SCDqqRby1YMgRG9c"
	wifKey, _ := btcutil.DecodeWIF(privateKey)

	script, err := ioutil.ReadFile("../demo/data/hello.lua")
	if err != nil {
		t.Error("Read contract script file err: ", err)
	}
	var tx WaykiRegisterContractTx

	tx.TxType = REG_CONT_TX
	tx.Version = 1
	tx.ValidHeight = 20999
	tx.UserId = NewRegUidByStr("7849-1")
	tx.Script = script
	tx.Description = "My hello contract!!!"

	tx.Fees = 110000001
	hash := tx.SignTx(wifKey)
	println(hash)
}

func TestSignRewardTx(t *testing.T) {
	privateKey := "Y9sx4Y8sBAbWDAqAWytYuUnJige3ZPwKDZp1SCDqqRby1YMgRG9c"
	wifKey, _ := btcutil.DecodeWIF(privateKey)

	var tx WaykiRewardTx

	tx.TxType = REWARD_TX
	tx.Version = 1
	tx.ValidHeight = 14599
	tx.UserId = NewRegUidByStr("7849-1")
	tx.Values = 123456
	hash := tx.SignTx(wifKey)
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
	ss1 := int64(COMMON_TX) //txtype
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
