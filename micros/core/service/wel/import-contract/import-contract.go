package importcontract

import (
	"bridge/libs"
	welImport "bridge/micros/core/abi/wel"
	welLogic "bridge/micros/core/blogic/wel"
	"bridge/micros/core/dao"
	welDAO "bridge/micros/core/dao/wel-account"
	welService "bridge/micros/core/service/wel"
	welethModel "bridge/micros/weleth/model"
	welethService "bridge/micros/weleth/temporal"
	"bridge/service-managers/logger"
	"context"
	"math/big"
	"sort"
	"time"

	welclient "github.com/Clownsss/gotron-sdk/pkg/client"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

//const (
//	ImportContractQueue = "WelImportContractService"
//
//	WithdrawWorkflow = "Withdraw"
//	IssueWorkflow    = "Issue"
//
//	// signal
//	BatchIssueSignal = "BatchedIssueSignal"
//
//	BatchIssueID = "BatchIssueWFOnlyInstance"
//)

type ImportContractService struct {
	imp             *welImport.WelImport
	dao             welDAO.IWelDAO
	cli             *welclient.GrpcClient
	tempCli         client.Client
	worker          worker.Worker
	defaultFeelimit int64
	batchIssueID    string
	batchIssueRunID string
}

const (
	ImportContractQueue = welService.ImportContractQueue

	Withdraw = welService.Withdraw
	Issue    = welService.Issue

	WatchForTx2TreasuryWF = welService.WatchForTx2TreasuryWF
	// signal
	BatchIssueSignal = welService.BatchIssueSignal

	BatchIssueID = welService.BatchIssueID
)

func MkImportContractService(client *welclient.GrpcClient, tempCli client.Client, daos *dao.DAOs, contractAddr string) (*ImportContractService, error) {
	imp := welImport.MkWelImport(client, contractAddr)

	return &ImportContractService{cli: client, tempCli: tempCli, imp: imp, dao: daos.Wel, defaultFeelimit: 8000000}, nil
}

func (ctr *ImportContractService) Issue(ctx context.Context, tokenAddr string, receivers []string, values []*big.Int) (string, error) {
	callerkey, err := welLogic.GetAuthenticatorKey()
	if err != nil {
		logger.Get().Err(err).Msgf("Unable to get authenticator's key")
		return "", err
	}
	pkey, err := crypto.HexToECDSA(callerkey)
	if err != nil {
		logger.Get().Err(err).Msgf("Unable to parse hexstring to ECDSA key")
		return "", err
	}

	caller, err := libs.KeyToB58Addr(callerkey)
	if err != nil {
		logger.Get().Err(err).Msgf("Unable to parse caller address")
		return "", err
	}
	opts := &welImport.CallOpts{
		From:      caller,
		Prikey:    pkey,
		Fee_limit: ctr.defaultFeelimit,
		T_amount:  0,
	}

	tx, err := ctr.imp.Issue(opts, tokenAddr, receivers, values)
	if err != nil {
		logger.Get().Err(err).Msgf("Unable to trigger import contract")
		return "", err
	}
	logger.Get().Info().Msgf("Contract call done with tx: %+v", tx)
	return common.Bytes2Hex(tx.GetTxid()), nil
}

//type txQueue = []welethModel.EthCashinWelTrans
type txQueue struct {
	lastIssuance time.Time
	queue        []welethModel.EthCashinWelTrans
}

func (ctr *ImportContractService) BatchIssueWF(ctx workflow.Context) error {
	ao := workflow.ActivityOptions{
		TaskQueue:              ImportContractQueue,
		ScheduleToCloseTimeout: time.Second * 60,
		ScheduleToStartTimeout: time.Second * 60,
		StartToCloseTimeout:    time.Second * 60,
		HeartbeatTimeout:       time.Second * 10,
		WaitForCancellation:    false,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumInterval: time.Second * 100,
			MaximumAttempts: 10,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	log := workflow.GetLogger(ctx)
	// iteration no.
	iterN := 0
	done := false
	//lastIssuance := workflow.Now(ctx) // last issuance, zeroth = now when there's no issuance already done
	// channel to receive tx to batch
	signalChan := workflow.GetSignalChannel(ctx, BatchIssueSignal)
	//txQueue := make([]welethModel.EthCashinWelTrans, 0)
	allTxQueues := make(map[string]*txQueue)

	// selector
	// main loop
	for !done {
		iterN++

		selector := workflow.NewSelector(ctx)
		selector.AddReceive(ctx.Done(), func(channel workflow.ReceiveChannel, more bool) { done = true })
		selector.AddReceive(signalChan, func(channel workflow.ReceiveChannel, more bool) {
			var tx = welethModel.EthCashinWelTrans{}
			channel.Receive(ctx, &tx)
			welToken := tx.WelTokenAddr
			_, ok := allTxQueues[welToken]
			if !ok {
				allTxQueues[welToken] = &txQueue{
					lastIssuance: workflow.Now(ctx), // zeroth lastIssuance set to the time the first tx of this token arrived
					queue:        make([]welethModel.EthCashinWelTrans, 0),
				}
			}
			//txQueue = append(txQueue, tx)
			allTxQueues[welToken].queue = append(allTxQueues[welToken].queue, tx)

			if len(allTxQueues[welToken].queue) >= 16 {
				receivers := libs.Map(
					func(tx welethModel.EthCashinWelTrans) string {
						return tx.WelWalletAddr
					}, allTxQueues[welToken].queue)
				values := libs.Map(
					func(tx welethModel.EthCashinWelTrans) *big.Int {
						ret := &big.Int{}
						ret.SetString(tx.Total, 10)
						return ret
					}, allTxQueues[welToken].queue)
				// issue
				log.Info("Contract call...")
				var txhash string
				res := workflow.ExecuteActivity(ctx, ctr.Issue, welToken, receivers, values)
				if err := res.Get(ctx, &txhash); err != nil {
					log.Error("Failed to call issue on import contract")
				} else {
					log.Info("Contract call succeeded")
				}
				// update e2wcashin txs with wel issue txhash
				for _, tran := range allTxQueues[welToken].queue {
					res = workflow.ExecuteActivity(ctx, welethService.UpdateEthCashinWelTrans, tran)
					if err := res.Get(ctx, nil); err != nil {
						log.Error("Failed to update E2W cashin trans")
					} else {
						log.Info("update E2W cashin trans succeeded")
					}
				}

				// update lastIssuance
				allTxQueues[welToken].lastIssuance = workflow.Now(ctx)
				//lastIssuance = allTxQueues[welToken].lastIssuance
				//reset txqueue
				allTxQueues[welToken].queue = make([]welethModel.EthCashinWelTrans, 0)
			}
		})

		//workflow.Sleep(ctx, 3*time.Minute)
		//selector.AddFuture(workflow.NewTimer(ctx, lastIssuance.Add(2*time.Minute).Sub(workflow.Now(ctx))), func(f workflow.Future) {
		selector.AddFuture(workflow.NewTimer(ctx, 1*time.Minute), func(f workflow.Future) {
			// deterministism
			welTokens := []string{}
			for k, _ := range allTxQueues {
				welTokens = append(welTokens, k)
			}
			sort.Strings(welTokens)

			for _, welToken := range welTokens {
				if workflow.Now(ctx).Sub(allTxQueues[welToken].lastIssuance).Seconds() > 120.0 {
					receivers := libs.Map(
						func(tx welethModel.EthCashinWelTrans) string {
							return tx.WelWalletAddr
						}, allTxQueues[welToken].queue)
					values := libs.Map(
						func(tx welethModel.EthCashinWelTrans) *big.Int {
							ret := &big.Int{}
							ret.SetString(tx.Total, 10)
							return ret
						}, allTxQueues[welToken].queue)
					// issue
					log.Info("Contract call...")
					var txhash string
					res := workflow.ExecuteActivity(ctx, ctr.Issue, welToken, receivers, values)
					if err := res.Get(ctx, &txhash); err != nil {
						log.Error("Failed to call issue on import contract")
					} else {
						log.Info("Contract call succeeded")
					}
					// update e2wcashin txs with wel issue txhash
					for _, tran := range allTxQueues[welToken].queue {
						res = workflow.ExecuteActivity(ctx, welethService.UpdateEthCashinWelTrans, tran)
						if err := res.Get(ctx, nil); err != nil {
							log.Error("Failed to update E2W cashin trans")
						} else {
							log.Info("update E2W cashin trans succeeded")
						}
					}

					// update lastIssuance
					allTxQueues[welToken].lastIssuance = workflow.Now(ctx)
					//lastIssuance = allTxQueues[welToken].lastIssuance
					//reset txqueue
					allTxQueues[welToken].queue = make([]welethModel.EthCashinWelTrans, 0)
				}
			}
		})

		// select...
		selector.Select(ctx)
		if done {
			break
		}

		if iterN > 6000 {
			log.Info("[BatchIssueWF] iteration number passed 6000, continuing as new WF insance...")
			return workflow.NewContinueAsNewError(ctx, ctr.BatchIssueWF)
		}
	}

	return nil
}

func (ctr *ImportContractService) WatchForTx2Treasury(ctx workflow.Context, from, to, treasury, netid, token, amount string) error {
	log := workflow.GetLogger(ctx)

	ao := workflow.ActivityOptions{
		TaskQueue:              welethService.WelethServiceQueue,
		ScheduleToCloseTimeout: time.Second * 60,
		ScheduleToStartTimeout: time.Second * 60,
		StartToCloseTimeout:    time.Second * 60,
		HeartbeatTimeout:       time.Second * 10,
		WaitForCancellation:    false,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumInterval: time.Second * 30,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	getTx2Treasury := func() error {
		var tx welethModel.TxToTreasury
		res := workflow.ExecuteActivity(ctx, welethService.GetTx2Treasury, from, treasury, token, amount)
		err := res.Get(ctx, &tx)

		if err != nil {
			log.Error("[WatchForTx2Treasury] error while getting tx2treasury", err)
			return err
		}
		cashinTx := welethModel.EthCashinWelTrans{
			EthTxHash: tx.TxID,

			EthTokenAddr: token,
			WelTokenAddr: welethModel.WelTokenFromEth[token],

			EthWalletAddr: from,
			WelWalletAddr: to,

			NetworkID: netid,
			Total:     tx.Amount,

			//CommissionFee:
			Status: welethModel.EthCashinWelUnconfirmed,
		}
		res = workflow.ExecuteActivity(ctx, welethService.CreateEthCashinWelTrans, cashinTx)
		err = res.Get(ctx, &(cashinTx.ID))
		if err != nil {
			log.Error("[WatchForTx2Treasury] error while createing W2ECashin trans", err)
			return err
		}

		se := workflow.SignalExternalWorkflow(ctx, ctr.batchIssueID, "", BatchIssueSignal, tx)
		err = se.Get(ctx, nil)
		if err != nil {
			log.Error("[WatchForTx2Treasury] error while sending tx to BatchIssue", err)
			return err
		}
		return nil
	}

	// txQueue := ...
	timer := workflow.NewTimer(ctx, 2*time.Minute) // check pending queue every 2 min
	selector := workflow.NewSelector(ctx)

	selector.AddReceive(ctx.Done(), func(channel workflow.ReceiveChannel, more bool) {})
	selector.AddFuture(timer, func(f workflow.Future) {
		if err := getTx2Treasury(); err == welethModel.ErrTx2TreasuryNotFound {
			log.Error("[WatchForTx2Treasury] tx2treasury not found")
			return
		}
		log.Info("[WatchForTx2Treasury] tx enqueued for BatchIssueWF")
	})

	// main
	if err := getTx2Treasury(); err == welethModel.ErrTx2TreasuryNotFound {
		selector.Select(ctx)
	}

	return nil
}

// Worker
func (ctr *ImportContractService) registerService(w worker.Worker) {
	//w.RegisterActivity(ctr.Withdraw)
	//w.RegisterActivity(ctr.Issue)

	//w.RegisterWorkflowWithOptions(ctr.WithdrawWorkflow, workflow.RegisterOptions{Name: WithdrawWorkflow})
	//w.RegisterWorkflowWithOptions(ctr.IssueWorkflow, workflow.RegisterOptions{Name: IssueWorkflow})
}

func (ctr *ImportContractService) StartService() error {
	w := worker.New(ctr.tempCli, ImportContractQueue, worker.Options{})
	ctr.registerService(w)

	ctr.worker = w
	logger.Get().Info().Msgf("Starting ImportContractService")
	if err := w.Start(); err != nil {
		logger.Get().Err(err).Msgf("Error while starting ImportContractService")
		return err
	}

	logger.Get().Info().Msgf("ImportContractService started")

	// start batch issue WF
	ctx := context.Background()
	wo := client.StartWorkflowOptions{
		TaskQueue: ImportContractQueue,
		ID:        BatchIssueID, // only one workflow ID allowed at all time
	}
	we, err := ctr.tempCli.ExecuteWorkflow(ctx, wo, ctr.BatchIssueWF)
	if err != nil {
		logger.Get().Err(err).Msgf("Error while starting long-running workflow BatchIssue")
		return err
	}
	ctr.batchIssueID = we.GetID()
	ctr.batchIssueRunID = we.GetRunID()

	return nil
}

func (ctr *ImportContractService) StopService() {
	if ctr.worker != nil {
		ctr.worker.Stop()
	}
}
