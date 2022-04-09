package welethService

import (
	"bridge/micros/weleth/dao"
	"bridge/micros/weleth/model"
	"bridge/service-managers/logger"
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"gitlab.com/rwxrob/uniq"
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

	CreateW2ECashinClaimRequest = "CreateW2ECashinClaimRequest"
	UpdateClaimWelCashinEth     = "UpdateClaimWelCashinEth"

	CreateE2WCashoutClaimRequest = "CreateE2WCashoutClaimRequest"
	UpdateClaimEthCashoutWel     = "UpdateClaimEthCashoutWel"
)

type BridgeTx struct {
	FromChainTxHash     string
	FromChainTokenAddr  string
	ToChainTokenAddr    string
	ToChainReceiverAddr string
	RequestID           string
	Amount              string
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

func (s *WelethBridgeService) GetWelToEthCashinByTxHash(ctx context.Context, txhash string) (tx model.WelCashinEthTrans, err error) {
	log := logger.Get()
	log.Info().Msgf("[W2E transaction get] getting cashin transaction")
	ct, err := s.CashinTransDAO.SelectTransByDepositTxHash(txhash)
	if err != nil {
		log.Err(err).Msg("[W2E transaction get] failed to get cashin transaction: " + txhash)
		return
	}
	return *ct, nil
}

func (s *WelethBridgeService) CreateW2ECashinClaimRequest(ctx context.Context, cashinTxHash string, inTokenAddr string, userAddr string, amount string, contractVersion string) (tx model.WelCashinEthTrans, err error) {
	log := logger.Get()
	log.Info().Msgf("[W2E claim request] getting cashin transaction")
	ct, err := s.CashinTransDAO.SelectTransByDepositTxHash(cashinTxHash)
	if err != nil {
		log.Err(err).Msg("[W2E claim request] failed to get cashin transaction: " + cashinTxHash)
		return
	}
	switch ct.ClaimStatus {
	case model.StatusSuccess:
		err = model.ErrAlreadyClaimed
		log.Err(err).Msgf("[W2E claim request] %s already claimed " + cashinTxHash)
		return

	case model.StatusPending:
		err = model.ErrRequestPending
		log.Err(err).Msgf("[W2E claim request] %s already pending for a request" + cashinTxHash)
		return
	case model.StatusUnknown:
		tx = *ct
		// validate
		if tx.EthWalletAddr != userAddr {
			err = fmt.Errorf("Inconsistent receiver address: %s != %s", userAddr, tx.EthWalletAddr)
			log.Err(err).Msg("[W2E claim request] Inconsistent request")
			return model.WelCashinEthTrans{}, err
		}
		if tx.EthTokenAddr != inTokenAddr {
			err = fmt.Errorf("Inconsistent receiver address: %s != %s", inTokenAddr, tx.EthTokenAddr)
			log.Err(err).Msg("[W2E claim request] Inconsistent request")
			return model.WelCashinEthTrans{}, err
		}
		if tx.Amount != amount {
			err = fmt.Errorf("Inconsistent receiver address: %s != %s", amount, tx.Amount)
			log.Err(err).Msg("[W2E claim request] Inconsistent request")
			return model.WelCashinEthTrans{}, err
		}

		tx.ReqID = crypto.Keccak256Hash(uniq.Bytes(32)).Big().String()
		if err := s.CashinTransDAO.CreateClaimRequest(tx.ReqID, tx.ID, model.StatusPending); err != nil {
			log.Err(err).Msgf("[W2E claim request] couldn't create claim request for %s" + cashinTxHash)
			return model.WelCashinEthTrans{}, err
		}
	default:
		err = model.ErrUnrecognizedStatus
		log.Err(err).Msgf("[W2E claim request] unrecognized claim request status for %s" + cashinTxHash)
		return
	}
	return
}

func (s *WelethBridgeService) UpdateClaimWelCashinEth(ctx context.Context, id int64, reqID string, reqStatus string, claimTxHash string, status string) error {
	log := logger.Get()
	log.Info().Msgf("[W2E update claim request] updating cashin transaction")
	err := s.CashinTransDAO.UpdateClaimWelCashinEth(id, reqID, reqStatus, claimTxHash, status)
	if err != nil {
		log.Err(err).Msg("[W2E update claim request] failed to update cashin request ")
		return err
	}
	return nil
}

func (s *WelethBridgeService) GetEthToWelCashoutByTxHash(ctx context.Context, txhash string) (tx model.EthCashoutWelTrans, err error) {
	log := logger.Get()
	log.Info().Msgf("[E2W transaction get] getting cashout transaction")
	ct, err := s.CashoutTransDAO.SelectTransByDepositTxHash(txhash)
	if err != nil {
		log.Err(err).Msg("[E2W transaction get] failed to get cashout transaction: " + txhash)
		return
	}
	return *ct, nil
}

func (s *WelethBridgeService) CreateE2WCashoutClaimRequest(ctx context.Context, cashoutTxHash string, outTokenAddr string, userAddr string, amount string, contractVersion string) (tx model.EthCashoutWelTrans, err error) {
	log := logger.Get()
	log.Info().Msgf("[E2W claim request] getting cashout transaction")
	ct, err := s.CashoutTransDAO.SelectTransByDepositTxHash(cashoutTxHash)
	if err != nil {
		log.Err(err).Msg("[E2W claim request] failed to get cashout transaction: " + cashoutTxHash)
		return
	}
	switch ct.ClaimStatus {
	case model.StatusSuccess:
		err = model.ErrAlreadyClaimed
		log.Err(err).Msgf("[E2W claim request] %s already claimed " + cashoutTxHash)
		return

	case model.StatusPending:
		err = model.ErrRequestPending
		log.Err(err).Msgf("[E2W claim request] %s already pending for a request" + cashoutTxHash)
		return
	case model.StatusUnknown:
		tx = *ct
		// validate
		if tx.WelWalletAddr != userAddr {
			err = fmt.Errorf("Inconsistent receiver address: %s != %s", userAddr, tx.WelWalletAddr)
			log.Err(err).Msg("[E2W claim request] Inconsistent request")
			return model.EthCashoutWelTrans{}, err
		}
		if tx.WelTokenAddr != outTokenAddr {
			err = fmt.Errorf("Inconsistent receiver address: %s != %s", outTokenAddr, tx.WelTokenAddr)
			log.Err(err).Msg("[E2W claim request] Inconsistent request")
			return model.EthCashoutWelTrans{}, err
		}
		if tx.Amount != amount {
			err = fmt.Errorf("Inconsistent receiver address: %s != %s", amount, tx.Amount)
			log.Err(err).Msg("[E2W claim request] Inconsistent request")
			return model.EthCashoutWelTrans{}, err
		}

		tx.ReqID = crypto.Keccak256Hash(uniq.Bytes(32)).Big().String()
		if err := s.CashoutTransDAO.CreateClaimRequest(tx.ReqID, tx.ID, model.StatusPending); err != nil {
			log.Err(err).Msgf("[E2W claim request] couldn't create claim request for %s" + cashoutTxHash)
			return model.EthCashoutWelTrans{}, err
		}
	default:
		err = model.ErrUnrecognizedStatus
		log.Err(err).Msgf("[E2W claim request] unrecognized claim request status for %s" + cashoutTxHash)
		return
	}
	return
}

func (s *WelethBridgeService) UpdateClaimEthCashoutWel(ctx context.Context, id int64, reqID string, reqStatus string, claimTxHash string, fee string, status string) error {
	log := logger.Get()
	log.Info().Msgf("[E2W update claim request] updating cashin transaction")
	err := s.CashoutTransDAO.UpdateClaimEthCashoutWel(id, reqID, reqStatus, claimTxHash, fee, status)
	if err != nil {
		log.Err(err).Msg("[E2W update claim request] failed to update cashin request ")
		return err
	}
	return nil
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

	w.RegisterActivityWithOptions(s.CreateW2ECashinClaimRequest, activity.RegisterOptions{Name: CreateW2ECashinClaimRequest})
	w.RegisterActivityWithOptions(s.UpdateClaimWelCashinEth, activity.RegisterOptions{Name: UpdateClaimWelCashinEth})

	w.RegisterActivityWithOptions(s.CreateE2WCashoutClaimRequest, activity.RegisterOptions{Name: CreateE2WCashoutClaimRequest})
	w.RegisterActivityWithOptions(s.UpdateClaimEthCashoutWel, activity.RegisterOptions{Name: UpdateClaimEthCashoutWel})
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
