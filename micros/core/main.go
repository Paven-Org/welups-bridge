package main

import (
	"bridge/common/consts"
	"bridge/micros/core/config"
	router "bridge/micros/core/http"
	manager "bridge/service-managers"
	"bridge/service-managers/logger"

	_ "github.com/lib/pq"
	//"https://github.com/rs/zerolog/log"
)

func main() {
	//ctx := context.Background()

	config.Load()
	//// layer 1 setup: foundation
	// logger
	logger.Init(config.Get().Structured)
	logger := logger.Get()
	logger.Info().Msgf("Initialize system with config: %+v ", *config.Get())

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
	// Redis

	// Message queue

	// Mailer

	// ETH chain stuff: contract address, prkey, contract event watcher...

	// WEL chain stuff

	// GRPC server/client

	/// HTTP server
	// Router setup
	// middlewares: TLS, CORS, JWT, secure cookie, json resp body, URL normalization...
	httpServ, err := router.MkHttpServer(&config.Get().HttpConfig)
	httpServ.ListenAndServe()
	// system validity check

	//// layer 2 setup:
	// load handlers for HTTP server
	// load handlers for GRPC server/client

}
