package main

import (
	"bridge/common/consts"
	"bridge/micros/weleth/config"
	manager "bridge/service-managers"
	"bridge/service-managers/logger"

	_ "github.com/lib/pq"

	"github.com/ethereum/go-ethereum/ethclient"
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
		for err := db.Close(); err != nil; {
			logger.Err(err).Msg("[main] Unable to close db, attempt to retry...")
		}
		logger.Info().Msg("[main] db closed successfully")
	}()
	// Redis

	// Message queue

	// Mailer

	// ETH chain stuff: contract address, prkey, contract event watcher...
	ethClient, err := ethclient.Dial(config.Get().EtherumConf.BlockchainRPC)
	if err != nil {
		logger.Err(err).Msg("[main] Etherum client initialization failed")
		panic(err)
	}
	defer ethClient.Close()

	// WEL chain stuff

	// GRPC server/client

	// system validity check

	//// layer 2 setup:
	// load handlers for HTTP server
	// load handlers for GRPC server/client

}
