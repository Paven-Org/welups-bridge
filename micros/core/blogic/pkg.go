package blogic

import (
	"bridge/libs"
	welABI "bridge/micros/core/abi/wel"
	ethLogic "bridge/micros/core/blogic/eth"
	userLogic "bridge/micros/core/blogic/user"
	welLogic "bridge/micros/core/blogic/wel"
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

	WelInquirer *welABI.WelInquirer
}

func Init(iv InitV) {
	userLogic.Init(iv.DAOs, iv.RedisManager, iv.TokenService)
	ethLogic.Init(iv.DAOs, iv.TemporalCli)
	welLogic.Init(iv.DAOs, iv.TemporalCli, iv.WelInquirer)
}
