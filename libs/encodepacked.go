package libs

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
)

// https://gist.github.com/miguelmota/bc4304bb21a8f4cc0a37a0f9347b8bbb
func EncodePacked(input ...[]byte) []byte {
	return bytes.Join(input, nil)
}

func EncodeString(str string) []byte {
	return []byte(str)
}

func EncodeBytesString(v string) []byte {
	decoded, err := hex.DecodeString(v)
	if err != nil {
		panic(err)
	}
	return decoded
}

func EncodeUint256(v string) []byte {
	bn := new(big.Int)
	bn.SetString(v, 10)
	return math.U256Bytes(bn)
}

func EncodeUint256Array(arr []string) []byte {
	var res [][]byte
	for _, v := range arr {
		b := EncodeUint256(v)
		res = append(res, b)
	}

	return bytes.Join(res, nil)
}

func ToEthSignedMessageHash(_token string, _user string, _amount *big.Int, _requestID *big.Int, _version string) []byte {
	// https://gist.github.com/trmaphi/04b5790328dd71693b591973e07a943a
	token := common.HexToAddress(_token).Bytes()
	user := common.HexToAddress(_user).Bytes()
	amount := _amount.Bytes()

	requestID := _requestID.Bytes()
	version := []byte(_version)

	hash := crypto.Keccak256Hash(token, user, amount, requestID, version)
	// ECDSA.toEthSignedMessageHash
	prefixedHash := crypto.Keccak256Hash(
		[]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(hash))),
		hash.Bytes(),
	).Bytes()

	return prefixedHash
}

func StdSignedMessageHash(_token string, _user string, _amount *big.Int, _requestID *big.Int, _version string, prikey string) ([]byte, error) {
	hash := ToEthSignedMessageHash(_token, _user, _amount, _requestID, _version)
	return SignerK256(hash, prikey)
}
