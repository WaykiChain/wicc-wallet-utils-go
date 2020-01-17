package common

import (
	"errors"
	"github.com/JKinGH/go-hdwallet"
	"github.com/WaykiChain/wicc-wallet-utils-go/common/ec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
	"github.com/tyler-smith/go-bip39"
	"strings"
)

var (
	BTCW 			  *Wallet
	BTCSegwitW  	  *Wallet
	BTCTestnetW       *Wallet
	BTCTestnetSegwitW *Wallet
	ETHW 			  *Wallet
	ETHLedgerW 		  *Wallet
	WICCW 			  *Wallet
	WICCTestnetW	  *Wallet
)

func init(){
	BTCW = NewWallet(hdwallet.BTC,false,false, &hdwallet.BTCParams)
	BTCSegwitW = NewWallet(hdwallet.BTC,true,false, &hdwallet.BTCParams)
	BTCTestnetW = NewWallet(hdwallet.BTC_TESTNET,false,false, &hdwallet.BTCTestnetParams)
	BTCTestnetSegwitW = NewWallet(hdwallet.BTC_TESTNET,true,false, &hdwallet.BTCTestnetParams)
	ETHW = NewWallet(hdwallet.ETH,false,false, nil)
	ETHLedgerW = NewWallet(hdwallet.ETH,false,true, nil)
	WICCW = NewWallet(hdwallet.WICC,false,false, &hdwallet.WICCParams)
	WICCTestnetW = NewWallet(hdwallet.WICC_TESTNET,false,false, &hdwallet.WICCTestnetParams)
}

type Wallet struct{
	CoinType uint32  			//BIP44 coinType
	IsSegwit bool    			//Is segwit address(BIP49) for BTC or not
	IsLedger bool    			//Is Ledger wallet address for ETH or not
	NetParam *chaincfg.Params   //mainnet or testnet
}

func NewWallet(cointype uint32, issegwit, isledger bool ,netParam *chaincfg.Params) *Wallet {
	return &Wallet{cointype,issegwit,isledger,netParam}
}

// New mnemonic follow the wordlists
func (w *Wallet) GenerateMnemonic(length int) (string, error) {
	mnemonic, err := hdwallet.NewMnemonic(length)
	if err != nil {
		return "", err
	}

	return mnemonic, nil
}

func (w *Wallet) ChangeMnemonicLanguage(englishMnemonic string, language string) (string, error){
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

// Generate seed from mnemonic and pass( optional )
func GenerateSeed(mnemonic ,password,language string) ([]byte, error) {
	return hdwallet.NewSeed(mnemonic, password,language)
}

// Generate address By Mnemonic for BTC/ETH/WICC ...
func (w *Wallet) GenerateAddressFromMnemonic(mnemonic,language string) (string, error) {

	//get master key
	master, err := hdwallet.NewKey(w.IsLedger,hdwallet.Mnemonic(mnemonic),hdwallet.Language(language))
	if err != nil {
		return "",err
	}

	wallet, err := master.GetWallet(hdwallet.CoinType(w.CoinType),IsSegwit(w.IsSegwit))
	if err != nil {
		return "",err
	}

	//Not Segwit address
	address,err := wallet.GetAddress()
	if err != nil {
		return "",err
	}
	//Segwit address
	if w.IsSegwit == true {
		address,err = wallet.GetKey().AddressP2WPKHInP2SH()
		if err != nil {
			return "",err
		}
	}

	return address, nil
}

//Import PrivateKey return address except ETH
func (w *Wallet) GenerateAddressFromPrivateKey(privateKey string) (string, error){

	WIF, err := btcutil.DecodeWIF(privateKey)
	if (err != nil) {
		return "", err
	}

	_, publicKey:= ec.PrivKeyFromBytes(WIF.PrivKey.Serialize())
	if  publicKey == nil{
		return "", nil
	}
	//publicKey Compressed 33 bytes
	//publicKeyStr := hex.EncodeToString(publicKey.SerializeCompressed())

	//Not Segwit address
	addr1, err := btcutil.NewAddressPubKey(publicKey.SerializeCompressed(),w.NetParam)
	if (err != nil) {
		return "", err
	}
	address := addr1.EncodeAddress()

	//Segwit address
	if w.IsSegwit == true {
		addr, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(publicKey.SerializeCompressed()), w.NetParam)
		if err != nil {
			return "", err
		}
		script, err := txscript.PayToAddrScript(addr)
		if err != nil {
			return "", err
		}
		addr2, err := btcutil.NewAddressScriptHash(script, w.NetParam)
		if err != nil {
			return "", err
		}
		address = addr2.EncodeAddress()
	}

	return address, nil
}

//Export PrivateKey From Mnemonics
func (w *Wallet) ExportPrivateKeyFromMnemonic(mnemonic ,language string ) (string,error){

	master, err := hdwallet.NewKey(w.IsLedger,hdwallet.Mnemonic(mnemonic),hdwallet.Language(language))
	if err != nil {
		return "" , err
	}
	wallet, err := master.GetWallet(hdwallet.CoinType(w.CoinType),IsSegwit(w.IsSegwit))
	if err != nil {
		return "" , err
	}

	//ETH privatekey not WIF
	if w.CoinType == hdwallet.ETH{
		return wallet.GetKey().PrivateHex(), nil
	}

	WIFPrivateKey, err := wallet.GetKey().PrivateWIF(true)
	if err != nil {
		return "" , err
	}
	return WIFPrivateKey, nil
}

//Check address isvalid expect ETH
func CheckAddress(address string, netParams *chaincfg.Params) (bool, error) {

	_,err := btcutil.DecodeAddress(address,netParams)
	if err != nil{
		return false ,err
	}
	return true ,nil
}

//Check privatekey isvalid expect ETH
func CheckPrivateKey(privateKey string, netParams *chaincfg.Params) (bool, error) {
	_, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return false, ERR_INVALID_PRIVATEKEY
	}
	versionAndDataBytes := base58.Decode(privateKey)
	if (len(versionAndDataBytes) < 1) {
		return false, ERR_INVALID_PRIVATEKEY
	}
	version := versionAndDataBytes[0] & 0xFF
	if (netParams.PrivateKeyID == version) {
		return true, nil
	} else {
		return false, ERR_INVALID_PRIVATEKEY
	}
	return false,ERR_INVALID_PRIVATEKEY
}

func IsSegwit(segwit bool) hdwallet.Option {
	if segwit == true {
		return hdwallet.Purpose(hdwallet.BIP49Purpose)
	}

	return  hdwallet.Purpose(hdwallet.DefaultPurpose)
}

