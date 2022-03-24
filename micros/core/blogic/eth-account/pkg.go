package ethLogic

import (
	"bridge/libs"
	"bridge/micros/core/dao"
	ethdao "bridge/micros/core/dao/eth-account"
	userdao "bridge/micros/core/dao/user"
	"bridge/micros/core/model"
	manager "bridge/service-managers"
	"bridge/service-managers/logger"
	"fmt"
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
	mailer  *manager.Mailer
	tempcli client.Client
	log     *zerolog.Logger
)

func Init(d *dao.DAOs, m *manager.Mailer, tmpcli client.Client) {
	log = logger.Get()
	ethDAO = d.Eth
	userDAO = d.User
	mailer = m
	tempcli = tmpcli
}

type ethSysAccounts struct {
	sync.RWMutex
	superAdmin    model.EthAccount
	authenticator model.EthAccount
}

var sysAccounts ethSysAccounts

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

func sendNotificationToRole(role string, subject string, body string) error {
	users, err := userDAO.GetUsersWithRole(role, 0, 1000)
	if err != nil {
		log.Err(err).Msgf("[ethAccount logic internal] Unable to fetch users with role %s", role)
		return err
	}
	fmt.Println("Users: ", users)

	mails := libs.Map(func(u model.User) string { return u.Email },
		libs.Filter(func(u model.User) bool { return u.Status == "ok" }, users))
	fmt.Println("Mails: ", mails)
	for _, mail := range mails {
		mess := mailer.MkPlainMessage(mail, subject, body)
		err := mailer.Send(mess)
		if err != nil {
			log.Err(err).Msgf("[ethAccount logic internal] unable to send mail to address %s", mail) // best effort lol
		}
	}

	return nil
}
