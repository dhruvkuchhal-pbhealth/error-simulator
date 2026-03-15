package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/your-org/error-simulator/models"
	"gorm.io/gorm"
)

// UserRepository performs user lookups against the database.
// Uses the same users table as face-recognition-service (face_recognition DB).
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository returns a repository with real DB connection.
// Uses DATABASE_URL (default: postgresql://postgres:postgres@localhost:5432/face_recognition).
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetUserByID fetches a user by ID from the users table.
func (r *UserRepository) GetUserByID(id string) (*models.User, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return nil, err // invalid UUID format
	}
	var u models.UserModel
	if err := r.db.Where("id = ?", parsed).First(&u).Error; err != nil {
		return nil, err
	}
	return u.ToUser(), nil
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
