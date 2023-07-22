package main

import (
	"github.com/unluckythoughts/go-microservice/tools/logger"
	"github.com/unluckythoughts/go-microservice/tools/sqlite"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
	"github.com/unluckythoughts/manga-reader/reader/repository"
	"gorm.io/gorm"
)

var (
	l = logger.New(logger.Options{
		LogLevel: "error",
	})
)

func getSqliteDB() *gorm.DB {
	opts := sqlite.Options{
		Filepath: "./db.sqlite",
		Logger:   l,
	}
	return sqlite.New(opts)
}

func getRepo() *repository.Repository {
	db := getSqliteDB()
	return repository.New(db)
}

func getFavs() ([]models.MangaFavorite, error) {
	r := getRepo()
	return r.GetFavorites(web.NewContext(l))
}
