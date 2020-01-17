// Package ec provides ...
package ec

import (
	"errors"
	"fmt"
	"math/big"
	"github.com/WaykiChain/wicc-wallet-utils-go/common/bytes"
)

const (
	PubKeyBytesLenCompressed   = 33
	PubKeyBytesLenUncompressed = 65
	PubKeyBytesLenHybrid       = 65

	pubKeyCompressed   byte = 0x2
	pubKeyUncompressed byte = 0x4
	pubKeyHybrid       byte = 0x6
)

// PublicKey struct
type PublicKey struct {
	X, Y *big.Int
}

func isOdd(a *big.Int) bool {
	return a.Bit(0) == 1
}

// decompressPoint decompresses a point on the given curve given the X point and
// the solution to use.
func decompressPoint(x *big.Int, ybit bool) (*big.Int, error) {
	// TODO: This will probably only work for secp256k1 due to
	// optimizations.

	// Y = +-sqrt(x^3 + B)
	x3 := new(big.Int).Mul(x, x)
	x3.Mul(x3, x)
	x3.Add(x3, secp256k1.Params().B)
	x3.Mod(x3, secp256k1.Params().P)

	// Now calculate sqrt mod p of x^3 + B
	// This code used to do a full sqrt based on tonelli/shanks,
	// but this was replaced by the algorithms referenced in
	// https://bitcointalk.org/index.php?topic=162805.msg1712294#msg1712294
	y := new(big.Int).Exp(x3, secp256k1.QPlus1Div4(), secp256k1.Params().P)

	if ybit != isOdd(y) {
		y.Sub(secp256k1.Params().P, y)
	}

	// Check that y is a square root of x^3 + B.
	y2 := new(big.Int).Mul(y, y)
	y2.Mod(y2, secp256k1.Params().P)
	if y2.Cmp(x3) != 0 {
		return nil, fmt.Errorf("invalid square root")
	}

	// Verify that y-coord has expected parity.
	if ybit != isOdd(y) {
		return nil, fmt.Errorf("ybit doesn't match oddness")
	}

	return y, nil
}

// ParsePubKey parses a public key for a koblitz curve from a bytestring into a
// ecdsa.Publickey, verifying that it is valid. It supports compressed,
// uncompressed and hybrid signature formats.
func ParsePubKey(pubKeyStr []byte) (key *PublicKey, err error) {
	pubkey := PublicKey{}

	if len(pubKeyStr) == 0 {
		return nil, errors.New("pubkey string is empty")
	}

	format := pubKeyStr[0]
	ybit := (format & 0x1) == 0x1
	format &= ^byte(0x1)

	switch len(pubKeyStr) {
	case PubKeyBytesLenUncompressed:
		if format != pubKeyUncompressed && format != pubKeyHybrid {
			return nil, fmt.Errorf("invalid magic in pubkey str: "+
				"%d", pubKeyStr[0])
		}

		pubkey.X = new(big.Int).SetBytes(pubKeyStr[1:33])
		pubkey.Y = new(big.Int).SetBytes(pubKeyStr[33:])
		// hybrid keys have extra information, make use of it.
		if format == pubKeyHybrid && ybit != isOdd(pubkey.Y) {
			return nil, fmt.Errorf("ybit doesn't match oddness")
		}
	case PubKeyBytesLenCompressed:
		// format is 0x2 | solution, <X coordinate>
		// solution determines which solution of the curve we use.
		/// y^2 = x^3 + Curve.B
		if format != pubKeyCompressed {
			return nil, fmt.Errorf("invalid magic in compressed "+
				"pubkey string: %d", pubKeyStr[0])
		}
		pubkey.X = new(big.Int).SetBytes(pubKeyStr[1:33])
		pubkey.Y, err = decompressPoint(pubkey.X, ybit)
		if err != nil {
			return nil, err
		}
	default: // wrong!
		return nil, fmt.Errorf("invalid pub key length %d",
			len(pubKeyStr))
	}

	if pubkey.X.Cmp(secp256k1.Params().P) >= 0 {
		return nil, fmt.Errorf("pubkey X parameter is >= to P")
	}
	if pubkey.Y.Cmp(secp256k1.Params().P) >= 0 {
		return nil, fmt.Errorf("pubkey Y parameter is >= to P")
	}
	if !secp256k1.IsOnCurve(pubkey.X, pubkey.Y) {
		return nil, fmt.Errorf("pubkey isn't on secp256k1 curve")
	}
	return &pubkey, nil
}

// SerializeCompressed serializes a public key 33-byte compressed format
func (key *PublicKey) SerializeCompressed() []byte {
	b := make([]byte, 0, PubKeyBytesLenCompressed)
	format := pubKeyCompressed
	if isOdd(key.Y) {
		format |= 0x1
	}
	b = append(b, format)
	return bytes.PaddedAppend(b, 32, key.X.Bytes())
}

// SerializeCompressed serializes a public key 33-byte compressed format
func (p *PublicKey) SerializeUnCompressed() []byte {
	b := make([]byte, 0, PubKeyBytesLenUncompressed)
	b = append(b, pubKeyUncompressed)
	b = paddedAppend(32, b, p.X.Bytes())
	return paddedAppend(32, b, p.Y.Bytes())
}

// paddedAppend appends the src byte slice to dst, returning the new slice.
// If the length of the source is smaller than the passed size, leading zero
// bytes are appended to the dst slice before appending src.
func paddedAppend(size uint, dst, src []byte) []byte {
	for i := 0; i < int(size)-len(src); i++ {
		dst = append(dst, 0)
	}
	return append(dst, src...)
}