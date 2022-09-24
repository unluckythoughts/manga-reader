package worker

import (
	"context"

	"github.com/unluckythoughts/manga-reader/connector"
	"github.com/unluckythoughts/manga-reader/models"
	"github.com/unluckythoughts/manga-reader/reader/repository"
	"gorm.io/gorm"

	_ "gitlab.cobalt.rocks/coderdojo/sqlite-regexp.git"
)

type Worker struct {
	db repository.Repository
}

func setupConfig(db *gorm.DB) {
	r := repository.New(db)

	err := db.Exec("SELECT '' REGEXP '';").Error
	if err != nil {
		db.Logger.Error(nil, "could not update sources")
		panic(err)
	}

	sources := []models.Source{}
	for _, conn := range connector.GetAllMangaConnectors() {
		sources = append(sources, conn.GetSource())
	}

	err = r.CreateSources(context.Background(), &sources)
	if err != nil {
		db.Logger.Error(nil, "could not update sources")
		panic(err)
	}

	novelSources := []models.NovelSource{}
	for _, conn := range connector.GetAllNovelConnectors() {
		novelSources = append(novelSources, conn.GetSource())
	}

	err = r.CreateNovelSources(context.Background(), &novelSources)
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
