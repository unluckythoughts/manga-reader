package worker

import (
	"time"

	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
	"go.uber.org/zap"
)

func (w *Worker) UpdateSourceMangas(ctx web.Context, domain string, mangas []models.Manga) {
	go func() {
		ctx = web.NewContext(ctx.Logger().Desugar())
		source, err := w.db.FindSourceByDomain(ctx, domain)
		if err != nil {
			ctx.Logger().With(zap.Error(err)).Errorf("could not find source for %+s", domain)
			return
		}

		for i := 0; i < len(mangas); i++ {
			mangas[i].SourceID = source.ID
		}

		err = w.db.UpdateMangas(ctx, &mangas)
		if err != nil {
			ctx.Logger().With(zap.Error(err)).Errorf("could not update manga for %+s", domain)
			return
		}

		source.UpdatedAt = time.Now().Format("2006-01-02")

		err = w.db.SaveSource(ctx, &source)
		if err != nil {
			ctx.Logger().With(zap.Error(err)).Errorf("could not update timestamp for %+s", domain)
			return
		}
	}()
}

func (w *Worker) UpdateSourceManga(ctx web.Context, domain string, manga models.Manga) {
	go func() {
		ctx = web.NewContext(ctx.Logger().Desugar())
		source, err := w.db.FindSourceByDomain(ctx, domain)
		if err != nil {
			ctx.Logger().With(zap.Error(err)).Errorf("could not find source for %+s", domain)
			return
		}

		manga.SourceID = source.ID
		for i := range manga.Chapters {
			if manga.Chapters[i].UploadDate == "" {
				manga.Chapters[i].UploadDate = time.Now().Format("2006-01-02")
			}
		}

		err = w.db.UpdateMangas(ctx, &[]models.Manga{manga})
		if err != nil {
			ctx.Logger().With(zap.Error(err)).Errorf("could not update manga for %+s", domain)
			return
		}

		// err = w.db.UpdateChapters(ctx, &manga.Chapters)
		// if err != nil {
		// 	ctx.Logger().With(zap.Error(err)).Errorf("could not update chapters of %+s for", manga.Title, domain)
		// 	return
		// }
	}()
}
