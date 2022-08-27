package userLogic

import (
	"bridge/libs"
	"bridge/micros/core/model"
	manager "bridge/service-managers"
	"bridge/service-managers/logger"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func Login(username string, password string) (string, string, string, time.Duration, error) {
	log.Info().Msgf("[user logic] Preparing to login user %s", username)
	redis, err := rm.GetRedisClient(manager.StdAuthDBName)
	if err != nil {
		log.Err(err).Msgf("[user logic] Failed to get redis connection %s's info", username)
		return "", "", "", -1, err
	}
	log.Info().Msg("[user logic] connected to redis server")
	ctx := context.Background()

	tk, sessionSecret, err := login(username, password)
	if err != nil {
		log.Err(err).Msgf("[user logic] User %s's login failed", username)
		return "", "", "", -1, err
	}

	signedTk, err := ts.SignToken(tk)
	if err != nil {
		log.Err(err).Msgf("[user logic] Error while creating user %s's credential", username)
		return "", "", "", -1, err
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
		log.Err(err).Msgf("[user logic] Error while saving session for user %s", username)
		return "", "", "", -1, err

	}

	return signedTk, sessionID, sessionSecret, expDur, nil
}

func login(username string, password string) (*jwt.Token, string, error) {
	log.Info().Msgf("[user logic internal] Logging user %s in...", username)
	user, err := userDAO.GetUserByName(username)
	if err != nil {
		log.Err(err).Msgf("[user logic internal] Failed to retrieve user %s's info", username)
		return nil, "", err
	}

	log.Info().Msgf("[user logic internal] Checking user %s's password...", username)
	if !libs.ValidatePasswd(user.Password, password) {
		err := model.ErrWrongPasswd
		log.Err(err).Msgf("[user logic internal] Wrong password")
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
		log.Err(err).Msgf("[user logic internal] User %s is not available", username)
		return nil, "", err
	}

	log.Info().Msgf("[user logic internal] Creating user %s's credential...", username)
	tk := ts.MkToken(user.Id, user.Username, time.Hour*24*30)

	// used as a httponly secure cookie to guard against XXS
	sessionSecret := libs.Uniq()

	return tk, sessionSecret, nil
}

func Logout(token string, cookie string) error {
	tk, err := ts.ValidateToken(token)
	if err != nil {
		log.Err(err).Msgf("[user logic] Failed to parse token %s", token)
		return err
	}
	if tk == nil {
		err := fmt.Errorf("[user logic] JWT parse result: nil")
		log.Err(err).Msgf("[user logic] Failed to parse token %s", token)
		return err
	}

	mClaims, _ := tk.Claims.(jwt.MapClaims)
	sessionID := mClaims["session"].(string)
	username := mClaims["username"].(string)
	logger.Get().Debug().Msgf("sessionID: %s", sessionID)
	logger.Get().Debug().Msgf("username: %s", username)

	log.Info().Msgf("[user logic] Preparing to logout user %s", username)
	redis, err := rm.GetRedisClient(manager.StdAuthDBName)
	if err != nil {
		log.Err(err).Msgf("[user logic] Failed to get redis connection")
		return err
	}
	ctx := context.Background()

	sessionSecret, err := redis.
		Get(ctx,
			fmt.Sprintf("session:user_%s:%s", username, sessionID)).
		Result()

	if err != nil {
		log.Err(err).Msgf("[user logic] Error while logging out user %s", username)
		return err
	}

	if cookie != sessionSecret {
		err := model.ErrInconsistentCredentials
		log.Err(err).Msgf("[user logic] Error while logging out user %s", username)
		return err
	}

	if err := redis.Del(ctx, fmt.Sprintf("session:user_%s:%s", username, sessionID)).Err(); err != nil {
		log.Err(err).Msgf("[user logic] Error while removing session for user %s", username)
		return err
	}

	return nil
}

// change password
func Passwd(username string, oldpasswd string, newpasswd string) error {
	log.Info().Msgf("[user logic internal] Changing password for user %s", username)
	user, err := userDAO.GetUserByName(username)
	if err != nil {
		log.Err(err).Msgf("[user logic internal] Failed to retrieve user %s's info", username)
		return err
	}

	log.Info().Msgf("[user logic internal] Checking user %s's old password...", username)
	if !libs.ValidatePasswd(user.Password, oldpasswd) {
		err := model.ErrWrongPasswd
		log.Err(err).Msgf("[user logic internal] old password doesn't match")
		return err
	}

	log.Info().Msgf("[user logic internal] Checking user %s's new password's strength...", username)
	if !libs.StrongPasswd(user.Password) {
		err := model.ErrWeakPasswd
		log.Err(err).Msgf("[user logic internal] Weak password")
		return err
	}

	log.Info().Msgf("[user logic internal] Hashing new password for %s", username)
	user.Password, err = libs.HashPasswd(newpasswd)
	if err != nil {
		log.Err(err).Msgf("[user logic internal] Unable to hash password")
		return err
	}

	log.Info().Msgf("[user logic internal] Updating database...")
	if err := userDAO.UpdateUser(&user); err != nil {
		log.Err(err).Msgf("[user logic internal] Failed to update user %s's info in database")
		return err
	}

	return nil
}

// Update user's info
func UpdateUserInfo(username, email string) error {
	return generalUpdateUserInfo(username, "", email, "", "")
}
