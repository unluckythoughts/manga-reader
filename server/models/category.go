package models

import "gorm.io/gorm"

// Category represents a user-defined category for organizing books
type Category struct {
	gorm.Model
	Name string `json:"name,omitempty"`
}

// TableName specifies the table name for Category model
func (Category) TableName() string {
	return "categories"
}
