package welLogic

import (
	"bridge/libs"
	"bridge/micros/core/dao"
	userdao "bridge/micros/core/dao/user"
	weldao "bridge/micros/core/dao/wel-account"
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
	welDAO  weldao.IWelDAO
	userDAO userdao.IUserDAO
	mailer  *manager.Mailer
	tempcli client.Client
	log     *zerolog.Logger
)

func Init(d *dao.DAOs, m *manager.Mailer, tmpcli client.Client) {
	log = logger.Get()
	welDAO = d.Wel
	userDAO = d.User
	mailer = m
	tempcli = tmpcli
}

type welSysAccounts struct {
	sync.RWMutex
	superAdmin    model.WelAccount
	authenticator model.WelAccount
}

var sysAccounts welSysAccounts

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

func sendNotificationToRole(role string, subject string, body string) error {
	users, err := userDAO.GetUsersWithRole(role, 0, 1000)
	if err != nil {
		log.Err(err).Msgf("[Wel logic internal] Unable to fetch users with role %s", role)
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
			log.Err(err).Msgf("[Wel logic internal] unable to send mail to address %s", mail) // best effort lol
		}
	}

	return nil
}
