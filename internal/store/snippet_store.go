package store

import (
	"database/sql"
	"time"
)

type Snippet struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Code        string    `json:"code"`
	Language    string    `json:"language"`
	UserID      int       `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PostgresSnippetStore struct {
	db *sql.DB
}

func NewPostgresSnippetStore(db *sql.DB) *PostgresSnippetStore {
	return &PostgresSnippetStore{
		db: db,
	}
}

type SnippetStore interface {
	CreateSnippet(*Snippet) (*Snippet, error)
	GetSnippetsByUserID(userID int) ([]*Snippet, error)
	DeleteSnippet(snippetID int) error
	UpdateSnippet(snippet *Snippet) (*Snippet, error)
}

func (pg *PostgresSnippetStore) CreateSnippet(snippet *Snippet) (*Snippet, error) {

	query := `
		INSERT INTO snippets (title, description, code, language, user_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	err := pg.db.QueryRow(query, snippet.Title, snippet.Description, snippet.Code, snippet.Language, snippet.UserID).Scan(&snippet.ID, &snippet.CreatedAt, &snippet.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return snippet, nil
}

func (pg *PostgresSnippetStore) GetSnippetsByUserID(userID int) ([]*Snippet, error) {
	query := `
		SELECT id, title, description, code, language, user_id, created_at, updated_at
		FROM snippets
		WHERE user_id = $1
	`

	rows, err := pg.db.Query(query, userID)
	if err != nil {
		return nil, err
	}

	snippets := []*Snippet{}
	for rows.Next() {
		snippet := &Snippet{}
		err := rows.Scan(&snippet.ID, &snippet.Title, &snippet.Description, &snippet.Code, &snippet.Language, &snippet.UserID, &snippet.CreatedAt, &snippet.UpdatedAt)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, snippet)
	}

	return snippets, nil
}

func (pg *PostgresSnippetStore) DeleteSnippet(snippetID int) error {
	query := `DELETE FROM snippets WHERE id = $1`
	_, err := pg.db.Exec(query, snippetID)
	return err
}

func (pg *PostgresSnippetStore) UpdateSnippet(snippet *Snippet) (*Snippet, error) {
	query := `
		UPDATE snippets
		SET title = $1, description = $2, code = $3, language = $4, updated_at = NOW()
		WHERE id = $5
  	RETURNING id, title, description, code, language, user_id, created_at, updated_at
	`
	err := pg.db.QueryRow(query, snippet.Title, snippet.Description, snippet.Code, snippet.Language, snippet.ID).Scan(&snippet.ID, &snippet.Title, &snippet.Description, &snippet.Code, &snippet.Language, &snippet.UserID, &snippet.CreatedAt, &snippet.UpdatedAt)
	return snippet, err
}
