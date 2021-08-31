package mock

import (
	"time"

	"github.com/terdia/snippetbox/pkg/models"
	"github.com/terdia/snippetbox/pkg/repository"
)

type userRepository struct{}

func NewUserRepository() repository.UserRepository {
	return &userRepository{}
}

var mockUser = &models.User{
	ID:             1,
	Name:           "Terry",
	Email:          "terry@yahoo.com",
	Active:         true,
	Created:        time.Now(),
	HashedPassword: []byte("$2a$12$MlU36qi0aP80m4J3rrhuj.AWU7gM96Co9Ch.FH/5cXlOGAn.HNoyu"),
}

func (repo *userRepository) GetById(id int) (*models.User, error) {
	switch id {
	case 1:

		return mockUser, nil
	default:

		return nil, models.ErrNoRecord
	}
}

func (repo *userRepository) Insert(name, email, password string) error {

	switch email {
	case "duplicate@example.com":

		return models.ErrDuplicateEmail
	default:

		return nil
	}
}

func (repo *userRepository) Unique(field, value string) bool {

	if value == "duplicate@example.com" {
		return false
	}

	return true
}

func (repo *userRepository) FindByEmail(email string) (*models.User, error) {
	switch email {
	case "terry@yahoo.com":

		return mockUser, nil
	default:

		return nil, models.ErrInvalidCredentials
	}
}
