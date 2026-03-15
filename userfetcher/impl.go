package userfetcher

import (
	"github.com/google/uuid"
	"github.com/your-org/error-simulator/models"
	"gorm.io/gorm"
)

// Impl implements usersvc.Fetcher. Fetches from the real users DB (same as face-recognition-service).
type Impl struct {
	db *gorm.DB
}

// NewImpl returns an impl that fetches from the real users DB.
func NewImpl(db *gorm.DB) *Impl {
	return &Impl{db: db}
}

// FetchUser implements Fetcher. Fetches from DB using GORM.
func (i *Impl) FetchUser(id string) (*models.User, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	var u models.UserModel
	if err := i.db.Where("id = ?", parsed).First(&u).Error; err != nil {
		return nil, err
	}
	return u.ToUser(), nil
}
