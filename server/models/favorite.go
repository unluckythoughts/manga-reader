package models

import (
	"time"

	"github.com/unluckythoughts/go-microservice/tools/auth"
)

// Favorite represents a user's favorite book with progress tracking
type Favorite struct {
	ID         int        `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     int        `gorm:"column:user_id" json:"user_id,omitempty"`
	BookID     int        `gorm:"column:book_id" json:"book_id,omitempty"`
	Progress   List       `json:"progress,omitempty"`
	Categories List       `json:"categories,omitempty"`
	UpdatedAt  time.Time  `gorm:"not null" json:"updated_at"`
	DeletedAt  *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	User *auth.User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Book *Book      `gorm:"foreignKey:BookID" json:"book,omitempty"`
}

// TableName specifies the table name for Favorite model
func (Favorite) TableName() string {
	return "favorite"
}
