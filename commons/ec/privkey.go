// Package ec provides ...
package ec

import (
	"math/big"
	"github.com/WaykiChain/wicc-wallet-utils-go/commons/bytes"
)

const PrivKeyBytesLen = 32

type PrivateKey struct {
	PubKey PublicKey
	D      *big.Int
}

func PrivKeyFromBytes(bytes []byte) (*PrivateKey, *PublicKey) {
	x, y := secp256k1.ScalarBaseMult(bytes)
	privKey := &PrivateKey{
		PubKey: PublicKey{
			X: x,
			Y: y,
		},
		D: new(big.Int).SetBytes(bytes),
	}
	return privKey, &privKey.PubKey
}

func (privKey *PrivateKey) Serialize() []byte {
	return bytes.PaddedBytes(32, privKey.D.Bytes())
}
