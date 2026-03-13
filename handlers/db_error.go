package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
)

// User represents a simple user record.
type User struct {
	ID   int
	Name string
}

// UserRepository provides access to the users store.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository constructs a UserRepository with the provided DB.
// It does not open the DB itself; callers should pass a non-nil *sql.DB.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetUserByID returns a user by ID. If the repository was not properly
// initialized with a non-nil DB, it returns an error instead of panicking.
func (r *UserRepository) GetUserByID(id int) (User, error) {
	var u User
	if r == nil || r.db == nil {
		return u, errors.New("database is not initialized")
	}

	row := r.db.QueryRow("SELECT id, name FROM users WHERE id = ?", id)
	if err := row.Scan(&u.ID, &u.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return u, fmt.Errorf("user not found: %w", err)
		}
		return u, fmt.Errorf("query scan error: %w", err)
	}
	return u, nil
}

// DBError is an HTTP handler demonstrating DB usage and error handling.
func DBError(repo *UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := repo.GetUserByID(1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "user: %+v", user)
	}
}
