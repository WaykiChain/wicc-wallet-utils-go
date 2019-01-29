package commons

import (
	"errors"
	"wiccwallet/wordslists"
)

// Please refer the link: https://iancoleman.io/bip39/ for purpose double check result

type Mnemonic struct {
	EntropySize int
	Password    string
}

func NewMnemonicWithDefaultOption() *Mnemonic {
	return &Mnemonic{EntropySize: DefaultEntropySize, Password: DefaultSeedPass}
}

func NewMnemonicWithLanguage(language MnemonicLanguage) *Mnemonic {
	SetWordList(loadWordList(language))
	return &Mnemonic{EntropySize: DefaultEntropySize, Password: DefaultSeedPass}
}

// New mnemonic follow the wordlists
func (m *Mnemonic) GenerateMnemonic() (string, error) {
	entropy, err := NewEntropy(m.EntropySize)
	if err != nil {
		return "", err
	}

	return NewMnemonic(entropy)
}

// Generate seed from mnemonic and pass( optional )
func (m *Mnemonic) GenerateSeed(mnemonic string) ([]byte, error) {
	if !IsMnemonicValid(mnemonic) {
		return nil, errors.New("invalidate mnemonic")
	}
	return NewSeed(mnemonic, m.Password), nil
}

// Get word list
func (m *Mnemonic) ListWord() []string {
	return GetWordList()
}

// loadWordList returns word lists base on language setting in the configuration
func loadWordList(language MnemonicLanguage) []string {
	switch language {
	case JAPANESE:
		return wordlists.Japanese
	case ITALIAN:
		return wordlists.Italian
	case KOREAN:
		return wordlists.Korean
	case SPANISH:
		return wordlists.Spanish
	case FRENCH:
		return wordlists.French
	default:
		return wordlists.English
	}
}