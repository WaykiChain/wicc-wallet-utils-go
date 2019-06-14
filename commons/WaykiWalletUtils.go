package commons

import (
	"fmt"
	"github.com/WaykiChain/wicc-wallet-utils-go/commons/wif"
	"github.com/btcsuite/btcutil"
	"github.com/WaykiChain/wicc-wallet-utils-go/commons/ec"
)

func GetAddressFromMnemonic(words string, net Network) string {
	seed := NewSeed(words, DefaultSeedPass)
	netParams, err := NetworkToChainConfig(net)
	if (err != nil) {
		fmt.Errorf("invalid network")
		return ""
	}
	key, err := NewMasterKey(seed, netParams)
	if err != nil {
		fmt.Errorf("Failed to new masterKey")
		return ""
	}
	childKey, err := key.DerivePath(BIP44PATH)
	if err != nil {
		fmt.Errorf("Failed to derive address node")
		return ""
	}
	privKey, err := childKey.ECPrivKey()
	if err != nil {
		fmt.Errorf("Failed to generate privatekey")
		return ""
	}
	_,pubk:=ec.PrivKeyFromBytes(privKey.Serialize())
	address, err := btcutil.NewAddressPubKey(pubk.SerializeCompressed(),netParams)
	if err != nil {
		fmt.Errorf("Failed to generate Address")
		return ""
	}

	return address.EncodeAddress()
}

func GetPrivateKeyFromMnemonic(words string, net Network) string {
	seed,err:= NewSeedWithErrorChecking(words, DefaultSeedPass)
	if (err != nil) {
		fmt.Errorf("invalid seed")
		return ""
	}
	netParams, err := NetworkToChainConfig(net)
	if (err != nil) {
		fmt.Errorf("invalid network")
		return ""
	}
	key, err := NewMasterKey(seed, netParams)
	if err != nil {
		fmt.Errorf("Failed to new masterKey")
		return ""
	}

	childKey, err := key.DerivePath(BIP44PATH)
	if err != nil {
		fmt.Errorf("Failed to derive address node")
		return ""
	}
	privKey, err := childKey.ECPrivKey()
	if err != nil {
		fmt.Errorf("Failed to generate privatekey")
		return ""
	}
	wif := wif.NewWIF(privKey, netParams.PrivateKeyID)
	if (err != nil) {
		fmt.Errorf("Failed to generate privatekey")
		return ""
	}
	return wif.EncodeCompressed()
}

func GetAddressFromPrivateKey(priv string ,net Network) string{
	netParams, err := NetworkToChainConfig(net)
	if (err != nil) {
		fmt.Errorf("invalid network")
		return ""
	}
	wifpr, err := btcutil.DecodeWIF(priv)
	if (err != nil) {
		fmt.Errorf("decode wif error")
		return ""
	}
	_,pub:=ec.PrivKeyFromBytes(wifpr.PrivKey.Serialize())
	address, err := btcutil.NewAddressPubKey(pub.SerializeCompressed(),netParams)
	if (err != nil) {
		fmt.Errorf("Failed to generate address")
		return ""
	}
	return address.EncodeAddress()
}
