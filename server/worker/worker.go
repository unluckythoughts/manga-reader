package worker

import (
	"github.com/unluckythoughts/book-reader/server/db"
	"github.com/unluckythoughts/go-microservice/tools/worker"
)

type ServiceWorker struct {
	w  *worker.Worker
	db *db.DB
}

func (s *ServiceWorker) registerTasks() {
	// Run checkDBSources once at startup
	s.w.RunInBackground("check_sources", s.checkDBSources)
	// Schedule updateFavorites every 5 minutes
	s.w.ScheduleCron("update_favorites", "0 */5 * * * *", s.updateFavoritesTask)
	// Schedule updateSources every day at midnight
	s.w.ScheduleCron("update_sources", "0 0 0 * * *", s.updateSources)
}

func New(w *worker.Worker, db *db.DB) *ServiceWorker {
	worker := &ServiceWorker{
		w:  w,
		db: db,
	}
	worker.registerTasks()

	return worker
}
