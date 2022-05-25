package blogic

import (
	"bridge/libs"
	ethLogic "bridge/micros/core/blogic/eth"
	userLogic "bridge/micros/core/blogic/user"
	welLogic "bridge/micros/core/blogic/wel"
	"bridge/micros/core/dao"
	manager "bridge/service-managers"

	welclient "github.com/Paven-Org/gotron-sdk/pkg/client"
	"github.com/ethereum/go-ethereum/ethclient"

	"go.temporal.io/sdk/client"
)

type InitV struct {
	DAOs         *dao.DAOs
	RedisManager *manager.RedisManager
	Mailer       *manager.Mailer
	Httpcli      *manager.HttpClient
	TokenService libs.ITokenService
	TemporalCli  client.Client
	WelCli       *welclient.GrpcClient
	EthCli       *ethclient.Client
}

func Init(iv InitV) {
	userLogic.Init(iv.DAOs, iv.RedisManager, iv.TokenService)
	ethLogic.Init(iv.DAOs, iv.TemporalCli, iv.EthCli)
	welLogic.Init(iv.DAOs, iv.TemporalCli, iv.WelCli)
}
