package main

import (
	"bridge/common/consts"
	"bridge/libs"
	"bridge/micros/core/blogic"
	"bridge/micros/core/config"
	"bridge/micros/core/dao"
	router "bridge/micros/core/http"
	"bridge/micros/core/microservices/weleth/mswelethImp"
	"bridge/micros/core/middlewares"
	ethService "bridge/micros/core/service/eth"
	ethMulsend "bridge/micros/core/service/eth/mulsend"
	"bridge/micros/core/service/notifier"
	welService "bridge/micros/core/service/wel"
	importcontract "bridge/micros/core/service/wel/import-contract"
	manager "bridge/service-managers"
	ethListener "bridge/service-managers/listener/eth"
	welListener "bridge/service-managers/listener/wel"
	"bridge/service-managers/logger"
	"context"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	patchedWelclient "github.com/Paven-Org/gotron-sdk/pkg/client"
	welclient "github.com/Paven-Org/gotron-sdk/pkg/client"

	"github.com/casbin/casbin/v2"
	"github.com/ethereum/go-ethereum/ethclient"
	_ "github.com/lib/pq"
	//"https://github.com/rs/zerolog/log"
)

func main() {
	manager.SetOSParams()
	//ctx := context.Background()

	config.Load()

	cnf := config.Get()
	//// layer 1 setup: foundation
	// logger
	logger.Init(cnf.Structured)
	logger := logger.Get()
	//logger.Info().Msgf("[main] Initialize system with config: %+v ", *cnf)
	defer logger.Info().Msg("[main] Core exited")

	//	ctx = logger.WithContext(ctx)
	//	zerolog.Ctx(ctx).Info().Msgf("getting log from context: ", ctx)
	// loading config, secret, key
	// DB
	db, err := manager.MkDB(cnf.DBconfig)
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
		cnf.RedisConfig,
		manager.StdDbMap)
	defer func() {
		if cnf.HttpConfig.Mode == "debug" {
			rm.Flush(manager.StdAuthDBName)
		}
		rm.CloseAll()
	}()

	// token service
	ts := libs.MkTokenServ(cnf.Secrets.JwtSecret)

	// Mailer
	mailer := manager.MkMailer(cnf.Mailerconf)

	// Temporal
	tempCli, err := manager.MkTemporalClient(cnf.TemporalCliConfig, []string{"callerkey", "signerkey"})
	if err != nil {
		logger.Err(err).Msg("[main] Unable to connect to Temporal cluster")
		return
	}
	defer tempCli.Close()

	// ETH chain stuff: contract address, prkey, contract event watcher...
	ethCli, err := ethclient.Dial(cnf.EthereumConfig.BlockchainRPC)
	if err != nil {
		logger.Err(err).Msgf("Unable to connect to ethereum RPC server")
		return
	}
	defer ethCli.Close()

	ethGovService, err := ethService.MkGovContractService(ethCli, tempCli, daos, cnf.EthGovContract)
	if err != nil {
		logger.Err(err).Msgf("Unable to initialize ethererum GovContractService")
		return
	}

	ethGovService.StartService()
	defer ethGovService.StopService()

	ethMulsendService, err := ethMulsend.MkMulsendContractService(ethCli, tempCli, daos, cnf.EthMulsendContract)
	if err != nil {
		logger.Err(err).Msgf("Unable to initialize ethererum MulsendContractService")
		return
	}

	ethMulsendService.StartService()
	defer ethMulsendService.StopService()

	// WEL chain stuff
	welCli := welclient.NewGrpcClient(cnf.WelupsConfig.Nodes[0])
	defer welCli.Stop()
	if err := welCli.Start(); err != nil {
		logger.Err(err).Msgf("Unable to start welCli's GRPC connection")
		return
	}

	welGovService, err := welService.MkGovContractService(welCli, tempCli, daos, cnf.WelGovContract)
	if err != nil {
		logger.Err(err).Msgf("Unable to initialize welups GovContractService")
		return
	}

	welGovService.StartService()
	defer welGovService.StopService()

	patchedWelCli := patchedWelclient.NewGrpcClient(cnf.WelupsConfig.Nodes[0])
	defer patchedWelCli.Stop()
	if err := patchedWelCli.Start(); err != nil {
		logger.Err(err).Msgf("Unable to start patchedWelCli's GRPC connection")
		return
	}
	welImportService, err := importcontract.MkImportContractService(patchedWelCli, tempCli, daos, cnf.WelImportContract)
	if err != nil {
		logger.Err(err).Msgf("Unable to initialize welups WelImportService")
		return
	}

	welImportService.StartService()
	defer welImportService.StopService()
	// Bridge microservices
	//weleth
	msWelEth := mswelethImp.MkWeleth(tempCli)
	msWelEth.StartService()
	defer msWelEth.StopService()

	notifierS := notifier.MkNotifier(tempCli, daos, mailer)
	notifierS.StartService()
	defer notifierS.StopService()

	// Core business logic init
	initVector := blogic.InitV{
		DAOs:         daos,
		RedisManager: rm,
		Mailer:       mailer,
		//Httpcli: nil,
		TokenService: ts,
		TemporalCli:  tempCli,
		WelCli:       welCli,
		EthCli:       ethCli,
	}

	blogic.Init(initVector)

	// sync stuff
	wg := sync.WaitGroup{}
	ctx := context.Background()

	/// event listeners
	// ethereum side
	ethblockdao := daos.EthBlockDAO
	ethListen := ethListener.NewEthListener(ethblockdao, ethCli, cnf.EthereumConfig.BlockTime, cnf.EthereumConfig.BlockOffSet, logger)

	ethGovEvConsumer := ethService.NewGovEvConsumer(cnf.EthGovContract, daos, tempCli)
	ethListen.RegisterConsumer(ethGovEvConsumer)

	wg.Add(1)
	go func() {
		ethListen.Start(ctx)
		wg.Done()
	}()

	// welups side
	welExtcli := welListener.NewExtNodeClientFromCli(welCli, time.Minute*time.Duration(cnf.WelupsConfig.ClientTimeout))
	welTransHandler := welListener.NewTransHandler(welExtcli, cnf.WelupsConfig.BlockOffSet)
	welblockdao := daos.WelBlockDAO
	welListen := welListener.NewWelListener(welblockdao, welTransHandler, cnf.WelupsConfig.BlockTime, cnf.WelupsConfig.BlockOffSet, logger)

	welEvtConsumer := welService.NewGovEvConsumer(cnf.WelGovContract, daos, tempCli)
	welListen.RegisterConsumer(welEvtConsumer)

	wg.Add(1)
	go func() {
		welListen.Start(ctx)
		wg.Done()
	}()

	/// HTTP server
	// RBAC enforcer
	enforcer, err := casbin.NewEnforcer(cnf.Casbin.ModelPath, cnf.Casbin.PolicyPath)
	if err != nil {
		logger.Err(err).Msg("[main] constructing casbin enforcer failed")
		return
	}
	authMW := middlewares.MkAuthMW(enforcer, rm)
	// Router setup
	// middlewares: TLS, CORS, JWT, secure cookie, json resp body, URL normalization...
	mainRouter := router.InitMainRouter(cnf.HttpConfig, authMW)
	httpServ := manager.MkHttpServer(cnf.HttpConfig, mainRouter)
	go func() {
		if err := httpServ.ListenAndServeTLS(cnf.HttpConfig.X509CertFile, cnf.HttpConfig.X509KeyFile); err != nil && err != http.ErrServerClosed {
			logger.Err(err).Msg("[main] Failed to start HTTP server")
			return
		}
	}()
	// system validity check

	// shutdown & cleanup
	logger.Info().Msg("[main] Waiting for daemons to stop...")

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)
	defer stop()
	<-ctx.Done()
	stop()

	logger.Info().Msg("[main] Shutting down HTTP server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServ.Shutdown(ctx); err != nil {
		logger.Err(err).Msg("[main] Failed to gracefully shutdown HTTP server")
	}

	wg.Wait()
	logger.Info().Msg("[main] Closing core process, cleaning up...")

}
