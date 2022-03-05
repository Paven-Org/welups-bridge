package libs

import (
	"bridge/logger"
	"crypto/md5"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type ITokenService interface {
	MkToken(uid uint64, password string, expiration time.Duration) *jwt.Token
	SignToken(token *jwt.Token) (string, error)
	ValidateToken(string) (*jwt.Token, error)
}

type tokenServ struct {
	jwtSecret string
}

var _ ITokenService = &tokenServ{}

func MkTokenServ(jwtSecret string) ITokenService {
	return &tokenServ{
		jwtSecret: jwtSecret,
	}
}

func (t *tokenServ) MkToken(uid uint64, username string, expiration time.Duration) *jwt.Token {
	created := time.Now().Unix()
	exp := time.Now().Add(expiration).Unix()
	claims := jwt.MapClaims{
		"exp":      fmt.Sprintf("%d", exp),
		"iat":      fmt.Sprintf("%d", created),
		"iss":      "welbridge",
		"uid":      fmt.Sprintf("%d", uid),
		"username": username,
		"session":  fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%d%s%d", uid, username, created)))),
	}

	logger.Get().Debug().Msgf("sessionID: %s", claims["session"])
	logger.Get().Debug().Msgf("exp: %s", time.Unix(exp, 0).String())

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token

}

func (t *tokenServ) SignToken(token *jwt.Token) (string, error) {

	signedToken, err := token.SignedString([]byte(t.jwtSecret))
	if err != nil {
		logger.Get().Err(err).Msg("[GenerateToken] Failed to generate JWT token")
	}

	return signedToken, err
}

func (t *tokenServ) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected JWT signing method: %v", token.Header["alg"])
		}
		return []byte(t.jwtSecret), nil
	})
}
