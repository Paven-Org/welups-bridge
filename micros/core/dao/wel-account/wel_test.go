package welDAO

import (
	"bridge/micros/core/config"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var welDao = &welDAO{}

func TestMain(m *testing.M) {
	config.Load()
	cnf := config.Get().DBconfig

	connString := fmt.Sprintf("host='%s' port=%d user='%s' password='%s' dbname='%s' sslmode=%s", cnf.Host, cnf.Port, cnf.Username, cnf.Password, cnf.DBname, cnf.SSLMode)

	// mock DB
	txdb.Register("psql_txdb", "postgres", connString)
	sqlx.BindDriver("psql_txdb", sqlx.DOLLAR)
	db, _ := sqlx.Open("psql_txdb", "test")
	defer db.Close()

	// DAOs initialization
	welDao.db = db

	m.Run()
}

func TestAddWelAccount(t *testing.T) {
	if err := welDao.AddWelAccount("efef", "ok"); err != nil {
		t.Fatal("Error: ", err.Error())
	}
}
func TestGetAllWelAccounts(t *testing.T) {
	if accounts, err := welDao.GetAllWelAccounts(0, 20); err != nil {
		t.Fatal("Error: ", err.Error())
	} else {
		fmt.Println("All accounts: ", accounts)
	}
}

func TestGetAllRoles(t *testing.T) {
	if roles, err := welDao.GetAllRoles(); err != nil {
		t.Fatal("Error: ", err.Error())
	} else {
		fmt.Println("All roles: ", roles)
	}
}

func TestGetWelAccount(t *testing.T) {
	if account, err := welDao.GetWelAccount("efef"); err != nil {
		t.Fatal("Error: ", err.Error())
	} else {
		fmt.Println("Account: ", account)
	}
}

func TestGetWelAccountRoles(t *testing.T) {
	if roles, err := welDao.GetWelAccountRoles("efef"); err != nil {
		t.Fatal("Error: ", err.Error())
	} else {
		fmt.Println("roles: ", roles)
	}
}

func TestGetWelAccountsWithRole(t *testing.T) {
	if accounts, err := welDao.GetWelAccountsWithRole("super_admin", 0, 20); err != nil {
		t.Fatal("Error: ", err.Error())
	} else {
		fmt.Println("All accounts with role 'super_admin': ", accounts)
	}
}
func TestGetWelPrikeyIfExists(t *testing.T) {}
func TestGrantRole(t *testing.T)            {}
func TestRemoveWelAccount(t *testing.T)     {}
func TestRevokeRole(t *testing.T)           {}
func TestSetWelAccountStatus(t *testing.T)  {}
func TestSetPriKey(t *testing.T)            {}
func TestUnsetPrikey(t *testing.T)          {}
