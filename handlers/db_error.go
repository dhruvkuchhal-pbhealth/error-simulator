package handlers

import (
	"database/sql"
	"net/http"

	"github.com/your-org/error-simulator/models"
)

// UserRepository performs user lookups against the database.
// The bug: db is never initialized (nil); GetUserByID calls r.db.QueryRow and panics.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository returns a repository. In this test target, db is left nil
// to simulate a failed connection pool initialization in production.
func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: nil, // simulate failed DB init
	}
}

// GetUserByID fetches a user by ID. If the repository's db connection was never
// initialized, r.db is nil and r.db.QueryRow causes a nil pointer dereference.
func (r *UserRepository) GetUserByID(id string) (*models.User, error) {
	query := `SELECT id, email, first_name, last_name, created_at FROM users WHERE id = $1`
	row := r.db.QueryRow(query, id)
	var u models.User
	err := row.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// DBError handles GET /error/db.
// It calls GetUserByID on a repository with nil db to trigger the panic.
func DBError(repo *UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := repo.GetUserByID("user-abc-123")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}
