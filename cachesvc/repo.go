package cachesvc

import (
	"database/sql"

	"github.com/your-org/error-simulator/models"
)

// Repo fetches users from DB. BUG: db is nil; FindByID panics (multi-file stack).
type Repo struct {
	db *sql.DB
}

// NewRepo returns a repo with nil db (simulates failed DB init).
func NewRepo() *Repo {
	return &Repo{db: nil}
}

// FindByID runs a query. Panics when r.db is nil.
func (r *Repo) FindByID(id string) (*models.User, error) {
	row := r.db.QueryRow("SELECT id, email FROM users WHERE id = $1", id)
	var u models.User
	err := row.Scan(&u.ID, &u.Email)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
