package models

import "time"

// BookType represents the type of book content
type BookType string

const (
	BookTypeManga BookType = "manga"
	BookTypeNovel BookType = "novel"
)

// Book represents a book (manga or novel) entity
type Book struct {
	ID        int        `gorm:"primaryKey;autoIncrement" json:"id"`
	URL       string     `gorm:"not null" json:"url"`
	Title     string     `gorm:"not null" json:"title"`
	Type      BookType   `gorm:"type:text;not null;check:type IN ('manga', 'novel')" json:"type"`
	ImageURL  string     `gorm:"column:image_url" json:"image_url,omitempty"`
	Synopsis  string     `json:"synopsis,omitempty"`
	Slug      string     `json:"slug,omitempty"`
	OtherID   string     `gorm:"column:other_id" json:"other_id,omitempty"`
	SourceID  int        `gorm:"column:source_id" json:"source_id,omitempty"`
	UpdatedAt time.Time  `gorm:"not null" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	Source   *Source   `gorm:"foreignKey:SourceID" json:"source,omitempty"`
	Chapters []Chapter `gorm:"foreignKey:BookID" json:"chapters,omitempty"`
}

// TableName specifies the table name for Book model
func (Book) TableName() string {
	return "book"
}
