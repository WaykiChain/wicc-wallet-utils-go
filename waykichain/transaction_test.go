package waykichain

import (
	"encoding/hex"
	"encoding/json"
	"github.com/WaykiChain/wicc-wallet-utils-go/common"
	"io/ioutil"
	"testing"
)

/*
  * 多币种转账交易 ,支持多种币种转账
  * Test nUniversal Coin Transfer Tx
  * fee Minimum 100000 sawi
  * */
func TestCreateUCoinTransferRawTx(t *testing.T) {
	privateKey := "Y6amwxjHqUM37UrquokPsbCXTRNughoM27gDUGfbXhJikS39i9h1"
	var txParam UCoinTransferTxParam
	txParam.FeeSymbol = common.WICCSYM
	txParam.Fees = 1000000
	txParam.ValidHeight = 2644318
	txParam.SrcRegId = ""
	txParam.Dests=NewDestArr()
	dest:=Dest{common.WICCSYM,1000000, "wPYU7FA2Y3WQ5TCAA84aaAmbYDMo8CHV2T"}
	txParam.Dests.Add(&dest)
	txParam.PubKey = "031b27286c65b81ac13cfd4067b030398a19eb147f439c094fbb19a2f3ab9ec10b"
	txParam.Memo = ""
	rawTx,txid, err := txParam.CreateRawTx(privateKey)
	if err != nil {
		t.Error("UCoinTransferTx err: ", err)
	}
	t.Log("rawTx=",rawTx)
	t.Log("txid=",txid)
}

/*
  * 多币种合约调用交易 ,支持多种币种转账
  * Test nUniversal Coin Contract Tx
  * fee Minimum 1000000 sawi
  * */
func TestCreateUCoinCallContractRawTx(t *testing.T) {
	privateKey := "Y6J4aK6Wcs4A3Ex4HXdfjJ6ZsHpNZfjaS4B9w7xqEnmFEYMqQd13"
	var txParam UCoinContractTxParam
	txParam.FeeSymbol = common.WICCSYM
	txParam.CoinSymbol = common.WICCSYM
	txParam.CoinAmount = 1000000
	txParam.Fees = 1000000
	txParam.ValidHeight = 297449
	txParam.SrcRegId = "0-1"
	txParam.AppId = "0-1"
	txParam.PubKey = "036c5397f3227a1e209952829d249b7ad0f615e43b763ac15e3a6f52627a10df21"
	txParam.ContractHex = "f017"
	rawTx,txid, err := txParam.CreateRawTx(privateKey)
	if err != nil {
		t.Error("SignUCoinCallContractTx err: ", err)
	}
	t.Log("rawTx=",rawTx)
	t.Log("txid=",txid)
}


/*
	多币种部署合约,推荐使用
	fee Minimum: 110000000 sawi
*/
func TestCreateUCoinDeployContractTx(t *testing.T) {

	privateKey := "YAa1wFCfFnZ5bt4hg9MDeDevTMd1Nu874Mn83hEXwtfAL2vkQE9f" //"Y9sx4Y8sBAbWDAqAWytYuUnJige3ZPwKDZp1SCDqqRby1YMgRG9c"

	script, err := ioutil.ReadFile("./../demo/data/hello.lua")
	if err != nil {
		t.Error("Read contract script file err: ", err)
	}

	var txParam UCoinRegisterContractTxParam
	txParam.ValidHeight = 630314 //20999
	txParam.SrcRegId = "0-1"     //"7849-1"
	txParam.Fees = 110000000
	txParam.Script = script
	txParam.Description = "My hello contract!!!"
	txParam.FeeSymbol = common.WICCSYM

	rawTx,txid, err := txParam.CreateRawTx(privateKey)
	if err != nil {
		t.Error("UCoinRegisterContractTx err: ", err)
	}
	t.Log("rawTx=",rawTx)
	t.Log("txid=",txid)
}

/*
Asset Create 创建资产

fee Minimum 0.01 WICC
make sure account balance >550.01 wicc
发布资产需要扣除550wicc的发布费用
*/
func TestCreateAssetCreateRawTx(t *testing.T) {
	privateKey := "Y6J4aK6Wcs4A3Ex4HXdfjJ6ZsHpNZfjaS4B9w7xqEnmFEYMqQd13"
	var txParam AssetIssueTxParam
	txParam.FeeSymbol = common.WICCSYM
	txParam.Fees = 1000000
	txParam.ValidHeight = 11625
	txParam.SrcRegId = "0-1"
	txParam.AssetOwner="0-1"//only support regid
	txParam.AssetTotal=100*100000000
	txParam.AssetSymbol="SSSSSSS"
	txParam.AssetName="SK Token"
	txParam.MinTable=true
	rawTx,txid, err := txParam.CreateRawTx(privateKey)
	if err != nil {
		t.Error("SignAssetIssueTx err: ", err)
	}
	t.Log("rawTx=",rawTx)
	t.Log("txid=",txid)
}

/*
Asset Update 更新资产
fee Minimum 0.01WICC
make sure account balance >110.01 wicc
发布资产需要扣除110wicc的发布费用
*/
func TestCreateAssetUpdateRawTx(t *testing.T) {
	privateKey := "Y6J4aK6Wcs4A3Ex4HXdfjJ6ZsHpNZfjaS4B9w7xqEnmFEYMqQd13"
	var txParam AssetUpdateTxParam
	txParam.FeeSymbol = common.WICCSYM
	txParam.Fees = 1000000
	txParam.ValidHeight = 688709
	txParam.SrcRegId = "0-1"
	txParam.UpdateType=int(ASSET_OWNER_UID)
	txParam.AssetSymbol="STOKEN"
	txParam.AssetOwner="111-1" //only support regid
	//txParam.AssetTotal=100*100000000
	//txParam.AssetName="SK Token"
	rawTx,txid, err := txParam.CreateRawTx(privateKey)
	if err != nil {
		t.Error("AssetUpdateTx err: ", err)
	}
	t.Log("rawTx=",rawTx)
	t.Log("txid=",txid)
}


/*
* 投票交易
* Voting transaction
* fee Minimum 0.01 wicc
* */
func TestCreateDelegateRawTx(t *testing.T) {
	privateKey := "Y9sx4Y8sBAbWDAqAWytYuUnJige3ZPwKDZp1SCDqqRby1YMgRG9c"
	var txParams DelegateTxParam
	txParams.ValidHeight = 95728
	txParams.SrcRegId = ""
	txParams.Fees = 1000000
	txParams.PubKey = "03e93e7d870ce6f1c9997076c56fc24e6381c612662cd9a5a59294fac9ba7d21d7"
	txParams.Votes = NewOperVoteFunds()
	pubKey, _ := hex.DecodeString("025a37cb6ec9f63bb17e562865e006f0bafa9afbd8a846bd87fc8ff9e35db1252e")
	vote := OperVoteFund{PubKey: pubKey, VoteValue: 10000}
	txParams.Votes.Add(&vote)

	rawTx,txid, err := txParams.CreateRawTx(privateKey)
	if err != nil {
		t.Error("SignDelegateTx err: ", err)
	}
	t.Log("rawTx=",rawTx)
	t.Log("txid=",txid)
}


/*
    * 创建,追加cdp交易
    * Create or append an  cdp transaction
    * fee Minimum 0.01 wicc
    * */
func TestCreateCdpStakeRawTx(t *testing.T) {
	privateKey := "Y6J4aK6Wcs4A3Ex4HXdfjJ6ZsHpNZfjaS4B9w7xqEnmFEYMqQd13"
	var txParam CdpStakeTxParam
	txParam.CdpTxid = "0b9734e5db3cfa38e76bb273dba4f65a210cc76ca2cf739f3c131d0b24ff89c1"
	txParam.ScoinSymbol = common.WUSDSYM
	txParam.FeeSymbol = common.WICCSYM
	AssetSymbol:= common.WICCSYM
	AssetAmount:= 1000
	model:=AssetModel{int64(AssetAmount),AssetSymbol}
	txParam.Assets=NewCdpAssets()
	txParam.Assets.Add(&model)
	txParam.ScoinMint = 10
	txParam.Fees = 10000000
	txParam.ValidHeight = 283308
	txParam.SrcRegId = "0-1"
	txParam.PubKey = "03e93e7d870ce6f1c9997076c56fc24e6381c612662cd9a5a59294fac9ba7d21d7"

	rawTx,txid, err := txParam.CreateRawTx(privateKey)
	if err != nil {
		t.Error("CdpStakeTx err: ", err)
	}
	t.Log("rawTx=",rawTx)
	t.Log("txid=",txid)
}

/*
  * 赎回cdp交易
  * Redeem cdp transaction
  * fee Minimum 0.01 wicc
  * */
func TestCreateCdpRedeemRawTx(t *testing.T) {
	privateKey := "Y6J4aK6Wcs4A3Ex4HXdfjJ6ZsHpNZfjaS4B9w7xqEnmFEYMqQd13"
	var txParam CdpRedeemTxParam
	txParam.CdpTxid = "0b9734e5db3cfa38e76bb273dba4f65a210cc76ca2cf739f3c131d0b24ff89c1"
	txParam.FeeSymbol = common.WICCSYM
	txParam.ScoinsToRepay = 0
	AssetSymbol:= common.WICCSYM
	AssetAmount:= 1000
	model:=AssetModel{int64(AssetAmount),AssetSymbol}
	txParam.Assets=NewCdpAssets()
	txParam.Assets.Add(&model)
	txParam.Fees = 1000000
	txParam.ValidHeight = 25
	txParam.SrcRegId = "0-1"
	txParam.PubKey = "03e93e7d870ce6f1c9997076c56fc24e6381c612662cd9a5a59294fac9ba7d21d7"

	rawTx,txid, err := txParam.CreateRawTx(privateKey)
	if err != nil {
		t.Error("CdpRedeemTx err: ", err)
	}
	t.Log("rawTx=",rawTx)
	t.Log("txid=",txid)
}

/*
  * 清算cdp交易
  * Liquidate cdp transaction
  * fee Minimum 0.01 wicc
  * */
func TestCreateCdpLiquidateRawTx(t *testing.T) {
	privateKey := "Y6J4aK6Wcs4A3Ex4HXdfjJ6ZsHpNZfjaS4B9w7xqEnmFEYMqQd13"
	var txParam CdpLiquidateTxParam
	txParam.CdpTxid = "0b9734e5db3cfa38e76bb273dba4f65a210cc76ca2cf739f3c131d0b24ff89c1"
	txParam.FeeSymbol = common.WICCSYM
	txParam.ScoinsLiquidate = 100
	txParam.Fees = 1000000
	txParam.ValidHeight = 25
	txParam.SrcRegId = "0-1"
	txParam.PubKey = "03e93e7d870ce6f1c9997076c56fc24e6381c612662cd9a5a59294fac9ba7d21d7"
	txParam.AssetSymbol=common.WICCSYM

	rawTx, txid,err := txParam.CreateRawTx(privateKey)
	if err != nil {
		t.Error("CdpLiquidateTx err: ", err)
	}
	t.Log("rawTx=",rawTx)
	t.Log("txid=",txid)
}

/*
 * Dex 限价买单交易
 * Dex limit price transaction
 * fee Minimum 0.001 wicc
 * */
func TestCreateDexLimitBuyRawTx(t *testing.T) {
	privateKey := "Y6J4aK6Wcs4A3Ex4HXdfjJ6ZsHpNZfjaS4B9w7xqEnmFEYMqQd13"
	var txParam DexLimitBuyTxParam
	txParam.FeeSymbol = common.WICCSYM
	txParam.Fees = 1000000
	txParam.CoinSymbol = common.WUSDSYM
	txParam.AssetSymbol = common.WICCSYM
	txParam.AssetAmount = 10000
	txParam.ValidHeight = 25
	txParam.Price = 25
	txParam.SrcRegId = "0-1"
	txParam.PubKey = "03e93e7d870ce6f1c9997076c56fc24e6381c612662cd9a5a59294fac9ba7d21d7"

	rawTx, txid,err := txParam.CreateRawTx(privateKey)
	if err != nil {
		t.Error("DexBuyLimitTx err: ", err)
	}
	t.Log("rawTx=",rawTx)
	t.Log("txid=",txid)
}

/*
   * Dex 限价卖单交易
   * Dex limit sell price transaction
   * fee Minimum 0.001 wicc
  * */
func TestCreateDexLimitSellRawTx(t *testing.T) {
	privateKey := "Y6J4aK6Wcs4A3Ex4HXdfjJ6ZsHpNZfjaS4B9w7xqEnmFEYMqQd13"
	var txParam DexLimitSellTxParam
	txParam.FeeSymbol = common.WICCSYM
	txParam.Fees = 100000
	txParam.CoinSymbol = common.WUSDSYM
	txParam.AssetSymbol = common.WICCSYM
	txParam.AssetAmount = 1000000
	txParam.ValidHeight = 282956
	txParam.Price = 200000000
	txParam.SrcRegId = "0-1"
	txParam.PubKey = "03e93e7d870ce6f1c9997076c56fc24e6381c612662cd9a5a59294fac9ba7d21d7"

	rawTx,txid, err := txParam.CreateRawTx(privateKey)
	if err != nil {
		t.Error("DexSellLimitTx err: ", err)
	}
	t.Log("rawTx=",rawTx)
	t.Log("txid=",txid)
}

/*
     *  Dex 市价卖单交易
     * Dex market sell price transaction
     * fee Minimum 0.001 wicc
    * */
func TestCreateDexMarketSellRawTx(t *testing.T) {
	privateKey := "Y6J4aK6Wcs4A3Ex4HXdfjJ6ZsHpNZfjaS4B9w7xqEnmFEYMqQd13"
	var txParam DexMarketSellTxParam
	txParam.FeeSymbol = common.WICCSYM
	txParam.Fees = 1000000
	txParam.CoinSymbol = common.WUSDSYM
	txParam.AssetSymbol = common.WICCSYM
	txParam.AssetAmount = 10000
	txParam.ValidHeight = 25
	txParam.SrcRegId = "0-1"
	txParam.PubKey = "03e93e7d870ce6f1c9997076c56fc24e6381c612662cd9a5a59294fac9ba7d21d7"

	rawTx,txid, err := txParam.CreateRawTx(privateKey)
	if err != nil {
		t.Error("DexMarketSellTx err: ", err)
	}
	t.Log("rawTx=",rawTx)
	t.Log("txid=",txid)
}

/*
     *  Dex 市价买单交易
     * Dex market buy price transaction
     * fee Minimum 0.001 wicc
    * */
func TestCreateDexMarketBuyRawTx(t *testing.T) {
	privateKey := "Y6J4aK6Wcs4A3Ex4HXdfjJ6ZsHpNZfjaS4B9w7xqEnmFEYMqQd13"
	var txParam DexMarketBuyTxParam
	txParam.FeeSymbol = common.WICCSYM
	txParam.Fees = 1000000
	txParam.CoinSymbol = common.WUSDSYM
	txParam.AssetSymbol = common.WICCSYM
	txParam.AssetAmount = 10000
	txParam.ValidHeight = 25
	txParam.SrcRegId = "0-1"
	txParam.PubKey = "03e93e7d870ce6f1c9997076c56fc24e6381c612662cd9a5a59294fac9ba7d21d7"

	rawTx,txid, err := txParam.CreateRawTx(privateKey)
	if err != nil {
		t.Error("DexMarketBuyTx err: ", err)
	}
	t.Log("rawTx=",rawTx)
	t.Log("txid=",txid)
}

/*
 *  Dex 取消挂单交易
 * Dex cancel order tx
 * fee Minimum 0.001 wicc
* */
func TestCreateDexCancelRawTx(t *testing.T) {
	privateKey := "Y6J4aK6Wcs4A3Ex4HXdfjJ6ZsHpNZfjaS4B9w7xqEnmFEYMqQd13"
	var txParam DexCancelTxParam
	txParam.FeeSymbol = common.WICCSYM
	txParam.Fees = 100000
	txParam.DexTxid = "009c0e665acdd9e8ae754f9a51337b85bb8996980a93d6175b61edccd3cdc144"
	txParam.ValidHeight = 25
	txParam.SrcRegId = "0-1"
	txParam.PubKey = "03e93e7d870ce6f1c9997076c56fc24e6381c612662cd9a5a59294fac9ba7d21d7"

	rawTx,txid, err := txParam.CreateRawTx(privateKey)
	if err != nil {
		t.Error("DexCancelTx err: ", err)
	}
	t.Log("rawTx=",rawTx)
	t.Log("txid=",txid)
}

// Support UCOIN_TRANSFER_TX and UCOIN_CONTRACT_INVOKE_TX only
func TestDecodeRawTx(t *testing.T){
	//rawTx := "0b01df390684f0c10c82480457494343bc834002141c758724cc60db35dd387bcf619a478ec3c065f20457494343bc8340142af03ec43eb893039b5dd5bab612d73034cf1b610457555344858c1f00473045022100d68782ebf4059ac26b169ae035ca2a8c1533c4f5639c9fd64445f205d86fbf2c022008b7ed1467ec9321382284ce9d762967a604602a26295f4d569f9a15b643e1db"
	//rawTx := "0b01df3921036c5397f3227a1e209952829d249b7ad0f615e43b763ac15e3a6f52627a10df210457494343bc83400214079b9296a00a2b655787fa90e66ec3cde4bf1c8c0457494343bc834014079b9296a00a2b655787fa90e66ec3cde4bf1c8c0457494343866700473045022100a80f36f5b260bdbb76a46f0d9563bdfe0c79a8c2f7b5935b5e79a3c280040551022001ba6dee175619b2c3ba19976690d848ffdbfb4df596bdf866932a4b7befb022"
	rawTx := "0b01c7a10803de4901045749434382dbea93000114079b9296a00a2b655787fa90e66ec3cde4bf1c8c045749434382dbea930006e8bdace8b4a647304502210091922890a5ccc26fa2c6c404378d3684a414f952018f0dfa414ee463f81e6cfa022057babb41679557218c6bdf6fa24f495423289639aac71643566d1e8bb6c0f5f8"
	netParams := common.WICCTestnetParams //testnet
//	netParams := common.WICCParams  //mainnet
	result ,err := DecodeRawTx(rawTx,netParams)
	if err != nil {
		t.Error("Umarshal failed:", err)
	}
	jsonBytes, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		t.Error("Umarshal failed:", err)
	}

	t.Log("DecodeRawTx result=\n",string(jsonBytes))
}