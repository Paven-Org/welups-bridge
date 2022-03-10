package main

import (
	"bridge/common/consts"
	"bridge/micros/weleth/config"
	"bridge/micros/weleth/dao"
	manager "bridge/service-managers"
	ethListener "bridge/service-managers/listener/eth"
	welListener "bridge/service-managers/listener/wel"
	"bridge/service-managers/logger"
	"context"
	"time"

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

	// create parent context
	ctx := context.Background()

	// ETH chain stuff: contract address, prkey, contract event watcher...
	ethClient, err := ethclient.Dial(config.Get().EtherumConf.BlockchainRPC)
	if err != nil {
		logger.Err(err).Msg("[main] Etherum client initialization failed")
		panic(err)
	}
	defer ethClient.Close()
	ethSysDAO := dao.MkEthSysDao(db)
	ethListen := ethListener.NewEthListener(ethSysDAO, ethClient, config.Get().EtherumConf.BlockTime, config.Get().EtherumConf.BlockOffSet, logger)
	ethListen.Start(ctx)

	// WEL chain stuff
	var welClient = &welListener.ExtNodeClient{}
	for _, addr := range config.Get().WelupsConf.Nodes {
		welClient := welListener.NewExtNodeClientWithTimeout(addr, time.Duration(config.Get().WelupsConf.ClientTimeout)*time.Minute)
		if welClient.Start() == nil {
			break
		}
	}
	defer welClient.Stop()

	welTransHandler := welListener.NewTransHandler(welClient, config.Get().WelupsConf.BlockOffSet)
	welSysDAO := dao.MkWelSysDao(db)
	welListen := welListener.NewWelListener(welSysDAO, welTransHandler, config.Get().WelupsConf.BlockTime, config.Get().WelupsConf.BlockOffSet, logger)
	welListen.Start(ctx)

	// GRPC server/client

	// system validity check

	//// layer 2 setup:
	// load handlers for HTTP server
	// load handlers for GRPC server/client

}
