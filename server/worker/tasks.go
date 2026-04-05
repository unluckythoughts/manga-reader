package worker

import (
	"github.com/unluckythoughts/book-reader/server/connector"
	"github.com/unluckythoughts/book-reader/server/models"
	"github.com/unluckythoughts/go-microservice/v2/tools/context"
	"go.uber.org/zap"
)

func (w *ServiceWorker) updateFavoritesTask(ctx context.Context) error {
	ctx.Logger().Info("Updating favorites...")
	// Implementation for updating favorites
	return nil
}

func (w *ServiceWorker) checkDBSources(ctx context.Context) error {
	ctx.Logger().Info("Checking DB sources for all connectors...")

	sources, err := w.db.GetAllSources()
	if err != nil {
		ctx.Logger().Error("Failed to get sources: ", zap.Error(err))
		return err
	}

	conns := connector.GetAllConnectors()
	ctx.Logger().Sugar().Infof("found %d connectors", len(conns))
	for _, conn := range conns {
		found := false
		for _, source := range sources {
			if source.Name == conn.GetName() {
				found = true
			}
		}

		if !found {
			ctx.Logger().Sugar().Infof("Registering new source from connector: %s", conn.GetName())
			newSource := &models.Source{
				Name:    conn.GetName(),
				Domain:  conn.GetDomain(),
				IconURL: conn.GetIconURL(),
			}
			err := w.db.CreateSource(newSource)
			if err != nil {
				ctx.Logger().Sugar().Errorf("Failed to create source: %v", err)
			}
		}
	}

	return nil
}

func (w *ServiceWorker) updateSourceBooks(ctx context.Context, source models.Source, conn models.IConnector) error {
	dbBooksCount, err := w.db.GetBookCountBySourceID(source.ID)
	if err != nil {
		ctx.Logger().Sugar().Errorf("Failed to get DB book count for source %s: %v", source.Name, err)
		return err
	}
	ctx.Logger().Sugar().Infof("Source %s has %d books in DB", source.Name, dbBooksCount)

	if dbBooksCount > 0 && conn.SupportsLastPageSelector() {
		// books will be nil if connector supports only count fetching
		count, _, err := conn.GetBookCount()
		if err != nil {
			ctx.Logger().Sugar().Errorf("Failed to get book count for source %s: %v", source.Name, err)
			return err
		}

		if count <= int(dbBooksCount) {
			ctx.Logger().Sugar().Infof("No new books for source %s. DB count: %d, Connector count: %d", source.Name, dbBooksCount, count)
			return nil
		}
	}

	bookChan, err := conn.GetBooksAsync()
	if err != nil {
		ctx.Logger().Sugar().Errorf("Failed to get all books for source %s: %v", source.Name, err)
		return err
	}

	for book := range bookChan {
		book.SourceID = source.ID
		err := w.db.CreateBook(&book)
		if err != nil {
			ctx.Logger().Sugar().Errorf("Failed to create book for source %s: %v", source.Name, err)
			continue
		}
		ctx.Logger().Sugar().Debugf("Added or Updated book '%s' for source %s", book.Title, source.Name)
	}

	return nil
}

func (w *ServiceWorker) updateSources(ctx context.Context) error {
	ctx.Logger().Info("Running updateSources...")

	sources, err := w.db.GetAllSources()
	if err != nil {
		ctx.Logger().Error("Failed to get sources: ", zap.Error(err))
		return err
	}
	conns := connector.GetAllConnectors()

	for _, source := range sources {
		conn, ok := conns[source.Name]
		if !ok {
			ctx.Logger().Sugar().Errorf("No connector found for source: %s", source.Name)
			continue
		}

		err := w.updateSourceBooks(ctx, source, conn)
		if err != nil {
			ctx.Logger().Sugar().Errorf("Failed to update books for source %s: %v", source.Name, err)
			continue
		}
	}

	return nil
}
