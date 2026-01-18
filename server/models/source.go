package models

import "time"

// Source represents a book source website entity
type Source struct {
	ID        int        `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string     `gorm:"unique;not null" json:"name"`
	Domain    string     `gorm:"unique;not null" json:"domain"`
	IconURL   string     `gorm:"column:icon_url" json:"icon_url,omitempty"`
	UpdatedAt time.Time  `gorm:"not null" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	Books []Book `gorm:"foreignKey:SourceID" json:"books,omitempty"`
}

// TableName specifies the table name for Source model
func (Source) TableName() string {
	return "source"
}
