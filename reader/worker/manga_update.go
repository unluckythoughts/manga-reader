package worker

import (
	"time"

	"github.com/pkg/errors"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/connector"
	"github.com/unluckythoughts/manga-reader/models"
	"github.com/unluckythoughts/manga-reader/utils"
)

func (w *Worker) UpdateFavorites(ctx web.Context) error {
	favorites, err := w.db.GetFavorites(ctx)
	if err != nil {
		return err
	}

	workerFn := func(favorite models.MangaFavorite, errChan chan<- error) {
		ctx.Logger().Debugf("Updating manga %s", favorite.Manga.Title)

		conn, err := connector.GetMangaConnector(ctx, favorite.Manga.Source.Domain)
		if err != nil {
			errChan <- errors.Wrapf(err, "error getting connector for %s", favorite.Manga.Source.Domain)
			return
		}
		manga, err := conn.GetMangaInfo(ctx, favorite.Manga.URL)
		if err != nil {
			errChan <- errors.Wrapf(err, "error getting manga info for %s", favorite.Manga.Title)
			return
		}

		manga.SourceID = favorite.Manga.SourceID
		err = w.db.UpdateMangas(ctx, &[]models.Manga{manga})
		if err != nil {
			errChan <- errors.Wrapf(err, "error updating db for %s", favorite.Manga.Title)
			return
		}

		for i := range manga.Chapters {
			manga.Chapters[i].MangaID = favorite.MangaID
			if manga.Chapters[i].UploadDate == "" {
				manga.Chapters[i].UploadDate = time.Now().Format("2006-01-02")
			}
		}

		err = w.db.UpdateChapters(ctx, &manga.Chapters)
		if err != nil {
			errChan <- errors.Wrapf(err, "error updating chapters in db for %s", favorite.Manga.Title)
		}

		return
	}

	workerCount := 5
	errorList := utils.RunParallel[models.MangaFavorite, error](workerCount, favorites, workerFn)

	for _, err := range errorList {
		ctx.Logger().
			With("error", err).
			Warn("error updating manga")
	}

	if len(errorList) > 0 {
		return errorList[0]
	}

	return nil
}
