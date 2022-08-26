package connector

import (
	"net/url"
	"strconv"

	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
	"github.com/unluckythoughts/manga-reader/scrapper"
)

type leviatan models.Connector

func GetLeviatanScansConnector() models.IConnector {
	return &leviatan{
		Source: models.Source{
			Name:    "Leviatan Scans",
			Domain:  "leviatanscans.com",
			IconURL: "https://styles.redditmedia.com/t5_2hfywp/styles/communityIcon_qdo3swk6vzl41.png",
		},
		BaseURL:       "https://en.leviatanscans.com/api/",
		MangaListPath: "comics/",
	}
}

func (l *leviatan) GetSource() models.Source {
	return l.Source
}

func (l *leviatan) GetMangaList(ctx web.Context) ([]models.Manga, error) {
	type apiRespBody struct {
		Data []struct {
			ImagePath string `json:"imagePath" manga-reader:"manga.image-url"`
			Title     string `json:"title" manga-reader:"manga.title"`
		} `json:"data" manga-reader:"manga-list"`
		MetaData struct {
			CurrentPage int `json:"currentPage" manga-reader:"manga-list.current-page"`
			PerPage     int `json:"perPage"`
			Total       int `json:"total"`
			LastPage    int `json:"lastPage" manga-reader:"manga-list.last-page"`
		} `json:"metadata"`
	}

	params := map[string]string{
		"page":            "1",
		"limit":           "100",
		"bilibiliEnabled": "false",
	}

	mangas := []models.Manga{}

	for {
		apiResp := apiRespBody{}
		q := models.APIQueryData{
			URL:         l.BaseURL + l.MangaListPath,
			QueryParams: params,
			Response:    &apiResp,
		}
		err := scrapper.GetAPIResponse(ctx, q)
		if err != nil {
			return []models.Manga{}, err
		}

		pageMangas, err := scrapper.GetMangaListFromTags(apiResp)
		if err != nil {
			return mangas, err
		}

		mangas = append(mangas, pageMangas...)

		currentPage, lastPage, _ := scrapper.GetMangaListPageData(apiResp)
		if currentPage >= lastPage {
			break
		}

		params["page"] = strconv.Itoa(currentPage + 1)
	}

	for i, m := range mangas {
		mangas[i].ImageURL = "https://en.leviatanscans.com/" + m.ImageURL
		mangas[i].URL = l.BaseURL + "comics-title/" + url.PathEscape(m.Title)
	}

	return mangas, nil
}

func (l *leviatan) GetMangaInfo(ctx web.Context, mangaURL string) (models.Manga, error) {
	type apiResonseBody struct {
		Data struct {
			ImagePath string `json:"imagePath" manga-reader:"manga.image-url"`
			Title     string `json:"title" manga-reader:"manga.title"`
			Synopsis  string `json:"synopsis" manga-reader:"manga.synopsis"`
			Chapters  []struct {
				Number models.StrFloat `json:"number" manga-reader:"manga.chapter.number"`
				Title  string          `json:"title" manga-reader:"manga.chapter.title"`
			} `json:"chapters" manga-reader:"manga.chapter-list"`
		} `json:"data"`
	}

	apiResp := apiResonseBody{}
	q := models.APIQueryData{
		URL:      mangaURL,
		Response: &apiResp,
	}

	err := scrapper.GetAPIResponse(ctx, q)
	if err != nil {
		return models.Manga{}, err
	}

	manga, err := scrapper.GetMangaInfoFromTags(apiResp)
	if err != nil {
		return manga, err
	}

	manga.URL = l.BaseURL + "comics-title/" + url.PathEscape(manga.Title)
	manga.ImageURL = "https://en.leviatanscans.com/" + manga.ImageURL

	for i, c := range manga.Chapters {
		manga.Chapters[i].URL = l.BaseURL + "chapters-title/" + c.Title
		manga.Chapters[i].Number = scrapper.GetChapterNumber(string(c.Number))
	}

	return manga, nil
}

func (l *leviatan) GetChapterPages(ctx web.Context, chapterURL string) (models.Pages, error) {
	type apiResonseBody struct {
		Data struct {
			Content []string `json:"content" manga-reader:"manga.chapter.pages-list"`
		} `json:"data"`
	}

	apiResp := apiResonseBody{}
	q := models.APIQueryData{
		URL:      chapterURL,
		Response: &apiResp,
	}

	err := scrapper.GetAPIResponse(ctx, q)
	if err != nil {
		return models.Pages{}, err
	}

	pages, err := scrapper.GetChapterPagesFromTags(apiResp)
	if err != nil {
		return pages, err
	}

	for i, url := range pages.URLs {
		pages.URLs[i] = "https://en.leviatanscans.com/" + url
	}

	return pages, nil
}
