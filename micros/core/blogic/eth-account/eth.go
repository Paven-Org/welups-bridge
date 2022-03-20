package ethLogic

import (
	"bridge/libs"
	"bridge/micros/core/model"
)

//AddEthAccount(address string, status string)
func AddEthAccount(address, status string) error {
	log.Info().Msgf("[ethAccount logic internal] Creating ethAccount %s...", address)

	if !verifyAddress(address) {
		err := model.ErrEthInvalidAddress
		log.Err(err).Msgf("[ethAccount logic internal] Address %s invalid", address)
		return err
	}

	err := ethDAO.AddEthAccount(address, status)
	if err != nil {
		log.Err(err).Msgf("[ethAccount logic internal] Failed to create ethAccount %s", address)
		return err
	}
	return nil
}

//RemoveEthAccount(address string)
func RemoveEthAccount(address string) error {
	log.Info().Msgf("[ethAccount logic internal] Start removing ethAccount %s...", address)
	_, err := ethDAO.GetEthAccount(address)
	if err != nil {
		log.Err(err).Msgf("[ethAccount logic internal] Failed to retrieve ethAccount %s", address)
		return err
	}

	log.Info().Msgf("[ethAccount logic internal] Removing ethAccount %s...", address)
	if err := ethDAO.RemoveEthAccount(address); err != nil {
		log.Err(err).Msgf("[ethAccount logic internal] Failed to remove ethAccount %s", address)
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
	log.Info().Msgf("[ethAccount logic internal] Getting ethAccounts...")
	ethAccounts, err := ethDAO.GetAllEthAccounts(offset, size)
	if err != nil {
		log.Err(err).Msg("[ethAccount logic internal] Failed to retrieve ethAccounts")
		return nil, err
	}

	return ethAccounts, nil
}

//GetAllRoles()
func GetAllRoles() ([]string, error) {
	log.Info().Msgf("[ethAccount logic internal] Getting all roles...")
	roles, err := ethDAO.GetAllRoles()
	if err != nil {
		log.Err(err).Msg("[ethAccount logic internal] Failed to retrieve roles")
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
	log.Info().Msgf("[ethAccount logic internal] Getting ethAccounts with role %s...", role)
	ethAccounts, err := ethDAO.GetEthAccountsWithRole(role, offset, size)
	if err != nil {
		log.Err(err).Msg("[ethAccount logic internal] Failed to retrieve ethAccounts with role " + role)
		return nil, err
	}

	return ethAccounts, nil
}

//GrantRole(address string, role string)
func GrantRole(address, role string) error {
	log.Info().Msgf("[ethAccount logic internal] Start granting role %s to ethAccount %s...", role, address)
	//ethAccount, err := ethDAO.GetEthAccount(address)
	//if err != nil {
	//	log.Err(err).Msgf("[ethAccount logic internal] Failed to retrieve ethAccount %s", address)
	//	return err
	//}

	log.Info().Msgf("[ethAccount logic internal] Granting role %s to ethAccount %s...", role, address)
	if err := ethDAO.GrantRole(address, role); err != nil {
		log.Err(err).Msgf("[ethAccount logic internal] Failed to grant role %s to ethAccount %s", role, address)
		return err
	}

	return nil
}

//RevokeRole(address string, role string)
func RevokeRole(address, role string) error {
	log.Info().Msgf("[ethAccount logic internal] Start revoking role %s from ethAccount %s...", role, address)
	//ethAccount, err := ethDAO.GetEthAccount(address)
	//if err != nil {
	//	log.Err(err).Msgf("[ethAccount logic internal] Failed to retrieve ethAccount %s", address)
	//	return err
	//}

	log.Info().Msgf("[ethAccount logic internal] Revoking role %s from ethAccount %s...", role, address)
	if err := ethDAO.RevokeRole(address, role); err != nil {
		log.Err(err).Msgf("[ethAccount logic internal] Failed to revoke role %s from ethAccount %s", role, address)
		return err
	}

	return nil
}

//SetEthAccountStatus(address string, status string)
func SetEthAccountStatus(address, status string) error {
	log.Info().Msgf("[ethAccount logic internal] Start setting status %s to ethAccount %s...", status, address)
	//ethAccount, err := ethDAO.GetEthAccount(address)
	//if err != nil {
	//	log.Err(err).Msgf("[ethAccount logic internal] Failed to retrieve ethAccount %s", address)
	//	return err
	//}

	log.Info().Msgf("[ethAccount logic internal] setting status %s to ethAccount %s...", status, address)
	if err := ethDAO.SetEthAccountStatus(address, status); err != nil {
		log.Err(err).Msgf("[ethAccount logic internal] Failed to set status %s to ethAccount %s", status, address)
		return err
	}

	return nil
}

// system keys
func SetCurrentSuperAdmin(address string, prikey string) error {
	sysAccounts.Lock()
	defer sysAccounts.Unlock()

	log.Info().Msgf("[ethAccount logic internal] Set current super admin to %s", address)

	if !verifyKeyAndAddress(prikey, address) {
		err := model.ErrEthKeyAndAddressMismatch
		log.Err(err).Msgf("[ethAccount logic internal] Key and address mismatch for account %s", address)
		return err
	}

	accs, err := ethDAO.GetEthAccountsWithRole(model.EthAccountRoleSuperAdmin, 0, 1000)
	if err != nil {
		log.Err(err).Msgf("[ethAccount logic internal] couldn't retrieve super admin accounts")
		return err
	}
	match := libs.DropWhile(func(a model.EthAccount) bool { return a.Address != address }, accs)
	if len(match) < 1 {
		err = model.ErrEthAccountNotFound
		log.Err(err).Msgf("[ethAccount logic internal] super admin %s not found", address)
		return err
	}

	if match[0].Status != "ok" {
		err = model.ErrEthAccountLocked
		log.Err(err).Msgf("[ethAccount logic internal] super admin %s is locked", address)
		return err
	}

	sysAccounts.superAdmin.Address = address
	sysAccounts.superAdmin.Prikey = prikey
	sysAccounts.superAdmin.Status = match[0].Status
	return nil
}

func SetCurrentAuthenticator(address string, prikey string) error {
	sysAccounts.Lock()
	defer sysAccounts.Unlock()

	log.Info().Msgf("[ethAccount logic internal] Set current authenticator to %s", address)

	if !verifyKeyAndAddress(prikey, address) {
		err := model.ErrEthKeyAndAddressMismatch
		log.Err(err).Msgf("[ethAccount logic internal] Key and address mismatch for account %s", address)
		return err
	}

	accs, err := ethDAO.GetEthAccountsWithRole(model.EthAccountRoleAuthenticator, 0, 1000)
	if err != nil {
		log.Err(err).Msgf("[ethAccount logic internal] couldn't retrieve authenticator accounts")
		return err
	}
	match := libs.DropWhile(func(a model.EthAccount) bool { return a.Address != address }, accs)
	if len(match) < 1 {
		err = model.ErrEthAccountNotFound
		log.Err(err).Msgf("[ethAccount logic internal] authenticator %s not found", address)
		return err
	}

	if match[0].Status != "ok" {
		err = model.ErrEthAccountLocked
		log.Err(err).Msgf("[ethAccount logic internal] authenticator %s is locked", address)
		return err
	}

	sysAccounts.authenticator.Address = address
	sysAccounts.authenticator.Prikey = prikey
	sysAccounts.authenticator.Status = match[0].Status
	return nil
}

func CheckSuperAdminKey() bool {

	return true
}

func CheckAuthenticatorKey() bool {

	return true
}

//GetEthPrikeyIfExists(address string)

//SetPriKey(address string, key string)

//UnsetPrikey(address string)
