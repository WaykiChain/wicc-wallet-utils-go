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

package bitcoin

import (
	"github.com/WaykiChain/wicc-wallet-utils-go/log"
	"github.com/blocktree/openwallet/common"
	"github.com/shopspring/decimal"
	"net/url"
	"testing"
)



func TestListUnspentByBitcore(t *testing.T) {
	list, err := tw.BitcoreClient.listUnspentByBitcore(0, "midiZgstuUWHJhNgpWxX9unT7g9chBvzwd")
	if err != nil {
		t.Errorf("listUnspentByExplorer failed unexpected error: %v\n", err)
		return
	}
	for i, unspent := range list {
		t.Logf("listUnspentByExplorer[%d] = %v \n", i, unspent)
	}

}




func TestEstimateFeeRateByExplorer(t *testing.T) {
	feeRate, _ := tw.BitcoreClient.EstimateFeeRateByBitcore()
	t.Logf("EstimateFee feeRate = %s\n", feeRate.String())
	fees, _ := tw.EstimateFee(10, 2, feeRate)
	t.Logf("EstimateFee fees = %v\n", fees)
}

func TestURLParse(t *testing.T) {
	apiUrl, err := url.Parse("http://192.168.32.107:20003/insight-api/")
	if err != nil {
		t.Errorf("url.Parse failed unexpected error: %v\n", err)
		return
	}
	domain := apiUrl.Hostname()
	port := common.NewString(apiUrl.Port()).Int()
	t.Logf("%s : %d", domain, port)
}

func TestDecimalAdd(t *testing.T) {
	unconfirmBalance, _ := decimal.NewFromString("-5.1")
	confirmBalance, _ := decimal.NewFromString("5.1")
	balance := confirmBalance.Add(unconfirmBalance)
	log.Infof("balance = %s", balance.String())
}