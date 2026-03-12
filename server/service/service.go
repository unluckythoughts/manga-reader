package service

import (
	"github.com/unluckythoughts/book-reader/server/db"
	localWorker "github.com/unluckythoughts/book-reader/server/worker"
	"github.com/unluckythoughts/go-microservice/v2/tools/worker"
	"go.uber.org/zap"
)

type ReaderService struct {
	db *db.DB
	w  *localWorker.ServiceWorker
	l  *zap.Logger
}

// New creates a new instance of ReaderService
func New(db *db.DB, w *worker.Worker, l *zap.Logger) *ReaderService {
	sw := localWorker.New(w, db)
	return &ReaderService{db: db, w: sw, l: l.Named("reader-service")}
}
