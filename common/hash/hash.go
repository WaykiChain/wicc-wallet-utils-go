// Package utils provides ...
package hash

import "crypto/sha256"
import "golang.org/x/crypto/ripemd160"

// Hash256 sha256(b)
func Hash256(bytes []byte) []byte {
	hash := sha256.Sum256(bytes)
	return hash[:]
}

// DoubleHash256 double sha256(b)
func DoubleHash256(bytes []byte) []byte {
	first := sha256.Sum256(bytes)
	second := sha256.Sum256(first[:])
	return second[:]
}

// Hash160  https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki#key-identifiers
// Hash160 (RIPEMD160 after SHA256)
func Hash160(bytes []byte) []byte {
	hash1 := sha256.Sum256(bytes)
	return hashRipeMD160(hash1[:])
}

func hashRipeMD160(data []byte) []byte {
	hasher := ripemd160.New()
	hasher.Write(data)
	return hasher.Sum(nil)
}
