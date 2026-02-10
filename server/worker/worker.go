package worker

import (
	"time"

	"github.com/unluckythoughts/book-reader/server/db"
	"github.com/unluckythoughts/go-microservice/v2/tools/worker"
	"github.com/unluckythoughts/go-microservice/v2/utils"
)

type ServiceWorker struct {
	w  *worker.Worker
	db *db.DB
}

type JOB_OPTIONS struct {
	JOB_UPDATE_SOURCES_SCHEDULE   string `env:"JOB_SOURCES_UPDATE_SCHEDULE" envDefault:"0 0 0 * * *"`     // every day at midnight
	JOB_UPDATE_FAVORITES_SCHEDULE string `env:"JOB_FAVORITES_UPDATE_SCHEDULE" envDefault:"0 */5 * * * *"` // every 5 minutes
}

func (s *ServiceWorker) registerTasks(opts *JOB_OPTIONS) {
	// Schedule updateFavorites every 5 minutes
	s.w.ScheduleCron("update_favorites", opts.JOB_UPDATE_FAVORITES_SCHEDULE, s.updateFavoritesTask)
	// Schedule updateSources every day at midnight
	s.w.ScheduleCron("update_sources", opts.JOB_UPDATE_SOURCES_SCHEDULE, s.updateSources)
	// Run checkDBSources once at startup
	s.w.RunInBackground("check_sources", s.checkDBSources)
	// Initial delay before first run of updateSources
	time.Sleep(2 * time.Second)
	// Run updateSources once at startup
	s.w.RunInBackground("update_sources", s.updateSources)
}

func New(w *worker.Worker, db *db.DB) *ServiceWorker {
	worker := &ServiceWorker{
		w:  w,
		db: db,
	}

	opts := &JOB_OPTIONS{}
	utils.ParseEnvironmentVars(opts)
	worker.registerTasks(opts)

	return worker
}
