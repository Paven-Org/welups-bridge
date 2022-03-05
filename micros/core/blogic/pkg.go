package blogic

import (
	"bridge/libs"
	userLogic "bridge/micros/core/blogic/user"
	"bridge/micros/core/dao"
	manager "bridge/service-managers"
)

type InitV struct {
	DAOs         *dao.DAOs
	RedisManager *manager.RedisManager
	Mailer       *manager.Mailer
	Httpcli      *manager.HttpClient
	TokenService libs.ITokenService
}

func Init(iv InitV) {
	userLogic.Init(iv.DAOs, iv.RedisManager, iv.TokenService)
}
