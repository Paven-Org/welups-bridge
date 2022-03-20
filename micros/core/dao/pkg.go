package dao

import (
	ethDAO "bridge/micros/core/dao/eth-account"
	userDAO "bridge/micros/core/dao/user"

	"github.com/jmoiron/sqlx"
)

// sort of a locator for DAOs
type DAOs struct {
	User userDAO.IUserDAO
	Eth  ethDAO.IEthDAO
}

func MkDAOs(db *sqlx.DB) *DAOs {
	return &DAOs{
		User: userDAO.MkUserDAO(db),
		Eth:  ethDAO.MkEthDAO(db),
	}
}
