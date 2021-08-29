package services

import (
	"github.com/terdia/snippetbox/pkg/forms"
	"github.com/terdia/snippetbox/pkg/models"
)

type UserServiceInterface interface {
	GetById(id int) (*models.User, error)
	SignupUser(form *forms.Form) (*forms.Form, error)
	Authenticate(form *forms.Form) AuthResponse
}
