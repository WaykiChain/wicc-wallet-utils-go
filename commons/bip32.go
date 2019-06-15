// Package bip32 provides ...
package commons

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"github.com/WaykiChain/wicc-wallet-utils-go/commons/bytes"
	"github.com/WaykiChain/wicc-wallet-utils-go/commons/hash"
	"github.com/WaykiChain/wicc-wallet-utils-go/commons/base58"
	"github.com/WaykiChain/wicc-wallet-utils-go/commons/ec"
	"github.com/btcsuite/btcd/chaincfg"
)

// https://github.com/btcsuite/btcutil/blob/master/hdkeychain/extendedkey.go
// https://github.com/tyler-smith/go-bip32/blob/master/bip32.go

const (
	//HardenedKeyStart hardended key starts.
	HardenedKeyStart uint32= 0x80000000 // 2^31

	// version(4 bytes) || depth(1 byte) || parent fingerprint(4 bytes) || child number(4 bytes) || chaincode(32 bytes) || pub/pri key data(33 bytes)
	serializedKeyLen = 78

	// max Depth
	maxDepth = 0xFF

	minSeedBytes = 16 //128 bits

	maxSeedBytes = 64 // 512 bits
)

var (
	// ErrInvalidSeedLen seed Len error
	ErrInvalidSeedLen = errors.New("seed lenght must be between 128 and 512 bits")
	// ErrUnusableSeed describes an error in which the provided seed is not
	// usable due to the derived key falling outside of the valid range for
	// secp256k1 private keys.  This error indicates the caller must choose
	// another seed.
	ErrUnusableSeed = errors.New("unusable seed")

	// ErrKeyByteSize error extended key bytes
	ErrKeyByteSize = errors.New("error extended key bytes")

	//ErrDeriveBeyondMaxDepth max 255 indices
	ErrDeriveBeyondMaxDepth = errors.New("cannot derive a key with more than 255 indices in its path")

	// ErrDeriveHardFromPublic cannot derive a hardened key from a public key
	ErrDeriveHardFromPublic = errors.New("cannot derive a hardened key from a public key")

	//ErrInvalidChild  child index invalid
	ErrInvalidChild = errors.New("the extended key at this index is invalid")

	// ErrNotPrivExtKey  not Private Key
	ErrNotPrivExtKey = errors.New("cant't create private keys from public extended key")
)

var (
	masterKey = []byte("Bitcoin seed")
)

// ExtendedKey private/public key data
type ExtendedKey struct {
	coinParams     * chaincfg.Params
	version        []byte // 4 bytes
	depth          byte   // 1 byte
	parentFP       []byte // 4 bytes
	childNum       uint32 // 4 bytes
	chainCode      []byte // 32 bytes
	key            []byte // will be the pubkey for extended pub keys
	pubKey         []byte // only for extended pri keys
	isPrivate      bool
	DerivationPath string
}

//NewMasterKey create a new master key data from seed
func NewMasterKey(seed []byte, params * chaincfg.Params) (*ExtendedKey, error) {
	if len(seed) < minSeedBytes || len(seed) > maxSeedBytes {
		return nil, ErrInvalidSeedLen
	}
	// I = HMAC-SHA512(Key = "Bitcoin seed", Data = S)
	hmac512 := hmac.New(sha512.New, masterKey)
	hmac512.Write(seed)
	lr := hmac512.Sum(nil)

	// Split "I" into two 32-byte sequences Il and Ir where:
	//   Il = master secret key
	//   Ir = master chain code
	secretKey := lr[:32]
	chainCode := lr[32:]

	secretKeyNum := new(big.Int).SetBytes(secretKey)
	if secretKeyNum.Cmp(ec.Secp265k1().Params().N) >= 0 || secretKeyNum.Sign() == 0 {
		return nil, ErrUnusableSeed
	}

	return &ExtendedKey{
		coinParams:     params,
		version:        params.HDPrivateKeyID[:],
		depth:          0,
		parentFP:       []byte{0x00, 0x00, 0x00, 0x00},
		childNum:       0,
		chainCode:      chainCode,
		key:            secretKey,
		isPrivate:      true,
		DerivationPath: "m",
	}, nil
}

// pubKeyBytes returns bytes for the serialized compressed public key associated
// with this extended key in an efficient manner including memoization as
// necessary.
//
// When the extended key is already a public key, the key is simply returned as
// is since it's already in the correct form.  However, when the extended key is
// a private key, the public key will be calculated and memoized so future
// accesses can simply return the cached result.
func (key *ExtendedKey) pubKeyBytes() []byte {
	// Just return the key if it's already an extended public key.
	if !key.isPrivate {
		return key.key
	}

	// This is a private extended key, so calculate and memoize the public
	// key if needed.
	if len(key.pubKey) == 0 {
		_, pubKey := ec.PrivKeyFromBytes(key.key)
		key.pubKey = pubKey.SerializeCompressed()
	}
	return key.pubKey
}

// HardenedChild derivation hardened child
func (key *ExtendedKey) HardenedChild(i uint32) (*ExtendedKey, error) {
	i += HardenedKeyStart
	return key.Child(i)
}

// Child create extended child key
// https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki#child-key-derivation-ckd-functions
func (key *ExtendedKey) Child(i uint32) (*ExtendedKey, error) {
	if key.depth == maxDepth {
		return nil, ErrDeriveBeyondMaxDepth
	}

	// four scenarios
	// #1 Private parent key -> private child key (hardened)
	// #2 Private parent key -> private child key (normal)
	// #3 Public parent key -> public child key (normal)
	// #4 Public parent key -> public child key (hardened) failure!!

	isHardenedChild := i >= HardenedKeyStart

	// #4
	if !key.isPrivate && isHardenedChild {
		return nil, ErrDeriveHardFromPublic
	}

	// data:
	// #1 hardened: 0x00 || ser256(parentKey) || ser32(i)
	//     1+32+4 = 37
	// #2 normal: serP(point(parentKey)) || ser32(i)
	//   point(parentKey) = parentPubKey
	// #3 normal: serP(parentPubKey) || ser32(i)
	//   P=(x,y)  serP(P) = (0x02 or 0x03) || ser256(x) = Compressed PubKey
	// 33 +4 = 37
	keyIdentifier := ""
	data := make([]byte, 37)
	if isHardenedChild {
		// #1
		copy(data[1:], key.key)
		keyIdentifier = fmt.Sprintf("%s/%d'", key.DerivationPath, i%HardenedKeyStart)
	} else {
		// #2 #3
		copy(data, key.pubKeyBytes())
		keyIdentifier = fmt.Sprintf("%s/%d", key.DerivationPath, i)
	}
	binary.BigEndian.PutUint32(data[33:], i)

	// I = HMAC-SHA512(Key = chainCode, Data=data)
	hmac512 := hmac.New(sha512.New, key.chainCode)
	hmac512.Write(data)
	ilr := hmac512.Sum(nil)

	il := ilr[:32]
	childChainCode := ilr[32:]

	ilNum := new(big.Int).SetBytes(il)
	if ilNum.Cmp(ec.Secp265k1().Params().N) >= 0 || ilNum.Sign() == 0 {
		return nil, ErrInvalidChild
	}

	var isPrivate bool
	var childKey []byte
	if key.isPrivate {
		// #1 #2
		// child key ki is parse256(IL) + kpar (mod n)
		keyNum := new(big.Int).SetBytes(key.key)
		ilNum.Add(ilNum, keyNum)
		ilNum.Mod(ilNum, ec.Secp265k1().Params().N)
		// childKey = ilNum.Bytes()
		childKey = bytes.PaddedBytes(32, ilNum.Bytes())
		isPrivate = true
	} else {
		// #3
		// child key Ki is point(parse256(IL)) + Kpar.
		ilx, ily := ec.Secp265k1().ScalarBaseMult(il)
		if ilx.Sign() == 0 || ily.Sign() == 0 {
			return nil, ErrInvalidChild
		}

		parentPubKey, err := ec.ParsePubKey(key.key)
		if err != nil {
			return nil, err
		}

		childX, childY := ec.Secp265k1().Add(ilx, ily, parentPubKey.X, parentPubKey.Y)
		pk := ec.PublicKey{X: childX, Y: childY}
		childKey = pk.SerializeCompressed()
	}
	parentFP := hash.Hash160(key.pubKeyBytes())[:4]
	return &ExtendedKey{
		coinParams:     key.coinParams,
		version:        key.version,
		depth:          key.depth + 1,
		parentFP:       parentFP,
		childNum:       i,
		chainCode:      childChainCode,
		key:            childKey,
		pubKey:         nil,
		isPrivate:      isPrivate,
		DerivationPath: keyIdentifier,
	}, nil
}

// Neuter https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki#private-parent-key--public-child-key
func (key *ExtendedKey) Neuter() *ExtendedKey {
	// N((k, c)) → (K, c)
	// N(CKDpriv((kpar, cpar), i)) (works always).
	// CKDpub(N(kpar, cpar), i) (works only for non-hardened child keys).
	if !key.isPrivate {
		return key
	}
	return &ExtendedKey{
		coinParams:     key.coinParams,
		version:        key.coinParams.HDPublicKeyID[:],
		depth:          key.depth,
		parentFP:       key.parentFP,
		childNum:       key.childNum,
		chainCode:      key.chainCode,
		key:            key.pubKeyBytes(),
		pubKey:         nil,
		isPrivate:      false,
		DerivationPath: key.DerivationPath,
	}
}

// Address return pay-to-pubkey-has (P2PKH) address
func (key *ExtendedKey) Address() (*ec.AddressPubKeyHash, error) {
	pkHash := hash.Hash160(key.pubKeyBytes())
	return ec.NewAddressPubKeyHash(pkHash, key.coinParams.PubKeyHashAddrID)
}

func (key *ExtendedKey) ECPrivKey() (*ec.PrivateKey, error) {
	if !key.isPrivate {
		return nil, ErrNotPrivExtKey
	}

	privKey, _ := ec.PrivKeyFromBytes(key.key)
	return privKey, nil
}

func (key *ExtendedKey) Serialize() ([]byte, error) {
	if len(key.key) == 0 {
		return nil, ErrKeyByteSize
	}

	var childNumBytes [4]byte
	binary.BigEndian.PutUint32(childNumBytes[:], key.childNum)

	// The serialized format is:
	//   version (4) || depth (1) || parent fingerprint (4)) ||
	//   child num (4) || chain code (32) || key data (33) || checksum (4)
	serializedBytes := make([]byte, 0, serializedKeyLen+4)
	serializedBytes = append(serializedBytes, key.version...)
	serializedBytes = append(serializedBytes, key.depth)
	serializedBytes = append(serializedBytes, key.parentFP...)
	serializedBytes = append(serializedBytes, childNumBytes[:]...)
	serializedBytes = append(serializedBytes, key.chainCode...)
	if key.isPrivate {
		serializedBytes = append(serializedBytes, 0x00)
		serializedBytes = bytes.PaddedAppend(serializedBytes, 32, key.key)
	} else {
		serializedBytes = append(serializedBytes, key.pubKeyBytes()...)
	}
	checkSum := hash.DoubleHash256(serializedBytes)[:4]
	serializedBytes = append(serializedBytes, checkSum...)
	return serializedBytes, nil
}

func (key *ExtendedKey) B58Serialize() string {
	serializeKey, err := key.Serialize()
	if err != nil {
		return ""
	}
	return base58.Encode(serializeKey)
}
