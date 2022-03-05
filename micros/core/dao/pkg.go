package dao

import (
	userDAO "bridge/micros/core/dao/user"

	"github.com/jmoiron/sqlx"
)

// sort of a locator for DAOs
type DAOs struct {
	User userDAO.IUserDAO
}

func MkDAOs(db *sqlx.DB) *DAOs {
	return &DAOs{
		User: userDAO.MkUserDAO(db),
	}
}
