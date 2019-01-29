package commons

import (
	"fmt"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/btcsuite/btcutil"
)


func GenerateAddress(words string, net Network)  string {
	wallet, err := CreateWalletFromMnemonic(words, DefaultSeedPass)
	if (err != nil) {
		fmt.Errorf("invalid mnemonic")
		return ""
	}
	netParams, err := NetworkToChainConfig(net)
	if (err != nil) {
		fmt.Errorf("invalid network")
		return ""
	}
	wallet.Initialize([]uint32{HardenedKeyZeroIndex})
	btc := wallet.Coins[0]
	acc, err := btc.DeriveAccountNode(hdkeychain.HardenedKeyStart + 0)
	if err != nil {
		fmt.Errorf("Failed to derive account node")
		return ""
	}
	ch, err := acc.DeriveChainNode(0)
	if err != nil {
		fmt.Errorf("Failed to derive chain node")
		return ""
	}
	addr, err := ch.DeriveAddressNode(0)
	if err != nil {
		fmt.Errorf("Failed to derive address node")
		return ""
	}
	address, err := addr.Address.Address(netParams)
	if (err != nil) {
		fmt.Errorf("Failed to generate address")
		return ""
	}
	return address.String()
}

func GeneratePrivateKey(words string ,net Network) string{
	wallet, err := CreateWalletFromMnemonic(words, DefaultSeedPass)
	if (err != nil) {
		fmt.Errorf("invalid mnemonic")
		return ""
	}
	netParams, err := NetworkToChainConfig(net)
	if (err != nil) {
		fmt.Errorf("invalid network")
		return ""
	}
	wallet.Initialize([]uint32{HardenedKeyZeroIndex})
	btc := wallet.Coins[0]
	acc, err := btc.DeriveAccountNode(hdkeychain.HardenedKeyStart + 0)
	if err != nil {
		fmt.Errorf("Failed to derive account node")
		return ""
	}
	ch, err := acc.DeriveChainNode(0)
	if err != nil {
		fmt.Errorf("Failed to derive chain node")
		return ""
	}
	addr, err := ch.DeriveAddressNode(0)
	if err != nil {
		fmt.Errorf("Failed to derive address node")
		return ""
	}
	priv,_:=addr.Address.ECPrivKey()
	wif1, err := btcutil.NewWIF(priv, netParams,true)
	if (err != nil) {
		fmt.Errorf("Failed to generate address")
		return ""
	}
	return wif1.String()
}
