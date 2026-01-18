package models

import "time"

// User represents a user entity
type User struct {
	ID        int        `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string     `json:"name,omitempty"`
	UpdatedAt time.Time  `gorm:"not null" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	Favorites []Favorite `gorm:"foreignKey:UserID" json:"favorites,omitempty"`
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "user"
}
