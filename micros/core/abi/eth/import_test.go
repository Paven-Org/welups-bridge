package eth

import (
	"bridge/libs"
	"bridge/service-managers/logger"
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	importC *EthImportC
)

func TestClaim(t *testing.T) {
	tokenaddr := "0x4272ffC0682d68aCF5eEbD2ABFDc38d721BCF55a"
	signerkey := "ce0d51b2062e5694d28a21ad64b7efd583856ba20afe437ae4c4ad7d7a5ae34a"

	//toaddr := "0x0b49cfbc07542c39d95a6b079b0e821e2cbfbb1e5c4b3a6e85fc562d590b8de6"
	//prikey := "AC91B3A0E2EDB0C692D753018277D6D1869242F6666A3B58B58F7593E8A0CE35"
	//reqID := "104554513604985853153126454866643575863446308692812117996278932301483034333893"
	//amount := "99"

	prikey := "ce0d51b2062e5694d28a21ad64b7efd583856ba20afe437ae4c4ad7d7a5ae34a"
	reqID := "33520334248965224490069560844488943606812912433996205144170613492011902220912"
	amount := "0"
	contractVersion := "IMPORTS_ETH_v1"

	ctx := context.Background()
	pkey, err := crypto.HexToECDSA(prikey)
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
	caller := crypto.PubkeyToAddress(pkey.PublicKey)
	address := caller.Hex()
	log.Info().Msgf("[Eth logic internal] caller address: %s", address)

	_requestID := &big.Int{}
	_requestID.SetString(reqID, 10)

	_amount := &big.Int{}
	_amount.SetString(amount, 10)

	signature, err := libs.StdSignedMessageHash(tokenaddr, address, _amount, _requestID, contractVersion, signerkey)
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
	//fmt.Printf("signature hex: %x\n", signature)

	nonce, err := ethCli.PendingNonceAt(ctx, caller)
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}

	gasPrice, err := ethCli.SuggestGasPrice(context.Background())
	if err != nil {
		logger.Get().Err(err).Msgf("Unable to get recommended gas price, set to default")
		gasPrice = big.NewInt(3000000000)
	}

	opts := bind.NewKeyedTransactor(pkey)
	opts.GasLimit = uint64(300000)
	opts.Value = big.NewInt(0)
	opts.GasPrice = gasPrice
	opts.Nonce = big.NewInt(int64(nonce))

	tokenAddr := common.HexToAddress(tokenaddr)
	tx, err := importC.EthImportCTransactor.Claim(opts, tokenAddr, _requestID, _amount, signature)

	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
	fmt.Println("tx: ", tx)
}
