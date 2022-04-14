package mswelethImp

import (
	ethLogic "bridge/micros/core/blogic/eth"
	welLogic "bridge/micros/core/blogic/wel"
	msweleth "bridge/micros/core/microservices/weleth"
	"bridge/micros/weleth/model"
	welethService "bridge/micros/weleth/temporal"
	"bridge/service-managers/logger"
	"context"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

const (
	TaskQueue = welethService.WelethServiceQueue

	GetWelToEthCashinByTxHash  = msweleth.GetWelToEthCashinByTxHash
	GetEthToWelCashoutByTxHash = msweleth.GetEthToWelCashoutByTxHash

	GetEthToWelCashinByTxHash  = msweleth.GetEthToWelCashinByTxHash
	GetWelToEthCashoutByTxHash = msweleth.GetWelToEthCashoutByTxHash

	CreateW2ECashinClaimRequestWF  = msweleth.CreateW2ECashinClaimRequestWF
	CreateE2WCashoutClaimRequestWF = msweleth.CreateE2WCashoutClaimRequestWF

	WaitForPendingW2ECashinClaimRequestWF  = msweleth.WaitForPendingW2ECashinClaimRequestWF
	WaitForPendingE2WCashoutClaimRequestWF = msweleth.WaitForPendingE2WCashoutClaimRequestWF
)

type Weleth struct {
	tempCli client.Client
	worker  worker.Worker
}

func MkWeleth(cli client.Client) *Weleth {
	return &Weleth{
		tempCli: cli,
	}
}

func (cli *Weleth) CreateW2ECashinClaimRequestWF(ctx workflow.Context, txhash string, userAddr string) (tx model.WelCashinEthTrans, err error) {
	log := workflow.GetLogger(ctx)
	log.Info("[Core MSWeleth] Creating cashin claim request from wel to eth with wel's side txhash: " + txhash)

	ao := workflow.ActivityOptions{
		TaskQueue:              welethService.WelethServiceQueue,
		ScheduleToCloseTimeout: time.Second * 60,
		ScheduleToStartTimeout: time.Second * 60,
		StartToCloseTimeout:    time.Second * 60,
		HeartbeatTimeout:       time.Second * 10,
		WaitForCancellation:    false,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumInterval: time.Second * 15,
			MaximumAttempts: 10,
		},
	}

	ctx = workflow.WithActivityOptions(ctx, ao)

	// call weleth
	log.Info("[Core MSWeleth] Call weleth...")
	res := workflow.ExecuteActivity(ctx, welethService.CreateW2ECashinClaimRequest, txhash, userAddr)
	if err = res.Get(ctx, &tx); err != nil {
		log.Error("[Core MSWeleth] Error while executing activity CreateW2ECashinClaimRequest in weleth microservice", err.Error())
		return
	}

	// outro...
	log.Info("[Core MSWeleth] Call weleth successfully, result: ", tx)

	return tx, nil
}

func (cli *Weleth) InvalidateW2ECashinClaim(ctx context.Context, tokenAddr, reqid string) error {
	return ethLogic.InvalidateRequestClaim(tokenAddr, "0", reqid, "EXPORT_WELUPS_v1")
}
func (cli *Weleth) WaitForPendingW2ECashinClaimRequestWF(ctx workflow.Context, txhash string) error {
	log := workflow.GetLogger(ctx)

	log.Info("[Core MSWeleth] Waiting for claim request...")
	workflow.Sleep(ctx, 3*time.Minute)

	log.Info("[Core MSWeleth] Pending duration expired, checking claim request status...")
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
	ctx = workflow.WithActivityOptions(ctx, ao)

	var tx model.WelCashinEthTrans
	res := workflow.ExecuteActivity(ctx, welethService.GetWelToEthCashinByTxHash, txhash)
	if err := res.Get(ctx, &tx); err != nil {
		log.Info("[Temporal BG] Error while processing pending claim request: ", err.Error())
		return err
	}
	if tx.ClaimStatus == model.StatusPending { // if still pending after 1 minute
		// TODO: add a deliberate fail claim contract call here to invalidate the ReqID
		//ethLogic.InvalidateRequestClaim(tx.EthTokenAddr, "0", tx.ReqID, "IMPORTS_ETH_v1")
		if err := workflow.ExecuteActivity(ctx, cli.InvalidateW2ECashinClaim, tx.WelTokenAddr, tx.ReqID).Get(ctx, nil); err != nil {
			log.Info("[Temporal BG] Error while processing pending claim request: ", err.Error())
			return err
		}
		if err := workflow.ExecuteActivity(ctx, welethService.UpdateClaimWelCashinEth, tx.ID, tx.ReqID, model.RequestExpired, tx.ClaimTxHash, model.StatusUnknown).Get(ctx, nil); err != nil {
			log.Info("[Temporal BG] Error while processing pending claim request: ", err.Error())
			return err
		}
	}
	log.Info("[Temporal] nothing to do")
	return nil
}
func (cli *Weleth) CreateE2WCashoutClaimRequestWF(ctx workflow.Context, txhash string, userAddr string) (tx model.EthCashoutWelTrans, err error) {
	log := workflow.GetLogger(ctx)
	log.Info("[Core MSWeleth] Creating cashout claim request from eth to wel with eth's side txhash: " + txhash)

	ao := workflow.ActivityOptions{
		TaskQueue:              welethService.WelethServiceQueue,
		ScheduleToCloseTimeout: time.Second * 60,
		ScheduleToStartTimeout: time.Second * 60,
		StartToCloseTimeout:    time.Second * 60,
		HeartbeatTimeout:       time.Second * 10,
		WaitForCancellation:    false,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumInterval: time.Second * 15,
			MaximumAttempts: 10,
		},
	}

	ctx = workflow.WithActivityOptions(ctx, ao)

	// call weleth
	log.Info("[Core MSWeleth] Call weleth...")
	res := workflow.ExecuteActivity(ctx, welethService.CreateE2WCashoutClaimRequest, txhash, userAddr)
	if err = res.Get(ctx, &tx); err != nil {
		log.Error("[Core MSWeleth] Error while executing activity CreateE2WCashoutClaimRequest in weleth microservice", err.Error())
		return
	}

	// outro...
	log.Info("[Core MSWeleth] Call weleth successfully, result: ", tx)

	return tx, nil
}

func (cli *Weleth) InvalidateE2WCashoutClaim(ctx context.Context, tokenAddr, reqid string) error {
	return welLogic.InvalidateRequestClaim(tokenAddr, "0", reqid, "EXPORT_WELUPS_v1")
}

func (cli *Weleth) WaitForPendingE2WCashoutClaimRequestWF(ctx workflow.Context, txhash string) error {
	log := workflow.GetLogger(ctx)

	log.Info("[Core MSWeleth] Waiting for claim request...")
	workflow.Sleep(ctx, 3*time.Minute)

	log.Info("[Core MSWeleth] Pending duration expired, checking claim request status...")
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
	ctx = workflow.WithActivityOptions(ctx, ao)

	var tx model.EthCashoutWelTrans
	res := workflow.ExecuteActivity(ctx, welethService.GetEthToWelCashoutByTxHash, txhash)
	if err := res.Get(ctx, &tx); err != nil {
		log.Info("[Temporal BG] Error while processing pending claim request: ", err.Error())
		return err
	}
	if tx.ClaimStatus == model.StatusPending { // if still pending after 1 minute
		// TODO: add a deliberate fail claim contract call here to invalidate the ReqID
		if err := workflow.ExecuteActivity(ctx, cli.InvalidateE2WCashoutClaim, tx.WelTokenAddr, tx.ReqID).Get(ctx, nil); err != nil {
			log.Info("[Temporal BG] Error while processing pending claim request: ", err.Error())
			return err
		}
		if err := workflow.ExecuteActivity(ctx, welethService.UpdateClaimEthCashoutWel, tx.ID, tx.ReqID, model.RequestExpired, tx.ClaimTxHash, tx.Fee, model.StatusUnknown).Get(ctx, nil); err != nil {
			log.Info("[Temporal BG] Error while processing pending claim request: ", err.Error())
			return err
		}
	}
	log.Info("[Temporal] nothing to do")
	return nil
}
func (cli *Weleth) GetWelToEthCashinByTxHashWF(ctx workflow.Context, txhash string) (tx welethService.BridgeTx, err error) {
	log := workflow.GetLogger(ctx)
	log.Info("[Core MSWeleth] Getting cashin transaction from wel to eth with wel's side txhash: " + txhash)

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

	ctx = workflow.WithActivityOptions(ctx, ao)

	// call weleth
	log.Info("[Core MSWeleth] Call weleth...")
	res := workflow.ExecuteActivity(ctx, welethService.GetWelToEthCashinByTxHash, txhash)
	if err = res.Get(ctx, &tx); err != nil {
		log.Error("[Core MSWeleth] Error while executing activity GetWelToEthCashinByTxHash in weleth microservice", err.Error())
		return
	}

	log.Info("[Core MSWeleth] Call weleth successfully, result: ", tx)
	return tx, nil
}

func (cli *Weleth) GetEthToWelCashoutByTxHashWF(ctx workflow.Context, txhash string) (tx welethService.BridgeTx, err error) {
	log := workflow.GetLogger(ctx)
	log.Info("[Core MSWeleth] Getting cashout transaction from eth to wel with eth's side txhash: " + txhash)

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

	ctx = workflow.WithActivityOptions(ctx, ao)

	// call weleth
	log.Info("[Core MSWeleth] Call weleth...")
	res := workflow.ExecuteActivity(ctx, welethService.GetEthToWelCashoutByTxHash, txhash)
	if err = res.Get(ctx, &tx); err != nil {
		log.Error("[Core MSWeleth] Error while executing activity GetEthToWelCashoutByTxHash in weleth microservice", err.Error())
		return
	}

	log.Info("[Core MSWeleth] Call weleth successfully, result: ", tx)
	return tx, nil
}

func (cli *Weleth) GetEthToWelCashinByTxHashWF(ctx workflow.Context, txhash string) (tx welethService.BridgeTx, err error) {
	log := workflow.GetLogger(ctx)
	log.Info("[Core MSWeleth] Getting cashin transaction from eth to wel with eth's side txhash: " + txhash)

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

	ctx = workflow.WithActivityOptions(ctx, ao)

	// call weleth
	log.Info("[Core MSWeleth] Call weleth...")
	res := workflow.ExecuteActivity(ctx, welethService.GetEthToWelCashinByTxHash, txhash)
	if err = res.Get(ctx, &tx); err != nil {
		log.Error("[Core MSWeleth] Error while executing activity GetEthToWelCashinByTxHash in weleth microservice", err.Error())
		return
	}

	log.Info("[Core MSWeleth] Call weleth successfully, result: ", tx)
	return tx, nil
}

func (cli *Weleth) GetWelToEthCashoutByTxHashWF(ctx workflow.Context, txhash string) (tx welethService.BridgeTx, err error) {
	log := workflow.GetLogger(ctx)
	log.Info("[Core MSWeleth] Getting cashout transaction from wel to eth with wel's side txhash: " + txhash)

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

	ctx = workflow.WithActivityOptions(ctx, ao)

	// call weleth
	log.Info("[Core MSWeleth] Call weleth...")
	res := workflow.ExecuteActivity(ctx, welethService.GetWelToEthCashoutByTxHash, txhash)
	if err = res.Get(ctx, &tx); err != nil {
		log.Error("[Core MSWeleth] Error while executing activity GetWelToEthCashoutByTxHash in weleth microservice", err.Error())
		return
	}

	log.Info("[Core MSWeleth] Call weleth successfully, result: ", tx)
	return tx, nil
}

func (cli *Weleth) registerService(w worker.Worker) {
	// register workflow an activities
	w.RegisterWorkflowWithOptions(cli.GetWelToEthCashinByTxHashWF, workflow.RegisterOptions{Name: GetWelToEthCashinByTxHash})
	w.RegisterWorkflowWithOptions(cli.GetEthToWelCashoutByTxHashWF, workflow.RegisterOptions{Name: GetEthToWelCashoutByTxHash})
	w.RegisterWorkflowWithOptions(cli.GetEthToWelCashinByTxHashWF, workflow.RegisterOptions{Name: GetEthToWelCashinByTxHash})
	w.RegisterWorkflowWithOptions(cli.GetWelToEthCashoutByTxHashWF, workflow.RegisterOptions{Name: GetWelToEthCashoutByTxHash})

	w.RegisterWorkflowWithOptions(cli.CreateW2ECashinClaimRequestWF, workflow.RegisterOptions{Name: CreateW2ECashinClaimRequestWF})
	w.RegisterWorkflowWithOptions(cli.CreateE2WCashoutClaimRequestWF, workflow.RegisterOptions{Name: CreateE2WCashoutClaimRequestWF})

	w.RegisterActivity(cli.InvalidateE2WCashoutClaim)
	w.RegisterWorkflowWithOptions(cli.WaitForPendingW2ECashinClaimRequestWF, workflow.RegisterOptions{Name: WaitForPendingW2ECashinClaimRequestWF})

	w.RegisterActivity(cli.InvalidateW2ECashinClaim)
	w.RegisterWorkflowWithOptions(cli.WaitForPendingE2WCashoutClaimRequestWF, workflow.RegisterOptions{Name: WaitForPendingE2WCashoutClaimRequestWF})

}

func (cli *Weleth) StartService() error {
	w := worker.New(cli.tempCli, TaskQueue, worker.Options{})
	cli.registerService(w)

	cli.worker = w
	logger.Get().Info().Msgf("Starting Weleth")
	if err := w.Start(); err != nil {
		logger.Get().Err(err).Msgf("Error while starting Weleth")
		return err
	}

	logger.Get().Info().Msgf("Weleth started")
	return nil
}

func (cli *Weleth) StopService() {
	if cli.worker != nil {
		cli.worker.Stop()
	}
}
