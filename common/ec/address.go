// Package address provides ...
package ec

import (
	"errors"

	"golang.org/x/crypto/ripemd160"
	"github.com/WaykiChain/wicc-wallet-utils-go/common/base58"
)

// AddressPubKeyHash pay-to-pubkey-hash (P2PKH)
type AddressPubKeyHash struct {
	hash    [ripemd160.Size]byte
	version byte
}

func NewAddressPubKeyHash(pkHash []byte, version byte) (*AddressPubKeyHash, error) {
	return newAddressPubKeyHash(pkHash, version)
}

func newAddressPubKeyHash(pkHash []byte, version byte) (*AddressPubKeyHash, error) {
	if len(pkHash) != ripemd160.Size {
		return nil, errors.New("pkHash must be 20 bytes")
	}
	addr := &AddressPubKeyHash{}
	addr.version = version
	copy(addr.hash[:], pkHash)
	return addr, nil
}

// EncodeAddress return P2PKH address
func (a *AddressPubKeyHash) EncodeAddress() string {
	return encodeAddress(a.hash[:], a.version)
}

// P2PKH P2SH address encoding
func encodeAddress(hash160 []byte, version byte) string {
	input := make([]byte, 21)
	input[0] = version
	copy(input[1:], hash160)
	return base58.CheckEncode(input)
}

// Hash160 return hash160
func (a *AddressPubKeyHash) Hash160() []byte {
	return a.hash[:]
}
