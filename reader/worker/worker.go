package worker

import (
	"context"

	"github.com/unluckythoughts/manga-reader/connector"
	"github.com/unluckythoughts/manga-reader/models"
	"github.com/unluckythoughts/manga-reader/reader/repository"
	"gorm.io/gorm"
)

type Worker struct {
	db repository.Repository
}

func setupConfig(db *gorm.DB) {
	r := repository.New(db)
	sources := []models.Source{}
	for _, conn := range connector.GetAllConnectors() {
		sources = append(sources, conn.GetSource())
	}

	err := r.CreateSources(context.Background(), &sources)
	if err != nil {
		db.Logger.Error(nil, "could not update sources")
		panic(err)
	}
}

func New(db *gorm.DB) *Worker {
	setupConfig(db)

	return &Worker{
		db: *repository.New(db),
	}
}
