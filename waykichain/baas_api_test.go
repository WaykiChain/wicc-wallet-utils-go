package waykichain

import (
	"encoding/json"
	"testing"
)

func TestBaaSGetInfo(t *testing.T) {

	r, err := tw.BaaSClient.GetInfo()
	if err != nil {
		t.Errorf("GetInfo failed: %v\n", err)
	}
	jsonBytes, err := json.MarshalIndent(r, "", "    ")
	if err != nil {
		t.Error("MarshalIndent failed:", err)
	}

	t.Log("getinfo=\n",string(jsonBytes))
}

func TestBaaSGetSynBlockHeight(t *testing.T) {

	synBlockHeight, err := tw.WalletClient.GetSynBlockHeight()
	if err != nil {
		t.Errorf("GetSynBlockHeight failed: %v\n", err)
	}

	t.Log("synBlockHeight=",synBlockHeight)
}

func TestBaaSSubmitTxRaw(t *testing.T) {

	rawTx := "0b01809aac1621031b27286c65b81ac13cfd4067b030398a19eb147f439c094fbb19a2f3ab9ec10b0457494343bc834001140a0145a6f4d60eed4ac3ddd34489898bb8e9ab500457494343bc834000463044022067e69d527d4c7ef69c7b31656bba3cf79926ef8bc9a1f876b34a84b5d53a1491022078bc3f787b43b77bb82852fe4fade822d1c2465924bf3f0d8b5c2eb449dfaed0"
	txid, err := tw.BaaSClient.SubmitTxRaw(rawTx)
	if err != nil {
		t.Errorf("SubmitTxRaw failed: %v\n", err)
	}

	//a9a62c7c5e48162ce62dc11eb7cce7646dd9cb895c99be4515ae420740fd2ea7
	t.Log("txid=",txid)
}
