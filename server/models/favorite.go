package models

import (
	"github.com/unluckythoughts/go-microservice/v2/tools/auth"
	"gorm.io/gorm"
)

// Favorite represents a user's favorite book with progress tracking
type Favorite struct {
	gorm.Model
	UserID     uint `gorm:"column:user_id" json:"user_id,omitempty"`
	BookID     uint `gorm:"column:book_id" json:"book_id,omitempty"`
	Progress   List `json:"progress,omitempty"`
	Categories List `json:"categories,omitempty"`

	// Relationships
	User *auth.User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Book *Book      `gorm:"foreignKey:BookID" json:"book,omitempty"`
}

// TableName specifies the table name for Favorite model
func (Favorite) TableName() string {
	return "favorites"
}
