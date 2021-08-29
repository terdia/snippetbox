package services

import (
	"github.com/terdia/snippetbox/pkg/forms"
	"github.com/terdia/snippetbox/pkg/models"
)

type SnippetServiceInterface interface {
	CreateSnippet(form *forms.Form) (int, *forms.Form, error)
	GetById(id int) (*models.Snippet, error)
	GetLatest() ([]*models.Snippet, error)
}
