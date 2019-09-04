package wiccwallet

import (
	"testing"
	"encoding/hex"
	"io/ioutil"
	"github.com/WaykiChain/wicc-wallet-utils-go/commons"
)

/*
    * 账户注册交易,新版本已基本废弃，可改用公钥注册，免注册费用
    * Account registration transaction, the new version has been abandoned, you can use public key registration, free registration fee
    * */
func TestSignRegisterAccountTx(t *testing.T) {

	privateKey := "YAa1wFCfFnZ5bt4hg9MDeDevTMd1Nu874Mn83hEXwtfAL2vkQE9f" //"Y9sx4Y8sBAbWDAqAWytYuUnJige3ZPwKDZp1SCDqqRby1YMgRG9c"
	var txParam RegisterAccountTxParam
	txParam.ValidHeight = 630314
	txParam.Fees = 10000

	hash, err := SignRegisterAccountTx(privateKey, &txParam)
	if err != nil {
		t.Error("SignRegisterAccountTx err: ", err)
	}
	println(hash)
}

/*
  * 转账交易
  * common transfer
  * fee Minimum 0.0001 wicc
  * */
func TestSignCommonTx(t *testing.T) {

	privateKey := "YAa1wFCfFnZ5bt4hg9MDeDevTMd1Nu874Mn83hEXwtfAL2vkQE9f" //"Y7V1jwCRr8D3tyPTkcsjgBTHwZN45b1U3ueZfJ5oWVJqwcKpArou"
	var txParams CommonTxParam
	txParams.ValidHeight = 630314
	txParams.SrcRegId = "158-1"
	txParams.DestAddr = "wSSbTePArv6BkDsQW9gpGCTX55AXVxVKbd"
	txParams.Values = 10000
	txParams.Fees = 10000
	txParams.PubKey = "03e93e7d870ce6f1c9997076c56fc24e6381c612662cd9a5a59294fac9ba7d21d7"
    txParams.Memo="test transfer"
	hash, err := SignCommonTx(privateKey, &txParams)
	if err != nil {
		t.Error("SignCommonTx err: ", err)
	}
	println(hash)
}

/*
   * 合约调用交易
   * Contract transaction sample
   * fee Minimum 0.0001 wicc
   * */
func TestSignCallContractTx(t *testing.T) {

	privateKey := "Y9sx4Y8sBAbWDAqAWytYuUnJige3ZPwKDZp1SCDqqRby1YMgRG9c"

	var txParam CallContractTxParam
	txParam.ValidHeight = 22365
	txParam.SrcRegId = "0-1"
	txParam.AppId = "20988-1"
	txParam.Fees = 100000
	txParam.Values = 10000
	txParam.ContractHex = "f017"
	txParam.PubKey = "03e93e7d870ce6f1c9997076c56fc24e6381c612662cd9a5a59294fac9ba7d21d7"
	hash, err := SignCallContractTx(privateKey, &txParam)
	if err != nil {
		t.Error("SignCallContractTx err: ", err)
	}
	println(hash)
}

/*
* 投票交易
* Voting transaction
* fee Minimum 0.0001 wicc
* */
func TestSignDelegateTx(t *testing.T) {
	privateKey := "Y9sx4Y8sBAbWDAqAWytYuUnJige3ZPwKDZp1SCDqqRby1YMgRG9c"
	var txParams DelegateTxParam
	txParams.ValidHeight = 95728
	txParams.SrcRegId = ""
	txParams.Fees = 10000
	txParams.PubKey = "03e93e7d870ce6f1c9997076c56fc24e6381c612662cd9a5a59294fac9ba7d21d7"
	txParams.Votes = NewOperVoteFunds()
	pubKey, _ := hex.DecodeString("025a37cb6ec9f63bb17e562865e006f0bafa9afbd8a846bd87fc8ff9e35db1252e")
	vote := OperVoteFund{PubKey: pubKey, VoteValue: 10000}
	txParams.Votes.Add(&vote)

	hash, err := SignDelegateTx(privateKey, &txParams)
	if err != nil {
		t.Error("SignDelegateTx err: ", err)
	}
	println(hash)
}

func TestSignRegisterContractTx(t *testing.T) {

	privateKey := "YAa1wFCfFnZ5bt4hg9MDeDevTMd1Nu874Mn83hEXwtfAL2vkQE9f" //"Y9sx4Y8sBAbWDAqAWytYuUnJige3ZPwKDZp1SCDqqRby1YMgRG9c"

	script, err := ioutil.ReadFile("./demo/data/hello.lua")
	if err != nil {
		t.Error("Read contract script file err: ", err)
	}

	var txParam RegisterContractTxParam
	txParam.ValidHeight = 630314 //20999
	txParam.SrcRegId = "0-1"     //"7849-1"
	txParam.Fees = 110000000
	txParam.Script = script
	txParam.Description = "My hello contract!!!"
	hash, err := SignRegisterContractTx(privateKey, &txParam)
	if err != nil {
		t.Error("SignRegisterContractTx err: ", err)
	}
	println(hash)
}

/*
  * 多币种转账交易 ,支持多种币种转账
  * Test nUniversal Coin Transfer Tx
  * fee Minimum 0.0001 wicc
  * */
func TestSignUCoinTransferTx(t *testing.T) {
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
	println(hash)
}

/*
    * 创建,追加cdp交易
    * Create or append an  cdp transaction
    * fee Minimum 0.001 wicc
    * */
func TestCdpStakeTx(t *testing.T) {
	privateKey := "Y6J4aK6Wcs4A3Ex4HXdfjJ6ZsHpNZfjaS4B9w7xqEnmFEYMqQd13"
	var txParam CdpStakeTxParam
	txParam.CdpTxid = "0b9734e5db3cfa38e76bb273dba4f65a210cc76ca2cf739f3c131d0b24ff89c1"
	txParam.BcoinSymbol = string(commons.WICC)
	txParam.ScoinSymbol = string(commons.WUSD)
	txParam.FeeSymbol = string(commons.WICC)
	txParam.BcoinStake = 10
	txParam.ScoinMint = 10
	txParam.Fees = 100000
	txParam.ValidHeight = 283308
	txParam.SrcRegId = "0-1"
	txParam.PubKey = "03e93e7d870ce6f1c9997076c56fc24e6381c612662cd9a5a59294fac9ba7d21d7"

	hash, err := SignCdpStakeTx(privateKey, &txParam)
	if err != nil {
		t.Error("SignCdpStakeTx err: ", err)
	}
	println(hash)
}

/*
  * 赎回cdp交易
  * Redeem cdp transaction
  * fee Minimum 0.001 wicc
  * */
func TestSignCdpRedeemTx(t *testing.T) {
	privateKey := "Y6J4aK6Wcs4A3Ex4HXdfjJ6ZsHpNZfjaS4B9w7xqEnmFEYMqQd13"
	var txParam CdpRedeemTxParam
	txParam.CdpTxid = "0b9734e5db3cfa38e76bb273dba4f65a210cc76ca2cf739f3c131d0b24ff89c1"
	txParam.FeeSymbol = string(commons.WICC)
	txParam.ScoinsToRepay = 0
	txParam.BcoinsToRedeem = 10
	txParam.Fees = 100000
	txParam.ValidHeight = 25
	txParam.SrcRegId = "0-1"
	txParam.PubKey = "03e93e7d870ce6f1c9997076c56fc24e6381c612662cd9a5a59294fac9ba7d21d7"

	hash, err := SignCdpRedeemTx(privateKey, &txParam)
	if err != nil {
		t.Error("SignCdpStakeTx err: ", err)
	}
	println(hash)
}

/*
  * 清算cdp交易
  * Liquidate cdp transaction
  * fee Minimum 0.001 wicc
  * */
func TestSignCdpLiquidateTx(t *testing.T) {
	privateKey := "Y6J4aK6Wcs4A3Ex4HXdfjJ6ZsHpNZfjaS4B9w7xqEnmFEYMqQd13"
	var txParam CdpLiquidateTxParam
	txParam.CdpTxid = "0b9734e5db3cfa38e76bb273dba4f65a210cc76ca2cf739f3c131d0b24ff89c1"
	txParam.FeeSymbol = string(commons.WICC)
	txParam.ScoinsLiquidate = 100
	txParam.Fees = 100000
	txParam.ValidHeight = 25
	txParam.SrcRegId = "0-1"
	txParam.PubKey = "03e93e7d870ce6f1c9997076c56fc24e6381c612662cd9a5a59294fac9ba7d21d7"

	hash, err := SignCdpLiquidateTx(privateKey, &txParam)
	if err != nil {
		t.Error("SignCdpStakeTx err: ", err)
	}
	println(hash)
}

/*
 * Dex 限价买单交易
 * Dex limit price transaction
 * fee Minimum 0.0001 wicc
 * */
func TestSignDexBuyLimitTx(t *testing.T) {
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
	println(hash)
}

/*
   * Dex 限价卖单交易
   * Dex limit sell price transaction
   * fee Minimum 0.0001 wicc
  * */
func TestSignDexSellLimitTx(t *testing.T) {
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
	//txParam.PubKey = "03e93e7d870ce6f1c9997076c56fc24e6381c612662cd9a5a59294fac9ba7d21d7"

	hash, err := SignDexSellLimitTx(privateKey, &txParam)
	if err != nil {
		t.Error("SignCdpStakeTx err: ", err)
	}
	println(hash)
}

/*
     *  Dex 市价卖单交易
     * Dex market sell price transaction
     * fee Minimum 0.0001 wicc
    * */
func TestSignDexMarketSellTx(t *testing.T) {
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
	println(hash)
}

/*
     *  Dex 市价买单交易
     * Dex market buy price transaction
     * fee Minimum 0.0001 wicc
    * */
func TestSignDexMarketBuyTx(t *testing.T) {
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
	println(hash)
}

/*
 *  Dex 取消挂单交易
 * Dex cancel order tx
 * fee Minimum 0.0001 wicc
* */
func TestSignDexCancelTx(t *testing.T) {
	privateKey := "Y6J4aK6Wcs4A3Ex4HXdfjJ6ZsHpNZfjaS4B9w7xqEnmFEYMqQd13"
	var txParam DexCancelTxParam
	txParam.FeeSymbol = string(commons.WICC)
	txParam.Fees = 100000
	txParam.DexTxid = "009c0e665acdd9e8ae754f9a51337b85bb8996980a93d6175b61edccd3cdc144"
	txParam.ValidHeight = 25
	txParam.SrcRegId = "0-1"
	txParam.PubKey = "03e93e7d870ce6f1c9997076c56fc24e6381c612662cd9a5a59294fac9ba7d21d7"

	hash, err := SignDexCancelTx(privateKey, &txParam)
	if err != nil {
		t.Error("SignCdpStakeTx err: ", err)
	}
	println(hash)
}
