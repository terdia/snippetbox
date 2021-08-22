package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/terdia/snippetbox/pkg/models"
)

type UserRepository interface {
	GetById(id int) (*models.User, error)
	Insert(name, email, password string) error
	Unique(field, value string) bool
	FindByEmail(email string) (*models.User, error)
}

type repository struct {
	*sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &repository{db}
}

func (repo *repository) GetById(id int) (*models.User, error) {
	return nil, nil
}

func (repo *repository) Insert(name, email, password string) error {
	stmt := `INSERT INTO users (name, email, hashed_password, created)
			VALUES(?, ?, ?, UTC_TIMESTAMP())`

	_, err := repo.DB.Exec(stmt, name, email, string(password))
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

func (repo *repository) Unique(field, value string) bool {

	stmt := fmt.Sprintf("SELECT %s FROM users WHERE %s = ?", field, field)
	var email string
	row := repo.DB.QueryRow(stmt, value)

	_ = row.Scan(&email)

	return len(email) == 0
}

func (repo *repository) FindByEmail(email string) (*models.User, error) {

	stmt := `SELECT id, hashed_password FROM users WHERE email = ? AND active = TRUE`

	user := models.User{}

	row := repo.DB.QueryRow(stmt, email)

	err := row.Scan(&user.ID, &user.HashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrInvalidCredentials
		} else {
			return nil, err
		}
	}

	return &user, nil
}
