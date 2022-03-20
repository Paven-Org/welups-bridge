package ethDAO

import (
	"bridge/micros/core/config"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var ethDao = &ethDAO{}

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
	ethDao.db = db

	m.Run()
}

func TestAddEthAccount(t *testing.T) {
	if err := ethDao.AddEthAccount("efef", "ok"); err != nil {
		t.Fatal("Error: ", err.Error())
	}
}
func TestGetAllEthAccounts(t *testing.T) {
	if accounts, err := ethDao.GetAllEthAccounts(0, 20); err != nil {
		t.Fatal("Error: ", err.Error())
	} else {
		fmt.Println("All accounts: ", accounts)
	}
}

func TestGetAllRoles(t *testing.T) {
	if roles, err := ethDao.GetAllRoles(); err != nil {
		t.Fatal("Error: ", err.Error())
	} else {
		fmt.Println("All roles: ", roles)
	}
}

func TestGetEthAccount(t *testing.T) {
	if account, err := ethDao.GetEthAccount("1"); err != nil {
		t.Fatal("Error: ", err.Error())
	} else {
		fmt.Println("Account: ", account)
	}
}

func TestGetEthAccountRoles(t *testing.T) {
	if roles, err := ethDao.GetEthAccountRoles("1"); err != nil {
		t.Fatal("Error: ", err.Error())
	} else {
		fmt.Println("roles: ", roles)
	}
}

func TestGetEthAccountsWithRole(t *testing.T) {
	if accounts, err := ethDao.GetEthAccountsWithRole("vault", 0, 20); err != nil {
		t.Fatal("Error: ", err.Error())
	} else {
		fmt.Println("All accounts with role 'vault': ", accounts)
	}
}
func TestGetEthPrikeyIfExists(t *testing.T) {}
func TestGrantRole(t *testing.T)            {}
func TestRemoveEthAccount(t *testing.T)     {}
func TestRevokeRole(t *testing.T)           {}
func TestSetEthAccountStatus(t *testing.T)  {}
func TestSetPriKey(t *testing.T)            {}
func TestUnsetPrikey(t *testing.T)          {}
