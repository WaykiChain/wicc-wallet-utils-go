package wiccwallet

import (
	"errors"
	"wicc-wallet-utils-go/commons"
	"wicc-wallet-utils-go/wordlists"
)

// List Mnemonic language support
const (
	ENGLISH  string = "EN"
	JAPANESE        = "JP"
	FRENCH          = "FR"
	ITALIAN         = "IT"
	KOREAN          = "KR"
	SPANISH         = "ES"
)

// Please refer the link: https://iancoleman.io/bip39/ for purpose double check result

type Mnemonic struct {
	EntropySize int
	Password    string
}

func NewMnemonicWithDefaultOption() *Mnemonic {
	return &Mnemonic{EntropySize: commons.DefaultEntropySize, Password: commons.DefaultSeedPass}
}

func NewMnemonicWithLanguage(language string) *Mnemonic {
	commons.SetWordList(loadWordList(language))
	return &Mnemonic{EntropySize: commons.DefaultEntropySize, Password: commons.DefaultSeedPass}
}

// New mnemonic follow the wordlists
func (m *Mnemonic) GenerateMnemonic() (string, error) {
	entropy, err := commons.NewEntropy(m.EntropySize)
	if err != nil {
		return "", err
	}

	return commons.NewMnemonic(entropy)
}

// Generate seed from mnemonic and pass( optional )
func (m *Mnemonic) GenerateSeed(mnemonic string) ([]byte, error) {
	if !commons.IsMnemonicValid(mnemonic) {
		return nil, errors.New("invalidate mnemonic")
	}
	return commons.NewSeed(mnemonic, m.Password), nil
}

// Get word list
func (m *Mnemonic) ListWord() []string {
	return commons.GetWordList()
}

// loadWordList returns word lists base on language setting in the configuration
func loadWordList(language string) []string {
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
