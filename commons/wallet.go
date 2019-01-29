package commons

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"log"
)

type Address struct {
	Address *hdkeychain.ExtendedKey
	Network *chaincfg.Params
}

type Chain struct {
	Chain   *hdkeychain.ExtendedKey
	Network *chaincfg.Params
}

type Account struct {
	Account *hdkeychain.ExtendedKey
	Network *chaincfg.Params
}

type Coin struct {
	Name    string
	Coin    *hdkeychain.ExtendedKey
	Network *chaincfg.Params
}

type Wallet struct {
	Entropy     string
	Mnemonic    string
	Seed        string
	MasterNode  *hdkeychain.ExtendedKey
	PurposeNode *hdkeychain.ExtendedKey
	Coins       []*Coin
}

func CreateWalletWithPassword(password string) (w *Wallet, err error) {
	entropy, _ := NewEntropy(DefaultEntropySize)
	entropyToHexString := hex.EncodeToString(entropy)

	mnemonic, _ := NewMnemonic(entropy)

	seed, err := NewSeedWithErrorChecking(mnemonic, password)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	seedToHexString := hex.EncodeToString(seed)

	//@ToDo: create network params for FLO and LTC, etc

	masterKey, _ := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	purposeNode, _ := masterKey.Child(BIP44Purpose)

	wallet := Wallet{
		Mnemonic:    mnemonic,
		Seed:        seedToHexString,
		Entropy:     entropyToHexString,
		MasterNode:  masterKey,
		PurposeNode: purposeNode,
	}

	return &wallet, nil
}

func CreateWalletFromMnemonic(mnemonic, password string) (w *Wallet, err error) {
	entropy, _ := EntropyFromMnemonic(mnemonic)
	entropyToHexString := hex.EncodeToString(entropy)

	seed, err := NewSeedWithErrorChecking(mnemonic, password)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	seedToHexString := hex.EncodeToString(seed)

	//@ToDo: create network params for FLO and LTC, etc

	masterKey, _ := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	purposeNode, _ := masterKey.Child(BIP44Purpose)

	wallet := Wallet{
		Mnemonic: mnemonic,
		Seed:     seedToHexString,
		//ToDo: Derive entropy from mnemonic or seed
		Entropy:     entropyToHexString,
		MasterNode:  masterKey,
		PurposeNode: purposeNode,
	}

	return &wallet, nil
}

func CreateWalletFromSeed(seed []byte) (w *Wallet, err error) {
	//@ToDo: create network params for FLO and LTC, etc
	masterKey, _ := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	purposeNode, _ := masterKey.Child(hdkeychain.HardenedKeyStart + 44)

	seedToHexString := hex.EncodeToString(seed)

	wallet := Wallet{
		Seed: seedToHexString,
		//ToDo: Derive entropy from mnemonic or seed
		MasterNode:  masterKey,
		PurposeNode: purposeNode,
	}

	return &wallet, nil
}

func (w *Wallet) Initialize(bip44CoinConstants []uint32) (*Wallet, error) {

	for i := 0; i < len(bip44CoinConstants); i++ {
		//ToDo: make this dynamic to where it will choose the network configs based on the constant
		coin, err := w.DeriveCoinNode(&chaincfg.MainNetParams, bip44CoinConstants[i])
		if err != nil {
			log.Fatal("Failed to Derive coin node: terminate.")
		}

		w.Coins = append(w.Coins, coin)
	}

	return w, nil
}

//pkg/errors w errors.wrap
func (w *Wallet) DeriveCoinNode(network *chaincfg.Params, bip44CoinConstant uint32) (c *Coin, err error) {
	if bip44CoinConstant < hdkeychain.HardenedKeyStart {
		bip44CoinConstant += hdkeychain.HardenedKeyStart
	}
	coin, err := w.PurposeNode.Child(bip44CoinConstant)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	coin.SetNet(network)

	_coin := &Coin{
		Coin:    coin,
		Network: network,
	}

	return _coin, nil
}

//ToDo: create function to derive addresses at certain indices
//ToDo: Create GetCoin method for Wallet -> returns *Coin

func (c *Coin) DeriveAccountNode(index uint32) (account *Account, err error) {
	if index < hdkeychain.HardenedKeyStart {
		index += hdkeychain.HardenedKeyStart
	}
	a, err := c.Coin.Child(index)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	_account := &Account{
		Account: a,
		Network: c.Network,
	}

	return _account, nil
}

func (a *Account) DeriveChainNode(index uint32) (chain *Chain, err error) {
	c, err := a.Account.Child(index)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	_chain := &Chain{
		Chain:   c,
		Network: a.Network,
	}

	return _chain, nil
}

func (c *Chain) DeriveAddressNode(index uint32) (address *Address, err error) {
	a, err := c.Chain.Child(index)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	_address := &Address{
		Address: a,
		Network: c.Network,
	}

	return _address, nil
}
