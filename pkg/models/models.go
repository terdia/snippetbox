package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")

type Snippet struct {
	ID      int
	Title   string
	Content string // for nullable field use sql.NullString
	Created time.Time
	Expires time.Time
}
