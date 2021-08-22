package services

import (
	"errors"

	"github.com/terdia/snippetbox/pkg/forms"
	"github.com/terdia/snippetbox/pkg/models"
	"github.com/terdia/snippetbox/pkg/repository"
)

type AuthResponse struct {
	User  *models.User
	Form  *forms.Form
	Error error
}

type userService struct {
	repository      repository.UserRepository
	passwordService PasswordHashService
}

func NewUserService(repo repository.UserRepository, passwordService PasswordHashService) *userService {
	return &userService{
		repository:      repo,
		passwordService: passwordService,
	}
}

func (service *userService) SignupUser(form *forms.Form) (*forms.Form, error) {

	form.Required("name", "email", "password")
	form.MaxLength("name", 50)
	form.MaxLength("email", 255)
	form.MinLength("password", 6)
	form.MatchesPattern("email", forms.EmailRX)

	if form.Valid() {
		emailExists := service.repository.Unique("email", form.Get("email"))
		if !emailExists {
			form.Errors.Add("email", "Address is already in use")
		}
	}

	if !form.Valid() {
		return form, nil
	}

	hashPassword, err := service.passwordService.Hash(form.Get("password"))
	if err != nil {
		return nil, err
	}

	err = service.repository.Insert(form.Get("name"), form.Get("email"), hashPassword)

	return nil, err
}

func (service *userService) Authenticate(form *forms.Form) AuthResponse {

	form.Required("email", "password")
	if !form.Valid() {
		return AuthResponse{Form: form}
	}

	user, err := service.repository.FindByEmail(form.Get("email"))
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Errors.Add("generic", "Email or Password is incorrect")
			return AuthResponse{Form: form}
		}

		return AuthResponse{Error: err}
	}

	err = service.passwordService.VerifyHash(user.HashedPassword, form.Get("password"))
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			form.Errors.Add("generic", "Email or Password is incorrect")
			return AuthResponse{Form: form}
		}

		return AuthResponse{Error: err}
	}

	return AuthResponse{User: user}

}
