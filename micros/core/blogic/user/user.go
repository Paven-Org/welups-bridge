package userLogic

import (
	"bridge/libs"
	"bridge/micros/core/dao"
	userdao "bridge/micros/core/dao/user"
	"bridge/micros/core/model"
	manager "bridge/service-managers"
	"bridge/service-managers/logger"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	userDAO userdao.IUserDAO
	rm      *manager.RedisManager
	ts      libs.ITokenService
	log     = logger.Get()
)

func Init(d *dao.DAOs, r *manager.RedisManager, t libs.ITokenService) {
	userDAO = d.User
	rm = r
	ts = t
}

func Login(username string, password string) (string, string, string, error) {
	log.Info().Msgf("Preparing to login user %s", username)
	redis, err := rm.GetRedisClient(manager.StdAuthDBName)
	if err != nil {
		log.Err(err).Msgf("Failed to get redis connection %s's info", username)
		return "", "", "", err
	}
	ctx := context.Background()

	tk, sessionSecret, err := login(username, password)
	if err != nil {
		log.Err(err).Msgf("User %s's login failed", username)
		return "", "", "", err
	}

	signedTk, err := ts.SignToken(tk)
	if err != nil {
		log.Err(err).Msgf("Error while creating user %s's credential", username)
		return "", "", "", err
	}

	mClaims, _ := tk.Claims.(jwt.MapClaims)
	sessionID := mClaims["session"].(string)
	_exp, _ := strconv.ParseInt(mClaims["exp"].(string), 10, 64)
	exp := time.Unix(_exp, 0)
	expDur := exp.Sub(time.Now())

	logger.Get().Debug().Msgf("sessionID: %s", sessionID)
	logger.Get().Debug().Msgf("exp: %s", exp.String())
	logger.Get().Debug().Msgf("expDur: %s", expDur.String())

	if err := redis.SetNX(ctx, fmt.Sprintf("session:user_%s:%s", username, sessionID), sessionSecret, exp.Sub(time.Now())).Err(); err != nil {
		log.Err(err).Msgf("Error while saving session for user %s", username)
		return "", "", "", err

	}

	return signedTk, sessionID, sessionSecret, nil
}

func login(username string, password string) (*jwt.Token, string, error) {
	log.Info().Msgf("Logging user %s in...", username)
	user, err := userDAO.GetUserByName(username)
	if err != nil {
		log.Err(err).Msgf("Failed to retrieve user %s's info", username)
		return nil, "", err
	}

	log.Info().Msgf("Checking user %s's password...", username)
	if !libs.ValidatePasswd(user.Password, password) {
		err := model.ErrWrongPasswd
		log.Err(err).Msgf("Wrong password")
		return nil, "", err
	}

	if user.Status != model.UserStatusOK {
		var err error
		switch user.Status {
		case model.UserStatusLocked:
			err = model.ErrUserNotActivated
		case model.UserStatusBanned:
			err = model.ErrUserBanned
		case model.UserStatusPermabanned:
			err = model.ErrUserPermaBanned
		}
		log.Err(err).Msgf("User %s is not available", username)
		return nil, "", err
	}

	log.Info().Msgf("Creating user %s's credential...", username)
	tk := ts.MkToken(user.Id, user.Username, time.Hour*24*30)

	// used as a httponly secure cookie to guard against XXS
	sessionSecret := libs.Uniq()

	return tk, sessionSecret, nil
}

func Logout(token string, cookie string) error {
	tk, err := ts.ValidateToken(token)
	if err != nil {
		log.Err(err).Msgf("Failed to parse token %s", token)
		return err
	}
	if tk == nil {
		err := fmt.Errorf("JWT parse result: nil")
		log.Err(err).Msgf("Failed to parse token %s", token)
		return err
	}

	mClaims, _ := tk.Claims.(jwt.MapClaims)
	sessionID := mClaims["session"].(string)
	username := mClaims["username"].(string)
	logger.Get().Debug().Msgf("sessionID: %s", sessionID)
	logger.Get().Debug().Msgf("username: %s", username)

	log.Info().Msgf("Preparing to logout user %s", username)
	redis, err := rm.GetRedisClient(manager.StdAuthDBName)
	if err != nil {
		log.Err(err).Msgf("Failed to get redis connection")
		return err
	}
	ctx := context.Background()

	sessionSecret, err := redis.
		Get(ctx,
			fmt.Sprintf("session:user_%s:%s", username, sessionID)).
		Result()

	if err != nil {
		log.Err(err).Msgf("Error while logging out user %s", username)
		return err
	}

	if cookie != sessionSecret {
		err := model.ErrInconsistentCredentials
		log.Err(err).Msgf("Error while logging out user %s", username)
		return err
	}

	if err := redis.Del(ctx, fmt.Sprintf("session:user_%s:%s", username, sessionID)).Err(); err != nil {
		log.Err(err).Msgf("Error while removing session for user %s", username)
		return err
	}

	return nil
}

func GetUserRoles(id uint64) ([]string, error) {
	log.Info().Msgf("Getting user id %d's roles...", id)
	roles, err := userDAO.GetUserRoles(id)
	if err != nil {
		log.Err(err).Msgf("Failed to retrieve user id %d's roles", id)
		return nil, err
	}

	return roles, nil
}

func ParseToken(token string) (*jwt.Token, error) {
	tk, err := ts.ValidateToken(token)
	if err != nil {
		return tk, err
	}

	if tk == nil {
		err := fmt.Errorf("JWT parse result: nil")
		return tk, err
	}

	return tk, err
}

func ParseTokenToClaims(token string) (*model.Claims, error) {
	tk, err := ParseToken(token)
	if err != nil {
		log.Err(err).Msgf("Failed to parse token %s", token)
		return nil, err
	}

	claims, _ := tk.Claims.(jwt.MapClaims)

	//"session"
	sessionID := claims["session"].(string)
	//"exp":      exp,
	_exp, _ := strconv.ParseInt(claims["exp"].(string), 10, 64)
	exp := time.Unix(_exp, 0)
	//"iat":      created,
	_iat, _ := strconv.ParseInt(claims["iat"].(string), 10, 64)
	iat := time.Unix(_iat, 0)
	//"iss":      "welbridge",
	iss := claims["iss"].(string)
	//"uid":      fmt.Sprintf("%d", uid),
	uid, _ := strconv.ParseUint(claims["uid"].(string), 10, 64)
	//"username": username,
	username := claims["username"].(string)

	return &model.Claims{
			Exp:      exp,
			Iat:      iat,
			Iss:      iss,
			Uid:      uid,
			Username: username,
			Session:  sessionID,
		},
		nil
}
