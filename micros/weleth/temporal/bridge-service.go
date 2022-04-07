package welethService

import (
	"bridge/micros/weleth/dao"
	"bridge/service-managers/logger"
	"context"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

// API
const (
	WelethServiceQueue = "TEMPORAL_BRIDGE_QUEUE_WELETH"

	// Main API

	// Some elaboration:
	// ChainA -`CASHIN`-> ChainB <=> method `withdraw` was called on ChainA's Export
	// contract, waiting to be claimed by user as 1:1 equivalent (wrapped) tokens on ChainB,
	// via calling `claim` method on ChainB's Import contract.
	// ChainB -`CASHOUT`-> ChainA <=> method `withdraw` was called on ChainB's Import
	// contract, waiting to be claimed by user as original tokens/currencies on ChainA, via
	// calling `claim` method on ChainA's Export contract.
	GetWelToEthCashinByTxHash  = "WEL2ETH_CASHIN"  // original Wel values -> wrapped Eth tokens
	GetEthToWelCashoutByTxHash = "ETH2WEL_CASHOUT" // wrapped Eth tokens -> original Wel values

	GetEthToWelCashinByTxHash  = "ETH2WEL_CASHIN"  // original Eth values -> wrapped Wel tokens
	GetWelToEthCashoutByTxHash = "WEL2ETH_CASHOUT" // wrapped Wel tokens -> original Eth values

)

type BridgeTx struct {
	TriggerChainTxHash     string
	OtherChainToTokenAddr  string
	OtherChainReceiverAddr string
	Amount                 string
}

type WelethBridgeService struct {
	CashinTransDAO  dao.IWelCashinEthTransDAO
	CashoutTransDAO dao.IEthCashoutWelTransDAO
	tempCli         client.Client
	worker          worker.Worker
}

// Service implementation
func MkWelethBridgeService(cli client.Client, daos *dao.DAOs) *WelethBridgeService {
	return &WelethBridgeService{
		CashinTransDAO:  daos.WelCashinEthTransDAO,
		CashoutTransDAO: daos.EthCashoutWelTransDAO,
		tempCli:         cli,
	}
}

func (s *WelethBridgeService) GetWelToEthCashinByTxHash(ctx context.Context, txhash string) (tx BridgeTx, err error) {

	return
}

func (s *WelethBridgeService) GetEthToWelCashoutByTxHash(ctx context.Context, txhash string) (tx BridgeTx, err error) {

	return
}

func (s *WelethBridgeService) GetEthToWelCashinByTxHash(ctx context.Context, txhash string) (tx BridgeTx, err error) {
	// NOT IMPLEMENTED
	return
}

func (s *WelethBridgeService) GetWelToEthCashoutByTxHash(ctx context.Context, txhash string) (tx BridgeTx, err error) {
	// NOT IMPLEMENTED
	return
}

func (s *WelethBridgeService) registerService(w worker.Worker) {
	w.RegisterActivityWithOptions(s.GetWelToEthCashinByTxHash, activity.RegisterOptions{Name: GetWelToEthCashinByTxHash})
	w.RegisterActivityWithOptions(s.GetEthToWelCashoutByTxHash, activity.RegisterOptions{Name: GetEthToWelCashoutByTxHash})

	w.RegisterActivityWithOptions(s.GetEthToWelCashinByTxHash, activity.RegisterOptions{Name: GetEthToWelCashinByTxHash})
	w.RegisterActivityWithOptions(s.GetWelToEthCashoutByTxHash, activity.RegisterOptions{Name: GetWelToEthCashoutByTxHash})

}

func (s *WelethBridgeService) StartService() error {
	w := worker.New(s.tempCli, WelethServiceQueue, worker.Options{})
	s.registerService(w)

	s.worker = w
	logger.Get().Info().Msgf("Starting WelethBridgeService")
	if err := w.Start(); err != nil {
		logger.Get().Err(err).Msgf("Error while starting WelethBridgeService")
		return err
	}

	logger.Get().Info().Msgf("WelethBridgeService started")
	return nil
}

func (s *WelethBridgeService) StopService() {
	if s.worker != nil {
		s.worker.Stop()
	}
}
