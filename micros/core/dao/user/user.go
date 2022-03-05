package userDAO

import (
	"bridge/logger"
	"bridge/micros/core/model"
	"time"

	"github.com/jmoiron/sqlx"
)

type IUserDAO interface {
	GetUserById(id uint64) (model.User, error)
	GetUserByName(name string) (model.User, error)
	GetUserByEmail(email string) (model.User, error)
	UpdateUser(u *model.User) error
	RemoveUser(id uint64) error
	GrantRole(id uint64, role string) error
	GetUserRoles(id uint64) ([]string, error)
}

type userDAO struct {
	db *sqlx.DB
}

func MkUserDAO(db *sqlx.DB) IUserDAO {
	return &userDAO{db: db}
}

func (dao *userDAO) AddUser(name string, mail string, password string) (uint64, error) {
	db := dao.db
	log := logger.Get()
	var id uint64

	q := db.Rebind("INSERT INTO users(username, email, password) VALUES (?,?,?) RETURNING id")
	err := db.QueryRow(q, name, mail, password).Scan(&id)

	if err != nil {
		log.Err(err).Msgf("Error while inserting user %s", name)
		return 0, err
	}

	return id, nil
}

func (dao *userDAO) UpdateUser(u *model.User) error {
	db := dao.db
	log := logger.Get()

	q := db.Rebind("UPDATE users SET username = ?, email = ?, password = ?, updated_at = ? WHERE id = ?")
	_, err := db.Exec(q, u.Username, u.Email, u.Password, time.Now(), u.Id)

	if err != nil {
		log.Err(err).Msgf("Error while updating user %v", u)
		return err
	}

	return nil
}

func (dao *userDAO) GrantRole(id uint64, role string) error {
	db := dao.db
	log := logger.Get()

	q := db.Rebind("INSERT INTO user_roles(user_id, role) VALUES (?, ?)")
	_, err := db.Exec(q, id, role)
	if err != nil {
		log.Err(err).Msgf("Error while granting role %s to user %d", role, id)
	}
	return err
}

func (dao *userDAO) RemoveUser(id uint64) error {
	db := dao.db
	log := logger.Get()

	tx, err := db.Beginx() // begin tx
	if err != nil {
		log.Err(err).Msgf("Unable to begin transaction when deleting user %d", id)
		return err
	}

	qDeleteURoles := db.Rebind("DELETE FROM user_roles WHERE user_id = ?")
	_, err = tx.Exec(qDeleteURoles, id)
	if err != nil {
		log.Err(err).Msgf("Error while deleting user %d", id)
		for {
			if err := tx.Rollback(); err != nil {
				log.Err(err).Msg("Error while rolling back tx, retrying...")
			} else {
				break
			}
		}
		return err
	}

	qDeleteUser := db.Rebind("DELETE FROM users WHERE id = ?")
	_, err = tx.Exec(qDeleteUser, id)
	if err != nil {
		log.Err(err).Msgf("Error while deleting user %d", id)
		for {
			if err := tx.Rollback(); err != nil {
				log.Err(err).Msg("Error while rolling back tx, retrying...")
			} else {
				break
			}
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Err(err).Msgf("Error while deleting user %d", id)
		for {
			if err := tx.Rollback(); err != nil {
				log.Err(err).Msg("Error while rolling back tx, retrying...")
			} else {
				break
			}
		}
		return err
	}

	return nil
}

func (dao *userDAO) GetUserById(id uint64) (model.User, error) {
	db := dao.db
	log := logger.Get()

	user := model.User{}

	q := db.Rebind("SELECT * FROM users WHERE id = ?")
	err := db.QueryRowx(q, id).StructScan(&user)
	if err != nil {
		log.Err(err).Msgf("Error while querying for user id: %d", id)
	}
	return user, err
}

func (dao *userDAO) GetUserByName(name string) (model.User, error) {
	db := dao.db
	log := logger.Get()

	user := model.User{}

	q := db.Rebind("SELECT * FROM users WHERE username = ?")
	err := db.QueryRowx(q, name).StructScan(&user)
	if err != nil {
		log.Err(err).Msgf("Error while querying for username: %s", name)
	}
	return user, err
}

func (dao *userDAO) GetUserByEmail(email string) (model.User, error) {
	db := dao.db
	log := logger.Get()

	user := model.User{}

	q := db.Rebind("SELECT * FROM users WHERE email = ?")
	err := db.QueryRowx(q, email).StructScan(&user)
	if err != nil {
		log.Err(err).Msgf("Error while querying for email: %s", email)
	}
	return user, err
}

func (dao *userDAO) GetUserRoles(id uint64) ([]string, error) {
	db := dao.db
	log := logger.Get()

	var roles []string

	q := db.Rebind("SELECT role FROM user_roles WHERE user_roles.user_id = ?")
	err := db.Select(&roles, q, id)

	if err != nil {
		log.Err(err).Msgf("Error while querying for user id %d's roles", id)
		return nil, err
	}

	return roles, nil
}
