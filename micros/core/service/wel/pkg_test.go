package welService

import (
	"bridge/micros/core/config"
	"bridge/micros/core/dao"
	userdao "bridge/micros/core/dao/user"
	"bridge/micros/core/model"
	manager "bridge/service-managers"
	"bridge/service-managers/logger"
	"context"
	"fmt"
	"testing"

	welclient "github.com/Clownsss/gotron-sdk/pkg/client"

	"github.com/DATA-DOG/go-txdb"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"go.temporal.io/sdk/client"
)

const (
	testkey  = "ce0d51b2062e5694d28a21ad64b7efd583856ba20afe437ae4c4ad7d7a5ae34a"
	testAddr = "WDReBjymEzH5Bi4avyfUqXa4rhyvAT7DY2"
)

var GovService *GovContractService

var (
	log     *zerolog.Logger
	userDAO userdao.IUserDAO
	tempcli client.Client
)

func TestMain(m *testing.M) {
	config.Load()
	cnf := config.Get()
	dbCnf := cnf.DBconfig
	logger.Init(false)
	log = logger.Get()

	connString := fmt.Sprintf("host='%s' port=%d user='%s' password='%s' dbname='%s' sslmode=%s", dbCnf.Host, dbCnf.Port, dbCnf.Username, dbCnf.Password, dbCnf.DBname, dbCnf.SSLMode)

	// mock DB
	txdb.Register("psql_txdb", "postgres", connString)
	sqlx.BindDriver("psql_txdb", sqlx.DOLLAR)
	db, _ := sqlx.Open("psql_txdb", "test")
	defer db.Close()
	daos := dao.MkDAOs(db)
	userDAO = daos.User
	//ethDAO := daos.Eth
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

	GovService, err = MkGovContractService(welCli, tempcli, daos, cnf.WelGovContract)
	if err != nil {
		logger.Get().Err(err).Msgf("Unable to initialize GovContractService")
		return
	}

	GovService.StartService()
	defer GovService.StopService()
	m.Run()

}

func TestGrantRoleOnContract(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "callerkey", testkey)
	tx, err := GovService.GrantRoleOnContract(ctx, testAddr, "AUTHENTICATOR")
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
	fmt.Println(tx)
}

func TestRevokeRoleOnContract(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "callerkey", testkey)
	tx, err := GovService.RevokeRoleOnContract(ctx, testAddr, "AUTHENTICATOR")
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
	fmt.Println(tx)
}

func TestHasRole(t *testing.T) {
	res, err := GovService.HasRole(context.Background(), testAddr, "AUTHENTICATOR")
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
	fmt.Println(res)
}

func TestGrantRoleWF(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "callerkey", testkey)
	wo := client.StartWorkflowOptions{
		TaskQueue: GovContractQueue,
	}

	we, err := tempcli.ExecuteWorkflow(ctx, wo, GrantRoleWorkflow, testAddr, "AUTHENTICATOR")
	if err != nil {
		t.Fatal("Unable to call GrantRoleWorkflow: ", err.Error())
	}
	log.Info().Str("Workflow", we.GetID()).Str("runID=", we.GetRunID()).Msg("dispatched")

	var txhash string
	if err := we.Get(ctx, &txhash); err != nil {
		t.Fatal("GrantRoleWorkflow failed: ", err.Error())
	}
	fmt.Println(txhash)
}

func TestRevokeRoleWF(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "callerkey", testkey)
	wo := client.StartWorkflowOptions{
		TaskQueue: GovContractQueue,
	}

	we, err := tempcli.ExecuteWorkflow(ctx, wo, RevokeRoleWorkflow, testAddr, "AUTHENTICATOR")
	if err != nil {
		t.Fatal("Unable to call GrantRoleWorkflow: ", err.Error())
	}
	log.Info().Str("Workflow", we.GetID()).Str("runID=", we.GetRunID()).Msg("dispatched")

	var txhash string
	if err := we.Get(ctx, &txhash); err != nil {
		t.Fatal("GrantRoleWorkflow failed: ", err.Error())
	}
	fmt.Println(txhash)
}
func TestRoleBytes(t *testing.T) {
	var b [32]byte
	copy(b[:], common.Hex2Bytes("0x00"))
	fmt.Println("admin: ", b)
	fmt.Printf("a role: %x\n", crypto.Keccak256([]byte(model.EthAccountRoleAuthenticator)))
}
