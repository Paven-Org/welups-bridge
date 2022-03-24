package ethService

import (
	ethGov "bridge/micros/core/abi/eth"
	"bridge/micros/core/dao"
	ethDAO "bridge/micros/core/dao/eth-account"
	"bridge/micros/core/model"
	"bridge/service-managers/logger"
	"context"
	"fmt"
	"math/big"
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

const (
	GrantRoleWorkflow  = "GrantRole"
	RevokeRoleWorkflow = "RevokeRole"
	GovContractQueue   = "GovContractService"
)

type GovContractService struct {
	gov          *ethGov.EthGov
	dao          ethDAO.IEthDAO
	cli          *ethclient.Client
	tempCli      client.Client
	worker       worker.Worker
	lastGasPrice *big.Int
}

func MkGovContractService(client *ethclient.Client, tempCli client.Client, daos *dao.DAOs, contractAddr string) (*GovContractService, error) {
	contractAddress := common.HexToAddress(contractAddr)
	gov, err := ethGov.NewEthGov(contractAddress, client)
	if err != nil {
		logger.Get().Err(err).Msgf("Unable to create governance contract interface at address %s", contractAddr)
		return nil, err
	}

	return &GovContractService{cli: client, tempCli: tempCli, gov: gov, dao: daos.Eth, lastGasPrice: big.NewInt(1000000000)}, nil
}

func (ctr *GovContractService) GrantRoleOnContract(ctx context.Context, target string, role string) (string, error) {
	targetAddress := common.HexToAddress(target)

	gasPrice, err := ctr.cli.SuggestGasPrice(context.Background())
	if err != nil {
		logger.Get().Err(err).Msgf("Unable to get recommended gas price, set to default")
		gasPrice = ctr.lastGasPrice
	}
	ctr.lastGasPrice = gasPrice

	callerkey := ctx.Value("callerkey").(string)
	pkey, err := crypto.HexToECDSA(callerkey)
	if err != nil {
		logger.Get().Err(err).Msgf("Unable to parse hexstring to ECDSA key")
		return "", err
	}
	caller := crypto.PubkeyToAddress(pkey.PublicKey)

	nonce, err := ctr.cli.PendingNonceAt(ctx, caller)
	if err != nil {
		logger.Get().Err(err).Msgf("Unale to get last nonce of address %s", caller.Hex())
		return "", err
	}

	opts := bind.NewKeyedTransactor(pkey)
	opts.GasLimit = uint64(300000)
	opts.Value = big.NewInt(0)
	opts.GasPrice = gasPrice
	opts.Nonce = big.NewInt(int64(nonce))

	var brole [32]byte
	if role == model.EthAccountRoleSuperAdmin {
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
	return tx.Hash().Hex(), nil
}

func (ctr *GovContractService) FilterRoleGranted(ctx context.Context, target string, role string) (bool, error) {
	targetAddress := common.HexToAddress(target)

	var brole [32]byte
	if role == model.EthAccountRoleSuperAdmin {
		copy(brole[:], common.Hex2Bytes("0x00"))
	} else {
		copy(brole[:], crypto.Keccak256([]byte(role)))
	}

	callerkey := ctx.Value("callerkey").(string)
	pkey, err := crypto.HexToECDSA(callerkey)
	if err != nil {
		logger.Get().Err(err).Msgf("Unable to parse hexstring to ECDSA key")
		return false, err
	}
	caller := crypto.PubkeyToAddress(pkey.PublicKey)

	res, err := ctr.gov.EthGovFilterer.FilterRoleGranted(nil, [][32]byte{brole}, []common.Address{targetAddress}, []common.Address{caller})
	if err != nil {
		logger.Get().Err(err).Msgf("Unable to get RoleGranted event")
		return false, err
	}
	defer res.Close()

	cnt := 0
	for res.Next() {
		cnt++
		fmt.Println("event: ", res.Event)
	}

	return cnt > 0, nil
}
func (ctr *GovContractService) SaveRoleOnDB(ctx context.Context, address, role string) error {
	logger.Get().Info().Msgf("[ethAccount logic internal] Start saving role %s for ethAccount %s...", role, address)
	//ethAccount, err := ethDAO.GetEthAccount(address)
	//if err != nil {
	//	logger.Get().Err(err).Msgf("[ethAccount logic internal] Failed to retrieve ethAccount %s", address)
	//	return err
	//}
	//

	logger.Get().Info().Msgf("[Eth activity] Saving role %s for ethAccount %s...", role, address)
	if err := ctr.dao.GrantRole(address, role); err != nil {
		logger.Get().Err(err).Msgf("[ethAccount logic internal] Failed to save role %s for ethAccount %s", role, address)
		return err
	}

	return nil
}

func (ctr *GovContractService) GrantRoleWorkflow(ctx workflow.Context, target string, role string) (string, error) {
	log := workflow.GetLogger(ctx)
	log.Info("[Eth workflow] start granting role " + role + " for account " + target)
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
	log.Info("Confirm grantRole contract call...")

	//grntEvent := workflow.ExecuteActivity(ctx, ctr.FilterRoleGranted, target, role)
	//if err := grntEvent.Get(ctx, &res); err != nil {
	//	log.Error("Failed to confirm contract call with event")
	//	return txhash, err
	//}

	var has bool = false
	interval := time.Second * 1

	for !has {
		workflow.Sleep(ctx, interval)
		hasRolePromise := workflow.ExecuteActivity(ctx, ctr.HasRole, target, role)
		if err := hasRolePromise.Get(ctx, &has); err != nil {
			log.Error("Failed to call hasRole on governance contract, error: ", err.Error())
		}
		if has {
			break
		}

		interval = interval * 3
		if interval > time.Hour {
			break
		}
	}
	if !has {
		err := fmt.Errorf("Failed to confirm grantRole")
		log.Error(err.Error())
		return txhash, err
	}

	log.Info("grantRole confirmed")

	//save on DB
	log.Info("Save address-role pair from DB...")

	grntDB := workflow.ExecuteActivity(ctx, ctr.SaveRoleOnDB, target, role)
	if err := grntDB.Get(ctx, nil); err != nil {
		log.Error("Failed to save role for address in DB")
		return txhash, err
	}

	log.Info("Role saved for address in DB")

	return txhash, nil
}

func (ctr *GovContractService) RevokeRoleOnContract(ctx context.Context, target string, role string) (string, error) {
	targetAddress := common.HexToAddress(target)

	gasPrice, err := ctr.cli.SuggestGasPrice(context.Background())
	if err != nil {
		logger.Get().Err(err).Msgf("Unable to get recommended gas price, set to default")
		gasPrice = ctr.lastGasPrice
	}
	ctr.lastGasPrice = gasPrice

	callerkey := ctx.Value("callerkey").(string)
	pkey, err := crypto.HexToECDSA(callerkey)
	if err != nil {
		logger.Get().Err(err).Msgf("Unable to parse hexstring to ECDSA key")
		return "", err
	}
	caller := crypto.PubkeyToAddress(pkey.PublicKey)

	nonce, err := ctr.cli.PendingNonceAt(ctx, caller)
	if err != nil {
		logger.Get().Err(err).Msgf("Unale to get last nonce of address %s", caller.Hex())
		return "", err
	}

	opts := bind.NewKeyedTransactor(pkey)
	opts.GasLimit = uint64(300000)
	opts.Value = big.NewInt(0)
	opts.GasPrice = gasPrice
	opts.Nonce = big.NewInt(int64(nonce))

	var brole [32]byte
	if role == model.EthAccountRoleSuperAdmin {
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
	return tx.Hash().Hex(), nil
}

func (ctr *GovContractService) FilterRoleRevoked(ctx context.Context, target string, role string) (bool, error) {
	targetAddress := common.HexToAddress(target)

	var brole [32]byte
	if role == model.EthAccountRoleSuperAdmin {
		copy(brole[:], common.Hex2Bytes("0x00"))
	} else {
		copy(brole[:], crypto.Keccak256([]byte(role)))
	}

	callerkey := ctx.Value("callerkey").(string)
	pkey, err := crypto.HexToECDSA(callerkey)
	if err != nil {
		logger.Get().Err(err).Msgf("Unable to parse hexstring to ECDSA key")
		return false, err
	}
	caller := crypto.PubkeyToAddress(pkey.PublicKey)

	res, err := ctr.gov.EthGovFilterer.FilterRoleRevoked(nil, [][32]byte{brole}, []common.Address{targetAddress}, []common.Address{caller})
	if err != nil {
		logger.Get().Err(err).Msgf("Unable to get RoleRevoked event")
		return false, err
	}
	defer res.Close()

	cnt := 0
	for res.Next() {
		cnt++
		fmt.Println("event: ", res.Event)
	}

	return cnt > 0, nil
}

func (ctr *GovContractService) RemoveRoleOnDB(ctx context.Context, address, role string) error {
	logger.Get().Info().Msgf("[Eth activity] Start removing role %s from ethAccount %s...", role, address)
	//ethAccount, err := ethDAO.GetEthAccount(address)
	//if err != nil {
	//	logger.Get().Err(err).Msgf("[ethAccount logic internal] Failed to retrieve ethAccount %s", address)
	//	return err
	//}
	//

	logger.Get().Info().Msgf("[Eth activity] Removing role %s for ethAccount %s...", role, address)
	if err := ctr.dao.RevokeRole(address, role); err != nil {
		logger.Get().Err(err).Msgf("[Eth activity] Failed to remove role %s for ethAccount %s", role, address)
		return err
	}

	return nil
}

func (ctr *GovContractService) RevokeRoleWorkflow(ctx workflow.Context, target string, role string) (string, error) {
	log := workflow.GetLogger(ctx)
	log.Info("[Eth workflow] start revoking role " + role + " for account " + target)
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

	//listen for event
	log.Info("Confirm revokeRole call...")

	//rvkEvent := workflow.ExecuteActivity(ctx, ctr.FilterRoleRevoked, target, role)
	//if err := rvkEvent.Get(ctx, &res); err != nil {
	//	log.Error("Failed to confirm contract call with event")
	//	return txhash, err
	//}
	var has bool = true
	interval := time.Second * 1

	for has {
		workflow.Sleep(ctx, interval)
		hasRolePromise := workflow.ExecuteActivity(ctx, ctr.HasRole, target, role)
		if err := hasRolePromise.Get(ctx, &has); err != nil {
			log.Error("Failed to call hasRole on governance contract, error: ", err.Error())
		}
		if !has {
			break
		}

		interval = interval * 3
		if interval > time.Hour {
			break
		}
	}
	if has {
		err := fmt.Errorf("Failed to confirm revokeRole call")
		log.Error(err.Error())
		return txhash, err
	}

	log.Info("revokeRole confirmed")

	//save on DB
	log.Info("Remove address-role pair from DB...")

	rvkDB := workflow.ExecuteActivity(ctx, ctr.RemoveRoleOnDB, target, role)
	if err := rvkDB.Get(ctx, nil); err != nil {
		log.Error("Failed to remove role for address in DB")
		return txhash, err
	}

	log.Info("Role removed for address in DB")

	return txhash, nil
}

func (ctr *GovContractService) HasRole(ctx context.Context, target string, role string) (bool, error) {
	targetAddress := common.HexToAddress(target)

	opts := &bind.CallOpts{Context: ctx}

	var brole [32]byte
	if role == model.EthAccountRoleSuperAdmin {
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

// Worker
func (ctr *GovContractService) registerService(w worker.Worker) {
	w.RegisterActivity(ctr.GrantRoleOnContract)
	w.RegisterActivity(ctr.FilterRoleGranted)
	w.RegisterActivity(ctr.SaveRoleOnDB)

	w.RegisterActivity(ctr.HasRole)

	w.RegisterActivity(ctr.RevokeRoleOnContract)
	w.RegisterActivity(ctr.FilterRoleRevoked)
	w.RegisterActivity(ctr.RemoveRoleOnDB)

	w.RegisterWorkflowWithOptions(ctr.GrantRoleWorkflow, workflow.RegisterOptions{Name: GrantRoleWorkflow})
	w.RegisterWorkflowWithOptions(ctr.RevokeRoleWorkflow, workflow.RegisterOptions{Name: RevokeRoleWorkflow})
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
