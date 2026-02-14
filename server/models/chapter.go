package models

import (
	"time"

	"gorm.io/gorm"
)

// Chapter represents a book chapter entity
type Chapter struct {
	gorm.Model
	URL        string     `gorm:"not null" json:"url"`
	Title      string     `gorm:"not null" json:"title"`
	BookID     uint       `gorm:"column:book_id" json:"book_id,omitempty"`
	Number     string     `json:"number,omitempty"`
	Content    List       `gorm:"column:content;type:text" json:"content,omitempty"`
	UploadDate *time.Time `gorm:"column:upload_date" json:"upload_date,omitempty"`
	Completed  bool       `gorm:"not null;default:false" json:"completed"`
	Downloaded bool       `gorm:"not null;default:false" json:"downloaded"`
	OtherID    string     `gorm:"column:other_id" json:"other_id,omitempty"`

	// Relationships
	Book *Book `gorm:"foreignKey:BookID" json:"book,omitempty"`
}

// TableName specifies the table name for Chapter model
func (Chapter) TableName() string {
	return "chapters"
}
