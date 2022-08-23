package service

import (
	"github.com/unluckythoughts/manga-reader/reader/repository"
	"github.com/unluckythoughts/manga-reader/reader/worker"
	"gorm.io/gorm"
)

type Service struct {
	db repository.Repository
	w  *worker.Worker
}

func New(db *gorm.DB) *Service {
	return &Service{
		db: *repository.New(db),
		w:  worker.New(db),
	}
}
