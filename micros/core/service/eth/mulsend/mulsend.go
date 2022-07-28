package mulsend

import (
	"bridge/common/consts"
	"bridge/libs"
	ethMulsend "bridge/micros/core/abi/eth"
	ethLogic "bridge/micros/core/blogic/eth"
	"bridge/micros/core/config"
	"bridge/micros/core/dao"
	ethDAO "bridge/micros/core/dao/eth-account"
	ethService "bridge/micros/core/service/eth"
	welethModel "bridge/micros/weleth/model"
	welethService "bridge/micros/weleth/temporal"
	"bridge/service-managers/logger"
	"context"
	"fmt"
	"math/big"
	"sort"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

//const (
//	MulsendContractQueue = "EthMulsendContractService"
//
//	WithdrawWorkflow = "Withdraw"
//	DisperseWorkflow    = "Disperse"
//
//	// signal
//	BatchDisperseSignal = "BatchedDisperseSignal"
//
//	BatchDisperseID = "BatchDisperseWFOnlyInstance"
//)

type MulsendContractService struct {
	mulsend            *ethMulsend.EthMultiSenderC
	dao                ethDAO.IEthDAO
	cli                *ethclient.Client
	tempCli            client.Client
	worker             worker.Worker
	lastGasPrice       *big.Int
	batchDisperseID    string
	batchDisperseRunID string
}

const (
	MulsendContractQueue = ethService.MulsendContractQueue

	Disperse = ethService.Disperse

	// signal
	BatchDisperseSignal = ethService.BatchDisperseSignal

	BatchDisperseID = ethService.BatchDisperseID
)

func MkMulsendContractService(client *ethclient.Client, tempCli client.Client, daos *dao.DAOs, contractAddr string) (*MulsendContractService, error) {
	mulsend, err := ethMulsend.NewEthMultiSenderC(common.HexToAddress(contractAddr), client)
	if err != nil {
		logger.Get().Err(err).Msgf("Unable to create multisender contract interface at address %s", contractAddr)
		return nil, err
	}

	return &MulsendContractService{cli: client, tempCli: tempCli, mulsend: mulsend, dao: daos.Eth, lastGasPrice: big.NewInt(1000000000)}, nil
}

func (ctr *MulsendContractService) Disperse(ctx context.Context, tokenAddr string, receivers []string, values []*big.Int) (string, error) {
	callerkey, err := ethLogic.GetAuthenticatorKey()
	if err != nil {
		logger.Get().Err(err).Msgf("Unable to get authenticator's key")
		return "", err
	}
	pkey, err := crypto.HexToECDSA(callerkey)
	if err != nil {
		logger.Get().Err(err).Msgf("Unable to parse hexstring to ECDSA key")
		return "", err
	}
	caller := crypto.PubkeyToAddress(pkey.PublicKey)

	gasPrice, err := ctr.cli.SuggestGasPrice(context.Background())
	if err != nil {
		logger.Get().Err(err).Msgf("Unable to get recommended gas price, set to default")
		gasPrice = ctr.lastGasPrice
	}
	ctr.lastGasPrice = gasPrice

	nonce, err := ctr.cli.PendingNonceAt(ctx, caller)
	if err != nil {
		logger.Get().Err(err).Msgf("Unale to get last nonce of address %s", caller.Hex())
		return "", err
	}

	env := config.Get().Environment
	opts, err := bind.NewKeyedTransactorWithChainID(pkey, consts.EthChainFromEnv[env])
	if err != nil {
		logger.Get().Err(err).Msg("Unale to create call opts for Disperse method")
		return "", err
	}
	opts.GasLimit = 0
	opts.Value = big.NewInt(0)
	if tokenAddr == consts.EthereumTk {
		sumValues := libs.Reduce(
			func(accum *big.Int, this *big.Int) *big.Int {
				accum = accum.Add(accum, this)
				return accum
			},
			big.NewInt(0),
			values)
		opts.Value = sumValues
	}
	logger.Get().Info().Msgf("Total values: %s", opts.Value.String())
	opts.GasPrice = gasPrice
	opts.Nonce = big.NewInt(int64(nonce))

	tokenAddress := common.HexToAddress(tokenAddr)
	receiversAddress := libs.Map(
		func(addr string) common.Address {
			return common.HexToAddress(addr)
		}, receivers)
	tx, err := ctr.mulsend.Disperse(opts, tokenAddress, receiversAddress, values)
	if err != nil {
		logger.Get().Err(err).Msgf("Unable to trigger mulsend contract")
		return "", err
	}
	logger.Get().Info().Msgf("Contract call done with tx: %+v", tx)
	return tx.Hash().Hex(), nil
}

//type txQueue = []welethModel.EthCashoutEthTrans
type txQueue struct {
	lastIssuance time.Time
	queue        []welethModel.WelCashoutEthTrans
}

func (ctr *MulsendContractService) BatchDisperseWF(ctx workflow.Context) error {
	ao := workflow.ActivityOptions{
		TaskQueue:              MulsendContractQueue,
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
	signalChan := workflow.GetSignalChannel(ctx, BatchDisperseSignal)
	//txQueue := make([]welethModel.WelCashoutEthTrans, 0)
	allTxQueues := make(map[string]*txQueue)

	// selector
	// main loop
	for !done {
		iterN++

		selector := workflow.NewSelector(ctx)
		selector.AddReceive(ctx.Done(), func(channel workflow.ReceiveChannel, more bool) { done = true })
		selector.AddReceive(signalChan, func(channel workflow.ReceiveChannel, more bool) {
			var tx = welethModel.WelCashoutEthTrans{}
			channel.Receive(ctx, &tx)
			log.Info(fmt.Sprintf("BatchDisperse received tx: %+v", tx))
			ethToken := tx.EthTokenAddr
			_, ok := allTxQueues[ethToken]
			if !ok {
				allTxQueues[ethToken] = &txQueue{
					lastIssuance: workflow.Now(ctx), // zeroth lastIssuance set to the time the first tx of this token arrived
					queue:        make([]welethModel.WelCashoutEthTrans, 0),
				}
			}
			//txQueue = append(txQueue, tx)
			allTxQueues[ethToken].queue = append(allTxQueues[ethToken].queue, tx)
			log.Info(fmt.Sprintf("Current BatchDisperse TxQueues after append: %+v", allTxQueues))

			if len(allTxQueues[ethToken].queue) >= 16 {
				receivers := libs.Map(
					func(tx welethModel.WelCashoutEthTrans) string {
						return tx.EthWalletAddr
					}, allTxQueues[ethToken].queue)
				log.Info(fmt.Sprintf("BatchDisperse receivers: %+v", receivers))
				values := libs.Map(
					func(tx welethModel.WelCashoutEthTrans) *big.Int {
						ret := &big.Int{}
						ret.SetString(tx.Total, 10)
						return ret
					}, allTxQueues[ethToken].queue)
				log.Info(fmt.Sprintf("BatchDisperse values: %+v", values))
				// issue
				log.Info("Contract call...")
				var txhash string
				res := workflow.ExecuteActivity(ctx, ctr.Disperse, ethToken, receivers, values)
				if err := res.Get(ctx, &txhash); err != nil {
					log.Error("Failed to call issue on mulsend contract")
				} else {
					log.Info("Contract call succeeded")
				}
				// update e2wcashin txs with wel issue txhash
				for _, tran := range allTxQueues[ethToken].queue {
					ao := workflow.ActivityOptions{
						TaskQueue:              welethService.WelethServiceQueue,
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
					ctx := workflow.WithActivityOptions(ctx, ao)
					tran.EthDisperseTxHash = txhash
					res = workflow.ExecuteActivity(ctx, welethService.UpdateWelCashoutEthTrans, tran)
					if err := res.Get(ctx, nil); err != nil {
						log.Error("Failed to update E2W cashin trans")
					} else {
						log.Info("update E2W cashin trans succeeded")
					}
				}

				// update lastIssuance
				allTxQueues[ethToken].lastIssuance = workflow.Now(ctx)
				//lastIssuance = allTxQueues[ethToken].lastIssuance
				//reset txqueue
				allTxQueues[ethToken].queue = make([]welethModel.WelCashoutEthTrans, 0)
			}
		})

		//workflow.Sleep(ctx, 3*time.Minute)
		//selector.AddFuture(workflow.NewTimer(ctx, lastIssuance.Add(2*time.Minute).Sub(workflow.Now(ctx))), func(f workflow.Future) {
		selector.AddFuture(workflow.NewTimer(ctx, 1*time.Minute), func(f workflow.Future) {
			// deterministism
			ethTokens := []string{}
			for k, _ := range allTxQueues {
				ethTokens = append(ethTokens, k)
			}
			sort.Strings(ethTokens)

			for _, ethToken := range ethTokens {
				if workflow.Now(ctx).Sub(allTxQueues[ethToken].lastIssuance).Seconds() > 120.0 {
					if len(allTxQueues[ethToken].queue) < 1 {
						continue
					}
					receivers := libs.Map(
						func(tx welethModel.WelCashoutEthTrans) string {
							return tx.EthWalletAddr
						}, allTxQueues[ethToken].queue)
					values := libs.Map(
						func(tx welethModel.WelCashoutEthTrans) *big.Int {
							ret := &big.Int{}
							ret.SetString(tx.Total, 10)
							return ret
						}, allTxQueues[ethToken].queue)
					// issue
					log.Info("Contract call...")
					var txhash string
					res := workflow.ExecuteActivity(ctx, ctr.Disperse, ethToken, receivers, values)
					if err := res.Get(ctx, &txhash); err != nil {
						log.Error("Failed to call issue on mulsend contract")
					} else {
						log.Info("Contract call succeeded")
					}
					// update e2wcashin txs with wel issue txhash
					for _, tran := range allTxQueues[ethToken].queue {
						ao := workflow.ActivityOptions{
							TaskQueue:              welethService.WelethServiceQueue,
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
						ctx := workflow.WithActivityOptions(ctx, ao)
						tran.EthDisperseTxHash = txhash
						res = workflow.ExecuteActivity(ctx, welethService.UpdateWelCashoutEthTrans, tran)
						if err := res.Get(ctx, nil); err != nil {
							log.Error("Failed to update E2W cashin trans")
						} else {
							log.Info("update E2W cashin trans succeeded")
						}
					}

					// update lastIssuance
					allTxQueues[ethToken].lastIssuance = workflow.Now(ctx)
					//lastIssuance = allTxQueues[ethToken].lastIssuance
					//reset txqueue
					allTxQueues[ethToken].queue = make([]welethModel.WelCashoutEthTrans, 0)
				}
			}
		})

		// select...
		selector.Select(ctx)
		if done {
			break
		}

		if iterN > 6000 {
			log.Info("[BatchDisperseWF] iteration number passed 6000, continuing as new WF insance...")
			return workflow.NewContinueAsNewError(ctx, ctr.BatchDisperseWF)
		}
	}

	return nil
}

// Worker
func (ctr *MulsendContractService) registerService(w worker.Worker) {
	//w.RegisterActivity(ctr.Withdraw)
	w.RegisterActivity(ctr.Disperse)

	w.RegisterWorkflow(ctr.BatchDisperseWF)
}

func (ctr *MulsendContractService) StartService() error {
	w := worker.New(ctr.tempCli, MulsendContractQueue, worker.Options{})
	ctr.registerService(w)

	ctr.worker = w
	logger.Get().Info().Msgf("Starting MulsendContractService")
	if err := w.Start(); err != nil {
		logger.Get().Err(err).Msgf("Error while starting MulsendContractService")
		return err
	}

	logger.Get().Info().Msgf("MulsendContractService started")

	// start batch issue WF
	ctx := context.Background()
	wo := client.StartWorkflowOptions{
		TaskQueue: MulsendContractQueue,
		ID:        BatchDisperseID, // only one workflow ID allowed at all time
	}
	we, err := ctr.tempCli.ExecuteWorkflow(ctx, wo, ctr.BatchDisperseWF)
	if err != nil {
		logger.Get().Err(err).Msgf("Error while starting long-running workflow BatchDisperse")
		return err
	}
	ctr.batchDisperseID = we.GetID()
	ctr.batchDisperseRunID = we.GetRunID()

	return nil
}

func (ctr *MulsendContractService) StopService() {
	if ctr.worker != nil {
		ctr.worker.Stop()
	}
}
