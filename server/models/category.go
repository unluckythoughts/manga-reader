package models

import "time"

// Category represents a user-defined category for organizing books
type Category struct {
	ID        int        `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string     `json:"name,omitempty"`
	UpdatedAt time.Time  `gorm:"not null" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName specifies the table name for Category model
func (Category) TableName() string {
	return "category"
}
