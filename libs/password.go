package libs

import (
	"bridge/logger"

	"golang.org/x/crypto/bcrypt"
)

func HashPasswd(passwd string) (string, error) {
	bytePasswd := []byte(passwd)
	hPass, err := bcrypt.GenerateFromPassword(bytePasswd, bcrypt.DefaultCost)
	if err != nil {
		logger.Get().Err(err).Msg("[HashPasswd] Hash password failed")
	}
	return string(hPass), err
}

func ValidatePasswd(hashPasswd, passwd string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashPasswd), []byte(passwd)); err != nil {
		logger.Get().Err(err).Msgf("[ValidatePasswd] Password failed to validate")
		return false
	}
	return true
}
