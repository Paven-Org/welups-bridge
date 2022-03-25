package wel

import (
	"bridge/micros/core/config"
	"bridge/service-managers/logger"
	"fmt"
	"testing"

	welclient "github.com/Clownsss/gotron-sdk/pkg/client"
	"github.com/rs/zerolog"
)

var (
	inq      *WelInquirer
	log      *zerolog.Logger
	testAddr = "WDReBjymEzH5Bi4avyfUqXa4rhyvAT7DY2"
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

	m.Run()

}

func TestInq(t *testing.T) {
	balance, err := inq.GetNativeBalance(testAddr)
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
	fmt.Println(balance)

	tBalance, err := inq.WRC20balanceOf("WXXybedJRgXd6G675VFaGo14U6YzzfGY9A", testAddr)
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
	fmt.Println(tBalance.String())
}
