package repository

import (
	"database/sql"
	"errors"

	"github.com/terdia/snippetbox/pkg/models"
)

type SnippetRepository interface {
	Insert(title, content, expires string) (int, error)
	GetById(id int) (*models.Snippet, error)
	Latest() ([]*models.Snippet, error)
}

type snippetRepository struct {
	*sql.DB
}

func NewSnippetRepository(db *sql.DB) SnippetRepository {
	return &snippetRepository{db}
}

func (repo *snippetRepository) Insert(title, content, expires string) (int, error) {

	stmt := `INSERT INTO snippets (title, content, created, expires)
			 VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := repo.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (repo *snippetRepository) GetById(id int) (*models.Snippet, error) {

	stmt := `SELECT id, title, content, created, expires FROM snippets
			 WHERE expires > UTC_TIMESTAMP() AND id = ?`

	snippet := &models.Snippet{}

	err := repo.DB.QueryRow(stmt, id).Scan(
		&snippet.ID,
		&snippet.Title,
		&snippet.Content,
		&snippet.Created,
		&snippet.Expires,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return snippet, nil
}

func (repo *snippetRepository) Latest() ([]*models.Snippet, error) {

	stmt := `SELECT id, title, content, created, expires FROM snippets
			 WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`

	rows, err := repo.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	snippets := []*models.Snippet{}

	for rows.Next() {
		s := &models.Snippet{}

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
