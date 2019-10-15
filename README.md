# wicc-wallet-utils-go  
维基链go语言离线签名钱包库 
(WaykiChain Wallet Library for Offline Operation written in golang)

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
- **GenerateMnemonics.**  
生成12个助记词  
(You will get 12 words).  
- **GetPrivateKeyFromMnemonic.**  
 你提供你的助记词和网络类型（1 主网，2 测试网），函数会给你返回私钥，主网私钥大写P开头，测试网大写Y开头。  
（You should provide your mnemonic and network Type (1 MainNet,2 TestNet),function return private Key,MainNet Private key start with "P" ,TestNet private key start with "Y".）  
- **GetPubKeyFromPrivateKey.**   
（提供私钥获得公钥，获得公钥的16进制字符串）  
you should provide your Private Key,the function return wallet public key as hex string.  
- **GetAddressFromPrivateKey.**  
私钥获得钱包地址，地址是Base58编码的字符串，主网地址大写字母W开头，测试网地址小写字母w开头。  
you should provide your Private Key,the function return wallet Address as base58 encode string,MainNet Address start with "W",TestNet Address start with "w".  
- **GetAddressFromMnemonic.**  
从助记词获得地址。  
you should provide your mnemonic,the function return wallet Address as base58 encode string,MainNet Address start with "W",TestNet Address start with "w".

示例(Examples):

[生成助记词(Generate mnemonic:)](https://github.com/WaykiChain/wicc-wallet-utils-go/blob/master/WaykichainWallet_test.go)
```go
mnemonic := GenerateMnemonics()
```
[助记词生成钱包私钥(Get private key from mnemonic:)](https://github.com/WaykiChain/wicc-wallet-utils-go/blob/master/WaykichainWallet_test.go)
```go
mnemonic := "empty regular curve turtle student prize toy accuse develop spike scatter ginger"
privateKey := GetPrivateKeyFromMnemonic(mnemonic, WAYKI_MAINTNET)
```
[私钥获得公钥(Get public key from private key:)](https://github.com/WaykiChain/wicc-wallet-utils-go/blob/master/WaykichainWallet_test.go)
```go
publicKey,_:=GetPubKeyFromPrivateKey("Y6J4aK6Wcs4A3Ex4HXdfjJ6ZsHpNZfjaS4B9w7xqEnmFEYMqQd13")
```
[私钥获得钱包地址(Get address from private key:)](https://github.com/WaykiChain/wicc-wallet-utils-go/blob/master/WaykichainWallet_test.go)
```go
address := GetAddressFromPrivateKey(privateKey, WAYKI_MAINTNET)
```
[助记词获得钱包地址(Get address from mnemonic:)](https://github.com/WaykiChain/wicc-wallet-utils-go/blob/master/WaykichainWallet_test.go)
```go
mnemonic := "empty regular curve turtle student prize toy accuse develop spike scatter ginger"
address := GetAddressFromMnemonic(mnemonic, WAYKI_MAINTNET)
```
### 维基链签名交易(WaykiChain Sign Transaction)
用私钥签名交易，你可以使用bass提交钱包库生成的rawtx字符串。  
Signing a transaction with a private key,you can submit your offline signature rawtx transaction by bass.

|  BassNetwork |  Api | 
|-------------- |----------------------------------|
|   TestNet | https://baas-test.wiccdev.org/v2/api/swagger-ui.html#!/  |  
|   MainNet | https://baas.wiccdev.org/v2/api/swagger-ui.html#!/       |                                |

Submit raw string:

Mainnet <https://baas.wiccdev.org/v2/api/swagger-ui.html#!/transaction-controller/offlinTransactionUsingPOST> ,

TestNet <https://baas-test.wiccdev.org/v2/api/swagger-ui.html#!/transaction-controller/offlinTransactionUsingPOST>,

Get block height:

MainNet<https://baas.wiccdev.org/v2/api/swagger-ui.html#!/block-controller/getBlockCountUsingPOST>,

TestNet <https://baas-test.wiccdev.org/v2/api/swagger-ui.html#!/block-controller/getBlockCountUsingPOST>

#### 转账交易

#### 合约交易(部署、调用)

#### CDP交易(CDP Transaction)
用户可以通过抵押WICC获得WUSD,一个用户只能拥有一个cdp，除非之前的cdp已经关闭。  
Any user holding a WICC can send a WICC to the CDP (Collaterized Debt Position) to obtain a certain percentage of WUSD.a user can only have one cdp unless the previous cdp has been destroyed.

```go
func SignCdpStakeTx(privateKey string, param *CdpStakeTxParam) (string, error) 

func SignCdpRedeemTx(privateKey string, param *CdpRedeemTxParam) (string, error) 

func SignCdpLiquidateTx(privateKey string, param *CdpLiquidateTxParam) (string, error)
```
- **SignCdpStakeTx.**  
CDP抵押交易签名。  
sign cdp stake transaction with a private key , return the rawtx string.
- **SignCdpRedeemTx.**  
cdp赎回交易签名。  
sign cdp redeem transaction with a private key , return the rawtx string.
- **SignCdpLiquidateTx.**  
CDP清算交易签名。  
sign cdp liquidate transaction with a private key , return the rawtx string.

示例(Example:)

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
维基链去中心化交易所。  
WaykiChain decentralized exchange.

```go
func SignDexSellLimitTx(privateKey string, param *DexLimitTxParam) (string, error)

func SignDexMarketSellTx(privateKey string, param *DexMarketTxParam) (string, error)

func SignDexBuyLimitTx(privateKey string, param *DexLimitTxParam) (string, error)

func SignDexMarketBuyTx(privateKey string, param *DexMarketTxParam) (string, error)

func SignDexCancelTx(privateKey string, param *DexCancelTxParam) (string, error)
```
- **SignDexSellLimitTx.**  
限价卖单。 
sign dex sell limit price transaction with a private key , return the rawtx string.
- **SignDexMarketSellTx.**  
市价卖单。  
sign dex sell market price transaction with a private key , return the rawtx string.
- **SignDexBuyLimitTx.**  
限价买单。  
sign dex buy limit price transaction with a private key , return the rawtx string.
- **SignDexMarketBuyTx.**  
市价买单。  
sign dex buy market price transaction with a private key , return the rawtx string.
- **SignDexCancelTx.**  
取消挂单。  
sign cancel dex order transaction with a private key , return the rawtx string.

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







































