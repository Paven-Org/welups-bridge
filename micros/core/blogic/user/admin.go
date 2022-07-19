package userLogic

import (
	"bridge/libs"
	"bridge/micros/core/model"
	"fmt"
)

func AdminUpdateUserInfo(username, new_username, email, password, status string) error {
	return generalUpdateUserInfo(username, new_username, email, password, status)
}

func AddUser(username, email, password string) error {
	log.Info().Msgf("[user logic internal] Creating user %s...", username)
	log.Info().Msgf("[user logic internal] Hashing new password for %s", username)

	// check existing username and email
	log.Info().Msgf("[user logic internal] Checking username %s", username)

	{
		_, err := userDAO.GetUserByName(username) //
		if err != nil && err != model.ErrUserNotFound {
			log.Err(err).Msgf("[user logic internal] Failed to check username %s", username)
			return fmt.Errorf("Failed to check username %s", username)
		}
		if err == nil || err == model.ErrUserBanned || err == model.ErrUserNotActivated || err == model.ErrUserPermaBanned {
			log.Err(err).Msgf("[user logic internal] username %s exists", username)
			return fmt.Errorf("username %s exists", username)
		}
	}

	log.Info().Msgf("[user logic internal] Checking email %s", email)

	{
		_, err := userDAO.GetUserByEmail(email) //
		if err != nil && err != model.ErrUserNotFound {
			log.Err(err).Msgf("[user logic internal] Failed to check email %s", email)
			return fmt.Errorf("Failed to check email %s", email)
		}
		if err == nil || err == model.ErrUserBanned || err == model.ErrUserNotActivated || err == model.ErrUserPermaBanned {
			log.Err(err).Msgf("[user logic internal] email %s exists", email)
			return fmt.Errorf("email %s exists", email)
		}
	}

	password, err := libs.HashPasswd(password)
	if err != nil {
		log.Err(err).Msgf("[user logic internal] Unable to hash password")
		return err
	}
	_, err = userDAO.AddUser(username, email, password)
	if err != nil {
		log.Err(err).Msgf("[user logic internal] Failed to create user %s", username)
		return err
	}
	return nil
}

func RemoveUser(username string) error {
	log.Info().Msgf("[user logic internal] Start removing user %s...", username)
	user, err := userDAO.GetUserByName(username)
	if err != nil {
		log.Err(err).Msgf("[user logic internal] Failed to retrieve user %s", username)
		return err
	}

	log.Info().Msgf("[user logic internal] Removing user %s...", username)
	if err := userDAO.RemoveUser(user.Username); err != nil {
		log.Err(err).Msgf("[user logic internal] Failed to remove user %s", username)
		return err
	}

	return nil
}

func GrantRole(username, role string) error {
	log.Info().Msgf("[user logic internal] Start granting role %s to user %s...", role, username)
	user, err := userDAO.GetUserByName(username)
	if err != nil {
		log.Err(err).Msgf("[user logic internal] Failed to retrieve user %s", username)
		return err
	}

	log.Info().Msgf("[user logic internal] Granting role %s to user %s...", role, username)
	if err := userDAO.GrantRole(user.Username, role); err != nil {
		log.Err(err).Msgf("[user logic internal] Failed to grant role %s to user %s", role, username)
		return err
	}

	return nil
}

func RevokeRole(username, role string) error {
	log.Info().Msgf("[user logic internal] Start revoking role %s from user %s...", role, username)
	user, err := userDAO.GetUserByName(username)
	if err != nil {
		log.Err(err).Msgf("[user logic internal] Failed to retrieve user %s", username)
		return err
	}

	log.Info().Msgf("[user logic internal] Revoking role %s from user %s...", role, username)
	if err := userDAO.RevokeRole(user.Username, role); err != nil {
		log.Err(err).Msgf("[user logic internal] Failed to revoke role %s from user %s", role, username)
		return err
	}

	return nil
}

func GetUsers(offset uint, size uint) ([]model.User, error) {
	log.Info().Msgf("[user logic internal] Getting users...")
	users, err := userDAO.GetUsers(offset, size)
	if err != nil {
		log.Err(err).Msg("[user logic internal] Failed to retrieve users")
		return nil, err
	}

	return users, nil
}

func GetUsersWithRole(role string, offset uint, size uint) ([]model.User, error) {
	log.Info().Msgf("[user logic internal] Getting users with role %s...", role)
	users, err := userDAO.GetUsersWithRole(role, offset, size)
	if err != nil {
		log.Err(err).Msg("[user logic internal] Failed to retrieve users with role " + role)
		return nil, err
	}

	return users, nil
}

func GetAllRoles() ([]string, error) {
	log.Info().Msgf("[user logic internal] Getting all roles...")
	roles, err := userDAO.GetAllRoles()
	if err != nil {
		log.Err(err).Msg("[user logic internal] Failed to retrieve roles")
		return nil, err
	}

	return roles, nil
}

func TotalUsers() (uint64, error) {
	log.Info().Msgf("[user logic internal] Getting total number of users...")
	total, err := userDAO.TotalUsers()
	if err != nil {
		log.Err(err).Msg("[user logic internal] Failed to retrieve total number of users")
		return 0, err
	}

	return total, nil
}
