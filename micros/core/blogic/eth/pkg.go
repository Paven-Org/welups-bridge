package ethLogic

import (
	"bridge/micros/core/config"
	"bridge/micros/core/dao"
	ethdao "bridge/micros/core/dao/eth-account"
	userdao "bridge/micros/core/dao/user"
	"bridge/micros/core/model"
	"bridge/micros/core/service/notifier"
	"bridge/service-managers/logger"
	"context"
	"math/big"
	"regexp"
	"strings"
	"sync"

	"bridge/micros/core/abi/eth"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
	"go.temporal.io/sdk/client"
)

var (
	ethDAO  ethdao.IEthDAO
	userDAO userdao.IUserDAO
	//mailer  *manager.Mailer
	tempcli client.Client
	log     *zerolog.Logger
	importC *importContract
)

type importContract struct {
	impC         *eth.EthImportC
	cli          *ethclient.Client
	lastGasPrice *big.Int
}

func Init(d *dao.DAOs, tmpcli client.Client, ethcli *ethclient.Client) {
	log = logger.Get()
	ethDAO = d.Eth
	userDAO = d.User

	importContractAddress := common.HexToAddress(config.Get().EthImportContract)
	impC, err := eth.NewEthImportC(importContractAddress, ethcli)
	if err != nil {
		log.Err(err).Msg("Unable to initialize ethLogic")
		panic(err)
	}
	importC = &importContract{
		impC:         impC,
		cli:          ethcli,
		lastGasPrice: big.NewInt(1000000000),
	}
	//	mailer = m
	tempcli = tmpcli
	SetCurrentAuthenticator("ce0d51b2062e5694d28a21ad64b7efd583856ba20afe437ae4c4ad7d7a5ae34a") // quick and dirty for now

	if problem := Healthcheck(); problem != nil {
		ctx := context.Background()
		wo := client.StartWorkflowOptions{
			TaskQueue: notifier.NotifierQueue,
		}

		we, err := tempcli.ExecuteWorkflow(ctx, wo, notifier.NotifyProblemWF, problem.Error(), "admin")
		if err != nil {
			log.Err(err).Msg("[Eth logic init] Failed to notify admins of problem: " + problem.Error())
			return
		}
		log.Info().Str("Workflow", we.GetID()).Str("runID=", we.GetRunID()).Msg("dispatched")
		if err := we.Get(ctx, nil); err != nil {
			log.Err(err).Msg("[Eth logic init] Failed to notify admins of problem: " + problem.Error())
			return
		}
	}
}

type ethSysAccounts struct {
	sync.RWMutex
	superAdmin    model.EthAccount
	authenticator model.EthAccount
}

var sysAccounts ethSysAccounts

func Healthcheck() error {
	sysAccounts.RLock()
	defer sysAccounts.RUnlock()
	if sysAccounts.authenticator.Prikey == "" {
		logger.Get().Warn().Msg("[Eth logic] Ethereum authenticator key unavailable")
		return model.ErrEthAuthenticatorKeyUnavailable
	}
	return nil
}

func verifyKeyAndAddress(hexkey string, hexaddress string) bool {
	key, err := crypto.HexToECDSA(hexkey)
	if err != nil {
		return false // invalid key
	}
	keyAddr := crypto.PubkeyToAddress(key.PublicKey)
	keyaddr := strings.ToLower(keyAddr.Hex())

	return keyaddr == strings.ToLower(hexaddress)
}

func verifyAddress(hexaddress string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(hexaddress)
}
