package main

import (
	"bridge/common/consts"
	"bridge/libs"
	userLogic "bridge/micros/core/blogic/user"
	"bridge/micros/core/config"
	"bridge/micros/core/dao"
	router "bridge/micros/core/http"
	"bridge/micros/core/middlewares"
	manager "bridge/service-managers"
	"bridge/service-managers/logger"
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/casbin/casbin/v2"
	_ "github.com/lib/pq"
	//"https://github.com/rs/zerolog/log"
)

func main() {
	manager.SetOSParams()
	//ctx := context.Background()

	config.Load()
	//// layer 1 setup: foundation
	// logger
	logger.Init(config.Get().Structured)
	logger := logger.Get()
	logger.Info().Msgf("[main] Initialize system with config: %+v ", *config.Get())
	defer logger.Info().Msg("[main] Core exited")

	//	ctx = logger.WithContext(ctx)
	//	zerolog.Ctx(ctx).Info().Msgf("getting log from context: ", ctx)
	// loading config, secret, key
	// DB
	db, err := manager.MkDB(config.Get().DBconfig)
	if err != nil {
		logger.Err(err).Msg("[main] DB initialization failed")
		panic(err)
	}
	if db == nil {
		logger.Err(consts.ErrNilDB).Msg("[main] DB initialization failed")
		panic(err)
	}
	defer func() {
		for {
			if err := db.Close(); err != nil {
				logger.Err(err).Msg("[main] Unable to close db, attempt to retry...")
			} else {
				logger.Info().Msg("[main] db closed successfully")
				break
			}
		}
	}()

	// daos
	daos := dao.MkDAOs(db)

	// Redis
	rm := manager.MkRedisManager(
		config.Get().RedisConfig,
		manager.StdDbMap)
	defer func() {
		if config.Get().HttpConfig.Mode == "debug" {
			rm.Flush(manager.StdAuthDBName)
		}
		rm.CloseAll()
	}()

	// RBAC enforcer
	enforcer, err := casbin.NewEnforcer(config.Get().Casbin.ModelPath, config.Get().Casbin.PolicyPath)
	if err != nil {
		logger.Err(err).Msg("[main] constructing casbin enforcer failed")
		return
	}
	authMW := middlewares.MkAuthMW(enforcer, rm)
	// token service
	ts := libs.MkTokenServ(config.Get().Secrets.JwtSecret)

	// Mailer

	// ETH chain stuff: contract address, prkey, contract event watcher...

	// WEL chain stuff

	// Core business logic init
	userLogic.Init(daos, rm, ts)

	// Temporal

	/// HTTP server
	// Router setup
	// middlewares: TLS, CORS, JWT, secure cookie, json resp body, URL normalization...
	mainRouter := router.InitMainRouter(config.Get().HttpConfig, authMW)
	httpServ := manager.MkHttpServer(config.Get().HttpConfig, mainRouter)
	go func() {
		if err := httpServ.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Err(err).Msg("[main] Failed to start HTTP server")
			return
		}
	}()
	// system validity check

	// shutdown & cleanup
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)
	defer stop()
	<-ctx.Done()
	stop()

	logger.Info().Msg("[main] Shutting down HTTP server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServ.Shutdown(ctx); err != nil {
		logger.Err(err).Msg("[main] Failed to gracefully shutdown HTTP server")
	}

}
