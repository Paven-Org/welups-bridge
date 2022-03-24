package ethLogic

import (
	"bridge/common"
	"bridge/micros/core/config"
	"bridge/micros/core/dao"
	ethService "bridge/micros/core/service/eth"
	manager "bridge/service-managers"
	"bridge/service-managers/logger"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	testkey  = "ce0d51b2062e5694d28a21ad64b7efd583856ba20afe437ae4c4ad7d7a5ae34a"
	testaddr = "0x25e8370E0e2cf3943Ad75e768335c892434bD090"
)

var GovService *ethService.GovContractService

func TestMain(m *testing.M) {
	mCnf := common.Mailerconf{
		SmtpHost: "smtp.gmail.com",
		SmtpPort: 587,
		Address:  "bridgemail.welups@gmail.com",
		Password: "showmethemoney11!1",
	}
	mailer = manager.MkMailer(mCnf)

	config.Load()
	cnf := config.Get()
	dbCnf := cnf.DBconfig
	log = logger.Get()

	connString := fmt.Sprintf("host='%s' port=%d user='%s' password='%s' dbname='%s' sslmode=%s", dbCnf.Host, dbCnf.Port, dbCnf.Username, dbCnf.Password, dbCnf.DBname, dbCnf.SSLMode)

	// mock DB
	txdb.Register("psql_txdb", "postgres", connString)
	sqlx.BindDriver("psql_txdb", sqlx.DOLLAR)
	db, _ := sqlx.Open("psql_txdb", "test")
	defer db.Close()
	daos := dao.MkDAOs(db)
	userDAO = daos.User
	ethDAO = daos.Eth

	// temporal
	tcli, err := manager.MkTemporalClient(cnf.TemporalCliConfig, []string{"callerkey", "signerkey"})
	if err != nil {
		logger.Get().Err(err).Msgf("Unable to connect to temporal backend")
		return
	}
	defer tcli.Close()

	tempcli = tcli

	ethCli, err := ethclient.Dial(cnf.EthereumConfig.BlockchainRPC)
	if err != nil {
		logger.Get().Err(err).Msgf("Unable to connect to ethereum RPC server")
		return
	}
	defer ethCli.Close()

	GovService, err = ethService.MkGovContractService(ethCli, tempcli, daos, cnf.EthGovContract)
	if err != nil {
		logger.Get().Err(err).Msgf("Unable to initialize GovContractService")
		return
	}

	GovService.StartService()
	defer GovService.StopService()
	m.Run()

}

func TestKeyAndAddress(t *testing.T) {
	fmt.Println(verifyAddress("0x25e8370E0e2cf3943Ad75e768335c892434bD090"))
	fmt.Println(verifyKeyAndAddress("ce0d51b2062e5694d28a21ad64b7efd583856ba20afe437ae4c4ad7d7a5ae34a", "0x25e8370E0e2cf3943Ad75e768335c892434bD090"))
}

func TestSetAuthenticator(t *testing.T) {
	ethDAO.GrantRole("0x25e8370E0e2cf3943Ad75e768335c892434bD090", "AUTHENTICATOR")
	SetCurrentAuthenticator(testkey)
	fmt.Println(sysAccounts.authenticator)
	ethDAO.RevokeRole("0x25e8370E0e2cf3943Ad75e768335c892434bD090", "AUTHENTICATOR")
}

func TestGrantRole(t *testing.T) {
	tx, err := GrantRole(testaddr, "MANAGER_ROLE", testkey)
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
	fmt.Println(tx)
}

func TestRevokeRole(t *testing.T) {
	tx, err := RevokeRole(testaddr, "MANAGER_ROLE", testkey)
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
	fmt.Println(tx)
}

//func TestSendMailToAdmins(t *testing.T) {
//	if err := sendNotificationToRole("admin", "test", "test mail"+libs.Uniq()); err != nil {
//		t.Fatal("Error: ", err.Error())
//	}
//}
