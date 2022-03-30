package msweleth

import (
	welethService "bridge/micros/weleth/temporal"
	"bridge/service-managers/logger"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

const (
	TaskQueue = welethService.WelethServiceQueue

	GetWelToEthCashinByTxHash  = "GetWelToEthCashinByTxHashWF"
	GetEthToWelCashoutByTxHash = "GetEthToWelCashoutByTxHashWF"

	GetEthToWelCashinByTxHash  = "GetEthToWelCashinByTxHashWF"
	GetWelToEthCashoutByTxHash = "GetWelToEthCashoutByTxHashWF"
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
	if err = res.Get(ctx, tx); err != nil {
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
	if err = res.Get(ctx, tx); err != nil {
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
	if err = res.Get(ctx, tx); err != nil {
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
	if err = res.Get(ctx, tx); err != nil {
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
