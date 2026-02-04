package service

import (
	"github.com/unluckythoughts/book-reader/server/db"
	localWorker "github.com/unluckythoughts/book-reader/server/worker"
	"github.com/unluckythoughts/go-microservice/v2/tools/worker"
)

type ReaderService struct {
	db *db.DB
	w  *localWorker.ServiceWorker
}

// New creates a new instance of ReaderService
func New(db *db.DB, w *worker.Worker) *ReaderService {
	sw := localWorker.New(w, db)
	return &ReaderService{db: db, w: sw}
}
