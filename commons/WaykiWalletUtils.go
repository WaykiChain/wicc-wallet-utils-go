package commons

import (
	"fmt"
	"github.com/btcsuite/btcutil"
)

func ImportPrivateKey(priv string ,net Network) string{

	netParams, err := NetworkToChainConfig(net)
	if (err != nil) {
		fmt.Errorf("invalid network")
		return ""
	}
	wif1, err := btcutil.DecodeWIF(priv)
	if (err != nil) {
		fmt.Errorf("invalid privatekey")
		return ""
	}
	address, err := btcutil.NewAddressPubKey(PublicKeyForPrivateKey(wif1.PrivKey.Serialize()),netParams)

	if (err != nil) {
		fmt.Errorf("Failed to generate address")
		return ""
	}
	return address.EncodeAddress()
}
