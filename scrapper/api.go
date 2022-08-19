package scrapper

import (
	"net/url"

	"github.com/pkg/errors"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
)

func getData(ctx web.Context, q models.APIQueryData, cb func(interface{})) error {
	c := web.NewClient(q.URL)
	params := url.Values{}
	for k, v := range q.QueryParams {
		params.Add(k, v)
	}

	for {
		apiResp := q.Response
		path := "/"
		if len(params) > 0 {
			path = "?" + params.Encode()
		}
		status, err := c.GetResponse(path, apiResp)
		if err != nil {
			return errors.Wrapf(err, "error while get data from %s", q.URL+path)
		}

		if status != 200 {
			return errors.Errorf("unexpected status %d when get data from %s", status, q.URL+path)
		}

		cb(apiResp)

		if q.HasNextPage != nil && q.HasNextPage(apiResp) {
			pageParam := "page"
			if q.PageParam != "" {
				pageParam = q.PageParam
			}
			params.Set(pageParam, strAdd(params.Get(pageParam), 1))
		} else {
			break
		}
	}

	return nil
}

func GetMangaListAPI(ctx web.Context, q models.APIQueryData, t models.MangaListTransform) ([]models.Manga, error) {
	mangas := []models.Manga{}

	getData(ctx, q, func(apiResp interface{}) {
		mangas = append(mangas, t(apiResp)...)
	})

	return mangas, nil
}

func GetMangaInfoAPI(ctx web.Context, q models.APIQueryData, t models.MangaInfoTransform) (models.Manga, error) {
	manga := models.Manga{}

	getData(ctx, q, func(apiResp interface{}) {
		manga = t(apiResp)
	})

	return manga, nil
}

func GetChapterListAPI(ctx web.Context, q models.APIQueryData, t models.ChapterListTransform) ([]models.Chapter, error) {
	chapters := []models.Chapter{}

	getData(ctx, q, func(apiResp interface{}) {
		chapters = append(chapters, t(apiResp)...)
	})

	return chapters, nil
}

func GetPagesListAPI(ctx web.Context, q models.APIQueryData, t models.PagesListTransform) ([]string, error) {
	pages := []string{}

	getData(ctx, q, func(apiResp interface{}) {
		pages = append(pages, t(apiResp)...)
	})

	return pages, nil
}
