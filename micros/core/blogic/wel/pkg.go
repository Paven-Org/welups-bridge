package welLogic

import (
	"bridge/libs"
	welABI "bridge/micros/core/abi/wel"
	"bridge/micros/core/dao"
	userdao "bridge/micros/core/dao/user"
	weldao "bridge/micros/core/dao/wel-account"
	"bridge/micros/core/model"
	"bridge/micros/core/service/notifier"
	"bridge/service-managers/logger"
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog"
	"go.temporal.io/sdk/client"
)

var (
	welDAO  weldao.IWelDAO
	userDAO userdao.IUserDAO
	//mailer  *manager.Mailer
	tempcli client.Client
	welInq  *welABI.WelInquirer
	log     *zerolog.Logger
)

func Init(d *dao.DAOs, tmpcli client.Client, inq *welABI.WelInquirer) {
	log = logger.Get()
	welDAO = d.Wel
	userDAO = d.User
	//mailer = m
	tempcli = tmpcli
	welInq = inq
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

type welSysAccounts struct {
	sync.RWMutex
	superAdmin    model.WelAccount
	authenticator model.WelAccount
}

var sysAccounts welSysAccounts

func Healthcheck() error {
	sysAccounts.RLock()
	defer sysAccounts.RUnlock()
	if sysAccounts.authenticator.Prikey == "" {
		logger.Get().Warn().Msg("[Wel logic] Welups authenticator key unavailable")
		return model.ErrWelAuthenticatorKeyUnavailable
	}
	return nil
}

func verifyKeyAndAddress(hexkey string, address string) bool {
	key, err := crypto.HexToECDSA(hexkey)
	if err != nil {
		return false // invalid key
	}
	keyAddr := crypto.PubkeyToAddress(key.PublicKey)

	var hexaddress string
	if !strings.HasPrefix(address, "0x") {
		hexaddress, err = libs.B58toHex(address)
		if err != nil {
			return false
		}
	} else {
		hexaddress = address
	}
	_, hexaddress, _ = strings.Cut(hexaddress, "0x41")
	hexaddress = "0x" + hexaddress
	hexaddress = strings.ToLower(hexaddress)

	keyaddr := strings.ToLower(keyAddr.Hex())

	fmt.Println("keyaddr and hexaddr: ", keyaddr, hexaddress)

	return hexaddress == keyaddr
}

func verifyAddress(address string) bool { // change! Base58 -> hex -> address
	var hexaddress string
	var err error
	if !strings.HasPrefix(address, "0x") {
		hexaddress, err = libs.B58toHex(address)
		if err != nil {
			return false
		}
	} else {
		hexaddress = address
	}

	re := regexp.MustCompile("^0x41[0-9a-fA-F]{40}$")
	return re.MatchString(hexaddress)
}
