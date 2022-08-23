package connector

import (
	"net/url"

	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
	"github.com/unluckythoughts/manga-reader/scrapper"
	"go.uber.org/zap"
)

type zero models.Source

func (z *zero) GetSource() models.Source {
	return models.Source(*z)
}

func (z *zero) GetMangaList(ctx web.Context) ([]models.Manga, error) {
	type apiRespStruct struct {
		Data struct {
			Comics []struct {
				ID      strFloat `json:"id"`
				Slug    string   `json:"Slug"`
				Name    string   `json:"name"`
				Summary string   `json:"summary"`
				Cover   struct {
					Horizontal string `json:"horizontal"`
					Vertical   string `json:"vertical"`
					Full       string `json:"full"`
				} `json:"cover"`
			} `json:"comics"`
		} `json:"data"`
	}

	q := models.APIQueryData{
		URL:      "https://zeroscans.com/swordflake/comics",
		Response: &apiRespStruct{},
	}

	transform := func(r interface{}) (mangas []models.Manga) {
		resp, ok := r.(*apiRespStruct)
		if !ok {
			return mangas
		}

		for _, item := range resp.Data.Comics {
			imageURL := item.Cover.Full
			if imageURL == "" {
				imageURL = item.Cover.Horizontal
				if imageURL == "" {
					imageURL = item.Cover.Vertical
				}
			}

			mangas = append(mangas, models.Manga{
				Title:    item.Name,
				ImageURL: imageURL,
				URL:      "https://zeroscans.com/swordflake/comic/" + item.Slug,
				OtherID:  string(item.ID),
				Synopsis: item.Summary,
			})
		}

		return mangas
	}

	return scrapper.GetMangaListAPI(ctx, q, transform)
}

func (z *zero) GetMangaInfo(ctx web.Context, mangaURL string) (models.Manga, error) {
	type apiRespStruct struct {
		Data struct {
			ID      strFloat `json:"id"`
			Name    string   `json:"name"`
			Slug    string   `json:"slug"`
			Summary string   `json:"summary"`
			Cover   struct {
				Horizontal string `json:"horizontal"`
				Vertical   string `json:"vertical"`
				Full       string `json:"full"`
			} `json:"cover"`
		} `json:"data"`
	}

	transform := func(r interface{}) models.Manga {
		resp, ok := r.(*apiRespStruct)
		if !ok {
			return models.Manga{}
		}

		imageURL := resp.Data.Cover.Full
		if imageURL == "" {
			imageURL = resp.Data.Cover.Horizontal
			if imageURL == "" {
				imageURL = resp.Data.Cover.Vertical
			}
		}

		manga := models.Manga{
			URL:      "https://zeroscans.com/swordflake/comic/" + url.PathEscape(resp.Data.Slug),
			Title:    resp.Data.Name,
			ImageURL: imageURL,
			Slug:     resp.Data.Slug,
			Synopsis: resp.Data.Summary,
			OtherID:  string(resp.Data.ID),
		}

		return manga
	}

	q := models.APIQueryData{
		URL:      mangaURL,
		Response: &apiRespStruct{},
	}

	manga, err := scrapper.GetMangaInfoAPI(ctx, q, transform)
	if err != nil {
		return manga, err
	}

	chapterListURL := "https://zeroscans.com/swordflake/comic/" + manga.OtherID + "/chapters"
	manga.Chapters, err = z.getChapters(ctx, chapterListURL, manga.Slug)

	return manga, err
}

func (z *zero) getChapters(ctx web.Context, chapterListURL, slug string) ([]models.Chapter, error) {
	type apiRespStruct struct {
		Data struct {
			Data []struct {
				ID        strFloat `json:"id"`
				Name      strFloat `json:"name"`
				CreatedAt string   `json:"created_at"`
			} `json:"data"`
			CurrentPage int `json:"current_page"`
			LastPage    int `json:"last_page"`
		} `json:"data"`
	}

	transform := func(r interface{}) []models.Chapter {
		chapters := []models.Chapter{}
		resp, ok := r.(*apiRespStruct)
		if !ok {
			return chapters
		}

		for _, item := range resp.Data.Data {
			uploadDate, err := scrapper.ParseDate(item.CreatedAt, scrapper.HUMAN_READABLE_DATE_FORMAT)
			if err != nil {
				ctx.Logger().With(zap.Error(err)).Debugf("could not parse date %s", item.CreatedAt)
				uploadDate = item.CreatedAt
			}

			chapters = append(chapters, models.Chapter{
				URL:        "https://zeroscans.com/swordflake/comic/" + slug + "/chapters/" + string(item.ID),
				Number:     string(item.Name),
				Title:      "Chapter " + string(item.Name),
				UploadDate: uploadDate,
			})
		}

		return chapters
	}

	hasNextPage := func(r interface{}) bool {
		resp, ok := r.(*apiRespStruct)
		if !ok {
			return false
		}

		return resp.Data.CurrentPage < resp.Data.LastPage
	}

	q := models.APIQueryData{
		URL:         chapterListURL,
		Response:    &apiRespStruct{},
		HasNextPage: hasNextPage,
	}

	chapters, err := scrapper.GetChapterListAPI(ctx, q, transform)

	return uniqChapters(chapters), err
}

func (z *zero) GetChapterPages(ctx web.Context, pageListURL string) ([]string, error) {
	type apiRespStruct struct {
		Data struct {
			Chapter struct {
				HighQuality []string `json:"high_quality"`
				GoodQuality []string `json:"good_quality"`
			} `json:"chapter"`
		} `json:"data"`
	}

	transform := func(r interface{}) (imageURLs []string) {
		resp, ok := r.(*apiRespStruct)
		if !ok {
			return imageURLs
		}

		if len(resp.Data.Chapter.HighQuality) >= len(resp.Data.Chapter.GoodQuality) {
			return resp.Data.Chapter.HighQuality
		}

		return resp.Data.Chapter.GoodQuality
	}

	q := models.APIQueryData{
		URL:      pageListURL,
		Response: &apiRespStruct{},
	}

	return scrapper.GetPagesListAPI(ctx, q, transform)
}

func getZeroScansConnector() models.IConnector {
	return &zero{
		Name:    "Zero Scans",
		Domain:  "zeroscans.com",
		IconURL: "https://zeroscans.com/favicon.ico",
	}
}
