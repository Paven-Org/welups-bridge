package welLogic

import (
	"bridge/libs"
	welABI "bridge/micros/core/abi/wel"
	"bridge/micros/core/config"
	"bridge/micros/core/dao"
	"bridge/micros/core/model"
	welService "bridge/micros/core/service/wel"
	manager "bridge/service-managers"
	"bridge/service-managers/logger"
	"fmt"
	"testing"

	welclient "github.com/Paven-Org/gotron-sdk/pkg/client"

	"github.com/DATA-DOG/go-txdb"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	testkey     = "ce0d51b2062e5694d28a21ad64b7efd583856ba20afe437ae4c4ad7d7a5ae34a"
	testaddr    = "WDReBjymEzH5Bi4avyfUqXa4rhyvAT7DY2"
	testaddrhex = "0x4125e8370E0e2cf3943Ad75e768335c892434bD090"
)

var GovService *welService.GovContractService

func TestMain(m *testing.M) {
	//mCnf := common.Mailerconf{
	//	SmtpHost: "smtp.gmail.com",
	//	SmtpPort: 587,
	//	Address:  "bridgemail.welups@gmail.com",
	//	Password: "showmethemoney11!1",
	//}
	//mailer = manager.MkMailer(mCnf)

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
	welDAO = daos.Wel
	//ethDAO.GrantRole("0x25e8370E0e2cf3943Ad75e768335c892434bD090", "AUTHENTICATOR")

	// temporal
	tcli, err := manager.MkTemporalClient(cnf.TemporalCliConfig, []string{"callerkey", "signerkey"})
	if err != nil {
		logger.Get().Err(err).Msgf("Unable to connect to temporal backend")
		return
	}
	defer tcli.Close()

	tempcli = tcli

	welCli := welclient.NewGrpcClient(cnf.WelupsConfig.Nodes[0])
	defer welCli.Stop()
	if err := welCli.Start(); err != nil {
		logger.Get().Err(err).Msgf("Unable to start welCli's GRPC connection")
		return
	}

	exportContr := cnf.WelExportContract
	welExp = welABI.MkWelExport(welCli, exportContr)

	GovService, err = welService.MkGovContractService(welCli, tempcli, daos, cnf.WelGovContract)
	if err != nil {
		logger.Get().Err(err).Msgf("Unable to initialize GovContractService")
		return
	}

	GovService.StartService()
	defer GovService.StopService()
	m.Run()

}

func TestKeyAndAddress(t *testing.T) {
	fmt.Println(verifyAddress(testaddr))
	fmt.Println(verifyKeyAndAddress(testkey, testaddr))

	b58, err := libs.KeyToB58Addr(testkey)
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
	fmt.Println(b58)

	hex, err := libs.KeyToWelHexAddr(testkey)
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
	fmt.Println(hex)
}

func TestSetAuthenticator(t *testing.T) {
	SetCurrentAuthenticator(testkey)
	fmt.Println(sysAccounts.authenticator)
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

func TestInvalidate(t *testing.T) {
	reqID := "36520334248965224490069560844488943606812912433996205144170613492011902220912"
	tokenAddr := "W9yD14Nj9j7xAB4dbGeiX9h8unkKHxuTtb"
	testAddr := "WDReBjymEzH5Bi4avyfUqXa4rhyvAT7DY2"
	testKey := "ce0d51b2062e5694d28a21ad64b7efd583856ba20afe437ae4c4ad7d7a5ae34a"
	contractVersion := "EXPORT_WELS_v1"
	//pkey, err := crypto.HexToECDSA(testKey)
	//if err != nil {
	//	t.Fatal("Error: ", err.Error())
	//}
	sysAccounts.authenticator.Address = testAddr
	sysAccounts.authenticator.Prikey = testKey
	sysAccounts.authenticator.Status = model.WelAccountStatusOK

	if err := InvalidateRequestClaim(tokenAddr, "0", reqID, contractVersion); err != nil {
		t.Fatal("Error: ", err.Error())
	}
}

//func TestSendMailToAdmins(t *testing.T) {
//	if err := sendNotificationToRole("admin", "test", "test mail"+libs.Uniq()); err != nil {
//		t.Fatal("Error: ", err.Error())
//	}
//}
