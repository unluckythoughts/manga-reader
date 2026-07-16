package models

import (
	"github.com/unluckythoughts/go-microservice/v2/tools/auth"
)

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Items      interface{} `json:"items"`
	Pagination Pagination  `json:"pagination"`
}

// Pagination contains pagination metadata
type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// CalculateTotalPages calculates the total number of pages
func (p *Pagination) CalculateTotalPages() {
	if p.Limit > 0 {
		p.TotalPages = int((p.Total + int64(p.Limit) - 1) / int64(p.Limit))
	}
}

// BooksPaginatedResponse represents a paginated response for books
type BooksPaginatedResponse struct {
	Items      []Book     `json:"items"`
	Pagination Pagination `json:"pagination"`
}

// SourcesPaginatedResponse represents a paginated response for sources
type SourcesPaginatedResponse struct {
	Items      []Source   `json:"items"`
	Pagination Pagination `json:"pagination"`
}

// UsersPaginatedResponse represents a paginated response for users
type UsersPaginatedResponse struct {
	Items      []auth.User `json:"items"`
	Pagination Pagination  `json:"pagination"`
}

type FavoriteBookValue struct {
	Book     Book     `json:"book"`
	Progress Progress `json:"progress"`
}

type FavoriteResponse struct {
	Sources  []Source                   `json:"sources"`
	Books    map[uint]FavoriteBookValue `json:"books"`
	Settings map[string]string          `json:"settings"`
}

type ChaptersPaginatedResponse struct {
	Items      []Chapter  `json:"items"`
	Pagination Pagination `json:"pagination"`
}
