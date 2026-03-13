package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// UserRepository provides methods to access users in the DB.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository constructs a UserRepository with the provided DB. db may be nil,
// but callers should prefer providing an initialized *sql.DB.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetUserByID returns a user by its ID. If the repository's DB is not initialized,
// it returns an error instead of causing a panic.
func (r *UserRepository) GetUserByID(id int64) (string, error) {
	if r == nil {
		return "", errors.New("user repository is nil")
	}
	if r.db == nil {
		return "", errors.New("database is not initialized")
	}

	var name string
	// Use QueryRow on the non-nil DB
	row := r.db.QueryRow("SELECT name FROM users WHERE id = ?", id)
	if err := row.Scan(&name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return name, nil
}

// DBError demonstrates an HTTP handler that uses UserRepository.
func DBError(repo *UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		name, err := repo.GetUserByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if name == "" {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(name))
	}
}
