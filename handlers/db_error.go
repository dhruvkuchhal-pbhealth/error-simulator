package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
)

// UserRepository is a simple repository holding a database handle.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository constructs a UserRepository and validates the db is non-nil.
func NewUserRepository(db *sql.DB) (*UserRepository, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}
	return &UserRepository{db: db}, nil
}

// GetUserByID retrieves a user by id. It returns an error if the repository was not properly initialized.
func (r *UserRepository) GetUserByID(id int) (string, error) {
	if r == nil {
		return "", errors.New("user repository is nil")
	}
	if r.db == nil {
		return "", errors.New("database handle is nil")
	}

	var name string
	row := r.db.QueryRow("SELECT name FROM users WHERE id = ?", id)
	if err := row.Scan(&name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("user with id %d not found: %w", id, err)
		}
		return "", fmt.Errorf("query failed: %w", err)
	}
	return name, nil
}

// DBError is an http handler demonstrating usage of the repository.
func DBError(repo *UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// For demonstration, use id=1
		name, err := repo.GetUserByID(1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "user: %s", name)
	}
}
