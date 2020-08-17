package internal

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/jpmorganchase/quorum/crypto/secp256k1"
	"golang.org/x/crypto/sha3"
)

//TODO(cjh) these are duplicated in the quorum-account-hashicorp-vault-plugin project
// consider moving these and similar util functions to a separate project

func zeroKey(key *ecdsa.PrivateKey) {
	b := key.D.Bits()
	for i := range b {
		b[i] = 0
	}
}

func PrivateKeyToHexString(key *ecdsa.PrivateKey) (string, error) {
	byt, err := PrivateKeyToBytes(key)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(byt), nil
}

const (
	keyLen = 32 // the expected length in bytes of a S256 secp256k1 private key
)

// PrivateKeyToBytes returns the bytes for the private component of the key, and if necessary, left-0 pads them to 32 bytes.
//
// As outlined in https://github.com/openethereum/openethereum/issues/2263, 256 bit secp256k1 can generate valid keys that are shorter than 32 bytes.
// To protect against potential issues with variable lengths, key bytes should be left-0 padded.
//
// common/math.PaddedBigBytes accomplishes this using bitwise operations for maximal performance.
// Given the additional HTTP overhead already introduced when using Vault, and to improve code-readability, standard slice functions have been used here.
func PrivateKeyToBytes(key *ecdsa.PrivateKey) ([]byte, error) {
	// TODO(cjh) if performance of crypto operations needs to be improved, consider revisiting this func (see godoc)
	if key == nil {
		return nil, errors.New("nil key")
	}
	byt := key.D.Bytes()

	switch {
	case len(byt) > keyLen:
		return nil, fmt.Errorf("key cannot be longer than %v bytes", keyLen)
	case len(byt) < keyLen:
		padded := make([]byte, keyLen)
		return append(padded[:len(padded)-len(byt)], byt...), nil
	default:
		return byt, nil
	}
}

const addressLength = 20

type Address [addressLength]byte

func PrivateKeyToAddress(key *ecdsa.PrivateKey) (Address, error) {
	if key == nil || key.PublicKey.X == nil || key.PublicKey.Y == nil {
		return Address{}, errors.New("invalid key: unable to derive address")
	}
	pubBytes := elliptic.Marshal(secp256k1.S256(), key.PublicKey.X, key.PublicKey.Y)

	d := sha3.NewLegacyKeccak256()
	_, err := d.Write(pubBytes[1:])
	if err != nil {
		return Address{}, err
	}
	pubHash := d.Sum(nil)

	return NewAddress(pubHash[12:])
}

func NewAddress(byt []byte) (Address, error) {
	if len(byt) != addressLength {
		return Address{}, fmt.Errorf("account address must have length %v bytes", addressLength)
	}
	var addr Address
	copy(addr[:], byt)

	return addr, nil
}

func (a Address) ToBytes() []byte {
	return a[:]
}

// ToHexString encodes the Address as a hex string without the '0x' prefix
func (a Address) ToHexString() string {
	return hex.EncodeToString(a[:])
}
