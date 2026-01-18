package db

import (
	"gorm.io/gorm"
)

// DB is the global database instance
type DB struct {
	db *gorm.DB
}

// New returns the database instance
func New(db *gorm.DB) *DB {
	return &DB{db: db}
}
