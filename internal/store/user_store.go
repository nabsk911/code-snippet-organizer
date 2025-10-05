package store

import (
	"database/sql"
)

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{
		db: db,
	}
}

type UserStore interface {
	CreateUser(*User) error
	GetUserByUsername(username string) (*User, error)
}

func (pg *PostgresUserStore) CreateUser(user *User) error {
	query := `INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id`
	err := pg.db.QueryRow(query, user.Username, user.PasswordHash).Scan(&user.ID)
	return err
}

func (pg *PostgresUserStore) GetUserByUsername(username string) (*User, error) {
	query := `SELECT id, username, password_hash FROM users WHERE username = $1`
	user := &User{}
	err := pg.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.PasswordHash)
	return user, err
}
