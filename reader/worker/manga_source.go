package worker

import (
	"errors"
	"time"

	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/connector"
	"github.com/unluckythoughts/manga-reader/models"
	"go.uber.org/zap"
)

func (w *Worker) UpdateSourceMangasSync(ctx web.Context, source models.Source, mangas []models.Manga, fix bool) error {
	if len(mangas) == 0 {
		return errors.New("no mangas to update")
	}

	ctx = web.NewContext(ctx.Logger().Desugar())
	for i := 0; i < len(mangas); i++ {
		mangas[i].SourceID = source.ID
		mangas[i].Source.ID = source.ID
	}

	if fix {
		err := w.db.DeleteChaptersBySource(ctx, source.ID)
		if err != nil {
			ctx.Logger().With(zap.Error(err)).Errorf("could not delete chapters for %+s", source.Domain)
		}
	}

	if len(mangas) > 1000 {
		batch := 500
		for i := 0; i < len(mangas); i = i + batch {
			batchMangas := []models.Manga{}
			if i+batch > len(mangas) {
				batchMangas = mangas[i:]
			} else {
				batchMangas = mangas[i : i+batch]
			}

			err := w.db.UpdateMangas(ctx, &batchMangas)
			if err != nil {
				ctx.Logger().With(zap.Error(err)).Errorf("could not update manga for %+s", source.Domain)
				return err
			}
		}
	} else {
		err := w.db.UpdateMangas(ctx, &mangas)
		if err != nil {
			ctx.Logger().With(zap.Error(err)).Errorf("could not update manga for %+s", source.Domain)
			return err
		}
	}

	err := w.db.SaveSource(ctx, &source)
	if err != nil {
		ctx.Logger().With(zap.Error(err)).Errorf("could not update timestamp for %+s", source.Domain)
		return err
	}

	return nil
}

func (w *Worker) UpdateAllSourceMangas(ctx web.Context) error {
	sources, err := w.db.GetSources(ctx)
	if err != nil {
		return err
	}

	for _, source := range sources {
		conn, err := connector.GetMangaConnector(ctx, source.Domain)
		if err != nil {
			return err
		}

		go func(src models.Source) {
			var mangas []models.Manga
			latest, err := w.db.GetLatestMangaForSource(ctx, src.ID)
			if err != nil {
				mangas, err = conn.GetMangaList(ctx)
			} else {
				mangas, err = conn.GetLatestMangaList(ctx, latest.Title)
			}

			if err != nil {
				ctx.Logger().With(zap.Error(err)).Errorf("could not get manga list for %+s", src.Domain)
			}

			w.UpdateSourceMangasSync(ctx, src, mangas, false)
		}(source)
	}

	return nil
}

func (w *Worker) UpdateSourceMangas(ctx web.Context, source models.Source, mangas []models.Manga, fix bool) {
	go w.UpdateSourceMangasSync(ctx, source, mangas, fix)
}

func (w *Worker) UpdateSourceMangaSync(ctx web.Context, domain string, manga *models.Manga) error {
	ctx = web.NewContext(ctx.Logger().Desugar())
	source, err := w.db.FindSourceByDomain(ctx, domain)
	if err != nil {
		ctx.Logger().With(zap.Error(err)).Errorf("could not find source for %+s", domain)
		return err
	}

	manga.Source = source
	mangas := []models.Manga{*manga}
	err = w.db.UpdateMangas(ctx, &mangas)
	if err != nil {
		ctx.Logger().With(zap.Error(err)).Errorf("could not update manga for %+s", domain)
		return err
	}
	manga.ID = mangas[0].ID

	if len(manga.Chapters) > 0 {
		for i := range manga.Chapters {
			manga.Chapters[i].MangaID = manga.ID
			if manga.Chapters[i].UploadDate == "" {
				manga.Chapters[i].UploadDate = time.Now().Format("2006-01-02")
			}
		}

		err = w.db.UpdateChapters(ctx, &manga.Chapters)
		if err != nil {
			ctx.Logger().With(zap.Error(err)).Errorf("could not update chapters of %+s for", manga.Title, domain)
			return err
		}
	}

	return nil
}

func (w *Worker) UpdateSourceManga(ctx web.Context, domain string, manga *models.Manga) {
	go w.UpdateSourceMangaSync(ctx, domain, manga)
}
