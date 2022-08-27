package model

import (
	"fmt"
	"time"
)

type User struct {
	Id       uint64 `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"-" db:"password"`
	Email    string `json:"email" db:"email"`
	Status   string `json:"status" db:"status"`

	Created_at time.Time `json:"created_at" db:"created_at"`
	Updated_at time.Time `json:"updated_at,omitempty" db:"updated_at,omitempty"`
}

type Claims struct {
	Exp      time.Time
	Iat      time.Time
	Iss      string
	Uid      uint64
	Username string
	Session  string
}

const (
	UserStatusOK          = "ok"
	UserStatusLocked      = "locked"
	UserStatusBanned      = "banned"
	UserStatusPermabanned = "permabanned"
)

const (
	UserRoleRoot    = "root"
	UserRoleAdmin   = "admin"
	UserRoleService = "service"

	UserRoleDefault = UserRoleAdmin
)

var (
	ErrWrongPasswd             = fmt.Errorf("Wrong password")
	ErrWeakPasswd              = fmt.Errorf("Password must be a mix of uppercase character, lowercase character, symbols and numbers, and at least 8 charaters long")
	ErrUserNotActivated        = fmt.Errorf("User not activated")
	ErrUserBanned              = fmt.Errorf("User not banned")
	ErrUserNotFound            = fmt.Errorf("User not found")
	ErrRoleNotFound            = fmt.Errorf("Role not found")
	ErrUserPermaBanned         = fmt.Errorf("User was wiped out of existence")
	ErrInconsistentCredentials = fmt.Errorf("Cannot reconcile user's token and cookie")
)
