package hdwallet

import (
	"github.com/tyler-smith/go-bip39"
	"github.com/tyler-smith/go-bip39/wordlists"
)

func SetLanguage(language string) {
	switch language {
	case ChineseSimplified:
		bip39.SetWordList(wordlists.ChineseSimplified)
	case ChineseTraditional:
		bip39.SetWordList(wordlists.ChineseTraditional)
	case JAPANESE:
		bip39.SetWordList(wordlists.Japanese)
	case ITALIAN:
		bip39.SetWordList(wordlists.Italian)
	case KOREAN:
		bip39.SetWordList(wordlists.Korean)
	case SPANISH:
		bip39.SetWordList(wordlists.Spanish)
	case FRENCH:
		bip39.SetWordList(wordlists.French)
	default:
		bip39.SetWordList(wordlists.English)
	}
}

// NewMnemonic creates a random mnemonic
// Return the English Mnemonic default
func NewMnemonic(length int) (string, error) {
	//setLanguage(language)

	if length < 12 {
		length = 12
	}

	if length > 24 {
		length = 24
	}

	entropy, err := bip39.NewEntropy(length / 3 * 32)
	if err != nil {
		return "", err
	}

	return bip39.NewMnemonic(entropy)
}

// NewSeed creates a hashed seed
func NewSeed(mnemonic, password, language string) ([]byte, error) {
	SetLanguage(language)
	return bip39.NewSeedWithErrorChecking(mnemonic, password)
}
