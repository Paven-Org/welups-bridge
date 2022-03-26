package main

import (
	"bridge/common/consts"
	"bridge/micros/weleth/config"
	"bridge/micros/weleth/dao"
	"bridge/micros/weleth/service"
	manager "bridge/service-managers"
	ethListener "bridge/service-managers/listener/eth"
	welListener "bridge/service-managers/listener/wel"
	"bridge/service-managers/logger"
	"context"
	"sync"
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
	welEthDAO := dao.MkWelEthTransDao(db)

	ctx := context.Background()

	wg := sync.WaitGroup{}

	// ETH chain stuff: contract address, prkey, contract event watcher...
	ethClient, err := ethclient.Dial(config.Get().EtherumConf.BlockchainRPC)
	if err != nil {
		logger.Err(err).Msg("[main] Etherum client initialization failed")
		panic(err)
	}
	defer ethClient.Close()
	ethSysDAO := dao.MkEthSysDao(db)
	ethListen := ethListener.NewEthListener(ethSysDAO, ethClient, config.Get().EtherumConf.BlockTime, config.Get().EtherumConf.BlockOffSet, logger)

	ethEvtConsumer := service.NewEthConsumer(config.Get().EthContractAddress[0], welEthDAO)
	ethListen.RegisterConsumer(ethEvtConsumer)

	wg.Add(1)
	go func() {
		ethListen.Start(ctx)
		wg.Done()
	}()

	// WEL chain stuff
	var welClient = &welListener.ExtNodeClient{}
	for _, addr := range config.Get().WelupsConf.Nodes {
		welClient = welListener.NewExtNodeClientWithTimeout(addr, time.Duration(config.Get().WelupsConf.ClientTimeout)*time.Minute)
		if err := welClient.Start(); err == nil {
			break
		} else {
			logger.Err(err).Msg("Can't start wel listener")
		}
	}
	defer welClient.Stop()

	welTransHandler := welListener.NewTransHandler(welClient, config.Get().WelupsConf.BlockOffSet)
	welSysDAO := dao.MkWelSysDao(db)
	welListen := welListener.NewWelListener(welSysDAO, welTransHandler, config.Get().WelupsConf.BlockTime, config.Get().WelupsConf.BlockOffSet, logger)

	welEvtConsumer := service.NewWelConsumer(config.Get().WelContractAddress[0], welEthDAO)
	welListen.RegisterConsumer(welEvtConsumer)

	wg.Add(1)
	go func() {
		welListen.Start(ctx)
		wg.Done()
	}()

	logger.Info().Msg("[main] Waiting for daemons to stop...")
	wg.Wait()
	logger.Info().Msg("[main] Closing weleth process, cleaning up...")
	//// Temporal
	//c, err := manager.MkTemporalClient(config.Get().TemporalCliConfig)
	//if err != nil {
	//	logger.Get().Err(err).Msgf("Unable to connect to temporal backend")
	//	panic(err)
	//}
	//defer c.Close()

	//// Should adapt to use with daemon manager
	//ctx, cancel := context.WithCancel(ctx)
	//wg := sync.WaitGroup{}
	//// might spawn multiple temporal workers here
	//wg.Add(1)
	//go func() {
	//	if err := manager.SpawnTemporalWorker(ctx, c, welethQueue, worker.Options{}, RegisterWelethBridgeService); err != nil {
	//		logger.Get().Err(err).Msg("Unable to spawn worker")
	//	}
	//	wg.Done()
	//	return
	//}()
	//////blocking on some channel then cancel()
	//wg.Wait()
	// system validity check

	//// layer 2 setup:
	// load handlers for HTTP server
	// load handlers for GRPC server/client

}
