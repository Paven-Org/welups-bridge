package ethService

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

	"github.com/DATA-DOG/go-txdb"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"go.temporal.io/sdk/client"
)

const testkey = "ce0d51b2062e5694d28a21ad64b7efd583856ba20afe437ae4c4ad7d7a5ae34a"

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

	ethCli, err := ethclient.Dial(cnf.EthereumConfig.BlockchainRPC)
	if err != nil {
		logger.Get().Err(err).Msgf("Unable to connect to ethereum RPC server")
		return
	}
	defer ethCli.Close()

	GovService, err = MkGovContractService(ethCli, tempcli, daos, cnf.EthGovContract)
	if err != nil {
		logger.Get().Err(err).Msgf("Unable to initialize GovContractService")
		return
	}

	GovService.StartService()
	defer GovService.StopService()
	m.Run()

}

func TestGrantRole(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "callerkey", testkey)
	tx, err := GovService.GrantRoleOnContract(ctx, "0x25e8370E0e2cf3943Ad75e768335c892434bD090", "AUTHENTICATOR")
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
	fmt.Println(tx)
}

func TestFilterGranted(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "callerkey", testkey)
	res, err := GovService.FilterRoleGranted(ctx, "0x25e8370E0e2cf3943Ad75e768335c892434bD090", "AUTHENTICATOR")
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
	fmt.Println(res)
}

func TestRevokeRole(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "callerkey", testkey)
	tx, err := GovService.RevokeRoleOnContract(ctx, "0x25e8370E0e2cf3943Ad75e768335c892434bD090", "AUTHENTICATOR")
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
	fmt.Println(tx)
}

func TestFilterRevoked(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "callerkey", testkey)
	res, err := GovService.FilterRoleRevoked(ctx, "0x25e8370E0e2cf3943Ad75e768335c892434bD090", "AUTHENTICATOR")
	if err != nil {
		t.Fatal("Error: ", err.Error())
	}
	fmt.Println(res)
}

func TestHasRole(t *testing.T) {
	res, err := GovService.HasRole(context.Background(), "0x25e8370E0e2cf3943Ad75e768335c892434bD090", "AUTHENTICATOR")
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

	we, err := tempcli.ExecuteWorkflow(ctx, wo, GrantRoleWorkflow, "0x25e8370E0e2cf3943Ad75e768335c892434bD090", "AUTHENTICATOR")
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

	we, err := tempcli.ExecuteWorkflow(ctx, wo, RevokeRoleWorkflow, "0x25e8370E0e2cf3943Ad75e768335c892434bD090", "AUTHENTICATOR")
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
