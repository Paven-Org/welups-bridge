package utils

import "golang.org/x/crypto/sha3"

func HashKeccak256(in []byte) []byte {
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(in)
	b := hasher.Sum(nil)
	return b
}
