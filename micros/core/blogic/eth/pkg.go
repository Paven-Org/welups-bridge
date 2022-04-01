package ethLogic

import (
	"bridge/micros/core/dao"
	ethdao "bridge/micros/core/dao/eth-account"
	userdao "bridge/micros/core/dao/user"
	"bridge/micros/core/model"
	"bridge/micros/core/service/notifier"
	"bridge/service-managers/logger"
	"context"
	"regexp"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog"
	"go.temporal.io/sdk/client"
)

var (
	ethDAO  ethdao.IEthDAO
	userDAO userdao.IUserDAO
	//mailer  *manager.Mailer
	tempcli client.Client
	log     *zerolog.Logger
)

func Init(d *dao.DAOs, tmpcli client.Client) {
	log = logger.Get()
	ethDAO = d.Eth
	userDAO = d.User
	//	mailer = m
	tempcli = tmpcli
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
