/*
 * Copyright 2018 The openwallet Authors
 * This file is part of the openwallet library.
 *
 * The openwallet library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The openwallet library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 */

package wicc_wallet_utils_go

import (
	"encoding/json"
	"testing"
)

func TestEthGetBlockNumber(t *testing.T) {


	if r, err := tw.WalletClient.EthGetBlockNumber(); err != nil {
		t.Errorf("GetAccountNet failed: %v\n", err)
	} else {
		t.Logf("GetAccountNet return: \n\t%+v\n", r)
	}
}

func TestGetInfuraChainId(t *testing.T) {

	if r, err := tw.WalletClient.GetInfuraChainId(); err != nil {
		t.Errorf("EthGetChainId failed: %v\n", err)
	} else {
		t.Logf("EthGetChainId return: \n\t%+v\n", r)
	}
}

func TestGetTransactionCount(t *testing.T) {

	address := "0x232D23C22543144B988F738C701Df6dfd6eAcA4c"

	if r, err := tw.WalletClient.ethGetTransactionCount(address,LEATEST); err != nil {
		t.Errorf("ethGetTransactionCount failed: %v\n", err)
	} else {
		t.Logf("ethGetTransactionCount return: \n\t%+v\n", r)
	}
}

func TestGetAddrBalance(t *testing.T) {

	address := "0x232D23C22543144B988F738C701Df6dfd6eAcA4c"

	if r, err := tw.WalletClient.GetAddrBalance(address,LEATEST); err != nil {
		t.Errorf("GetAddrBalance failed: %v\n", err)
	} else {
		t.Logf("GetAddrBalance return: \n\t%+v\n", r)
	}
}


func TestERC20GetAddressBalance(t *testing.T) {

	address := "0x232D23C22543144B988F738C701Df6dfd6eAcA4c"
	contract_address := "0x8E1dA42EbC22F91d528ceB9865f241167Ebb8A0f"

	if r, err := tw.WalletClient.ERC20GetAddressBalance(address,contract_address); err != nil {
		t.Errorf("ERC20GetAddressBalance failed: %v\n", err)
	} else {

		t.Logf("ERC20GetAddressBalance return: \n\t%+v\n", r)
		t.Logf("ERC20GetAddressBalance return2: \n\t%+v\n", r.Int64())
	}
}

func TestEthGetTransactionByHash(t *testing.T) {

	txid := "0xab5a428263314427eb01cffa5ac63ecf4ad4d29601b6c9b4476b5e5c4df840cf"
	tx_submit, err := tw.WalletClient.EthGetTransactionByHash(txid)
	if err != nil {
		t.Errorf("get transaction by has failed, err=%v", err)
		return
	}
	jsonBytes, err := json.MarshalIndent(*tx_submit, "", "    ")
	if err != nil {
		t.Error("Umarshal failed:", err)
	}

	t.Log("tx_submit=\n",string(jsonBytes))
//	log.Infof("tx_submit: %+v", tx_submit)
}

