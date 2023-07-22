package worker

import (
	"time"

	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
	"go.uber.org/zap"
)

func (w *Worker) UpdateSourceNovels(ctx web.Context, domain string, novels []models.Novel) {
	go func() {
		if len(novels) == 0 {
			return
		}

		ctx = web.NewContext(ctx.Logger().Desugar())
		source, err := w.db.FindNovelSourceByDomain(ctx, domain)
		if err != nil {
			ctx.Logger().With(zap.Error(err)).Errorf("could not find source for %+s", domain)
			return
		}

		for i := 0; i < len(novels); i++ {
			novels[i].SourceID = source.ID
		}

		if len(novels) > 1000 {
			batch := 500
			for i := 0; i < len(novels); i = i + batch {
				batchNovels := []models.Novel{}
				if i+batch > len(novels) {
					batchNovels = novels[i:]
				} else {
					batchNovels = novels[i : i+batch]
				}

				err = w.db.UpdateNovels(ctx, &batchNovels)
				if err != nil {
					ctx.Logger().With(zap.Error(err)).Errorf("could not update novel for %+s", domain)
					return
				}
			}
		} else {
			err = w.db.UpdateNovels(ctx, &novels)
			if err != nil {
				ctx.Logger().With(zap.Error(err)).Errorf("could not update novel for %+s", domain)
				return
			}
		}

		err = w.db.SaveNovelSource(ctx, &source)
		if err != nil {
			ctx.Logger().With(zap.Error(err)).Errorf("could not update timestamp for %+s", domain)
			return
		}

	}()
}

func (w *Worker) UpdateSourceNovel(ctx web.Context, domain string, novel models.Novel) {
	go func() {
		ctx = web.NewContext(ctx.Logger().Desugar())
		source, err := w.db.FindNovelSourceByDomain(ctx, domain)
		if err != nil {
			ctx.Logger().With(zap.Error(err)).Errorf("could not find source for %+s", domain)
			return
		}

		novel.SourceID = source.ID
		novels := []models.Novel{novel}
		err = w.db.UpdateNovels(ctx, &novels)
		if err != nil {
			ctx.Logger().With(zap.Error(err)).Errorf("could not update novel for %+s", domain)
			return
		}
		novel.ID = novels[0].ID

		if len(novel.Chapters) > 1 {
			for i := range novel.Chapters {
				novel.Chapters[i].NovelID = novel.ID
				if novel.Chapters[i].UploadDate == "" {
					novel.Chapters[i].UploadDate = time.Now().Format("2006-01-02")
				}
			}

			err = w.db.UpdateNovelChapters(ctx, &novel.Chapters)
			if err != nil {
				ctx.Logger().With(zap.Error(err)).Errorf("could not update chapters of %+s for", novel.Title, domain)
				return
			}
		}
	}()
}
