package waykichain

import (
	"encoding/json"
	"testing"
)

func TestGetInfo(t *testing.T) {

	r, err := tw.WalletClient.GetInfo()
	if err != nil {
		t.Errorf("GetInfo failed: %v\n", err)
	}
	jsonBytes, err := json.MarshalIndent(r, "", "    ")
	if err != nil {
		t.Error("MarshalIndent failed:", err)
	}

	t.Log("getinfo=\n",string(jsonBytes))
}

func TestGetSynBlockHeight(t *testing.T) {

	synBlockHeight, err := tw.WalletClient.GetSynBlockHeight()
	if err != nil {
		t.Errorf("GetSynBlockHeight failed: %v\n", err)
	}

	t.Log("synBlockHeight=",synBlockHeight)
}

func TestSubmitTxRaw(t *testing.T) {

	rawTx := "0b018099fa4621031b27286c65b81ac13cfd4067b030398a19eb147f439c094fbb19a2f3ab9ec10b0457494343bc834001140a0145a6f4d60eed4ac3ddd34489898bb8e9ab500457494343bc83400046304402201c794cc99d97374086d10901d53acf18a5aee8f5dac0e47fdf76ac9a9708cb5902204da3543a02e5d64fe9868fe5284dba3bce521580aeced30d931e277ebddf227d"
	txid, err := tw.WalletClient.SubmitTxRaw(rawTx)
	if err != nil {
		t.Errorf("SubmitTxRaw failed: %v\n", err)
	}

	t.Log("txid=",txid)
}
