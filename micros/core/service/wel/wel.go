package welService

import (
	"bridge/libs"
	welGov "bridge/micros/core/abi/wel"
	"bridge/micros/core/dao"
	welDAO "bridge/micros/core/dao/wel-account"
	"bridge/micros/core/model"
	"bridge/service-managers/logger"
	"context"
	"time"

	welclient "github.com/Clownsss/gotron-sdk/pkg/client"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

const (
	GovContractQueue = "WelGovContractService"

	// going to be deprecated
	GrantRoleWorkflow  = "GrantRole"
	RevokeRoleWorkflow = "RevokeRole"

	SaveRoleWF   = "SaveRoleWF"
	RemoveRoleWF = "RemoveRoleWF"
)

type GovContractService struct {
	gov             *welGov.WelGov
	dao             welDAO.IWelDAO
	cli             *welclient.GrpcClient
	tempCli         client.Client
	worker          worker.Worker
	defaultFeelimit int64
}

func MkGovContractService(client *welclient.GrpcClient, tempCli client.Client, daos *dao.DAOs, contractAddr string) (*GovContractService, error) {
	gov := welGov.MkWelGov(contractAddr, client)

	return &GovContractService{cli: client, tempCli: tempCli, gov: gov, dao: daos.Wel, defaultFeelimit: 8000000}, nil
}

func (ctr *GovContractService) GrantRoleOnContract(ctx context.Context, targetAddress string, role string) (string, error) {
	//targetAddress, err := address.Base58ToAddress(target)
	//if err != nil {
	//	logger.Get().Err(err).Msgf("Unable to parse target address")
	//	return "", err
	//}

	callerkey := ctx.Value("callerkey").(string)
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
	opts := &welGov.CallOpts{
		From:      caller,
		Prikey:    pkey,
		Fee_limit: ctr.defaultFeelimit,
		T_amount:  0,
	}

	var brole [32]byte
	if role == model.WelAccountRoleSuperAdmin {
		copy(brole[:], common.Hex2Bytes("0x00")) // in case the default change
	} else {
		copy(brole[:], crypto.Keccak256([]byte(role)))
	}

	tx, err := ctr.gov.GrantRole(opts, brole, targetAddress)
	if err != nil {
		logger.Get().Err(err).Msgf("Unable to trigger governance contract")
		return "", err
	}
	logger.Get().Info().Msgf("Contract call done with tx: %+v", tx)
	return common.Bytes2Hex(tx.GetTxid()), nil
}

func (ctr *GovContractService) GrantRoleWorkflow(ctx workflow.Context, target string, role string) (string, error) {
	log := workflow.GetLogger(ctx)
	log.Info("[Wel workflow] start granting role " + role + " for account " + target)
	ao := workflow.ActivityOptions{
		TaskQueue:              GovContractQueue,
		ScheduleToCloseTimeout: time.Second * 60,
		ScheduleToStartTimeout: time.Second * 60,
		StartToCloseTimeout:    time.Second * 60,
		HeartbeatTimeout:       time.Second * 10,
		WaitForCancellation:    false,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumInterval: time.Second * 500, // huge maximum backoff, because we're dealing with slow blockchain here
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	//call contract
	log.Info("Contract call...")
	var txhash string
	grntContract := workflow.ExecuteActivity(ctx, ctr.GrantRoleOnContract, target, role)
	if err := grntContract.Get(ctx, &txhash); err != nil {
		log.Error("Failed to call grantRole on governance contract")
		return txhash, err
	}

	log.Info("Contract call succeeded")

	// confirmation
	log.Info("Delegate grantRole contract call confirmation to event listener...")

	return txhash, nil
}

func (ctr *GovContractService) RevokeRoleOnContract(ctx context.Context, targetAddress string, role string) (string, error) {
	//targetAddress, err := address.Base58ToAddress(target)
	//if err != nil {
	//	logger.Get().Err(err).Msgf("Unable to parse target address")
	//	return "", err
	//}

	callerkey := ctx.Value("callerkey").(string)
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
	opts := &welGov.CallOpts{
		From:      caller,
		Prikey:    pkey,
		Fee_limit: ctr.defaultFeelimit,
		T_amount:  0,
	}

	var brole [32]byte
	if role == model.WelAccountRoleSuperAdmin {
		copy(brole[:], common.Hex2Bytes("0x00"))
	} else {
		copy(brole[:], crypto.Keccak256([]byte(role)))
	}

	tx, err := ctr.gov.RevokeRole(opts, brole, targetAddress)
	if err != nil {
		logger.Get().Err(err).Msgf("Unable to trigger governance contract")
		return "", err
	}
	logger.Get().Info().Msgf("Contract call done with tx: %+v", tx)
	return common.Bytes2Hex(tx.GetTxid()), nil
}

func (ctr *GovContractService) RevokeRoleWorkflow(ctx workflow.Context, target string, role string) (string, error) {
	log := workflow.GetLogger(ctx)
	log.Info("[Wel workflow] start revoking role " + role + " for account " + target)
	ao := workflow.ActivityOptions{
		TaskQueue:              GovContractQueue,
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

	//call contract
	log.Info("Contract call...")
	var txhash string
	rvkContract := workflow.ExecuteActivity(ctx, ctr.RevokeRoleOnContract, target, role)
	if err := rvkContract.Get(ctx, &txhash); err != nil {
		log.Error("Failed to call revokeRole on governance contract")
		return txhash, err
	}

	log.Info("Contract call succeeded")

	// confirmation
	log.Info("Delegate revokeRole contract call confirmation to event listener...")

	return txhash, nil
}

func (ctr *GovContractService) HasRole(ctx context.Context, targetAddress string, role string) (bool, error) {
	//targetAddress, err := address.Base58ToAddress(target)
	//if err != nil {
	//	logger.Get().Err(err).Msgf("Unable to parse target address")
	//	return false, err
	//}

	opts := &welGov.CallOpts{From: targetAddress}

	var brole [32]byte
	if role == model.WelAccountRoleSuperAdmin {
		copy(brole[:], common.Hex2Bytes("0x00"))
	} else {
		copy(brole[:], crypto.Keccak256([]byte(role)))
	}

	has, err := ctr.gov.HasRole(opts, brole, targetAddress)
	if err != nil {
		logger.Get().Err(err).Msgf("Unable to trigger governance contract")
		return false, err
	}
	logger.Get().Info().Msgf("Contract call done with result: %t", has)
	return has, nil
}

func (ctr *GovContractService) SaveRole(ctx context.Context, address string, role string) error {
	_, err := ctr.dao.GetWelAccount(address)
	switch err {
	case nil:
		break
	case model.ErrWelAccountNotFound:
		if err := ctr.dao.AddWelAccount(address, model.WelAccountStatusOK); err != nil {
			logger.Get().Err(err).Msgf("Unable to add eth account: %s", address)
			return err
		}
	default:
		logger.Get().Err(err).Msgf("Unable to retrieve eth account: %s", address)
		return err
	}

	logger.Get().Info().Msgf("Saving role %s for ethAccount %s...", role, address)
	if err := ctr.dao.GrantRole(address, role); err != nil {
		logger.Get().Err(err).Msgf("Failed to save role %s for ethAccount %s", role, address)
		return err
	}
	return nil
}

func (ctr *GovContractService) RemoveRole(ctx context.Context, address string, role string) error {
	_, err := ctr.dao.GetWelAccount(address)
	switch err {
	case nil:
		break
	case model.ErrWelAccountNotFound:
		logger.Get().Info().Msgf("Wel account %s isn't recognized", address)
		return nil
	default:
		logger.Get().Err(err).Msgf("Unable to retrieve eth account: %s", address)
		return err
	}

	logger.Get().Info().Msgf("Removing role %s for ethAccount %s...", role, address)
	if err := ctr.dao.RevokeRole(address, role); err != nil {
		logger.Get().Err(err).Msgf("Failed to remove role %s for ethAccount %s", role, address)
		return err
	}

	return nil
}

func (ctr *GovContractService) SaveRoleWorkflow(ctx workflow.Context, target string, role string) error {
	log := workflow.GetLogger(ctx)
	log.Info("[Wel workflow] start saving role " + role + " for account " + target)
	ao := workflow.ActivityOptions{
		TaskQueue:              GovContractQueue,
		ScheduleToCloseTimeout: time.Second * 60,
		ScheduleToStartTimeout: time.Second * 60,
		StartToCloseTimeout:    time.Second * 60,
		HeartbeatTimeout:       time.Second * 10,
		WaitForCancellation:    false,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumInterval: time.Second * 100,
			MaximumAttempts: 20,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	res := workflow.ExecuteActivity(ctx, ctr.SaveRole, target, role)
	if err := res.Get(ctx, nil); err != nil {
		log.Error("Failed to save role")
		return err
	}

	log.Info("Role " + role + " saved for " + target)

	return nil
}

func (ctr *GovContractService) RemoveRoleWorkflow(ctx workflow.Context, target string, role string) error {
	log := workflow.GetLogger(ctx)
	log.Info("[Wel workflow] start revoking role " + role + " for account " + target)
	ao := workflow.ActivityOptions{
		TaskQueue:              GovContractQueue,
		ScheduleToCloseTimeout: time.Second * 60,
		ScheduleToStartTimeout: time.Second * 60,
		StartToCloseTimeout:    time.Second * 60,
		HeartbeatTimeout:       time.Second * 10,
		WaitForCancellation:    false,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumInterval: time.Second * 100,
			MaximumAttempts: 20,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	res := workflow.ExecuteActivity(ctx, ctr.RemoveRole, target, role)
	if err := res.Get(ctx, nil); err != nil {
		log.Error("Failed to remove role")
		return err
	}

	log.Info("Role " + role + " removed from " + target)

	return nil
}

// Worker
func (ctr *GovContractService) registerService(w worker.Worker) {
	w.RegisterActivity(ctr.GrantRoleOnContract)
	w.RegisterActivity(ctr.RevokeRoleOnContract)
	w.RegisterActivity(ctr.HasRole)

	w.RegisterWorkflowWithOptions(ctr.GrantRoleWorkflow, workflow.RegisterOptions{Name: GrantRoleWorkflow})
	w.RegisterWorkflowWithOptions(ctr.RevokeRoleWorkflow, workflow.RegisterOptions{Name: RevokeRoleWorkflow})
	//

	w.RegisterActivity(ctr.SaveRole)
	w.RegisterActivity(ctr.RemoveRole)

	w.RegisterWorkflowWithOptions(ctr.SaveRoleWorkflow, workflow.RegisterOptions{Name: SaveRoleWF})
	w.RegisterWorkflowWithOptions(ctr.RemoveRoleWorkflow, workflow.RegisterOptions{Name: RemoveRoleWF})
}

func (ctr *GovContractService) StartService() error {
	w := worker.New(ctr.tempCli, GovContractQueue, worker.Options{})
	ctr.registerService(w)

	ctr.worker = w
	logger.Get().Info().Msgf("Starting GovContractService")
	if err := w.Start(); err != nil {
		logger.Get().Err(err).Msgf("Error while starting GovContractService")
		return err
	}

	logger.Get().Info().Msgf("GovContractService started")
	return nil
}

func (ctr *GovContractService) StopService() {
	if ctr.worker != nil {
		ctr.worker.Stop()
	}
}
