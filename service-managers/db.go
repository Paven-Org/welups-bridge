package manager

import (
	"bridge/common"
	"bridge/service-managers/logger"
	"fmt"

	"github.com/jmoiron/sqlx"
)

//var db *sqlx.DB

//func Get() *sqlx.DB {
//	if db == nil {
//		return Init()
//	}
//	return db
//}

func MkDB(cnf common.DBconf) (*sqlx.DB, error) {
	var err error
	var dblog = logger.Get()

	//if db != nil {
	//	dblog.Info().Msg("db object already initialized")
	//	return db
	//}

	dblog.Info().Msg("Init db object")

	connString := fmt.Sprintf("host='%s' port=%d user='%s' password='%s' dbname='%s' sslmode=%s",
		cnf.Host, cnf.Port, cnf.Username, cnf.Password, cnf.DBname, cnf.SSLMode)
	dblog.Info().Msgf("connString: %s", connString)
	db, err := sqlx.Connect(cnf.DBbackend, connString)
	if err != nil {
		dblog.Err(err).Msg("Unable to connect to database, error: ")
		return db, err
	}

	if err = db.Ping(); err != nil {
		dblog.Err(err).Msg("Database pint failed, error: ")
		panic("db.Init() failed")
		return db, err
	}

	db.SetConnMaxIdleTime(cnf.ConnMaxLifetime)
	db.SetMaxOpenConns(cnf.MaxOpenConns)
	db.SetMaxIdleConns(cnf.MaxIdleConns)

	dblog.Info().Msg("db object initialized")
	return db, err
}
