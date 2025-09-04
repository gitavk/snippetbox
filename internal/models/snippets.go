package models

import (
	"context"
	"database/sql"
	"time"
)

// Define a Snippet type to hold the data for an individual snippet. Notice how
// the fields of the struct correspond to the fields in our MySQL snippets
// table?
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// Define a SnippetModel type which wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *sql.DB
}

// This will insert a new snippet into the database.
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	var id int
	stmt := `INSERT INTO snippets (title, content, created, expires)
          VALUES($1, $2, NOW(), NOW() + INTERVAL '1 day' * $3) RETURNING id`
	err := m.DB.QueryRowContext(context.Background(), stmt, title, content, expires).Scan(&id)
	if err != nil {
		return 0, err
	}
	// The ID returned has the type int64, so we convert it to an int type
	// before returning.
	return id, nil
}

// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (Snippet, error) {
	return Snippet{}, nil
}

// This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]Snippet, error) {
	return nil, nil
}
