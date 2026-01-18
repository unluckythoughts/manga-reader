package models

import "time"

// Chapter represents a book chapter entity
type Chapter struct {
	ID         int        `gorm:"primaryKey;autoIncrement" json:"id"`
	URL        string     `gorm:"not null" json:"url"`
	Title      string     `gorm:"not null" json:"title"`
	BookID     int        `gorm:"column:book_id" json:"book_id,omitempty"`
	Number     string     `json:"number,omitempty"`
	Content    Content    `gorm:"column:content;type:text" json:"content,omitempty"`
	UploadDate *time.Time `gorm:"column:upload_date" json:"upload_date,omitempty"`
	Completed  bool       `gorm:"not null;default:false" json:"completed"`
	Downloaded bool       `gorm:"not null;default:false" json:"downloaded"`
	OtherID    string     `gorm:"column:other_id" json:"other_id,omitempty"`
	UpdatedAt  time.Time  `gorm:"not null" json:"updated_at"`
	DeletedAt  *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	Book *Book `gorm:"foreignKey:BookID" json:"book,omitempty"`
}

// TableName specifies the table name for Chapter model
func (Chapter) TableName() string {
	return "chapter"
}
