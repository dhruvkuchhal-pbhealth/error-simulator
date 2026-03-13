package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
)

// UserRepository holds dependencies for user-related DB operations.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository constructs a UserRepository. It accepts a *sql.DB which may be nil,
// but callers should provide a valid DB. Keeping this simple to avoid changing call sites.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetUserByID fetches a user by ID from the database.
func (r *UserRepository) GetUserByID(id int) (string, error) {
	if r == nil {
		return "", errors.New("user repository is nil")
	}
	if r.db == nil {
		return "", errors.New("database is not initialized")
	}

	var name string
	row := r.db.QueryRow("SELECT name FROM users WHERE id = ?", id)
	if err := row.Scan(&name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("user with id %d not found: %w", id, err)
		}
		return "", fmt.Errorf("query scan error: %w", err)
	}
	return name, nil
}

// DBError is an HTTP handler that demonstrates using the repository.
func DBError(repo *UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// For example purposes, use ID 1
		name, err := repo.GetUserByID(1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "User: %s", name)
	}
}
