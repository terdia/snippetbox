package mock

import (
	"time"

	"github.com/terdia/snippetbox/pkg/models"
	"github.com/terdia/snippetbox/pkg/repository"
)

var mockSnippet = &models.Snippet{
	ID:      1,
	Title:   "An old silent pond",
	Content: "An old silent pond...",
	Created: time.Now(),
	Expires: time.Now(),
}

type snippetRepository struct{}

func NewSnippetRepository() repository.SnippetRepository {
	return &snippetRepository{}
}

func (repo *snippetRepository) Insert(title, content, expires string) (int, error) {
	return 2, nil
}

func (repo *snippetRepository) GetById(id int) (*models.Snippet, error) {
	switch id {
	case 1:

		return mockSnippet, nil
	default:

		return nil, models.ErrNoRecord
	}
}

func (repo *snippetRepository) Latest() ([]*models.Snippet, error) {
	return []*models.Snippet{mockSnippet}, nil
}
