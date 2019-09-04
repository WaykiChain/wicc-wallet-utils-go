# wicc-wallet-utils-go
WaykiChain Wallet Library for Offline Operation written in golang

## Install

* go get github.com/WaykiChain/wicc-wallet-utils-go
* go get github.com/btcsuite

## Usage

### WaykiChain Create Wallet
Generate mnemonics and private keys to manage your wallet.

```go
func GenerateMnemonics() string

func GetPrivateKeyFromMnemonic(words string, netType int) string

func GetPubKeyFromPrivateKey(privKey string) (string,error) 

func GetAddressFromPrivateKey(privateKey string, netType int) string

func GetAddressFromMnemonic(words string, netType int) string
```
- GenerateMnemonics. You will get 12 words.

- GetPrivateKeyFromMnemonic. You should provide your mnemonic and network Type (1 MainNet,2 TestNet),function return private Key,MainNet Private key start with "P" ,TestNet
private key start with "Y".

- GetPubKeyFromPrivateKey. you should provide your Private Key,the function return wallet public key as hex string.

- GetAddressFromPrivateKey. you should provide your Private Key,the function return wallet Address as base58 encode string,MainNet Address start with "W",TestNet Address start with "w".

- GetAddressFromMnemonic. you should provide your mnemonic,the function return wallet Address as base58 encode string,MainNet Address start with "W",TestNet Address start with "w".

Examples:

Generate mnemonic:
```go
mnemonic := GenerateMnemonics()
```
Get private key from mnemonic:
```go
mnemonic := "empty regular curve turtle student prize toy accuse develop spike scatter ginger"
privateKey := GetPrivateKeyFromMnemonic(mnemonic, WAYKI_MAINTNET)
```
Get public key from private key:
```go
publicKey,_:=GetPubKeyFromPrivateKey("Y6J4aK6Wcs4A3Ex4HXdfjJ6ZsHpNZfjaS4B9w7xqEnmFEYMqQd13")
```
Get address from private key:
```go
address := GetAddressFromPrivateKey(privateKey, WAYKI_MAINTNET)
```
Get address from mnemonic:
```go
mnemonic := "empty regular curve turtle student prize toy accuse develop spike scatter ginger"
address := GetAddressFromMnemonic(mnemonic, WAYKI_MAINTNET)
```
### WaykiChain sign transaction
Signing a transaction with a private key,you can submit your offline signature rawtx transaction via bass.

|  BassNetwork |  ApiAddr | 
|-------------- |----------------------------------|
|   TestNetwork | https://baas-test.wiccdev.org/v2/api/swagger-ui.html#!/  |  
|   ProdNetwork | https://baas.wiccdev.org/v2/api/swagger-ui.html#!/       |                                |

Submit raw string:
Mainnet <https://baas.wiccdev.org/v2/api/swagger-ui.html#!/transaction-controller/offlinTransactionUsingPOST> ,
TestNet <https://baas-test.wiccdev.org/v2/api/swagger-ui.html#!/transaction-controller/offlinTransactionUsingPOST>,

Get block height:
MainNet<https://baas.wiccdev.org/v2/api/swagger-ui.html#!/block-controller/getBlockCountUsingPOST>,
TestNet <https://baas-test.wiccdev.org/v2/api/swagger-ui.html#!/block-controller/getBlockCountUsingPOST>

#### common transaction
```go
func SignRegisterAccountTx(privateKey string, param *RegisterAccountTxParam) (string, error)

func SignCommonTx(privateKey string, param *CommonTxParam) (string, error)

func SignDelegateTx(privateKey string, param *DelegateTxParam) (string, error)

func SignCallContractTx(privateKey string, param *CallContractTxParam) (string, error)

func SignRegisterContractTx(privateKey string, param *RegisterContractTxParam) (string, error)
```
- SignRegisterAccountTx.sign registration transaction with a private key , return the rawtx string.
- SignCommonTx.sign transfer transaction with a private key , return the rawtx string.
- SignDelegateTx.sign delegate transaction with a private key , return the rawtx string.
- SignCallContractTx.sign invoke contract transaction with a private key , return the rawtx string.
- SignRegisterContractTx.sign deploy contract transaction with a private key , return the rawtx string.

Example:

Sign register account transaction:
```go
	privateKey := "YAa1wFCfFnZ5bt4hg9MDeDevTMd1Nu874Mn83hEXwtfAL2vkQE9f"
	var txParam RegisterAccountTxParam  
	txParam.ValidHeight = 630314 //WaykiChain block height 
	txParam.Fees = 10000         //Miner fee,minimum 1000sawi

	hash, err := SignRegisterAccountTx(privateKey, &txParam)
	if err != nil {
		t.Error("SignRegisterAccountTx err: ", err)
	}
```
Sign common transfer transaction:
```go
    privateKey := "YAa1wFCfFnZ5bt4hg9MDeDevTMd1Nu874Mn83hEXwtfAL2vkQE9f" 
	var txParams CommonTxParam
	txParams.ValidHeight = 630314
	txParams.SrcRegId = "158-1"                                //user regid
	txParams.DestAddr = "wSSbTePArv6BkDsQW9gpGCTX55AXVxVKbd"  //dest address
	txParams.Values = 10000                                   //transfer amount
	txParams.Fees = 10000
	txParams.PubKey = "03e93e7d870ce6f1c9997076c56fc24e6381c612662cd9a5a59294fac9ba7d21d7"  //wallet public key hex string
    txParams.Memo="test transfer"                                                           //transfer memo   
	hash, err := SignCommonTx(privateKey, &txParams)
	if err != nil {
		t.Error("SignCommonTx err: ", err)
	}
```
Sign Delegate transaction:
```go
    privateKey := "Y9sx4Y8sBAbWDAqAWytYuUnJige3ZPwKDZp1SCDqqRby1YMgRG9c"
	var txParams DelegateTxParam
	txParams.ValidHeight = 95728
	txParams.SrcRegId = "0-1"
	txParams.Fees = 10000
	txParams.PubKey = "03e93e7d870ce6f1c9997076c56fc24e6381c612662cd9a5a59294fac9ba7d21d7"
	txParams.Votes = NewOperVoteFunds()
	pubKey, _ := hex.DecodeString("025a37cb6ec9f63bb17e562865e006f0bafa9afbd8a846bd87fc8ff9e35db1252e") //Voted public key
	vote := OperVoteFund{PubKey: pubKey, VoteValue: 10000}
	txParams.Votes.Add(&vote)

	hash, err := SignDelegateTx(privateKey, &txParams)
	if err != nil {
		t.Error("SignDelegateTx err: ", err)
	}
```
Sign invoke contract transaction:
```go
	privateKey := "Y9sx4Y8sBAbWDAqAWytYuUnJige3ZPwKDZp1SCDqqRby1YMgRG9c"
	var txParam CallContractTxParam
	txParam.ValidHeight = 22365
	txParam.SrcRegId = "0-1"
	txParam.AppId = "20988-1"          //contract regid
	txParam.Fees = 100000
	txParam.Values = 10000
	txParam.ContractHex = "f017"      //call contract method
	txParam.PubKey = "03e93e7d870ce6f1c9997076c56fc24e6381c612662cd9a5a59294fac9ba7d21d7"
	hash, err := SignCallContractTx(privateKey, &txParam)
	if err != nil {
		t.Error("SignCallContractTx err: ", err)
	}
```
Sign deploy contract Transaction:
```go
	privateKey := "YAa1wFCfFnZ5bt4hg9MDeDevTMd1Nu874Mn83hEXwtfAL2vkQE9f"
	script, err := ioutil.ReadFile("./demo/data/hello.lua")
	if err != nil {
		t.Error("Read contract script file err: ", err)
	}
	var txParam RegisterContractTxParam
	txParam.ValidHeight = 630314 
	txParam.SrcRegId = "0-1"     
	txParam.Fees = 110000000
	txParam.Script = script                        //contract bytearray
	txParam.Description = "My hello contract!!!"  //contract description
	hash, err := SignRegisterContractTx(privateKey, &txParam)
	if err != nil {
		t.Error("SignRegisterContractTx err: ", err)
	}
```

#### cdp transaction
Any user holding a WICC can send a WICC to the CDP (Collaterized Debt Position) to obtain a certain percentage of WUSD.a user can only have one cdp unless the previous cdp has been destroyed.

```go
func SignCdpStakeTx(privateKey string, param *CdpStakeTxParam) (string, error) 

func SignCdpRedeemTx(privateKey string, param *CdpRedeemTxParam) (string, error) 

func SignCdpLiquidateTx(privateKey string, param *CdpLiquidateTxParam) (string, error)
```
- SignCdpStakeTx.sign cdp stake transaction with a private key , return the rawtx string.
- SignCdpRedeemTx.sign cdp redeem transaction with a private key , return the rawtx string.
- SignCdpLiquidateTx.sign cdp liquidate transaction with a private key , return the rawtx string.

Example:

Sign cdp stake transaction:
```go
    privateKey := "Y6J4aK6Wcs4A3Ex4HXdfjJ6ZsHpNZfjaS4B9w7xqEnmFEYMqQd13"
	var txParam CdpStakeTxParam
	txParam.CdpTxid = "0b9734e5db3cfa38e76bb273dba4f65a210cc76ca2cf739f3c131d0b24ff89c1"  //user cdp transaction hash,If the user has not created a cdp, then do not fill out
	txParam.BcoinSymbol = string(commons.WICC)   //pay WICC
	txParam.ScoinSymbol = string(commons.WUSD)   //get WUSD
	txParam.FeeSymbol = string(commons.WICC)    //fee symbol (WICC/WUSD)
	txParam.BcoinStake = 100000000
	txParam.ScoinMint = 50000000
	txParam.Fees = 100000
	txParam.ValidHeight = 283308
	txParam.SrcRegId = "0-1"
	txParam.PubKey = "03e93e7d870ce6f1c9997076c56fc24e6381c612662cd9a5a59294fac9ba7d21d7"

	hash, err := SignCdpStakeTx(privateKey, &txParam)
	if err != nil {
		t.Error("SignCdpStakeTx err: ", err)
	}
```
Sign cdp redeem transaction:
```go
    privateKey := "Y6J4aK6Wcs4A3Ex4HXdfjJ6ZsHpNZfjaS4B9w7xqEnmFEYMqQd13"
	var txParam CdpRedeemTxParam
	txParam.CdpTxid = "0b9734e5db3cfa38e76bb273dba4f65a210cc76ca2cf739f3c131d0b24ff89c1"//user cdp create transaction hash
	txParam.FeeSymbol = string(commons.WICC)
	txParam.ScoinsToRepay = 20000000
	txParam.BcoinsToRedeem = 100000000
	txParam.Fees = 100000
	txParam.ValidHeight = 25
	txParam.SrcRegId = "0-1"
	txParam.PubKey = "03e93e7d870ce6f1c9997076c56fc24e6381c612662cd9a5a59294fac9ba7d21d7"

	hash, err := SignCdpRedeemTx(privateKey, &txParam)
	if err != nil {
		t.Error("SignCdpStakeTx err: ", err)
	}
```
Sign cdp liquidate transaction:
```go
    privateKey := "Y6J4aK6Wcs4A3Ex4HXdfjJ6ZsHpNZfjaS4B9w7xqEnmFEYMqQd13"
	var txParam CdpLiquidateTxParam
	txParam.CdpTxid = "0b9734e5db3cfa38e76bb273dba4f65a210cc76ca2cf739f3c131d0b24ff89c1"//Liquidated cdp transaction hash id
	txParam.FeeSymbol = string(commons.WICC)
	txParam.ScoinsLiquidate = 100000000
	txParam.Fees = 100000
	txParam.ValidHeight = 25
	txParam.SrcRegId = "0-1"
	txParam.PubKey = "03e93e7d870ce6f1c9997076c56fc24e6381c612662cd9a5a59294fac9ba7d21d7"

	hash, err := SignCdpLiquidateTx(privateKey, &txParam)
	if err != nil {
		t.Error("SignCdpStakeTx err: ", err)
	}
```

#### dex transaction
WaykiChain decentralized exchange.









































