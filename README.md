# wicc-wallet-utils-go
（维基go语言离线签名钱包库）

WaykiChain Wallet Library for Offline Operation written in golang

## 下载(Install)

* go get github.com/WaykiChain/wicc-wallet-utils-go
* go get github.com/btcsuite

## 用法（Usage）

### 创建钱包（WaykiChain Create Wallet）
生成助记词和私钥管理你的钱包。

Generate mnemonics and private keys to manage your wallet.

```go
func GenerateMnemonics() string

func GetPrivateKeyFromMnemonic(words string, netType int) string

func GetPubKeyFromPrivateKey(privKey string) (string,error) 

func GetAddressFromPrivateKey(privateKey string, netType int) string

func GetAddressFromMnemonic(words string, netType int) string
```
- GenerateMnemonics. 

  生成12个助记词
  
  (You will get 12 words).

- GetPrivateKeyFromMnemonic. 

  你提供你的助记词和网络类型（1 主网，2 测试网），函数会给你返回私钥，主网私钥大写P开头，测试网大写Y开头。

 （You should provide your mnemonic and network Type (1 MainNet,2 TestNet),function return private Key,MainNet Private key start with "P" ,TestNet private key start with "Y".）

- GetPubKeyFromPrivateKey. 

 （提供私钥获得公钥，获得公钥的16进制字符串）
   
  you should provide your Private Key,the function return wallet public key as hex string.

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

func SignUCoinTransferTx(privateKey string, param *UCoinTransferTxParam) (string, error) 
```
- SignRegisterAccountTx.sign registration transaction with a private key , return the rawtx string.
- SignCommonTx.sign transfer transaction with a private key , return the rawtx string.
- SignDelegateTx.sign delegate transaction with a private key , return the rawtx string.
- SignCallContractTx.sign invoke contract transaction with a private key , return the rawtx string.
- SignRegisterContractTx.sign deploy contract transaction with a private key , return the rawtx string.
- SignUCoinTransferTx.sign Multi-coin transfer transaction with a private key , return the rawtx string.

Example:

**The register transaction is not required, you can activate wallet by public key in other transactions**

Sign register account transaction:
```go
	privateKey := "YAa1wFCfFnZ5bt4hg9MDeDevTMd1Nu874Mn83hEXwtfAL2vkQE9f"
	var txParam RegisterAccountTxParam  
	txParam.ValidHeight = 630314 //WaykiChain block height 
	txParam.Fees = 10000         //Miner fee,minimum 10000sawi

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
	txParam.Fees = 110000000                       //Miner fee,minimum 11000000sawi 
	txParam.Script = script                        //contract bytearray
	txParam.Description = "My hello contract!!!"  //contract description
	hash, err := SignRegisterContractTx(privateKey, &txParam)
	if err != nil {
		t.Error("SignRegisterContractTx err: ", err)
	}
```
Sign Multi-coin transfer transaction:
```go
    privateKey := "Y6J4aK6Wcs4A3Ex4HXdfjJ6ZsHpNZfjaS4B9w7xqEnmFEYMqQd13"
	var txParam UCoinTransferTxParam
	txParam.FeeSymbol = string(commons.WICC)
	txParam.CoinSymbol = string(commons.WICC)
	txParam.CoinAmount = 1000000
	txParam.Fees = 10000
	txParam.ValidHeight = 297449
	txParam.SrcRegId = "0-1"
	txParam.DestAddr = "wNDue1jHcgRSioSDL4o1AzXz3D72gCMkP6"
	txParam.PubKey = "036c5397f3227a1e209952829d249b7ad0f615e43b763ac15e3a6f52627a10df21"
	txParam.Memo = ""
	hash, err := SignUCoinTransferTx(privateKey, &txParam)
	if err != nil {
		t.Error("SignCdpStakeTx err: ", err)
	}
```

#### CDP Transaction
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
	txParam.Fees = 100000                       //Miner fee,minimum 100000sawi
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
	txParam.Fees = 100000                                     //Miner fee,minimum 100000sawi
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
	txParam.Fees = 100000                                                    //Miner fee,minimum 100000sawi
	txParam.ValidHeight = 25
	txParam.SrcRegId = "0-1"
	txParam.PubKey = "03e93e7d870ce6f1c9997076c56fc24e6381c612662cd9a5a59294fac9ba7d21d7"

	hash, err := SignCdpLiquidateTx(privateKey, &txParam)
	if err != nil {
		t.Error("SignCdpStakeTx err: ", err)
	}
```

#### DEX Transaction
WaykiChain decentralized exchange.

```go
func SignDexSellLimitTx(privateKey string, param *DexLimitTxParam) (string, error)

func SignDexMarketSellTx(privateKey string, param *DexMarketTxParam) (string, error)

func SignDexBuyLimitTx(privateKey string, param *DexLimitTxParam) (string, error)

func SignDexMarketBuyTx(privateKey string, param *DexMarketTxParam) (string, error)

func SignDexCancelTx(privateKey string, param *DexCancelTxParam) (string, error)
```
- SignDexSellLimitTx.sign dex sell limit price transaction with a private key , return the rawtx string.
- SignDexMarketSellTx.sign dex sell market price transaction with a private key , return the rawtx string.
- SignDexBuyLimitTx.sign dex buy limit price transaction with a private key , return the rawtx string.
- SignDexMarketBuyTx.sign dex buy market price transaction with a private key , return the rawtx string.
- SignDexCancelTx.sign cancel dex order transaction with a private key , return the rawtx string.

Sign dex Buy limit price transaction:
```go
	privateKey := "Y6J4aK6Wcs4A3Ex4HXdfjJ6ZsHpNZfjaS4B9w7xqEnmFEYMqQd13"
	var txParam DexLimitTxParam
	txParam.FeeSymbol = string(commons.WICC)
	txParam.Fees = 100000
	txParam.CoinSymbol = string(commons.WUSD)
	txParam.AssetSymbol = string(commons.WICC)
	txParam.AssetAmount = 10000
	txParam.ValidHeight = 25
	txParam.Price = 25                      
	txParam.SrcRegId = "0-1"
	txParam.PubKey = "03e93e7d870ce6f1c9997076c56fc24e6381c612662cd9a5a59294fac9ba7d21d7"

	hash, err := SignDexBuyLimitTx(privateKey, &txParam)
	if err != nil {
		t.Error("SignCdpStakeTx err: ", err)
	}
```
Sign dex sell limit price transaction:
```go
    privateKey := "Y6J4aK6Wcs4A3Ex4HXdfjJ6ZsHpNZfjaS4B9w7xqEnmFEYMqQd13"
	var txParam DexLimitTxParam
	txParam.FeeSymbol = string(commons.WICC)
	txParam.Fees = 10000
	txParam.CoinSymbol = string(commons.WUSD)
	txParam.AssetSymbol = string(commons.WICC)
	txParam.AssetAmount = 1000000
	txParam.ValidHeight = 282956
	txParam.Price = 200000000
	txParam.SrcRegId = "0-1"
	txParam.PubKey = "03e93e7d870ce6f1c9997076c56fc24e6381c612662cd9a5a59294fac9ba7d21d7"

	hash, err := SignDexSellLimitTx(privateKey, &txParam)
	if err != nil {
		t.Error("SignCdpStakeTx err: ", err)
	}
```
Sign dex sell market price transaction:
```go
    privateKey := "Y6J4aK6Wcs4A3Ex4HXdfjJ6ZsHpNZfjaS4B9w7xqEnmFEYMqQd13"
	var txParam DexMarketTxParam
	txParam.FeeSymbol = string(commons.WICC)
	txParam.Fees = 100000
	txParam.CoinSymbol = string(commons.WUSD)
	txParam.AssetSymbol = string(commons.WICC)
	txParam.AssetAmount = 10000
	txParam.ValidHeight = 25
	txParam.SrcRegId = "0-1"
	txParam.PubKey = "03e93e7d870ce6f1c9997076c56fc24e6381c612662cd9a5a59294fac9ba7d21d7"

	hash, err := SignDexMarketSellTx(privateKey, &txParam)
	if err != nil {
		t.Error("SignCdpStakeTx err: ", err)
	}
```

Sign dex buy market price transaction:
```go
    privateKey := "Y6J4aK6Wcs4A3Ex4HXdfjJ6ZsHpNZfjaS4B9w7xqEnmFEYMqQd13"
	var txParam DexMarketTxParam
	txParam.FeeSymbol = string(commons.WICC)
	txParam.Fees = 100000
	txParam.CoinSymbol = string(commons.WUSD)
	txParam.AssetSymbol = string(commons.WICC)
	txParam.AssetAmount = 10000
	txParam.ValidHeight = 25
	txParam.SrcRegId = "0-1"
	txParam.PubKey = "03e93e7d870ce6f1c9997076c56fc24e6381c612662cd9a5a59294fac9ba7d21d7"

	hash, err := SignDexMarketBuyTx(privateKey, &txParam)
	if err != nil {
		t.Error("SignCdpStakeTx err: ", err)
	}
```
Sign dex cancel order transaction:
```go
	privateKey := "Y6J4aK6Wcs4A3Ex4HXdfjJ6ZsHpNZfjaS4B9w7xqEnmFEYMqQd13"
	var txParam DexCancelTxParam
	txParam.FeeSymbol = string(commons.WICC)
	txParam.Fees = 100000
	txParam.DexTxid = "009c0e665acdd9e8ae754f9a51337b85bb8996980a93d6175b61edccd3cdc144"//dex transaction tx id
	txParam.ValidHeight = 25
	txParam.SrcRegId = "0-1"
	txParam.PubKey = "03e93e7d870ce6f1c9997076c56fc24e6381c612662cd9a5a59294fac9ba7d21d7"

	hash, err := SignDexCancelTx(privateKey, &txParam)
	if err != nil {
		t.Error("SignCdpStakeTx err: ", err)
	}
```







































