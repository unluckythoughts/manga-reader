package worker

import (
	"time"

	"github.com/pkg/errors"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/connector"
	"github.com/unluckythoughts/manga-reader/models"
	"github.com/unluckythoughts/manga-reader/utils"
)

func (w *Worker) UpdateNovelFavorites(ctx web.Context) error {
	favorites, err := w.db.GetNovelFavorites(ctx)
	if err != nil {
		return err
	}

	workerFn := func(favorite models.NovelFavorite, errChan chan<- error) {
		ctx.Logger().Debugf("Updating novel %s", favorite.Novel.Title)

		conn, err := connector.NewNovelConnector(ctx, favorite.Novel.Source.Domain)
		if err != nil {
			errChan <- errors.Wrapf(err, "error getting connector for %s", favorite.Novel.Source.Domain)
			return
		}
		novel, err := conn.GetNovelInfo(ctx, favorite.Novel.URL)
		if err != nil {
			errChan <- errors.Wrapf(err, "error getting novel info for %s", favorite.Novel.Title)
			return
		}

		novel.SourceID = favorite.Novel.SourceID
		err = w.db.UpdateNovels(ctx, &[]models.Novel{novel})
		if err != nil {
			errChan <- errors.Wrapf(err, "error updating db for %s", favorite.Novel.Title)
			return
		}

		for i := range novel.Chapters {
			novel.Chapters[i].NovelID = favorite.NovelID
			if novel.Chapters[i].UploadDate == "" {
				novel.Chapters[i].UploadDate = time.Now().Format("2006-01-02")
			}
		}

		err = w.db.UpdateNovelChapters(ctx, &novel.Chapters)
		if err != nil {
			errChan <- errors.Wrapf(err, "error updating chapters in db for %s", favorite.Novel.Title)
		}

		return
	}

	workerCount := 5
	errorList := utils.RunParallel[models.NovelFavorite, error](workerCount, favorites, workerFn)

	for _, err := range errorList {
		ctx.Logger().
			With("error", err).
			Warn("error updating novel")
	}

	if len(errorList) > 0 {
		return errorList[0]
	}

	return nil
}
