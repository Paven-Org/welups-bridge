package dao

import (
	"github.com/jmoiron/sqlx"
)

// sort of a locator for DAOs
type DAOs struct {
}

func MkDAOs(db *sqlx.DB) *DAOs {
	return &DAOs{}
}
