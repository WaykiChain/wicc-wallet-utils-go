package commons

import "github.com/btcsuite/btcd/chaincfg"

var WaykiTestNetParams = chaincfg.Params{
	PubKeyHashAddrID: 0x87,                            // 1
	PrivateKeyID:     0xd2,                            // 5(uncompressed) or K (compressed)
	HDPrivateKeyID:   [4]byte{0x04, 0x88, 0xad, 0xe4}, //xprv
	HDPublicKeyID:    [4]byte{0x04, 0x88, 0xb2, 0x1e}, //xpub
	HDCoinType:       99999,
}

var WaykiMainNetParams =  chaincfg.Params{
	PubKeyHashAddrID: 0x49,                            // 1
	PrivateKeyID:     0x99,                            // 5(uncompressed) or K (compressed)
	HDPrivateKeyID:   [4]byte{0x04, 0x88, 0xad, 0xe4}, //xprv
	HDPublicKeyID:    [4]byte{0x04, 0x88, 0xb2, 0x1e}, //xpub
	HDCoinType:       99999,
}