package services

import (
	"github.com/terdia/snippetbox/pkg/forms"
	"github.com/terdia/snippetbox/pkg/models"
	"github.com/terdia/snippetbox/pkg/repository"
)

type snippetService struct {
	repository repository.SnippetRepository
}

func NewSnippetService(repo repository.SnippetRepository) *snippetService {
	return &snippetService{repository: repo}
}

func (service *snippetService) CreateSnippet(form *forms.Form) (int, *forms.Form, error) {
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	if !form.Valid() {
		return 0, form, nil
	}

	id, err := service.repository.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		return 0, nil, err
	}

	return id, nil, nil
}

func (service *snippetService) GetById(id int) (*models.Snippet, error) {
	return service.repository.GetById(id)
}

func (service *snippetService) GetLatest() ([]*models.Snippet, error) {
	return service.repository.Latest()
}
