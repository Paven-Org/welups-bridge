package libs

import (
	"strings"

	"github.com/Clownsss/gotron-sdk/pkg/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func HexToB58(hex string) (string, error) {
	hexb, err := common.HexStringToBytes(hex)
	if err != nil {
		return "", err
	}
	return common.EncodeCheck(hexb), nil
}

func B58toHex(b58 string) (string, error) {
	hexb, err := common.DecodeCheck(b58)
	if err != nil {
		return "", err
	}
	return common.BytesToHexString(hexb), nil
}

func B58toStdHex(b58 string) (string, error) {
	hexb, err := common.DecodeCheck(b58)
	if err != nil {
		return "", err
	}
	hexAddr := common.BytesToHexString(hexb)
	_, hexAddr, _ = strings.Cut(hexAddr, "0x41")
	hexAddr = "0x" + hexAddr
	return hexAddr, nil
}

func KeyToB58Addr(hexkey string) (string, error) {
	key, err := crypto.HexToECDSA(hexkey)
	if err != nil {
		return "", err // invalid key
	}
	keyAddr := crypto.PubkeyToAddress(key.PublicKey)
	hexAddr := keyAddr.Hex()
	_, hexAddr, _ = strings.Cut(hexAddr, "0x")
	hexAddr = "0x" + "41" + hexAddr
	return HexToB58(hexAddr)
}

func KeyToHexAddr(hexkey string) (string, error) {
	key, err := crypto.HexToECDSA(hexkey)
	if err != nil {
		return "", err // invalid key
	}
	keyAddr := crypto.PubkeyToAddress(key.PublicKey)

	hexAddr := keyAddr.Hex()

	return hexAddr, nil
}

func KeyToWelHexAddr(hexkey string) (string, error) {
	key, err := crypto.HexToECDSA(hexkey)
	if err != nil {
		return "", err // invalid key
	}
	keyAddr := crypto.PubkeyToAddress(key.PublicKey)

	hexAddr := keyAddr.Hex()
	_, hexAddr, _ = strings.Cut(hexAddr, "0x")
	hexAddr = "0x" + "41" + hexAddr

	return hexAddr, nil
}
