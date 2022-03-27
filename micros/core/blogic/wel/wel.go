package welLogic

import (
	"bridge/libs"
	"bridge/micros/core/model"
	welService "bridge/micros/core/service/wel"
	"context"
	"database/sql"
	"strings"

	"go.temporal.io/sdk/client"
)

//AddWelAccount(address string, status string)
func AddWelAccount(address, status string) error {
	log.Info().Msgf("[Wel logic internal] Creating welAccount %s...", address)

	if !verifyAddress(address) {
		err := model.ErrWelInvalidAddress
		log.Err(err).Msgf("[Wel logic internal] Address %s invalid", address)
		return err
	}

	err := welDAO.AddWelAccount(address, status)
	if err != nil {
		log.Err(err).Msgf("[Wel logic internal] Failed to create welAccount %s", address)
		return err
	}
	return nil
}

//RemoveWelAccount(address string)
func RemoveWelAccount(address string) error {
	log.Info().Msgf("[Wel logic internal] Start removing welAccount %s...", address)
	_, err := welDAO.GetWelAccount(address)
	if err != nil {
		log.Err(err).Msgf("[Wel logic internal] Failed to retrieve welAccount %s", address)
		return err
	}

	log.Info().Msgf("[Wel logic internal] Removing welAccount %s...", address)
	if err := welDAO.RemoveWelAccount(address); err != nil {
		log.Err(err).Msgf("[Wel logic internal] Failed to remove welAccount %s", address)
		return err
	}

	return nil
}

//GetWelAccount(address string)
func GetWelAccount(address string) (*model.WelAccount, error) {
	log.Info().Msgf("[WelAccount logic internal] Getting WelAccount %s", address)
	welAccount, err := welDAO.GetWelAccount(address) // should eventually get by ID instead, but this is more convenient for now
	if err != nil {
		log.Err(err).Msgf("[WelAccount logic internal] Failed to retrieve WelAccount %s's info", address)
		return nil, err
	}

	return welAccount, nil
}

//GetAllWelAccounts(offset uint, size uint)
func GetAllWelAccounts(offset uint, size uint) ([]model.WelAccount, error) {
	log.Info().Msgf("[Wel logic internal] Getting welAccounts...")
	welAccounts, err := welDAO.GetAllWelAccounts(offset, size)
	if err != nil {
		log.Err(err).Msg("[Wel logic internal] Failed to retrieve welAccounts")
		return nil, err
	}

	return welAccounts, nil
}

//GetAllRoles()
func GetAllRoles() ([]string, error) {
	log.Info().Msgf("[Wel logic internal] Getting all roles...")
	roles, err := welDAO.GetAllRoles()
	if err != nil {
		log.Err(err).Msg("[Wel logic internal] Failed to retrieve roles")
		return nil, err
	}

	return roles, nil
}

//GetWelAccountRoles(address string)
func GetWelAccountRoles(address string) ([]string, error) {
	log.Info().Msgf("[WelAccount logic] Getting WelAccount %s's roles...", address)
	roles, err := welDAO.GetWelAccountRoles(address)
	if err != nil {
		log.Err(err).Msgf("[WelAccount logic] Failed to retrieve WelAccount %s's roles", address)
		return nil, err
	}

	return roles, nil
}

//GetWelAccountsWithRole(role string, offset uint, size uint)
func GetWelAccountsWithRole(role string, offset uint, size uint) ([]model.WelAccount, error) {
	log.Info().Msgf("[Wel logic internal] Getting welAccounts with role %s...", role)
	welAccounts, err := welDAO.GetWelAccountsWithRole(role, offset, size)
	if err != nil {
		log.Err(err).Msg("[Wel logic internal] Failed to retrieve welAccounts with role " + role)
		return nil, err
	}

	return welAccounts, nil
}

//GrantRole(address string, role string)
func GrantRole(address, role string, callerkey string) (string, error) {
	log.Info().Msgf("[Wel logic internal] Start granting role %s to welAccount %s...", role, address)

	callerAddress, err := libs.KeyToB58Addr(callerkey)

	if err != nil {
		log.Err(err).Msg("[Wel logic internal] Invalid private key")
		return "", err // invalid key
	}

	log.Info().Msgf("[Wel logic internal] caller address: %s", callerAddress)

	if strings.HasPrefix(address, "0x") {
		address, err = libs.HexToB58(address)
		if err != nil {
			log.Err(err).Msg("[Wel logic internal] Invalid address")
			return "", err
		}
	}

	acc, err := welDAO.GetWelAccount(callerAddress)
	if err != nil {
		log.Err(err).Msgf("[Wel logic internal] Unable to get wel account %s from DB", callerAddress)
		return "", err // invalid key
	}
	if acc.Status != "ok" {
		log.Info().Msgf("[Wel logic internal] account %s locked and cannot do administrative tasks", callerAddress)
		return "", model.ErrWelAccountLocked
	}

	if !verifyAddress(address) {
		err := model.ErrEthInvalidAddress
		log.Err(err).Msgf("[Wel logic internal] invalid address %s", address)
	}

	target, err := welDAO.GetWelAccount(address)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Info().Msgf("[Wel logic internal] address %s not in system", address)
			log.Info().Msgf("[Wel logic internal] Registering address %s", address)
			if err := AddWelAccount(address, model.WelAccountStatusOK); err != nil {
				log.Err(err).Msgf("[Wel logic internal] Unable to add welereum account %s to DB", address)
				return "", err
			}
		} else {
			log.Err(err).Msgf("[Wel logic internal] Unable to get welereum account %s from DB", address)
			return "", err // invalid key
		}
	}
	if target.Status != model.WelAccountStatusOK {
		log.Info().Msgf("[Wel logic internal] account %s locked", address)
		return "", model.ErrWelAccountLocked
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "callerkey", callerkey)
	// call workflow
	log.Info().Msgf("[Wel logic internal] Calling GovContractService workflow...")
	wo := client.StartWorkflowOptions{
		TaskQueue: welService.GovContractQueue,
	}

	we, err := tempcli.ExecuteWorkflow(ctx, wo, welService.GrantRoleWorkflow, address, role)
	if err != nil {
		log.Err(err).Msg("[Wel logic internal] Unable to call GrantRoleWorkflow")
		return "", err
	}
	log.Info().Str("Workflow", we.GetID()).Str("runID=", we.GetRunID()).Msg("dispatched")

	var txhash string
	if err := we.Get(ctx, &txhash); err != nil {
		log.Err(err).Msg("[Wel logic internal] GrantRoleWorkflow failed")
		return txhash, err
	}
	return txhash, nil
}

//RevokeRole(address string, role string)
func RevokeRole(address, role string, callerkey string) (string, error) {
	log.Info().Msgf("[Wel logic internal] Start revoking role %s to welAccount %s...", role, address)

	callerAddress, err := libs.KeyToB58Addr(callerkey)

	if err != nil {
		log.Err(err).Msg("[Wel logic internal] Invalid private key")
		return "", err // invalid key
	}

	log.Info().Msgf("[Wel logic internal] caller address: %s", callerAddress)

	if strings.HasPrefix(address, "0x") {
		address, err = libs.HexToB58(address)
		if err != nil {
			log.Err(err).Msg("[Wel logic internal] Invalid address")
			return "", err
		}
	}

	acc, err := welDAO.GetWelAccount(callerAddress)
	if err != nil {
		log.Err(err).Msgf("[Wel logic internal] Unable to get wel account %s from DB", callerAddress)
		return "", err // invalid key
	}
	if acc.Status != "ok" {
		log.Info().Msgf("[Wel logic internal] account %s locked and cannot do administrative tasks", callerAddress)
		return "", model.ErrWelAccountLocked
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "callerkey", callerkey)
	// call workflow
	log.Info().Msgf("[Wel logic internal] Calling GovContractService workflow...")
	wo := client.StartWorkflowOptions{
		TaskQueue: welService.GovContractQueue,
	}
	we, err := tempcli.ExecuteWorkflow(ctx, wo, welService.RevokeRoleWorkflow, address, role)
	if err != nil {
		log.Err(err).Msg("[Wel logic internal] Unable to call RevokeRoleWorkflow")
		return "", err
	}
	log.Info().Str("Workflow", we.GetID()).Str("runID=", we.GetRunID()).Msg("dispatched")

	var txhash string
	if err := we.Get(ctx, &txhash); err != nil {
		log.Err(err).Msg("[Wel logic internal] RevokeRoleWorkflow failed")
		return txhash, err
	}
	return txhash, nil
}

//SetWelAccountStatus(address string, status string)
func SetWelAccountStatus(address, status string) error {
	log.Info().Msgf("[Wel logic internal] Start setting status %s to welAccount %s...", status, address)

	log.Info().Msgf("[Wel logic internal] setting status %s to welAccount %s...", status, address)
	if err := welDAO.SetWelAccountStatus(address, status); err != nil {
		log.Err(err).Msgf("[Wel logic internal] Failed to set status %s to welAccount %s", status, address)
		return err
	}

	return nil
}

// system keys

func SetCurrentAuthenticator(prikey string) error {
	sysAccounts.Lock()
	defer sysAccounts.Unlock()

	address, err := libs.KeyToB58Addr(prikey)

	if err != nil {
		log.Err(err).Msg("[Wel logic internal] Invalid private key")
		return err // invalid key
	}

	log.Info().Msgf("[Wel logic internal] Set current authenticator to %s", address)

	accs, err := welDAO.GetWelAccountsWithRole(model.WelAccountRoleAuthenticator, 0, 1000) // should've made the DAO to branch out queries instead, but deadline
	if err != nil {
		log.Err(err).Msgf("[Wel logic internal] couldn't retrieve authenticator accounts")
		return err
	}
	match := libs.DropWhile(func(a model.WelAccount) bool { return a.Address != address }, accs)
	if len(match) < 1 {
		err = model.ErrWelAccountNotFound
		log.Err(err).Msgf("[Wel logic internal] authenticator %s not found", address)
		return err
	}

	if match[0].Status != "ok" {
		err = model.ErrWelAccountLocked
		log.Err(err).Msgf("[Wel logic internal] authenticator %s is locked", address)
		return err
	}

	sysAccounts.authenticator.Address = address
	sysAccounts.authenticator.Prikey = prikey
	sysAccounts.authenticator.Status = match[0].Status
	return nil
}

//GetWelPrikeyIfExists(address string)

//SetPriKey(address string, key string)

//UnsetPrikey(address string)
