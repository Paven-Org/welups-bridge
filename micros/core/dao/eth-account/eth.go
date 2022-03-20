package ethDAO

import (
	"bridge/micros/core/model"
	"bridge/service-managers/logger"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type IEthDAO interface {
	AddEthAccount(address string, status string) error
	GetAllEthAccounts(offset uint, size uint) ([]model.EthAccount, error)
	GetAllRoles() ([]string, error)
	GetEthAccount(address string) (*model.EthAccount, error)
	GetEthAccountRoles(address string) ([]string, error)
	GetEthAccountsWithRole(role string, offset uint, size uint) ([]model.EthAccount, error)
	GetEthPrikeyIfExists(address string) (string, error)
	GrantRole(address string, role string) error
	RemoveEthAccount(address string) error
	RevokeRole(address string, role string) error
	SetEthAccountStatus(address string, status string) error
	SetPriKey(address string, key string) error
	UnsetPrikey(address string) error
}

type ethDAO struct {
	db *sqlx.DB
}

func MkEthDAO(db *sqlx.DB) IEthDAO {
	return &ethDAO{db: db}
}

func (dao *ethDAO) GetEthAccount(address string) (*model.EthAccount, error) {
	var account model.EthAccount
	db := dao.db
	log := logger.Get()

	q := db.Rebind(`SELECT eth_sys_accounts.*, eth_sys_prikeys.prikey 
									FROM eth_sys_accounts LEFT JOIN eth_sys_prikeys 
									ON eth_sys_accounts.address = eth_sys_prikeys.address
									WHERE eth_sys_accounts.address  = ? 
									ORDER BY eth_sys_accounts.created_at`)
	err := db.Get(&account, q, address)
	if err != nil {
		log.Err(err).Msgf("Error while querying for account with address %s", address)
		return nil, err
	}
	return &account, nil
}

func (dao *ethDAO) AddEthAccount(address string, status string) error {
	db := dao.db
	log := logger.Get()

	if status == "" {
		status = "locked"
	}

	q := db.Rebind("INSERT INTO eth_sys_accounts(address, status) VALUES (?,?)")
	_, err := db.Exec(q, address, status)

	if err != nil {
		log.Err(err).Msgf("Error while inserting account %s", address)
		return err
	}

	return nil
}

func (dao *ethDAO) SetPriKey(address string, key string) error {
	db := dao.db
	log := logger.Get()

	q := db.Rebind(`INSERT INTO eth_sys_prikeys (address, prikey) VALUES (?,?)`)
	_, err := db.Exec(q, address, key)
	if err != nil {
		log.Err(err).Msgf("Error while assigning private key to address %s", address)
	}
	return err
}

func (dao *ethDAO) UnsetPrikey(address string) error {
	db := dao.db
	log := logger.Get()

	q := db.Rebind(`DELETE FROM eth_sys_prikeys WHERE address = ?`)
	_, err := db.Exec(q, address)
	if err != nil {
		log.Err(err).Msgf("Error while deleting private key of address %s", address)
	}
	return err
}

func (dao *ethDAO) SetEthAccountStatus(address string, status string) error {
	db := dao.db
	log := logger.Get()

	q := db.Rebind("UPDATE eth_sys_accounts SET status = ?, updated_at = ? WHERE address = ?")
	_, err := db.Exec(q, status, time.Now(), address)

	if err != nil {
		log.Err(err).Msgf("Error while updating address %s", address)
		return err
	}

	return nil
}

func (dao *ethDAO) GrantRole(address string, role string) error {
	db := dao.db
	log := logger.Get()

	q := db.Rebind(`INSERT INTO eth_sys_account_roles (address, role) VALUES (?,?)`)
	_, err := db.Exec(q, address, role)
	if err != nil {
		log.Err(err).Msgf("Error while granting role %s to address %s", role, address)
	}
	return err
}

func (dao *ethDAO) RevokeRole(address string, role string) error {
	db := dao.db
	log := logger.Get()

	q := db.Rebind(`DELETE FROM eth_sys_account_roles 
									WHERE address = ? AND role = ?`)
	_, err := db.Exec(q, address, role)
	if err != nil {
		log.Err(err).Msgf("Error while revoking role %s from address %s", role, address)
	}
	return err
}

func (dao *ethDAO) RemoveEthAccount(address string) error {
	db := dao.db
	log := logger.Get()

	tx, err := db.Beginx() // begin tx
	if err != nil {
		log.Err(err).Msgf("Unable to begin transaction when deleting address %s", address)
		return err
	}

	qDeleteEARoles := db.Rebind(`DELETE FROM eth_sys_account_roles WHERE address = ?`)
	_, err = tx.Exec(qDeleteEARoles, address)
	if err != nil {
		log.Err(err).Msgf("Error while deleting address %s", address)
		for {
			if err := tx.Rollback(); err != nil {
				log.Err(err).Msg("Error while rolling back tx, retrying...")
			} else {
				break
			}
		}
		return err
	}

	qDeleteEthAccountPrikey := db.Rebind("DELETE FROM eth_sys_prikeys WHERE address = ?")
	_, err = tx.Exec(qDeleteEthAccountPrikey, address)
	if err != nil {
		log.Err(err).Msgf("Error while deleting address %s", address)
		for {
			if err := tx.Rollback(); err != nil {
				log.Err(err).Msg("Error while rolling back tx, retrying...")
			} else {
				break
			}
		}
		return err
	}

	qDeleteEthAccount := db.Rebind("DELETE FROM eth_sys_accounts WHERE address = ?")
	_, err = tx.Exec(qDeleteEthAccount, address)
	if err != nil {
		log.Err(err).Msgf("Error while deleting address %s", address)
		for {
			if err := tx.Rollback(); err != nil {
				log.Err(err).Msg("Error while rolling back tx, retrying...")
			} else {
				break
			}
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Err(err).Msgf("Error while deleting address %s", address)
		for {
			if err := tx.Rollback(); err != nil {
				log.Err(err).Msg("Error while rolling back tx, retrying...")
			} else {
				break
			}
		}
		return err
	}

	return nil
}

func (dao *ethDAO) GetEthPrikeyIfExists(address string) (string, error) {
	db := dao.db
	log := logger.Get()

	var prikey string

	q := db.Rebind("SELECT prikey FROM eth_sys_prikeys WHERE address = ?")
	err := db.QueryRowx(q, address).Scan(&prikey)
	if err == sql.ErrNoRows {
		return "", model.ErrEthNoPrikey
	} else if err != nil {
		log.Err(err).Msgf("Error while querying for address' private key: %s", address)
	}
	return prikey, err
}

func (dao *ethDAO) GetEthAccountRoles(address string) ([]string, error) {
	db := dao.db
	log := logger.Get()

	var roles []string

	q := db.Rebind(`SELECT role FROM eth_sys_account_roles
									WHERE eth_sys_account_roles.address = ?`)
	err := db.Select(&roles, q, address)

	if err != nil {
		log.Err(err).Msgf("Error while querying for address %s's roles", address)
		return nil, err
	}

	return roles, nil
}

func (dao *ethDAO) GetAllEthAccounts(offset uint, size uint) ([]model.EthAccount, error) {
	var accounts []model.EthAccount
	db := dao.db
	log := logger.Get()

	q := db.Rebind(`SELECT eth_sys_accounts.*
									FROM eth_sys_accounts
									ORDER BY eth_sys_accounts.created_at 
									OFFSET ? LIMIT ?`)
	err := db.Select(&accounts, q, offset, size)
	if err != nil {
		log.Err(err).Msgf("Error while querying for accounts")
		return nil, err
	}

	qPrikey := db.Rebind("SELECT prikey FROM eth_sys_prikeys WHERE address = ?")
	var prikey string
	for i, acc := range accounts {
		err := db.QueryRowx(qPrikey, acc.Address).Scan(&prikey)
		if err == sql.ErrNoRows {
			accounts[i].Prikey = ""
		} else if err != nil {
			log.Err(err).Msgf("Error while querying for address' private key: %s", acc.Address)
			return nil, err
		}
		accounts[i].Prikey = prikey
	}
	return accounts, nil
}

func (dao *ethDAO) GetEthAccountsWithRole(role string, offset uint, size uint) ([]model.EthAccount, error) {
	var accounts []model.EthAccount
	db := dao.db
	log := logger.Get()

	qGetAccs := db.Rebind(`SELECT eth_sys_accounts.*
									FROM eth_sys_accounts
									INNER JOIN eth_sys_account_roles
									ON eth_sys_account_roles.address = eth_sys_accounts.address
									WHERE eth_sys_account_roles.role  = ? 
									ORDER BY eth_sys_accounts.created_at 
									OFFSET ? LIMIT ?`)
	err := db.Select(&accounts, qGetAccs, role, offset, size)
	if err != nil {
		log.Err(err).Msgf("Error while querying for accounts with role %s", role)
		return nil, err
	}

	qPrikey := db.Rebind("SELECT prikey FROM eth_sys_prikeys WHERE address = ?")
	var prikey string
	for i, acc := range accounts {
		err := db.QueryRowx(qPrikey, acc.Address).Scan(&prikey)
		if err == sql.ErrNoRows {
			accounts[i].Prikey = ""
		} else if err != nil {
			log.Err(err).Msgf("Error while querying for address' private key: %s", acc.Address)
			return nil, err
		}
		accounts[i].Prikey = prikey
	}

	return accounts, nil
}

// again, there won't be that many roles, so pagination isn't even required
func (dao *ethDAO) GetAllRoles() ([]string, error) {
	var roles []string
	db := dao.db
	log := logger.Get()

	q := db.Rebind("SELECT role FROM eth_sys_roles")
	err := db.Select(&roles, q)

	if err != nil {
		log.Err(err).Msg("Error while querying for eth_sys_roles")
		return nil, err
	}

	return roles, nil
}
