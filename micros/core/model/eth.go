package model

import (
	"fmt"
	"time"
)

type EthAccount struct {
	Address string `json:"address" db:"address"`
	Prikey  string `json:"-" db:"prikey"`
	Status  string `json:"status" db:"status"`

	Created_at time.Time `json:"created_at" db:"created_at"`
	Updated_at time.Time `json:"updated_at,omitempty" db:"updated_at,omitempty"`
}

const (
	EthAccountStatusLocked = "locked"
	EthAccountStatusOK     = "ok"
)

const (
	EthAccountRoleUnauthorized  = "unauthorized"
	EthAccountRoleSuperAdmin    = "super_admin"
	EthAccountRoleManager       = "MANAGER_ROLE"
	EthAccountRoleAuthenticator = "AUTHENTICATOR"
)

var (
	ErrEthAccountLocked               = fmt.Errorf("Ethereum account locked in the internal system")
	ErrEthNoPrikey                    = fmt.Errorf("Ethereum account doesn't have private key stored in the system")
	ErrEthInvalidAddress              = fmt.Errorf("Invalid Address")
	ErrEthKeyAndAddressMismatch       = fmt.Errorf("Key and address mismatch")
	ErrEthAccountNotFound             = fmt.Errorf("Account not found in internal system")
	ErrEthRoleNotFound                = fmt.Errorf("Role not found in internal system")
	ErrEthAuthenticatorKeyUnavailable = fmt.Errorf("Ethereum authenticator key unavailable")
)
