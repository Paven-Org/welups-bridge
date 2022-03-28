package dao

import (
	"github.com/jmoiron/sqlx"
)

// sort of a locator for DAOs
type DAOs struct {
	TransDAO  IWelEthTransDAO
	EthSysDAO *ethSysDAO
	WelSysDAO *welSysDAO
}

func MkDAOs(db *sqlx.DB) *DAOs {
	return &DAOs{TransDAO: MkWelEthTransDao(db),
		EthSysDAO: MkEthSysDao(db),
		WelSysDAO: MkWelSysDao(db)}
}
