package ethLogic

import (
	"bridge/micros/core/dao"
	ethdao "bridge/micros/core/dao/eth-account"
	"bridge/micros/core/model"
	"bridge/service-managers/logger"
	"regexp"
	"sync"

	"github.com/ethereum/go-ethereum/crypto"
)

var (
	ethDAO ethdao.IEthDAO
	log    = logger.Get()
)

func Init(d *dao.DAOs) {
	ethDAO = d.Eth
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

	return keyAddr.Hex() == hexaddress
}

func verifyAddress(hexaddress string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(hexaddress)
}
