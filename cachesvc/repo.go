package cachesvc

import (
	"github.com/google/uuid"
	"github.com/your-org/error-simulator/models"
	"gorm.io/gorm"
)

// Repo fetches users from DB. Uses same users table as face-recognition-service.
type Repo struct {
	db *gorm.DB
}

// NewRepo returns a repo with real DB connection.
func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

// FindByID fetches a user by ID using GORM.
func (r *Repo) FindByID(id string) (*models.User, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	var u models.UserModel
	if err := r.db.Where("id = ?", parsed).First(&u).Error; err != nil {
		return nil, err
	}
	return u.ToUser(), nil
}
