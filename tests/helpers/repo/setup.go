package repo

import (
	"github.com/unluckythoughts/go-microservice/tools/logger"
	"github.com/unluckythoughts/go-microservice/tools/sqlite"
	_ "gitlab.cobalt.rocks/coderdojo/sqlite-regexp.git"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func Setup() {
	opts := sqlite.Options{
		Logger:   logger.New(logger.Options{}),
		Filepath: "/home/vinay/workspace/personal/manga-reader/tests/helpers/repo/test.db",
	}

	db = sqlite.New(opts)
	err := db.Exec("SELECT '' REGEXP '';").Error
	if err != nil {
		db.Logger.Error(db.Statement.Context, "could not run REGEXP")
		panic(err)
	}
}
