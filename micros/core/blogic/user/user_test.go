package userLogic

import (
	"bridge/libs"
	"bridge/logger"
	"bridge/micros/core/config"
	"bridge/micros/core/dao"
	manager "bridge/service-managers"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func TestMain(m *testing.M) {
	config.Load()
	// should change to some test DB
	cnf := config.Get().DBconfig

	connString := fmt.Sprintf("host='%s' port=%d user='%s' password='%s' dbname='%s' sslmode=%s", cnf.Host, cnf.Port, cnf.Username, cnf.Password, cnf.DBname, cnf.SSLMode)

	txdb.Register("psql_txdb", "postgres", connString)
	sqlx.BindDriver("psql_txdb", sqlx.DOLLAR)
	db, _ := sqlx.Open("psql_txdb", "test")
	defer db.Close()

	daos := dao.MkDAOs(db)

	rm := manager.MkRedisManager(
		config.Get().RedisConfig,
		map[string]int{
			manager.StdAuthDBName: 15,
		})
	defer func() {
		rm.Flush(manager.StdAuthDBName)
		rm.CloseAll()
	}()

	ts := libs.MkTokenServ(config.Get().Secrets.JwtSecret)

	Init(daos, rm, ts)

	m.Run()
}

func TestLoginLogout(t *testing.T) {
	tk, ssid, ss, err := Login("root", "root")
	if err != nil {
		t.Fatalf("Root login failed, error: %s", err.Error())
	}
	logger.Get().Info().Msg("Root logged in")
	logger.Get().Info().Msgf("Token generated: %s", tk)
	logger.Get().Info().Msgf("sessionID: %s", ssid)
	logger.Get().Info().Msgf("Secret generated: %s", ss)

	claims, _ := ParseTokenToClaims(tk)
	roles, err := GetUserRoles(claims.Uid)
	if err != nil {
		t.Fatalf("Getting root's roles failed, error: %s", err.Error())
	}
	logger.Get().Info().Msgf("Roles: %v", roles)

	if err := Logout(tk, ss); err != nil {
		t.Fatalf("Root logout failed, error: %s", err.Error())
	}

}
