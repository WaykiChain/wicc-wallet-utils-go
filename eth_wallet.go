package wicc_wallet_utils_go

import (
	"encoding/hex"
	wicc_common "github.com/WaykiChain/wicc-wallet-utils-go/common"
	"github.com/WaykiChain/wicc-wallet-utils-go/common/ec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var(
	ETH 		 ="ETH"
	ETHLedger 	 ="ETHLedger"

	ETHW 			  *ETHWallet
	ETHLedgerW  	  *ETHWallet
)


func init(){
	ETHW = NewETHWallet(ETH)
	ETHLedgerW = NewETHWallet(ETHLedger)
}


//ETHWallet
type ETHWallet struct{
	wallet *wicc_common.Wallet
	mnemonicLen int
}

func NewETHWallet(wc string) *ETHWallet{

	newWallet := ETHWallet{}
	switch wc {
	case ETH:
		newWallet =  ETHWallet{
			wallet : wicc_common.NewWallet(wicc_common.ETH,false, false,nil),
			mnemonicLen : 12,
		}
	case ETHLedger:
		newWallet =  ETHWallet{
			wallet : wicc_common.NewWallet(wicc_common.ETH,false, false,nil),
			mnemonicLen : 12,
		}
	default:
		newWallet =  ETHWallet{
			wallet : wicc_common.NewWallet(wicc_common.ETH,false, false,nil),
			mnemonicLen : 12,
		}
	}
	return  &newWallet
}

func (ETHw *ETHWallet) GenerateAddressFromMnemonic(mnemonic,language string) (string, error) {
	return ETHw.wallet.GenerateAddressFromMnemonic(mnemonic,language)
}

func (ETHw *ETHWallet) GenerateAddressFromPrivateKey(privateKey string) (string, error){

	privateKeyBytes,_ := hex.DecodeString(privateKey)
	// (1) new  *PrivateKey、 *PublicKey
	_, publicKey:= ec.PrivKeyFromBytes(privateKeyBytes)
	if  publicKey == nil{
		return "", nil
	}
	/*publicKey Compressed 33 bytes
	publicKeyStr := hex.EncodeToString(publicKey.SerializeCompressed())*/

	// (2) pubBytes为04 开头的65字节公钥,去掉04后剩下64字节进行Keccak256运算
	pubBytes := crypto.Keccak256(publicKey.SerializeUnCompressed()[1:])
	// (3) 经过Keccak256运算后变成32字节，最终取这32字节的后20字节作为真正的地址
	address := common.BytesToAddress(pubBytes[12:])

	return address.Hex(), nil
}

func (ETHw *ETHWallet) ExportPrivateKeyFromMnemonic(mnemonic ,language string ) (string,error) {
	return ETHw.wallet.ExportPrivateKeyFromMnemonic(mnemonic ,language)
}

func (ETHw *ETHWallet) CheckAddress(address string) (bool,error){
	//去掉0x（如有）
	rm0xaddr := wicc_common.RemoveOxFromHex(address)
	//判断长度
	if len(rm0xaddr) != wicc_common.ETHADDRESSLEN{
		return false, wicc_common.ERR_ADDRESS_LEN
	}
	//判断stringTohex是否成功
	_,err := hex.DecodeString(rm0xaddr)
	if err != nil{
		return false, wicc_common.ERR_INVALID_ADDRESS
	}
	return true, nil
}

func (ETHw *ETHWallet) CheckPrivateKey(privateKey string) (bool,error){
	//去掉0x（如有）
	rm0xaddr := wicc_common.RemoveOxFromHex(privateKey)
	//判断长度
	if len(rm0xaddr) != wicc_common.ETHPRIVATEKEYLEN{
		return false, wicc_common.ERR_INVALID_PRIVATEKEY_LEN
	}
	//判断stringTohex是否成功
	_,err := hex.DecodeString(rm0xaddr)
	if err != nil{
		return false, wicc_common.ERR_INVALID_PRIVATEKEY
	}
	return true, nil
}

func (ETHw *ETHWallet) GetPubKeyFromPrivateKey(privateKey string) (string, error) {

	isValid, err := ETHw.CheckPrivateKey(privateKey)
	if isValid == false || err != nil{
		return "", err
	}

	privateKeyBytes, _ := hex.DecodeString(privateKey)
	// (1) new  *PrivateKey、 *PublicKey
	_, publicKey := ec.PrivKeyFromBytes(privateKeyBytes)
	if publicKey == nil {
		return "", nil
	}
	//publicKey Compressed 33 bytes
	return hex.EncodeToString(publicKey.SerializeCompressed()), nil
}




