package worker

import (
	"fmt"

	"github.com/unluckythoughts/book-reader/server/connector"
	"github.com/unluckythoughts/book-reader/server/models"
	"github.com/unluckythoughts/go-microservice/tools/web"
)

func (w *ServiceWorker) updateFavoritesTask(ctx web.Context) error {
	ctx.Logger().Info("Updating favorites...")
	// Implementation for updating favorites
	return nil
}

func (w *ServiceWorker) checkDBSources(ctx web.Context) error {
	ctx.Logger().Info("Checking DB sources for all connectors...")

	sources, err := w.db.GetAllSources()
	if err != nil {
		ctx.Logger().Error("Failed to get sources: ", err)
		return err
	}

	conns := connector.GetAllConnectors()

	for _, conn := range conns {
		found := false
		for _, source := range sources {
			if source.Name == conn.GetName() {
				found = true
			}
		}

		if !found {
			ctx.Logger().Infof("Registering new source from connector: %s", conn.GetName())
			newSource := &models.Source{
				Name:    conn.GetName(),
				Domain:  conn.GetDomain(),
				IconURL: conn.GetIconURL(),
			}
			err := w.db.CreateSource(newSource)
			if err != nil {
				ctx.Logger().Errorf("Failed to create source: %v", err)
			}
		}
	}

	return nil
}

func (w *ServiceWorker) updateSources(ctx web.Context) error {
	ctx.Logger().Info("Running updateSources...")

	sources, err := w.db.GetAllSources()
	if err != nil {
		ctx.Logger().Error("Failed to get sources: ", err)
		return err
	}
	conns := connector.GetAllConnectors()

	for _, source := range sources {
		conn, ok := conns[source.Name]
		if !ok {
			ctx.Logger().Errorf("No connector found for source: %s", source.Name)
			return fmt.Errorf("no connector found for source: %s", source.Name)
		}

		count, books, err := conn.GetBookCount()
		if err != nil {
			ctx.Logger().Errorf("Failed to get book count for source %s: %v", source.Name, err)
			return err
		}

		dbBooksCount, err := w.db.GetBookCountBySourceID(source.ID)
		if err != nil {
			ctx.Logger().Errorf("Failed to get DB book count for source %s: %v", source.Name, err)
			return err
		}
		if count <= int(dbBooksCount) {
			ctx.Logger().Infof("No new books for source %s. DB count: %d, Connector count: %d", source.Name, dbBooksCount, count)
			continue
		}

		ctx.Logger().Infof("New books found for source %s. DB count: %d, Connector count: %d", source.Name, dbBooksCount, count)
		for _, book := range books {
			exists, err := w.db.CheckBookExists(source.ID, book.URL)
			if err != nil {
				ctx.Logger().Errorf("Failed to check if book exists for source %s: %v", source.Name, err)
				return err
			}
			if !exists {
				book.SourceID = source.ID
				err := w.db.CreateBook(&book)
				if err != nil {
					ctx.Logger().Errorf("Failed to create book for source %s: %v", source.Name, err)
					return err
				}
				ctx.Logger().Infof("Added new book '%s' to source %s", book.Title, source.Name)
			}
		}
	}
	return nil
}
