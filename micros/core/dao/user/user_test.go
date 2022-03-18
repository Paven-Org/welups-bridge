package userDAO

import (
	"bridge/micros/core/config"
	"bridge/micros/core/model"
	"bridge/service-managers/logger"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var userDao = &userDAO{}

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
	userDao.db = db

	m.Run()
}

func TestAddUser(t *testing.T) {
	id, err := userDao.AddUser("abc", "nomail", "nopass")
	logger.Get().Info().Msgf("New user id: %d", id)
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
	logger.Get().Info().Msg("User abc added")

	u, err := userDao.GetUserByName("abc")
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
	logger.Get().Info().Msgf("Found user abc: %v", u)

	err = userDao.GrantRole("abc", "admin")
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
	logger.Get().Info().Msg("User abc granted role admin")

	err = userDao.UpdateUser(&model.User{Id: u.Id, Username: "def", Email: "def@abc.com", Status: u.Status, Password: u.Password})
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
	logger.Get().Info().Msg("User updated")

	err = userDao.RemoveUser("def")
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
	logger.Get().Info().Msg("User removed")

	u, err = userDao.GetUserByName("def")
	if err != nil {
		t.Log("Error: ", err.Error())
	}
	logger.Get().Info().Msgf("user def should not be found: %v", u)
}

func TestGetUserByID(t *testing.T) {
	u, err := userDao.GetUserById(1)
	logger.Get().Info().Msgf("user: %v", u)
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
}

func TestGetUserByName(t *testing.T) {
	u, err := userDao.GetUserByName("root")
	logger.Get().Info().Msgf("user: %v", u)
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
}

func TestGetUserByEmail(t *testing.T) {
	u, err := userDao.GetUserByEmail("nhatanh02@gmail.com")
	logger.Get().Info().Msgf("user: %v", u)
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
}

func TestGetUserRoles(t *testing.T) {
	id, err := userDao.AddUser("abc", "nomail", "nopass")
	logger.Get().Info().Msgf("New user id: %d", id)
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
	logger.Get().Info().Msg("User abc added")

	err = userDao.GrantRole("abc", "admin")
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
	logger.Get().Info().Msg("User abc granted role admin")

	err = userDao.GrantRole("abc", "root")
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
	logger.Get().Info().Msg("User abc granted role root")

	roles, err := userDao.GetUserRoles("abc")
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
	logger.Get().Info().Msgf("user roles: %v", roles)

	err = userDao.RevokeRole("abc", "root")
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
	logger.Get().Info().Msg("User abc revoked role root")

	roles, err = userDao.GetUserRoles("abc")
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
	logger.Get().Info().Msgf("user roles: %v", roles)
}

func TestGetUsers(t *testing.T) {
	users, err := userDao.GetUsers(0, 100)
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
	fmt.Println("Users got: ", users)
	users, err = userDao.GetUsersWithRole("root", 0, 100)
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
	fmt.Println("Users with role \"root\" got: ", users)
	roles, err := userDao.GetUserRoles("root")
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
	logger.Get().Info().Msgf("user roles: %v", roles)
}
