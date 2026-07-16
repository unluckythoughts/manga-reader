package models

import "github.com/unluckythoughts/go-microservice/v2/tools/auth"

const (
	UserRole  auth.Role = 1
	AdminRole auth.Role = 99
)

// CreateBookRequest represents the request body for creating a new book
type CreateBookRequest struct {
	URL      string   `json:"url" valid:"required"`
	Title    string   `json:"title" valid:"required"`
	Type     BookType `json:"type" valid:"required"`
	ImageURL string   `json:"image_url,omitempty"`
	Synopsis string   `json:"synopsis,omitempty"`
	Slug     string   `json:"slug,omitempty"`
	OtherID  string   `json:"other_id,omitempty"`
	SourceID uint     `json:"source_id,omitempty"`
}

// CreateChapterRequest represents the request body for creating a new chapter
type CreateChapterRequest struct {
	URL        string `json:"url" valid:"required"`
	Title      string `json:"title" valid:"required"`
	BookID     uint   `json:"book_id" valid:"required"`
	Number     string `json:"number,omitempty"`
	Content    string `json:"content,omitempty"`
	UploadDate string `json:"upload_date,omitempty"`
	Completed  bool   `json:"completed,omitempty"`
	Downloaded bool   `json:"downloaded,omitempty"`
	OtherID    string `json:"other_id,omitempty"`
}

// AddFavoriteRequest represents the request body for adding data to user favorite
type AddFavoriteRequest struct {
	SourceID uint   `json:"source_id" valid:"required"`
	BookID   uint   `json:"book_id" valid:"required"`
	Progress string `json:"progress,omitempty"`
}

// UpdateFavoriteProgressRequest represents the request body for updating a favorite's progress
type UpdateFavoriteProgressRequest struct {
	BookID  uint   `json:"book_id" valid:"required"`
	Chapter string `json:"chapter" valid:"required"`
	Level   int    `json:"level" valid:"optional,numeric,min=0"`
}

// CreateCategoryRequest represents the request body for creating a category
type CreateCategoryRequest struct {
	Name string `json:"name" valid:"required"`
}

// UpdateCategoryRequest represents the request body for updating a category
type UpdateCategoryRequest struct {
	Name string `json:"name" valid:"required"`
}
