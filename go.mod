module github.com/WaykiChain/wicc-wallet-utils-go

go 1.12

replace (
	golang.org/x/crypto v0.0.0-20191227163750-53104e6ec876 => github.com/golang/crypto v0.0.0-20191227163750-53104e6ec876
	golang.org/x/crypto v0.0.0-20200115085410-6d4e4cb37c7d => github.com/golang/crypto v0.0.0-20200115085410-6d4e4cb37c7d
	golang.org/x/net v0.0.0-20180821023952-922f4815f713 => github.com/golang/net v0.0.0-20180826012351-8a410e7b638d
	golang.org/x/sync v0.0.0-20181221193216-37e7f081c4d4 => github.com/golang/sync v0.0.0-20181221193216-37e7f081c4d4
	golang.org/x/sys v0.3.0 => github.com/golang/sys v0.3.0
	golang.org/x/text v0.3.0 => github.com/golang/text v0.3.0
)

require (
	github.com/JKinGH/go-hdwallet v0.0.0-20200117082521-b9fec2768008
	github.com/astaxie/beego v1.12.0
	github.com/blocktree/go-owcdrivers v1.2.0
	github.com/blocktree/openwallet v1.7.0
	github.com/btcsuite/btcd v0.20.1-beta
	github.com/btcsuite/btcutil v1.0.1
	github.com/cpacia/bchutil v0.0.0-20181003130114-b126f6a35b6c // indirect
	github.com/ethereum/go-ethereum v1.9.9
	github.com/imroc/req v0.2.4
	github.com/shopspring/decimal v0.0.0-20200105231215-408a2507e114
	github.com/tidwall/gjson v1.3.5
	github.com/tyler-smith/go-bip39 v1.0.2
	golang.org/x/crypto v0.0.0-20200115085410-6d4e4cb37c7d
)
