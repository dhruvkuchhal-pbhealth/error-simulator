package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

// User represents a simplified user model for the demo.
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// UserRepository provides DB operations for users.
type UserRepository struct {
	DB *sql.DB
}

// NewUserRepository constructs a UserRepository and returns an error if db is nil.
func NewUserRepository(db *sql.DB) (*UserRepository, error) {
	if db == nil {
		return nil, fmt.Errorf("db must not be nil")
	}
	return &UserRepository{DB: db}, nil
}

// GetUserByID fetches a user by id. It returns an error if the repository DB is nil.
func (r *UserRepository) GetUserByID(id int) (*User, error) {
	if r == nil {
		return nil, fmt.Errorf("user repository is nil")
	}
	if r.DB == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}

	row := r.DB.QueryRow("SELECT id, name FROM users WHERE id = ?", id)
	var u User
	if err := row.Scan(&u.ID, &u.Name); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

// DBError demonstrates an HTTP handler that uses UserRepository and reports DB errors.
func DBError(repo *UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// For demo purposes, we use a fixed id.
		user, err := repo.GetUserByID(1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if user == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("user not found"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(user)
	}
}
