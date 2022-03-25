package eth

import (
	"bridge/micros/core/config"
	"bridge/service-managers/logger"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
)

var (
	inq *EthInquirer
	log *zerolog.Logger
)

func TestMain(m *testing.M) {
	config.Load()
	cnf := config.Get()
	log = logger.Get()

	ethCli, err := ethclient.Dial(cnf.EthereumConfig.BlockchainRPC)
	if err != nil {
		logger.Get().Err(err).Msgf("Unable to connect to ethereum RPC server")
		return
	}
	defer ethCli.Close()

	inq = MkEthInquirer(ethCli)

	m.Run()

}

func TestInq(t *testing.T) {
	balance, err := inq.BalanceAt("0x25e8370E0e2cf3943Ad75e768335c892434bD090")
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
	fmt.Println(balance.String())

	balance, err = inq.BalanceOf("0x6A9865aDE2B6207dAAC49f8bCba9705dEB0B0e6D", "0x25e8370E0e2cf3943Ad75e768335c892434bD090")
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
	fmt.Println(balance.String())
}
