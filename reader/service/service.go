package service

import (
	"github.com/unluckythoughts/manga-reader/reader/repository"
	"gorm.io/gorm"
)

type Service struct {
	db repository.Repository
}

func New(db *gorm.DB) *Service {
	return &Service{
		db: *repository.New(db),
	}
}
