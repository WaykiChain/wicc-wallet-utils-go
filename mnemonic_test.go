package wicc_wallet_utils_go

import (
	"github.com/JKinGH/go-hdwallet"
	"testing"
)

func TestGenerateMnemonics(t *testing.T) {

	mnemonic,err := GenerateMnemonic(12)
	if err != nil {
		t.Errorf("Failed to GenerateMnemonic: %v",err)
	}

	t.Log("mnemonic=",mnemonic)
}

func TestChangeMnemonicLanguage(t *testing.T){
	//englishMnemonic := "jazz wine firm worth cry dumb glad foam viable knee pride purse"
	englishMnemonic := "wreck bullet carpet nerve belt border often trust exchange believe defense rebel"
	newMnemonic,err :=  ChangeMnemonicLanguage(englishMnemonic, hdwallet.ChineseSimplified)
	if err != nil {
		t.Errorf("Failed to ChangeMnemonicLanguage: %v",err)
	}

	t.Log("mnemonic=",newMnemonic)
}
