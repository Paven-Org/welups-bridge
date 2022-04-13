package dao

import (
	"bridge/micros/core/dao/blockscan"
	ethDAO "bridge/micros/core/dao/eth-account"
	userDAO "bridge/micros/core/dao/user"
	welDAO "bridge/micros/core/dao/wel-account"

	"github.com/jmoiron/sqlx"
)

// sort of a locator for DAOs
type DAOs struct {
	User        userDAO.IUserDAO
	Eth         ethDAO.IEthDAO
	Wel         welDAO.IWelDAO
	EthBlockDAO *blockscan.EthSysDAO
	WelBlockDAO *blockscan.WelSysDAO
}

func MkDAOs(db *sqlx.DB) *DAOs {
	return &DAOs{
		User:        userDAO.MkUserDAO(db),
		Eth:         ethDAO.MkEthDAO(db),
		Wel:         welDAO.MkWelDAO(db),
		EthBlockDAO: blockscan.MkEthSysDao(db),
		WelBlockDAO: blockscan.MkWelSysDao(db),
	}
}
