package libs

import (
	"crypto/sha256"

	"github.com/ethereum/go-ethereum/crypto"
)

type Hasher func([]byte) []byte

func H256(payload []byte) []byte {
	h256h := sha256.New()
	h256h.Write(payload)
	return h256h.Sum(nil)
}

func HKeccak(payload []byte) []byte {
	return crypto.Keccak256Hash(payload).Bytes()
}

func NoHash(payload []byte) []byte {
	return payload
}

type Signer func([]byte, string) ([]byte, error)

var SignerH256 = MkSigner(H256)
var SignerK256 = MkSigner(HKeccak)
var SignerNoHash = MkSigner(NoHash)

func MkSigner(h Hasher) Signer {
	return func(payload []byte, keyhex string) ([]byte, error) {
		prikey, err := crypto.HexToECDSA(keyhex)
		if err != nil {
			return nil, err
		}
		hash := h(payload)
		signature, err := crypto.Sign(hash, prikey)
		if err != nil {
			return nil, err
		}
		return signature, nil
	}
}
