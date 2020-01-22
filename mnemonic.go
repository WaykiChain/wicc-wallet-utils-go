package wicc_wallet_utils_go

import (
	"errors"
	"github.com/JKinGH/go-hdwallet"
	"github.com/tyler-smith/go-bip39"
	"strings"
)

// New mnemonic follow the wordlists
func GenerateMnemonic(length int) (string, error) {
	mnemonic, err := hdwallet.NewMnemonic(length)
	if err != nil {
		return "", err
	}

	return mnemonic, nil
}

func ChangeMnemonicLanguage(englishMnemonic string, language string) (string, error){
	//把英文助记词按空格分割
	words := strings.Fields(englishMnemonic)
	indexs := make([]int,0)

	//获取英文助记词在词库中的index
	for _, word := range words{
		index, ok := bip39.GetWordIndex(word)
		if !ok {
			return "", errors.New("GetWordIndex error!")
		}
		indexs = append(indexs,index)
	}

	//设置需要转换的语言并获取该词库内容
	newwords := make([]string,len(indexs))
	hdwallet.SetLanguage(language)
	wordList := bip39.GetWordList()

	//根据index在目标词库中获取对应助记词
	for i := len(indexs) - 1; i >= 0; i-- {
		newwords[i] = wordList[indexs[i]]
	}

	//将新获取的助记词以空格分隔序列号成string
	return strings.Join(newwords, " "), nil
}
