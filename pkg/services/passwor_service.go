package services

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("models: invalid credentials")
)

type PasswordHashService interface {
	Hash(plainText string) (string, error)
	VerifyHash(hash []byte, plainText string) error
}

type bcryptPasswordService struct {
}

func NewPasswordService() PasswordHashService {
	return &bcryptPasswordService{}
}

func (service *bcryptPasswordService) Hash(plainText string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainText), 12)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func (service *bcryptPasswordService) VerifyHash(hash []byte, plainText string) error {
	err := bcrypt.CompareHashAndPassword(hash, []byte(plainText))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidCredentials
		} else {
			return err
		}
	}

	return nil
}
