# wicc-wallet-utils-go
WaykiChain Wallet Library for Offline Operation written in golang

## Install

* go get github.com/WaykiChain/wicc-wallet-utils-go
* go get github.com/btcsuite

## Usage

### Create Wallet
Generate mnemonics and private keys to manage your wallet.

```go
func GenerateMnemonics() string

func GetPrivateKeyFromMnemonic(words string, netType int) string

func GetPubKeyFromPrivateKey(privKey string) (string,error) 

func GetAddressFromPrivateKey(privateKey string, netType int) string

func GetAddressFromMnemonic(words string, netType int) string
```
- GenerateMnemonics. You will get 12 words.
- GetPrivateKeyFromMnemonic. You should provide your mnemonic and network Type (1 MainNet 2 TestNet),then return private Key,Mainnet Private key start with "P" ,Testnet
private key start with "Y".
- GetPubKeyFromPrivateKey. You will get 12 words.
- GetAddressFromPrivateKey. You will get 12 words.
- GetAddressFromMnemonic. You will get 12 words.