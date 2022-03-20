package model

import (
	"fmt"
	"time"
)

type EthAccount struct {
	Address string
	Prikey  string
	Status  string

	Created_at time.Time
	Updated_at time.Time
}

const (
	EthAccountStatusLocked = "locked"
	EthAccountStatusOK     = "ok"
)

const (
	EthAccountRoleUnauthorized  = "unauthorized"
	EthAccountRoleSuperAdmin    = "super_admin"
	EthAccountRoleTreasury      = "treasury"
	EthAccountRoleVault         = "vault"
	EthAccountRoleAuthenticator = "authenticator"
)

var (
	ErrEthAccountLocked = fmt.Errorf("Ethereum account locked in the internal system")
	ErrEthNoPrikey      = fmt.Errorf("Ethereum account doesn't have private key stored in the system")
)
