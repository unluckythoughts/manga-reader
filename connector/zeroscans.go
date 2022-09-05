package connector

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
	"github.com/unluckythoughts/manga-reader/scrapper"
)

type zero models.Connector

type zeroAPIResponseBody struct {
	Data struct {
		Comics []struct {
			ID      models.StrFloat `json:"id" manga-reader:"manga.id"`
			Slug    string          `json:"Slug" manga-reader:"manga.slug"`
			Name    string          `json:"name" manga-reader:"manga.title"`
			Summary string          `json:"summary" manga-reader:"manga.synopsis"`
			Cover   struct {
				Horizontal string `json:"horizontal" manga-reader:"manga.image-url[2]"`
				Vertical   string `json:"vertical" manga-reader:"manga.image-url[3]"`
				Full       string `json:"full" manga-reader:"manga.image-url[1]"`
			} `json:"cover"`
		} `json:"comics" manga-reader:"manga-list"`
		ID      models.StrFloat `json:"id" manga-reader:"manga.id"`
		Name    string          `json:"name" manga-reader:"manga.title"`
		Slug    string          `json:"slug" manga-reader:"manga.slug"`
		Summary string          `json:"summary" manga-reader:"manga.synopsis"`
		Cover   struct {
			Horizontal string `json:"horizontal" manga-reader:"manga.image-url[2]"`
			Vertical   string `json:"vertical" manga-reader:"manga.image-url[3]"`
			Full       string `json:"full" manga-reader:"manga.image-url[1]"`
		} `json:"cover"`

		Data []struct {
			ID        models.StrFloat `json:"id" manga-reader:"manga.chapter.id"`
			Name      models.StrFloat `json:"name" manga-reader:"manga.chapter.title"`
			CreatedAt string          `json:"created_at" manga-reader:"manga.chapter.upload-date|"`
		} `json:"data" manga-reader:"manga.chapter-list"`
		CurrentPage int    `json:"current_page" manga-reader:"manga.chapter-list.current-page"`
		LastPage    int    `json:"last_page" manga-reader:"manga.chapter-list.last-page"`
		NextPageURL string `json:"next_page_url" manga-reader:"manga.chapter-list.next-page-url"`

		Chapter struct {
			HighQuality []string `json:"high_quality" manga-reader:"manga.chapter.pages-list[1]"`
			GoodQuality []string `json:"good_quality" manga-reader:"manga.chapter.pages-list[2]"`
		} `json:"chapter"`
	} `json:"data"`
}

func GetZeroScansConnector() models.IConnector {
	return &zero{
		Source: models.Source{
			Name:    "Zero Scans",
			Domain:  "zeroscans.com",
			IconURL: "https://zeroscans.com/favicon.ico",
		},
		BaseURL:       "https://zeroscans.com/swordflake/",
		MangaListPath: "comics/",
	}
}

func (z *zero) GetSource() models.Source {
	return z.Source
}

func (z *zero) GetMangaList(ctx web.Context) ([]models.Manga, error) {
	apiResp := zeroAPIResponseBody{}
	q := models.APIQueryData{
		URL:      z.BaseURL + z.MangaListPath,
		Response: &apiResp,
	}

	err := scrapper.GetAPIResponse(ctx, q)
	if err != nil {
		return []models.Manga{}, err
	}

	mangas, err := scrapper.GetMangaListFromTags(apiResp)
	if err != nil {
		return mangas, err
	}

	for i, m := range mangas {
		mangas[i].URL = z.BaseURL + "comic/" + m.Slug
	}

	return mangas, nil
}

func (z *zero) GetMangaInfo(ctx web.Context, mangaURL string) (models.Manga, error) {
	apiResp := zeroAPIResponseBody{}
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
	manga.URL = z.BaseURL + "comic/" + url.PathEscape(manga.Slug)

	chaptersURL := z.BaseURL + "comic/" + manga.OtherID + "/chapters"
	params := map[string]string{
		"page": "1",
	}
	for {
		apiResp := zeroAPIResponseBody{}
		q := models.APIQueryData{
			URL:         chaptersURL,
			Response:    &apiResp,
			QueryParams: params,
		}

		err := scrapper.GetAPIResponse(ctx, q)
		if err != nil {
			return models.Manga{}, err
		}

		mangaChapters, err := scrapper.GetMangaInfoFromTags(apiResp)
		if err != nil {
			return manga, err
		}

		manga.Chapters = append(manga.Chapters, mangaChapters.Chapters...)

		currentPage, lastPage, _ := scrapper.GetChapterListPageData(apiResp)
		if currentPage == lastPage {
			break
		}

		q.QueryParams["page"] = strconv.Itoa(currentPage + 1)
	}

	for i, c := range manga.Chapters {
		manga.Chapters[i].URL = z.BaseURL + "comic/" + manga.Slug + "/chapters/" + c.OtherID
		manga.Chapters[i].Title = "Chapter " + string(c.Title)
		manga.Chapters[i].Number = c.Title
	}

	return manga, err
}

func (z *zero) GetChapterPages(ctx web.Context, pageListURL string) (models.Pages, error) {
	headers := http.Header{}
	headers.Add("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")

	apiResp := zeroAPIResponseBody{}
	q := models.APIQueryData{
		URL:      pageListURL,
		Response: &apiResp,
		Headers:  headers,
	}

	err := scrapper.GetAPIResponse(ctx, q)
	if err != nil {
		return models.Pages{}, err
	}

	return scrapper.GetChapterPagesFromTags(apiResp)
}
