package models

import (
	"time"

	"github.com/google/uuid"
)

// UserModel is the GORM model for the users table.
// Matches face-recognition-service schema: id, email, name, password_hash, created_at.
type UserModel struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email        string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	Name         string    `gorm:"type:varchar(255)"`
	PasswordHash string    `gorm:"column:password_hash;type:varchar(255)"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

// TableName overrides GORM's default table name.
func (UserModel) TableName() string {
	return "users"
}

// ToUser converts the GORM model to the domain User.
func (u *UserModel) ToUser() *User {
	return &User{
		ID:        u.ID.String(),
		Email:     u.Email,
		FirstName: u.Name,
		CreatedAt: u.CreatedAt,
	}
}
