package model

import (
	"fmt"
	"time"
)

type User struct {
	Id       uint64
	Username string
	Password string
	Email    string
	Status   string

	Created_at time.Time
	Updated_at time.Time
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

var (
	ErrWrongPasswd             = fmt.Errorf("Wrong password")
	ErrUserNotActivated        = fmt.Errorf("User not activated")
	ErrUserBanned              = fmt.Errorf("User not banned")
	ErrUserPermaBanned         = fmt.Errorf("User was wiped out of existence")
	ErrInconsistentCredentials = fmt.Errorf("Cannot reconcile user's token and cookie")
)
