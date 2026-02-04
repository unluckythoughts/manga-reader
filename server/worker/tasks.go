package worker

import (
	"github.com/unluckythoughts/book-reader/server/connector"
	"github.com/unluckythoughts/book-reader/server/models"
	"github.com/unluckythoughts/go-microservice/v2/tools/web"
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

func (w *ServiceWorker) updateSourceBooks(ctx web.Context, source models.Source, conn models.IConnector) error {
	dbBooksCount, err := w.db.GetBookCountBySourceID(source.ID)
	if err != nil {
		ctx.Logger().Errorf("Failed to get DB book count for source %s: %v", source.Name, err)
		return err
	}

	books := []models.Book{}
	// If no books in DB, fetch all books
	if dbBooksCount <= 0 {
		books, err = conn.GetBooks()
		if err != nil {
			ctx.Logger().Errorf("Failed to get all books for source %s: %v", source.Name, err)
			return err
		}
	} else {
		var count int
		// books will be nil if connector supports only count fetching
		count, books, err = conn.GetBookCount()
		if err != nil {
			ctx.Logger().Errorf("Failed to get book count for source %s: %v", source.Name, err)
			return err
		}

		if count <= int(dbBooksCount) {
			ctx.Logger().Infof("No new books for source %s. DB count: %d, Connector count: %d", source.Name, dbBooksCount, count)
			return nil
		}

		// Fetch all books if books is empty
		// Books will be empty if connector supports last page selector
		if len(books) == 0 {
			books, err = conn.GetBooks()
			if err != nil {
				ctx.Logger().Errorf("Failed to get all books for source %s: %v", source.Name, err)
				return err
			}
		}
	}

	ctx.Logger().Infof("New books found for source %s. DB count: %d, Connector count: %d", source.Name, dbBooksCount, len(books))
	for _, book := range books {
		book.SourceID = source.ID
		err = w.db.CreateBook(&book)
		if err != nil {
			ctx.Logger().Errorf("Failed to create book for source %s: %v", source.Name, err)
			return err
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
			continue
		}

		err := w.updateSourceBooks(ctx, source, conn)
		if err != nil {
			ctx.Logger().Errorf("Failed to update books for source %s: %v", source.Name, err)
			continue
		}
	}

	return nil
}
