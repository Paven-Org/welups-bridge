package userDAO

import (
	"bridge/micros/core/model"
	"bridge/service-managers/logger"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type IUserDAO interface {
	GetUserById(id uint64) (model.User, error)
	GetUserByName(name string) (model.User, error)
	GetUserByEmail(email string) (model.User, error)
	GetUsers(offset uint, size uint) ([]model.User, error)
	GetUsersWithRole(role string, offset uint, size uint) ([]model.User, error)
	AddUser(name string, mail string, password string) (uint64, error)
	UpdateUser(u *model.User) error
	RemoveUser(name string) error
	GrantRole(name string, role string) error
	RevokeRole(name string, role string) error
	GetUserRoles(name string) ([]string, error)
	GetAllRoles() ([]string, error)

	TotalUsers() (uint64, error)
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

	q := db.Rebind("UPDATE users SET username = ?, email = ?, password = ?, status = ?, updated_at = ? WHERE id = ?")
	_, err := db.Exec(q, u.Username, u.Email, u.Password, u.Status, time.Now(), u.Id)

	if err != nil {
		log.Err(err).Msgf("Error while updating user %v", u)
		return err
	}

	return nil
}

func (dao *userDAO) GrantRole(name string, role string) error {
	db := dao.db
	log := logger.Get()

	q := db.Rebind(`INSERT INTO user_roles(user_id, role) 
									SELECT users.id, ? from users
									WHERE users.username = ?`)
	_, err := db.Exec(q, role, name)
	if err != nil {
		log.Err(err).Msgf("Error while granting role %s to user %s", role, name)
	}
	return err
}

func (dao *userDAO) RevokeRole(name string, role string) error {
	db := dao.db
	log := logger.Get()

	q := db.Rebind(`DELETE FROM user_roles 
									WHERE user_id = (SELECT id from users where username = ?) 
									AND role = ?`)
	_, err := db.Exec(q, name, role)
	if err != nil {
		log.Err(err).Msgf("Error while revoking role %s from user %s", role, name)
	}
	return err
}

func (dao *userDAO) RemoveUser(name string) error {
	db := dao.db
	log := logger.Get()

	tx, err := db.Beginx() // begin tx
	if err != nil {
		log.Err(err).Msgf("Unable to begin transaction when deleting user %s", name)
		return err
	}

	qDeleteURoles := db.Rebind(`DELETE FROM user_roles 
									WHERE user_id = (SELECT id from users where username = ?)`)
	_, err = tx.Exec(qDeleteURoles, name)
	if err != nil {
		log.Err(err).Msgf("Error while deleting user %s", name)
		for {
			if err := tx.Rollback(); err != nil && err != sql.ErrTxDone && err != sql.ErrConnDone {
				log.Err(err).Msg("Error while rolling back tx, retrying...")
			} else {
				break
			}
		}
		return err
	}

	qDeleteUser := db.Rebind("DELETE FROM users WHERE username = ?")
	_, err = tx.Exec(qDeleteUser, name)
	if err != nil {
		log.Err(err).Msgf("Error while deleting user %s", name)
		for {
			if err := tx.Rollback(); err != nil && err != sql.ErrTxDone && err != sql.ErrConnDone {
				log.Err(err).Msg("Error while rolling back tx, retrying...")
			} else {
				break
			}
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Err(err).Msgf("Error while deleting user %s", name)
		for {
			if err := tx.Rollback(); err != nil && err != sql.ErrTxDone && err != sql.ErrConnDone {
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
		if err != sql.ErrNoRows {
			log.Err(err).Msgf("Error while querying for user id: %d", id)
			return user, err
		}
		return user, model.ErrUserNotFound
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
		if err != sql.ErrNoRows {
			log.Err(err).Msgf("Error while querying for username: %s", name)
			return user, err
		}
		return user, model.ErrUserNotFound
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
		if err != sql.ErrNoRows {
			log.Err(err).Msgf("Error while querying for email: %s", email)
			return user, err
		}
		return user, model.ErrUserNotFound
	}
	return user, err
}

func (dao *userDAO) GetUserRoles(name string) ([]string, error) {
	db := dao.db
	log := logger.Get()

	var roles []string

	q := db.Rebind(`SELECT role FROM user_roles JOIN users 
									ON user_roles.user_id = users.id
									WHERE users.username = ?`)
	err := db.Select(&roles, q, name)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Err(err).Msgf("Error while querying for user %s's roles", name)
			return nil, err
		}
		return nil, model.ErrRoleNotFound
	}

	return roles, nil
}

// use the good old offset-limit pagination technique, because there won't be that many
// users in the internal system
func (dao *userDAO) GetUsers(offset uint, size uint) ([]model.User, error) {
	var users []model.User
	db := dao.db
	log := logger.Get()

	q := db.Rebind(`SELECT * FROM users 
									ORDER BY users.id 
									OFFSET ? LIMIT ?`)
	err := db.Select(&users, q, offset, size)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Err(err).Msg("Error while querying for users")
			return nil, err
		}
		return nil, model.ErrUserNotFound
	}

	return users, nil
}

func (dao *userDAO) GetUsersWithRole(role string, offset uint, size uint) ([]model.User, error) {
	var users []model.User
	db := dao.db
	log := logger.Get()

	q := db.Rebind(`SELECT users.* FROM users INNER JOIN user_roles 
									ON user_roles.user_id = users.id 
									WHERE user_roles.role  = ? 
									ORDER BY users.id 
									OFFSET ? LIMIT ?`)
	err := db.Select(&users, q, role, offset, size)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Err(err).Msgf("Error while querying for users with role %s", role)
			return nil, err
		}
		return nil, model.ErrUserNotFound
	}
	return users, nil
}

// again, there won't be that many roles, so pagination isn't even required
func (dao *userDAO) GetAllRoles() ([]string, error) {
	var roles []string
	db := dao.db
	log := logger.Get()

	q := db.Rebind("SELECT role FROM roles")
	err := db.Select(&roles, q)

	if err != nil {
		log.Err(err).Msg("Error while querying for roles")
		return nil, err
	}

	return roles, nil
}

func (dao *userDAO) TotalUsers() (uint64, error) {
	db := dao.db
	log := logger.Get()
	var total uint64

	q := db.Rebind("SELECT users FROM total_rows_of")
	err := db.QueryRow(q).Scan(&total)

	if err != nil {
		log.Err(err).Msgf("Error while getting total rows of users table")
		return 0, err
	}

	return total, nil
}
