package model

import (
	"fmt"
	"time"
)

type WelAccount struct {
	Address string `json:"address" db:"address"`
	Prikey  string `json:"-" db:"prikey"`
	Status  string `json:"status" db:"status"`

	Created_at time.Time `json:"created_at" db:"created_at"`
	Updated_at time.Time `json:"updated_at,omitempty" db:"updated_at,omitempty"`
}

const (
	WelAccountStatusLocked = "locked"
	WelAccountStatusOK     = "ok"
)

const (
	WelAccountRoleUnauthorized  = "unauthorized"
	WelAccountRoleSuperAdmin    = "super_admin"
	WelAccountRoleManager       = "MANAGER_ROLE"
	WelAccountRoleAuthenticator = "AUTHENTICATOR"
)

var (
	ErrWelAccountLocked               = fmt.Errorf("Welups account locked in the internal system")
	ErrWelNoPrikey                    = fmt.Errorf("Welups account doesn't have private key stored in the system")
	ErrWelInvalidAddress              = fmt.Errorf("Invalid Address")
	ErrWelKeyAndAddressMismatch       = fmt.Errorf("Key and address mismatch")
	ErrWelAccountNotFound             = fmt.Errorf("Account not found in internal system")
	ErrWelRoleNotFound                = fmt.Errorf("Role not found in internal system")
	ErrWelAuthenticatorKeyUnavailable = fmt.Errorf("Welups authenticator key unavailable")
)
