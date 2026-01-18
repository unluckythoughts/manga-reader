package service

import "github.com/unluckythoughts/book-reader/server/db"

type ReaderService struct {
	db *db.DB
}

// New creates a new instance of ReaderService
func New(db *db.DB) *ReaderService {
	return &ReaderService{db: db}
}
