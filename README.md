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

- GetPrivateKeyFromMnemonic. You should provide your mnemonic and network Type (1 MainNet,2 TestNet),function return private Key,MainNet Private key start with "P" ,TestNet
private key start with "Y".

- GetPubKeyFromPrivateKey. you should provide your Private Key,the function return wallet public key as hex string.

- GetAddressFromPrivateKey. you should provide your Private Key,the function return wallet Address as base58 encode string,MainNet Address start with "W",TestNet Address start with "w".

- GetAddressFromMnemonic. you should provide your mnemonic,the function return wallet Address as base58 encode string,MainNet Address start with "W",TestNet Address start with "w".

Examples:

Generate mnemonic:
```go
mnemonic := GenerateMnemonics()
```
Get private key from mnemonic
```go
mnemonic := "empty regular curve turtle student prize toy accuse develop spike scatter ginger"
privateKey := GetPrivateKeyFromMnemonic(mnemonic, WAYKI_MAINTNET)
```
Get public key from private key
```go
publicKey,_:=GetPubKeyFromPrivateKey("Y6J4aK6Wcs4A3Ex4HXdfjJ6ZsHpNZfjaS4B9w7xqEnmFEYMqQd13")
```