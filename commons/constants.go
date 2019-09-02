package commons

import (
	"errors"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
)

type WICCNet uint32
type WalletStatus int64
type ChangeType uint32
type Network int16

const (
	// MainNet represents the main wicc network.
	MainNet wire.BitcoinNet = 0xff421d1a

	// TestNet represents the  test wicc network.
	TestNet wire.BitcoinNet = 0xfd7d5cd7

	// Default entropy size for mnemonic
	DefaultEntropySize = 128
	// Default seed pass. it used to generate seed from mnemonic( BIP39 ). Don't change if determined
	DefaultSeedPass = ""

	HardenedKeyZeroIndex        = 0x8001869f
	BIP44Purpose         uint32 = 0x8000002C
	WICCCoinType         uint32 = 99999

	MAINNET Network = 1
	TESTNET Network = 2

	BIP44PATH = "m/44'/99999'/0'/0/0"
)

type MnemonicLanguage string

// List Mnemonic language support
const (
	ENGLISH  MnemonicLanguage = "EN"
	JAPANESE                  = "JP"
	FRENCH                    = "FR"
	ITALIAN                   = "IT"
	KOREAN                    = "KR"
	SPANISH                   = "ES"
)
type CoinType string
const (
	WICC CoinType = "WICC"
	WGRT          = "WGRT"
	WUSD          = "WUSD"
	WCNY          = "WCNY"
	WBTC          = "WBTC"
	WETH          = "WETH"
	WEOS          = "WEOS"
	USD           = "USD"
	CNY           = "CNY"
	EUR           = "EUR"
	BTC           = "BTC"
	USDT          = "USDT"
	GOLD          = "GOLD"
	KWH           = "KWH"
)

func NetworkToChainConfig(net Network) (*chaincfg.Params, error) {
	switch net {
	case 1:
		return &WaykiMainNetParams, nil

	case 2:
		return &WaykiTestNetParams, nil
	}

	return nil, errors.New("invalid network")
}
