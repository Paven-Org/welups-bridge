package libs

import (
	"crypto/sha256"

	"github.com/ethereum/go-ethereum/crypto"
)

func Sign(payload []byte, keyhex string) ([]byte, error) {
	prikey, err := crypto.HexToECDSA(keyhex)
	if err != nil {
		return nil, err
	}

	h256h := sha256.New()
	h256h.Write(payload)
	hash := h256h.Sum(nil)
	signature, err := crypto.Sign(hash, prikey)
	if err != nil {
		return nil, err
	}
	return signature, nil
}
