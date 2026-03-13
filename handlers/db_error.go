package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
)

// UserRepository wraps a *sql.DB
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository constructs a UserRepository and validates the db is non-nil
func NewUserRepository(db *sql.DB) (*UserRepository, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}
	return &UserRepository{db: db}, nil
}

// GetUserByID returns the user for the given id
func (repo *UserRepository) GetUserByID(id int) (string, error) {
	if repo == nil || repo.db == nil {
		return "", fmt.Errorf("database not initialized: repo or repo.db is nil")
	}

	row := repo.db.QueryRow("SELECT name FROM users WHERE id = ?", id)
	var name string
	if err := row.Scan(&name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return name, nil
}

// DBError is an HTTP handler demonstrating DB usage
func DBError(repo *UserRepository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name, err := repo.GetUserByID(1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if name == "" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("user not found"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(name))
	})
}
