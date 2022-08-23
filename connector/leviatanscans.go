package connector

import (
	"net/url"

	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
	"github.com/unluckythoughts/manga-reader/scrapper"
)

type leviatan models.Source

func (l *leviatan) GetSource() models.Source {
	return models.Source(*l)
}

func (l *leviatan) GetMangaList(ctx web.Context) ([]models.Manga, error) {
	type apiRespStruct struct {
		Data []struct {
			ImagePath string `json:"imagePath"`
			Title     string `json:"title"`
		} `json:"data"`
		MetaData struct {
			CurrentPage int `json:"currentPage"`
			PerPage     int `json:"perPage"`
			Total       int `json:"total"`
			LastPage    int `json:"lastPage"`
		} `json:"metadata"`
	}

	nextPage := func(r interface{}) bool {
		resp, ok := r.(*apiRespStruct)
		if !ok {
			return false
		}

		return resp.MetaData.CurrentPage < resp.MetaData.LastPage
	}

	q := models.APIQueryData{
		URL: "https://en.leviatanscans.com/api/comics",
		QueryParams: map[string]string{
			"page":            "1",
			"limit":           "100",
			"bilibiliEnabled": "false",
		},
		Response:    &apiRespStruct{},
		HasNextPage: nextPage,
	}

	transform := func(r interface{}) (mangas []models.Manga) {
		resp, ok := r.(*apiRespStruct)
		if !ok {
			return mangas
		}

		for _, item := range resp.Data {
			mangas = append(mangas, models.Manga{
				Title:    item.Title,
				ImageURL: "https://en.leviatanscans.com/" + item.ImagePath,
				URL:      "https://en.leviatanscans.com/api/comics-title/" + url.PathEscape(item.Title),
			})
		}

		return mangas
	}

	return scrapper.GetMangaListAPI(ctx, q, transform)
}

func (l *leviatan) GetMangaInfo(ctx web.Context, mangaURL string) (models.Manga, error) {
	type apiRespStruct struct {
		Data struct {
			ImagePath string `json:"imagePath"`
			Title     string `json:"title"`
			Synopsis  string `json:"synopsis"`
			Chapters  []struct {
				Number strFloat `json:"number"`
				Title  string   `json:"title"`
			} `json:"chapters"`
		} `json:"data"`
	}

	transform := func(r interface{}) models.Manga {
		resp, ok := r.(*apiRespStruct)
		if !ok {
			return models.Manga{}
		}

		manga := models.Manga{
			URL:      "https://en.leviatanscans.com/api/comics-title/" + url.PathEscape(resp.Data.Title),
			Title:    resp.Data.Title,
			ImageURL: "https://en.leviatanscans.com/" + resp.Data.ImagePath,
			Synopsis: resp.Data.Synopsis,
		}

		for _, item := range resp.Data.Chapters {
			manga.Chapters = append(manga.Chapters, models.Chapter{
				URL:    "https://en.leviatanscans.com/api/chapters-title/" + item.Title,
				Number: scrapper.GetChapterNumber(string(item.Number)),
				Title:  item.Title,
			})
		}

		return manga
	}

	q := models.APIQueryData{
		URL:      mangaURL,
		Response: &apiRespStruct{},
	}

	return scrapper.GetMangaInfoAPI(ctx, q, transform)
}

func (l *leviatan) GetChapterPages(ctx web.Context, pageListURL string) ([]string, error) {
	type apiRespStruct struct {
		Data struct {
			Content []string `json:"content"`
		} `json:"data"`
	}

	transform := func(r interface{}) (imageURLs []string) {
		resp, ok := r.(*apiRespStruct)
		if !ok {
			return imageURLs
		}

		for _, imageURL := range resp.Data.Content {
			imageURLs = append(imageURLs, "https://en.leviatanscans.com/"+imageURL)
		}

		return resp.Data.Content
	}

	q := models.APIQueryData{
		URL:      pageListURL,
		Response: &apiRespStruct{},
	}

	return scrapper.GetPagesListAPI(ctx, q, transform)
}

func getLeviatanScansConnector() models.IConnector {
	return &leviatan{
		Name:    "Leviatan Scans",
		Domain:  "leviatanscans.com",
		IconURL: "https://styles.redditmedia.com/t5_2hfywp/styles/communityIcon_qdo3swk6vzl41.png",
	}
}
