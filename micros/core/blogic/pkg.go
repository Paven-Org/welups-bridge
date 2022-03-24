package blogic

import (
	"bridge/libs"
	ethLogic "bridge/micros/core/blogic/eth-account"
	userLogic "bridge/micros/core/blogic/user"
	welLogic "bridge/micros/core/blogic/wel-account"
	"bridge/micros/core/dao"
	manager "bridge/service-managers"

	"go.temporal.io/sdk/client"
)

type InitV struct {
	DAOs         *dao.DAOs
	RedisManager *manager.RedisManager
	Mailer       *manager.Mailer
	Httpcli      *manager.HttpClient
	TokenService libs.ITokenService
	TemporalCli  client.Client
}

func Init(iv InitV) {
	userLogic.Init(iv.DAOs, iv.RedisManager, iv.TokenService)
	ethLogic.Init(iv.DAOs, iv.Mailer, iv.TemporalCli)
	welLogic.Init(iv.DAOs, iv.Mailer, iv.TemporalCli)
}
