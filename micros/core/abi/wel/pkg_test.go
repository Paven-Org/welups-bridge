package wel

import (
	"bridge/libs"
	"bridge/micros/core/config"
	"bridge/service-managers/logger"
	"fmt"
	"math/big"
	"testing"

	welclient "github.com/Clownsss/gotron-sdk/pkg/client"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog"
)

var (
	inq      *WelInquirer
	exp      *WelExport
	log      *zerolog.Logger
	testAddr = "WDReBjymEzH5Bi4avyfUqXa4rhyvAT7DY2"
	testKey  = "ce0d51b2062e5694d28a21ad64b7efd583856ba20afe437ae4c4ad7d7a5ae34a"
)

func TestMain(m *testing.M) {
	config.Load()
	cnf := config.Get()
	log = logger.Get()

	welCli := welclient.NewGrpcClient(cnf.WelupsConfig.Nodes[0])
	defer welCli.Stop()
	if err := welCli.Start(); err != nil {
		logger.Get().Err(err).Msgf("Unable to start welCli's GRPC connection")
		return
	}

	inq = MkWelInquirer(welCli)
	exp = MkWelExport(welCli, "WUbnXM9M4QYEkksG3ADmSan2kY5xiHTr1E")

	m.Run()

}

func TestGetAccount(t *testing.T) {
	account, err := inq.GetAccount(testAddr)
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
	fmt.Println("account 1: ", account)

	account, err = inq.GetAccount("WDgqRLQ3928bWkvxE655QJasotasgA5ANg")
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
	fmt.Println("account 2: ", account)
}

func TestActivateAccount(t *testing.T) {
	err := inq.ActivateAccountIfNotExist("WDgqRLQ3928bWkvxE655QJasotasgA5ANg", testAddr, testKey)
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
}
func TestBalance(t *testing.T) {
	balance, err := inq.GetNativeBalance(testAddr)
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
	fmt.Println("balance: ", balance)

	tBalance, err := inq.WRC20balanceOf("WXXybedJRgXd6G675VFaGo14U6YzzfGY9A", testAddr)
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
	fmt.Println(tBalance.String())
}

func TestExp(t *testing.T) {
	pkey, err := crypto.HexToECDSA(testKey)
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
	//target := "0x25e8370E0e2cf3943Ad75e768335c892434bD090"
	opts := &CallOpts{
		From:      testAddr,
		Prikey:    pkey,
		Fee_limit: 8000000,
		T_amount:  1,
	}
	tokenAddr := "W9yD14Nj9j7xAB4dbGeiX9h8unkKHxuTtb"
	//testAddr := "0x25e8370E0e2cf3943Ad75e768335c892434bD090"
	//tokenAddr := "0x00000000000000000000"
	fmt.Println(tokenAddr)
	fmt.Printf("%x\n", []byte{100})
	tx, err := exp.Withdraw(opts, tokenAddr, testAddr, big.NewInt(1), 1)
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
	fmt.Println(tx.Transaction)
}
func TestClaim(t *testing.T) {
	pkey, err := crypto.HexToECDSA(testKey)
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
	caller := crypto.PubkeyToAddress(pkey.PublicKey)
	//target := "0x25e8370E0e2cf3943Ad75e768335c892434bD090"
	reqID := "33520334248965224490069560844488943606812912433996205144170613492011902220912"
	contractVersion := "EXPORT_WELS_v1"
	opts := &CallOpts{
		From:      testAddr,
		Prikey:    pkey,
		Fee_limit: 8000000,
		T_amount:  0,
	}
	tokenAddr := "W9yD14Nj9j7xAB4dbGeiX9h8unkKHxuTtb"
	_token, _ := libs.B58toHex(tokenAddr)
	testAddr := "0x25e8370E0e2cf3943Ad75e768335c892434bD090"
	//tokenAddr := "0x00000000000000000000"
	fmt.Println(tokenAddr)

	_requestID := &big.Int{}
	_requestID.SetString(reqID, 10)

	_amount := big.NewInt(0)

	signature, err := libs.StdSignedMessageHash(_token, caller.Hex(), _amount, _requestID, contractVersion, testKey)

	fmt.Printf("%x\n", []byte{100})
	tx, err := exp.Claim(opts, _token, testAddr, _requestID, _amount, signature)
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
	fmt.Println(tx.Transaction)
}
