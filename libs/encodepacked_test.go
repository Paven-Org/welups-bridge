package libs

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func TestPacked(t *testing.T) {
	requestId, _ := big.NewInt(0).SetString("79242598130257478667448782863620113455545540517178919498485001773537412501089", 10)
	bytes := ToEthSignedMessageHash(
		"0x4272ffC0682d68aCF5eEbD2ABFDc38d721BCF55a", // token
		"0x4bb718Cb404787BF97bB012Bb08096602fb9544B", // user
		big.NewInt(99),   // amount
		requestId,        // request id
		"IMPORTS_ETH_v1", // version
	)

	// 0x10fcf96e8de6cfe6237e51af6786fa9c32d0763866f4275b16df8db8a571ccf9
	fmt.Printf("%+v\n", hexutil.Encode(bytes))
	signedMessage, _ := SignerNoHash(bytes, "ce0d51b2062e5694d28a21ad64b7efd583856ba20afe437ae4c4ad7d7a5ae34a")
	// 0xc333fded0b74fbb5a56d92262deb040c2cd3d6c241d7d1ad2e6d6a0501e56176522f6c2021c4e0298e0d5a7ec8487529a0349120c32c9331946e106b3a19f7501c
	fmt.Printf("%+v\n", hexutil.Encode(signedMessage))
}

func TestSign(t *testing.T) {
	signer := "ce0d51b2062e5694d28a21ad64b7efd583856ba20afe437ae4c4ad7d7a5ae34a"
	user := "0x5B38Da6a701c568545dCfcB03FcB875f56beddC4"
	token := "0xd8b934580fcE35a11B58C6D73aDeE468a2833fa8"
	requestID := &big.Int{}
	requestID.SetString("79242598130257478667448782863620113455545540517178919498485001773537412501089", 10)
	amount := &big.Int{}
	amount.SetString("1", 10)
	contractVersion := "IMPORTS_ETH_v1"

	signature, _ := StdSignedMessageHash(token, user, amount, requestID, contractVersion, signer)
	fmt.Println("RequestID: ", requestID.String())
	fmt.Printf("Sig: 0x%x\n", signature)
}
