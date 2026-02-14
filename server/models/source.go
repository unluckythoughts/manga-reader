package models

import "gorm.io/gorm"

// Source represents a book source website entity
type Source struct {
	gorm.Model
	Name    string `gorm:"unique;not null" json:"name"`
	Domain  string `gorm:"unique;not null" json:"domain"`
	IconURL string `gorm:"column:icon_url" json:"icon_url,omitempty"`

	// Relationships
	Books []Book `gorm:"foreignKey:SourceID" json:"books,omitempty"`
}

// TableName specifies the table name for Source model
func (Source) TableName() string {
	return "sources"
}
