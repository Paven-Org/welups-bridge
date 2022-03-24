package welDAO

import (
	"bridge/micros/core/model"
	"bridge/service-managers/logger"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type IWelDAO interface {
	AddWelAccount(address string, status string) error
	GetAllWelAccounts(offset uint, size uint) ([]model.WelAccount, error)
	GetAllRoles() ([]string, error)
	GetWelAccount(address string) (*model.WelAccount, error)
	GetWelAccountRoles(address string) ([]string, error)
	GetWelAccountsWithRole(role string, offset uint, size uint) ([]model.WelAccount, error)
	GetWelPrikeyIfExists(address string) (string, error)
	GrantRole(address string, role string) error
	RemoveWelAccount(address string) error
	RevokeRole(address string, role string) error
	SetWelAccountStatus(address string, status string) error
	SetPriKey(address string, key string) error
	UnsetPrikey(address string) error
}

type welDAO struct {
	db *sqlx.DB
}

func MkWelDAO(db *sqlx.DB) IWelDAO {
	return &welDAO{db: db}
}

func (dao *welDAO) GetWelAccount(address string) (*model.WelAccount, error) {
	var account model.WelAccount
	db := dao.db
	log := logger.Get()

	q := db.Rebind(`SELECT wel_sys_accounts.*
									FROM wel_sys_accounts
									WHERE wel_sys_accounts.address  = ?`)
	err := db.Get(&account, q, address)
	if err != nil {
		log.Err(err).Msgf("Error while querying for account with address %s", address)
		return nil, err
	}
	// get key if exists
	qPrikey := db.Rebind("SELECT prikey FROM wel_sys_prikeys WHERE address = ?")
	var prikey string
	err = db.Get(&prikey, qPrikey, address)
	if err == sql.ErrNoRows {
		prikey = ""
	} else if err != nil {
		log.Err(err).Msgf("Error while querying for address' private key: %s", address)
		return &account, err
	}

	account.Prikey = prikey

	return &account, nil
}

func (dao *welDAO) AddWelAccount(address string, status string) error {
	db := dao.db
	log := logger.Get()

	if status == "" {
		status = "locked"
	}

	q := db.Rebind("INSERT INTO wel_sys_accounts(address, status) VALUES (?,?)")
	_, err := db.Exec(q, address, status)

	if err != nil {
		log.Err(err).Msgf("Error while inserting account %s", address)
		return err
	}

	return nil
}

func (dao *welDAO) SetPriKey(address string, key string) error {
	db := dao.db
	log := logger.Get()

	q := db.Rebind(`INSERT INTO wel_sys_prikeys (address, prikey) VALUES (?,?)`)
	_, err := db.Exec(q, address, key)
	if err != nil {
		log.Err(err).Msgf("Error while assigning private key to address %s", address)
	}
	return err
}

func (dao *welDAO) UnsetPrikey(address string) error {
	db := dao.db
	log := logger.Get()

	q := db.Rebind(`DELETE FROM wel_sys_prikeys WHERE address = ?`)
	_, err := db.Exec(q, address)
	if err != nil {
		log.Err(err).Msgf("Error while deleting private key of address %s", address)
	}
	return err
}

func (dao *welDAO) SetWelAccountStatus(address string, status string) error {
	db := dao.db
	log := logger.Get()

	q := db.Rebind("UPDATE wel_sys_accounts SET status = ?, updated_at = ? WHERE address = ?")
	_, err := db.Exec(q, status, time.Now(), address)

	if err != nil {
		log.Err(err).Msgf("Error while updating address %s", address)
		return err
	}

	return nil
}

func (dao *welDAO) GrantRole(address string, role string) error {
	db := dao.db
	log := logger.Get()

	q := db.Rebind(`INSERT INTO wel_sys_account_roles (address, role) VALUES (?,?)`)
	_, err := db.Exec(q, address, role)
	if err != nil {
		log.Err(err).Msgf("Error while granting role %s to address %s", role, address)
	}
	return err
}

func (dao *welDAO) RevokeRole(address string, role string) error {
	db := dao.db
	log := logger.Get()

	q := db.Rebind(`DELETE FROM wel_sys_account_roles 
									WHERE address = ? AND role = ?`)
	_, err := db.Exec(q, address, role)
	if err != nil {
		log.Err(err).Msgf("Error while revoking role %s from address %s", role, address)
	}
	return err
}

func (dao *welDAO) RemoveWelAccount(address string) error {
	db := dao.db
	log := logger.Get()

	tx, err := db.Beginx() // begin tx
	if err != nil {
		log.Err(err).Msgf("Unable to begin transaction when deleting address %s", address)
		return err
	}

	qDeleteEARoles := db.Rebind(`DELETE FROM wel_sys_account_roles WHERE address = ?`)
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

	qDeleteWelAccountPrikey := db.Rebind("DELETE FROM wel_sys_prikeys WHERE address = ?")
	_, err = tx.Exec(qDeleteWelAccountPrikey, address)
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

	qDeleteWelAccount := db.Rebind("DELETE FROM wel_sys_accounts WHERE address = ?")
	_, err = tx.Exec(qDeleteWelAccount, address)
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

func (dao *welDAO) GetWelPrikeyIfExists(address string) (string, error) {
	db := dao.db
	log := logger.Get()

	var prikey string

	q := db.Rebind("SELECT prikey FROM wel_sys_prikeys WHERE address = ?")
	err := db.QueryRowx(q, address).Scan(&prikey)
	if err == sql.ErrNoRows {
		return "", model.ErrWelNoPrikey
	} else if err != nil {
		log.Err(err).Msgf("Error while querying for address' private key: %s", address)
	}
	return prikey, err
}

func (dao *welDAO) GetWelAccountRoles(address string) ([]string, error) {
	db := dao.db
	log := logger.Get()

	var roles []string

	q := db.Rebind(`SELECT role FROM wel_sys_account_roles
									WHERE wel_sys_account_roles.address = ?`)
	err := db.Select(&roles, q, address)

	if err != nil {
		log.Err(err).Msgf("Error while querying for address %s's roles", address)
		return nil, err
	}

	return roles, nil
}

func (dao *welDAO) GetAllWelAccounts(offset uint, size uint) ([]model.WelAccount, error) {
	var accounts []model.WelAccount
	db := dao.db
	log := logger.Get()

	q := db.Rebind(`SELECT wel_sys_accounts.*
									FROM wel_sys_accounts
									ORDER BY wel_sys_accounts.created_at 
									OFFSET ? LIMIT ?`)
	err := db.Select(&accounts, q, offset, size)
	if err != nil {
		log.Err(err).Msgf("Error while querying for accounts")
		return nil, err
	}

	qPrikey := db.Rebind("SELECT prikey FROM wel_sys_prikeys WHERE address = ?")
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

func (dao *welDAO) GetWelAccountsWithRole(role string, offset uint, size uint) ([]model.WelAccount, error) {
	var accounts []model.WelAccount
	db := dao.db
	log := logger.Get()

	qGetAccs := db.Rebind(`SELECT wel_sys_accounts.*
									FROM wel_sys_accounts
									INNER JOIN wel_sys_account_roles
									ON wel_sys_account_roles.address = wel_sys_accounts.address
									WHERE wel_sys_account_roles.role  = ? 
									ORDER BY wel_sys_accounts.created_at 
									OFFSET ? LIMIT ?`)
	err := db.Select(&accounts, qGetAccs, role, offset, size)
	if err != nil {
		log.Err(err).Msgf("Error while querying for accounts with role %s", role)
		return nil, err
	}

	qPrikey := db.Rebind("SELECT prikey FROM wel_sys_prikeys WHERE address = ?")
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
func (dao *welDAO) GetAllRoles() ([]string, error) {
	var roles []string
	db := dao.db
	log := logger.Get()

	q := db.Rebind("SELECT role FROM wel_sys_roles")
	err := db.Select(&roles, q)

	if err != nil {
		log.Err(err).Msg("Error while querying for wel_sys_roles")
		return nil, err
	}

	return roles, nil
}
