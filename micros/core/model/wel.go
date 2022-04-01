package model

import (
	"fmt"
	"time"
)

type WelAccount struct {
	Address string
	Prikey  string
	Status  string

	Created_at time.Time
	Updated_at time.Time
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
