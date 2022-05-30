package ethLogic

import (
	"bridge/libs"
	msweleth "bridge/micros/core/microservices/weleth"
	"bridge/micros/core/model"
	ethService "bridge/micros/core/service/eth"
	"bridge/micros/core/service/notifier"
	welService "bridge/micros/core/service/wel"
	welethModel "bridge/micros/weleth/model"
	"bridge/service-managers/logger"
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"go.temporal.io/sdk/client"
)

//AddEthAccount(address string, status string)
func AddEthAccount(address, status string) error {
	log.Info().Msgf("[Eth logic internal] Creating ethAccount %s...", address)

	if !verifyAddress(address) {
		err := model.ErrEthInvalidAddress
		log.Err(err).Msgf("[Eth logic internal] Address %s invalid", address)
		return err
	}

	err := ethDAO.AddEthAccount(address, status)
	if err != nil {
		log.Err(err).Msgf("[Eth logic internal] Failed to create ethAccount %s", address)
		return err
	}
	return nil
}

//RemoveEthAccount(address string)
func RemoveEthAccount(address string) error {
	log.Info().Msgf("[Eth logic internal] Start removing ethAccount %s...", address)
	_, err := ethDAO.GetEthAccount(address)
	if err != nil {
		log.Err(err).Msgf("[Eth logic internal] Failed to retrieve ethAccount %s", address)
		return err
	}

	log.Info().Msgf("[Eth logic internal] Removing ethAccount %s...", address)
	if err := ethDAO.RemoveEthAccount(address); err != nil {
		log.Err(err).Msgf("[Eth logic internal] Failed to remove ethAccount %s", address)
		return err
	}

	return nil
}

//GetEthAccount(address string)
func GetEthAccount(address string) (*model.EthAccount, error) {
	log.Info().Msgf("[EthAccount logic internal] Getting EthAccount %s", address)
	ethAccount, err := ethDAO.GetEthAccount(address) // should eventually get by ID instead, but this is more convenient for now
	if err != nil {
		log.Err(err).Msgf("[EthAccount logic internal] Failed to retrieve EthAccount %s's info", address)
		return nil, err
	}

	return ethAccount, nil
}

//GetAllEthAccounts(offset uint, size uint)
func GetAllEthAccounts(offset uint, size uint) ([]model.EthAccount, error) {
	log.Info().Msgf("[Eth logic internal] Getting ethAccounts...")
	ethAccounts, err := ethDAO.GetAllEthAccounts(offset, size)
	if err != nil {
		log.Err(err).Msg("[Eth logic internal] Failed to retrieve ethAccounts")
		return nil, err
	}

	return ethAccounts, nil
}

//GetAllRoles()
func GetAllRoles() ([]string, error) {
	log.Info().Msgf("[Eth logic internal] Getting all roles...")
	roles, err := ethDAO.GetAllRoles()
	if err != nil {
		log.Err(err).Msg("[Eth logic internal] Failed to retrieve roles")
		return nil, err
	}

	return roles, nil
}

//GetEthAccountRoles(address string)
func GetEthAccountRoles(address string) ([]string, error) {
	log.Info().Msgf("[EthAccount logic] Getting EthAccount %s's roles...", address)
	roles, err := ethDAO.GetEthAccountRoles(address)
	if err != nil {
		log.Err(err).Msgf("[EthAccount logic] Failed to retrieve EthAccount %s's roles", address)
		return nil, err
	}

	return roles, nil
}

//GetEthAccountsWithRole(role string, offset uint, size uint)
func GetEthAccountsWithRole(role string, offset uint, size uint) ([]model.EthAccount, error) {
	log.Info().Msgf("[Eth logic internal] Getting ethAccounts with role %s...", role)
	ethAccounts, err := ethDAO.GetEthAccountsWithRole(role, offset, size)
	if err != nil {
		log.Err(err).Msg("[Eth logic internal] Failed to retrieve ethAccounts with role " + role)
		return nil, err
	}

	return ethAccounts, nil
}

//GrantRole(address string, role string)
func GrantRole(address, role string, callerkey string) (string, error) {
	log.Info().Msgf("[Eth logic internal] Start granting role %s to ethAccount %s...", role, address)

	key, err := crypto.HexToECDSA(callerkey)
	if err != nil {
		log.Err(err).Msg("[Eth logic internal] Invalid private key")
		return "", err // invalid key
	}

	callerAddress := crypto.PubkeyToAddress(key.PublicKey).Hex()
	log.Info().Msgf("[Eth logic internal] caller address: %s", callerAddress)

	acc, err := ethDAO.GetEthAccount(callerAddress)
	if err != nil {
		log.Err(err).Msgf("[Eth logic internal] Unable to get ethereum account %s from DB", callerAddress)
		return "", err // invalid key
	}
	if acc.Status != model.EthAccountStatusOK {
		log.Info().Msgf("[Eth logic internal] account %s locked and cannot do administrative tasks", callerAddress)
		return "", model.ErrEthAccountLocked
	}

	if !verifyAddress(address) {
		err := model.ErrEthInvalidAddress
		log.Err(err).Msgf("[Eth logic internal] invalid address %s", address)
	}

	target, err := ethDAO.GetEthAccount(address)
	if err != nil {
		if err == model.ErrEthAccountNotFound {
			log.Info().Msgf("[Eth logic internal] address %s not in system", address)
			log.Info().Msgf("[Eth logic internal] Registering address %s", address)
			if err := AddEthAccount(address, model.EthAccountStatusOK); err != nil {
				log.Err(err).Msgf("[Eth logic internal] Unable to add ethereum account %s to DB", address)
				return "", err
			}
		} else {
			log.Err(err).Msgf("[Eth logic internal] Unable to get ethereum account %s from DB", address)
			return "", err // invalid key
		}
	}
	if target.Status != model.EthAccountStatusOK {
		log.Info().Msgf("[Eth logic internal] account %s locked", address)
		return "", model.ErrEthAccountLocked
	}

	// call contract & persist granted role in system DB via workflow
	// cross-system transactional semantics is needed, thus the use of workflow
	ctx := context.Background()
	ctx = context.WithValue(ctx, "callerkey", callerkey)
	// call workflow
	log.Info().Msgf("[Eth logic internal] Calling GovContractService workflow...")
	wo := client.StartWorkflowOptions{
		TaskQueue: ethService.GovContractQueue,
	}

	we, err := tempcli.ExecuteWorkflow(ctx, wo, ethService.GrantRoleWorkflow, address, role)
	if err != nil {
		log.Err(err).Msg("[Eth logic internal] Unable to call GrantRoleWorkflow")
		return "", err
	}
	log.Info().Str("Workflow", we.GetID()).Str("runID=", we.GetRunID()).Msg("dispatched")

	var txhash string
	if err := we.Get(ctx, &txhash); err != nil {
		log.Err(err).Msg("[Eth logic internal] GrantRoleWorkflow failed")
		return txhash, err
	}
	return txhash, nil
}

//RevokeRole(address string, role string)
func RevokeRole(address, role string, callerkey string) (string, error) {
	log.Info().Msgf("[Eth logic internal] Start revoking role %s to ethAccount %s...", role, address)

	key, err := crypto.HexToECDSA(callerkey)
	if err != nil {
		log.Err(err).Msg("[Eth logic internal] Invalid private key")
		return "", err // invalid key
	}

	callerAddress := crypto.PubkeyToAddress(key.PublicKey).Hex()
	log.Info().Msgf("[Eth logic internal] caller address: %s", callerAddress)

	acc, err := ethDAO.GetEthAccount(callerAddress)
	if err != nil {
		log.Err(err).Msgf("[Eth logic internal] Unable to get ethereum account %s from DB", callerAddress)
		return "", err // invalid key
	}
	if acc.Status != "ok" {
		log.Info().Msgf("[Eth logic internal] account %s locked and cannot do administrative tasks", callerAddress)
		return "", model.ErrEthAccountLocked
	}

	// call contract & remove revoked role from system DB via workflow
	// cross-system transactional semantics is needed, thus the use of workflow
	ctx := context.Background()
	ctx = context.WithValue(ctx, "callerkey", callerkey)
	// call workflow
	log.Info().Msgf("[Eth logic internal] Calling GovContractService workflow...")
	wo := client.StartWorkflowOptions{
		TaskQueue: ethService.GovContractQueue,
	}
	we, err := tempcli.ExecuteWorkflow(ctx, wo, ethService.RevokeRoleWorkflow, address, role)
	if err != nil {
		log.Err(err).Msg("[Eth logic internal] Unable to call RevokeRoleWorkflow")
		return "", err
	}
	log.Info().Str("Workflow", we.GetID()).Str("runID=", we.GetRunID()).Msg("dispatched")

	var txhash string
	if err := we.Get(ctx, &txhash); err != nil {
		log.Err(err).Msg("[Eth logic internal] RevokeRoleWorkflow failed")
		return txhash, err
	}
	return txhash, nil
}

//SetEthAccountStatus(address string, status string)
func SetEthAccountStatus(address, status string) error {
	log.Info().Msgf("[Eth logic internal] Start setting status %s to ethAccount %s...", status, address)
	//ethAccount, err := ethDAO.GetEthAccount(address)
	//if err != nil {
	//	log.Err(err).Msgf("[Eth logic internal] Failed to retrieve ethAccount %s", address)
	//	return err
	//}

	log.Info().Msgf("[Eth logic internal] setting status %s to ethAccount %s...", status, address)
	if err := ethDAO.SetEthAccountStatus(address, status); err != nil {
		log.Err(err).Msgf("[Eth logic internal] Failed to set status %s to ethAccount %s", status, address)
		return err
	}

	return nil
}

// system keys

func SetCurrentAuthenticator(prikey string) error {
	sysAccounts.Lock()
	defer sysAccounts.Unlock()

	key, err := crypto.HexToECDSA(prikey)
	if err != nil {
		log.Err(err).Msgf("[Eth logic internal] invalid private key")
		return err
	}

	address := crypto.PubkeyToAddress(key.PublicKey).Hex()
	log.Info().Msgf("[Eth logic internal] Set current authenticator to %s", address)

	accs, err := ethDAO.GetEthAccountsWithRole(model.EthAccountRoleAuthenticator, 0, 1000) // should've made the DAO to branch out queries instead, but deadline
	if err != nil {
		log.Err(err).Msgf("[Eth logic internal] couldn't retrieve authenticator accounts")
		return err
	}
	match := libs.DropWhile(func(a model.EthAccount) bool { return a.Address != address }, accs)
	if len(match) < 1 {
		err = model.ErrEthAccountNotFound
		log.Err(err).Msgf("[Eth logic internal] authenticator %s not found", address)
		return err
	}

	if match[0].Status != "ok" {
		err = model.ErrEthAccountLocked
		log.Err(err).Msgf("[Eth logic internal] authenticator %s is locked", address)
		return err
	}

	sysAccounts.authenticator.Address = address
	sysAccounts.authenticator.Prikey = prikey
	sysAccounts.authenticator.Status = match[0].Status
	return nil
}

func UnsetCurrentAuthenticator() error {
	sysAccounts.Lock()
	defer sysAccounts.Unlock()
	sysAccounts.authenticator = model.EthAccount{}
	// immediately send notification email to admin
	return nil
}

// Claim cashin = get wrapped tokens equivalent to another chain's original tokens
func ClaimWel2EthCashin(cashinTxId string, userAddr string, contractVersion string) (inTokenAddr string, amount string, requestID []byte, signature []byte, err error) {
	// Get tx info from weleth microservice
	// tmpCli.ExecuteWorkflow
	ctx := context.Background()
	// call workflow
	log.Info().Msgf("[Eth logic internal] Calling MSWeleth workflow...")
	wo := client.StartWorkflowOptions{
		TaskQueue: msweleth.TaskQueue,
	}
	we, err := tempcli.ExecuteWorkflow(ctx, wo, msweleth.CreateW2ECashinClaimRequestWF, cashinTxId, userAddr)
	if err != nil {
		log.Err(err).Msg("[Eth logic internal] Unable to call CreateW2ECashinClaimRequest workflow")
		return
	}
	log.Info().Str("Workflow", we.GetID()).Str("runID=", we.GetRunID()).Msg("dispatched")

	var tx welethModel.WelCashinEthTrans
	if err = we.Get(ctx, &tx); err != nil {
		log.Err(err).Msg("[Eth logic internal] CreateW2ECashinClaimRequest workflow failed")
		return
	}
	tempcli.ExecuteWorkflow(ctx, wo, msweleth.WaitForPendingW2ECashinClaimRequestWF, cashinTxId)
	// process

	log.Info().Msg("[Eth logic internal] Everything a-ok, proceeding to create signature and requestID")

	sysAccounts.RLock()
	defer sysAccounts.RUnlock()
	prikey := sysAccounts.authenticator.Prikey
	// if prikey == "", send notification mail to admin and return error
	if prikey == "" {
		problem := model.ErrEthAuthenticatorKeyUnavailable
		wo := client.StartWorkflowOptions{
			TaskQueue: notifier.NotifierQueue,
		}

		we, err := tempcli.ExecuteWorkflow(ctx, wo, notifier.NotifyProblemWF, problem.Error(), "admin")
		if err != nil {
			log.Err(err).Msg("[Eth logic internal] Failed to notify admins of problem: " + problem.Error())
			return "", "", nil, nil, err
		}
		log.Info().Str("Workflow", we.GetID()).Str("runID=", we.GetRunID()).Msg("dispatched")
		if err := we.Get(ctx, nil); err != nil {
			log.Err(err).Msg("[Eth logic internal] Failed to notify admins of problem: " + problem.Error())
			return "", "", nil, nil, err
		}
		err = problem
		return "", "", nil, nil, err
	}

	inTokenAddr = tx.EthTokenAddr
	toAddress := tx.EthWalletAddr
	amount = tx.Amount

	_requestID := &big.Int{}
	_requestID.SetString(tx.ReqID, 10)
	requestID = _requestID.Bytes()

	_amount := &big.Int{}
	_amount.SetString(tx.Amount, 10)

	signature, err = libs.StdSignedMessageHash(inTokenAddr, toAddress, _amount, _requestID, contractVersion, prikey)
	if err != nil {
		log.Err(err).Msg("[Eth logic internal] Failed to create claim signature for user")
		return
	}

	log.Info().Msg("[Eth logic internal] Successfully create claim signature for user")
	return
}

func GetAuthenticatorKey() (string, error) {
	sysAccounts.RLock()
	defer sysAccounts.RUnlock()
	prikey := sysAccounts.authenticator.Prikey
	// if prikey == "", send notification mail to admin and return error
	if prikey == "" {
		ctx := context.Background()
		problem := model.ErrEthAuthenticatorKeyUnavailable
		wo := client.StartWorkflowOptions{
			TaskQueue: notifier.NotifierQueue,
		}

		we, err := tempcli.ExecuteWorkflow(ctx, wo, notifier.NotifyProblemWF, problem.Error(), "admin")
		if err != nil {
			log.Err(err).Msg("[Eth logic internal] Failed to notify admins of problem: " + problem.Error())
			return "", err
		}
		log.Info().Str("Workflow", we.GetID()).Str("runID=", we.GetRunID()).Msg("dispatched")
		if err := we.Get(ctx, nil); err != nil {
			log.Err(err).Msg("[Eth logic internal] Failed to notify admins of problem: " + problem.Error())
			return "", err
		}
		err = problem
		return "", err
	}

	return prikey, nil
}

func InvalidateRequestClaim(inTokenAddr, amount, reqID, contractVersion string) error {
	ctx := context.Background()
	prikey, err := GetAuthenticatorKey()
	if err != nil {
		log.Err(err).Msgf("[Eth logic internal] Authenticator key not available")
		return err
	}
	pkey, err := crypto.HexToECDSA(prikey)
	if err != nil {
		log.Err(err).Msgf("[Eth logic internal] invalid private key")
		return err
	}

	caller := crypto.PubkeyToAddress(pkey.PublicKey)
	address := caller.Hex()
	log.Info().Msgf("[Eth logic internal] operator address: %s", address)

	_requestID := &big.Int{}
	_requestID.SetString(reqID, 10)

	_amount := &big.Int{}
	_amount.SetString(amount, 10)

	signature, err := libs.StdSignedMessageHash(inTokenAddr, address, _amount, _requestID, contractVersion, prikey)
	if err != nil {
		log.Err(err).Msg("[Eth logic internal] Failed to create claim signature")
		return err
	}

	log.Info().Msg("[Eth logic internal] Successfully create claim signature")

	log.Info().Msg("[Eth logic internal] Invalidating request ID " + reqID)

	nonce, err := importC.cli.PendingNonceAt(ctx, caller)
	if err != nil {
		logger.Get().Err(err).Msgf("Unale to get last nonce of address %s", address)
		return err
	}
	gasPrice, err := importC.cli.SuggestGasPrice(context.Background())
	if err != nil {
		logger.Get().Err(err).Msgf("Unable to get recommended gas price, set to default")
		gasPrice = importC.lastGasPrice
	}
	importC.lastGasPrice = gasPrice

	opts := bind.NewKeyedTransactor(pkey)
	opts.GasLimit = uint64(300000)
	opts.Value = big.NewInt(0)
	opts.GasPrice = gasPrice
	opts.Nonce = big.NewInt(int64(nonce))

	tokenAddr := common.HexToAddress(inTokenAddr)
	tx, err := importC.impC.EthImportCTransactor.Claim(opts, tokenAddr, _requestID, _amount, signature)
	//if err != nil {
	logger.Get().Err(err).Msgf("[Eth logic internal] failed tx: %v", tx)
	//	return err
	//}

	return nil
}

func GetE2WCashinTransByEthTxHash(txhash string) (*welethModel.EthCashinWelTrans, error) {
	wo := client.StartWorkflowOptions{
		TaskQueue: msweleth.TaskQueue,
	}

	var tx welethModel.EthCashinWelTrans
	ctx := context.Background()

	we, err := tempcli.ExecuteWorkflow(ctx, wo, msweleth.GetEthToWelCashinByTxHash, txhash)
	if err != nil {
		log.Err(err).Msgf("[Eth logic internal] Failed to execute Get E2W cashin tx workflow with eth txhash %s", txhash)
		return nil, err
	}
	if err = we.Get(ctx, &tx); err != nil {
		log.Err(err).Msgf("[Eth logic internal] Failed to get E2W cashin tx with eth txhash %s", txhash)
		return &tx, err
	}
	log.Info().Msgf("[Eth logic internal] Retrieved E2W cashin tx with txhash %s", txhash)
	return &tx, nil
}

func WatchTx2TreasuryRequest(from, to, treasury, netid, token, amount string) error {
	wo := client.StartWorkflowOptions{
		TaskQueue: welService.ImportContractQueue,
	}

	we, err := tempcli.ExecuteWorkflow(context.Background(), wo, welService.WatchForTx2TreasuryWF, from, to, treasury, netid, token, amount)
	if err != nil {
		log.Err(err).Msgf("[Eth logic internal] Failed to request BE to watch for transaction to treasury from %s", from)
		return err
	}
	log.Info().Msgf("[Eth logic internal] Request BE to watch for transaction to treasury, WF ID: %s, run ID: %s", we.GetID(), we.GetRunID())
	return nil
}

func WatchTx2TreasuryRequestByTxhash(txhash, to, netid, token string) error {
	wo := client.StartWorkflowOptions{
		TaskQueue: welService.ImportContractQueue,
	}

	we, err := tempcli.ExecuteWorkflow(context.Background(), wo, welService.WatchForTx2TreasuryByTxHashWF, txhash, to, netid, token)
	if err != nil {
		log.Err(err).Msgf("[Eth logic internal] Failed to request BE to watch for transaction to treasury with txhash %s", txhash)
		return err
	}
	log.Info().Msgf("[Eth logic internal] Request BE to watch for transaction to treasury, WF ID: %s, run ID: %s", we.GetID(), we.GetRunID())
	return nil
}

//GetEthPrikeyIfExists(address string)

//SetPriKey(address string, key string)

//UnsetPrikey(address string)
