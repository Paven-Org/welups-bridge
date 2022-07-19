package userLogic

import (
	"bridge/libs"
	"bridge/micros/core/dao"
	userdao "bridge/micros/core/dao/user"
	"bridge/micros/core/model"
	manager "bridge/service-managers"
	"bridge/service-managers/logger"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/rs/zerolog"
)

var (
	userDAO userdao.IUserDAO
	rm      *manager.RedisManager
	ts      libs.ITokenService
	log     *zerolog.Logger
)

func Init(d *dao.DAOs, r *manager.RedisManager, t libs.ITokenService) {
	log = logger.Get()
	userDAO = d.User
	rm = r
	ts = t
}

func ParseToken(token string) (*jwt.Token, error) {
	tk, err := ts.ValidateToken(token)
	if err != nil {
		return tk, err
	}

	if tk == nil {
		err := fmt.Errorf("[user logic] JWT parse result: nil")
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

func generalUpdateUserInfo(username, new_username, email, password, status string) error {
	log.Info().Msgf("[user logic internal] Updating user %s", username)
	user, err := userDAO.GetUserByName(username) // should eventually get by ID instead, but this is more convenient for now
	if err != nil {
		log.Err(err).Msgf("[user logic internal] Failed to retrieve user %s's info", username)
		return err
	}
	// updating fields
	if new_username != "" {
		log.Info().Msgf("[user logic internal] change username from %s to %s", username, new_username)

		user, err := userDAO.GetUserByName(new_username) //
		if err != nil && err != sql.ErrNoRows {
			log.Err(err).Msgf("[user logic internal] Failed to check new_username %s", new_username)
			return fmt.Errorf("Failed to check new_username %s", new_username)
		}
		if err == nil {
			log.Err(err).Msgf("[user logic internal] new_username %s exists", new_username)
			return fmt.Errorf("new_username %s exists", new_username)
		}

		user.Username = new_username
	}

	if email != "" {
		log.Info().Msgf("[user logic internal] change email from %s to %s", user.Email, email)

		user, err := userDAO.GetUserByEmail(email) //
		if err != nil && err != sql.ErrNoRows {
			log.Err(err).Msgf("[user logic internal] Failed to check new email %s", email)
			return fmt.Errorf("Failed to check new email %s", email)
		}
		if err == nil {
			log.Err(err).Msgf("[user logic internal] new email %s exists", email)
			return fmt.Errorf("new email %s exists", email)
		}

		user.Email = email
	}

	if status != "" {
		log.Info().Msgf("[user logic internal] change status from %s to %s", user.Status, status)
		user.Status = status
	}

	if password != "" {
		log.Info().Msgf("[user logic internal] Hashing new password for %s", username)
		user.Password, err = libs.HashPasswd(password)
		if err != nil {
			log.Err(err).Msgf("[user logic internal] Unable to hash password")
			return err
		}
	}

	log.Info().Msgf("[user logic internal] Updating database...")
	if err := userDAO.UpdateUser(&user); err != nil {
		log.Err(err).Msgf("[user logic internal] Failed to update user %s's info in database")
		return err
	}

	return nil
}

func GetUserByName(username string) (*model.User, error) {
	log.Info().Msgf("[user logic internal] Getting user %s", username)
	user, err := userDAO.GetUserByName(username) // should eventually get by ID instead, but this is more convenient for now
	if err != nil {
		log.Err(err).Msgf("[user logic internal] Failed to retrieve user %s's info", username)
		return nil, err
	}

	return &user, nil
}

func GetUserRoles(name string) ([]string, error) {
	log.Info().Msgf("[user logic] Getting user %s's roles...", name)
	roles, err := userDAO.GetUserRoles(name)
	if err != nil {
		log.Err(err).Msgf("[user logic] Failed to retrieve user %s's roles", name)
		return nil, err
	}

	return roles, nil
}
