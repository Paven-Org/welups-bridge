package libs

import (
	"bridge/service-managers/logger"
	"unicode"

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

// VerifyPasswdStrength...
// 8 characters minimum, mix of uppercase, lowercase, symbols and number
func StrongPasswd(plainPasswd string) bool {
	conditions := map[string]bool{}
	if len(plainPasswd) >= 8 {
		conditions["length"] = true
	}

	for _, c := range plainPasswd {
		switch {
		// case unicode.IsSpace(c):
		// 	return false
		case unicode.IsNumber(c):
			conditions["number"] = true
		case unicode.IsLower(c):
			conditions["lower"] = true
		case unicode.IsUpper(c):
			conditions["upper"] = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			conditions["special"] = true
		}
	}
	return conditions["length"] && 
				 conditions["number"] && 
				 conditions["lower"] && 
				 conditions["upper"] && 
				 conditions["special"]
}
