package connector

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	cloudflarebp "github.com/DaRealFreak/cloudflare-bp-go"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
	"github.com/unluckythoughts/manga-reader/scrapper"
)

var (
	accessToken = ""
)

type mangahub models.Source

func (m *mangahub) GetSource() models.Source {
	return models.Source(*m)
}

func (m *mangahub) _fetchAccessKey(ctx web.Context) string {
	url := "https://mangahub.io/?reloadKey=1"

	headers := scrapper.GetHeaders(ctx, url, m.Transport)
	headers.Add("cookie", headers.Get("set-cookie"))
	req := http.Request{Header: headers}
	cookies := req.Cookies()

	for _, c := range cookies {
		if c.Name == "mhub_access" {
			return c.Value
		}
	}

	return ""
}

func (m *mangahub) _getRequestBody(slugs ...string) []byte {
	query := "{search(x: m01,q: \"\",genre: \"all\",mod:ALPHABET,limit: 99999){rows{id,slug,title,image}}}"
	if len(slugs) == 1 {
		query = "{manga(x:m01,slug:\"" + slugs[0] + "\"){id,slug,title,image,description,chapters{id,number,title,slug,date}}}"
	} else if len(slugs) > 1 {
		query = "{chapter(x:m01,slug:\"" + slugs[0] + "\",number:" + slugs[1] + "){pages}}"
	}

	reqBody := map[string]string{
		"query": query,
	}

	data, _ := json.Marshal(reqBody)
	return data
}

func (m *mangahub) _getRequestHeaders(ctx web.Context) http.Header {
	if accessToken == "" {
		accessToken = m._fetchAccessKey(ctx)
	}

	headers := http.Header{}
	headers.Add("origin", "https://mangahub.io")
	headers.Add("referer", "https://mangahub.io/")
	headers.Add("content-type", "application/json")
	headers.Add("x-mhub-access", accessToken)

	return headers
}

type apiResp struct {
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
	Data struct {
		Search struct {
			Rows []struct {
				ID    int    `json:"id"`
				Title string `json:"title"`
				Slug  string `json:"slug"`
				Image string `json:"image"`
			} `json:"rows"`
		} `json:"search"`
		Manga struct {
			ID          int    `json:"id"`
			Title       string `json:"title"`
			Slug        string `json:"slug"`
			Image       string `json:"image"`
			Description string `json:"description"`
			Chapters    []struct {
				ID     int    `json:"id"`
				Title  string `json:"title"`
				Number int    `json:"number"`
				Date   string `json:"date"`
			} `json:"chapters"`
		} `json:"manga"`
		Chapter struct {
			Pages string `json:"pages"`
		} `json:"chapter"`
	} `json:"data"`
}

func (m *mangahub) GetMangaList(ctx web.Context) ([]models.Manga, error) {
	q := models.APIQueryData{
		URL:       "https://api.mghubcdn.com/graphql",
		Method:    http.MethodPost,
		Body:      m._getRequestBody(),
		Response:  &apiResp{},
		Headers:   m._getRequestHeaders(ctx),
		Transport: m.Transport,
	}

	transform := func(r interface{}) (mangas []models.Manga) {
		resp, ok := r.(*apiResp)
		if !ok {
			return []models.Manga{}
		}

		for _, row := range resp.Data.Search.Rows {
			mangas = append(mangas, models.Manga{
				OtherID:  strconv.Itoa(row.ID),
				ImageURL: "https://thumb.mghubcdn.com/" + row.Image,
				Slug:     row.Slug,
				Title:    row.Title,
				URL:      "https://mangahub.io/manga/" + row.Slug,
			})
		}

		return mangas
	}

	return scrapper.GetMangaListAPI(ctx, q, transform)
}

func (m *mangahub) GetMangaInfo(ctx web.Context, mangaURL string) (models.Manga, error) {
	slug := strings.Replace(mangaURL, "https://mangahub.io/manga/", "", -1)

	q := models.APIQueryData{
		URL:       "https://api.mghubcdn.com/graphql",
		Method:    http.MethodPost,
		Body:      m._getRequestBody(slug),
		Response:  &apiResp{},
		Headers:   m._getRequestHeaders(ctx),
		Transport: m.Transport,
	}

	transform := func(r interface{}) models.Manga {
		resp, ok := r.(*apiResp)
		if !ok {
			return models.Manga{}
		}

		manga := models.Manga{
			URL:      "https://mangahub.io/manga/" + resp.Data.Manga.Slug,
			Title:    resp.Data.Manga.Title,
			ImageURL: "https://thumb.mghubcdn.com/" + resp.Data.Manga.Image,
			Synopsis: resp.Data.Manga.Description,
		}

		for _, item := range resp.Data.Manga.Chapters {
			uploadDate, _ := scrapper.ParseDate(item.Date, "2006-01-02T15:04:05.000Z")
			manga.Chapters = append(manga.Chapters, models.Chapter{
				URL:        "https://mangahub.io/chapter/" + manga.Slug + "/chapter-" + strconv.Itoa(item.Number),
				Number:     strconv.Itoa(item.Number),
				Title:      item.Title,
				UploadDate: uploadDate,
			})
		}

		return manga
	}

	return scrapper.GetMangaInfoAPI(ctx, q, transform)
}

func (m *mangahub) GetChapterPages(ctx web.Context, pageListURL string) ([]string, error) {
	pageListURL = strings.Replace(pageListURL, "https://mangahub.io/chapter/", "", -1)
	slugs := strings.Split(strings.Replace(pageListURL, "/chapter-", ":", -1), ":")

	q := models.APIQueryData{
		URL:       "https://api.mghubcdn.com/graphql",
		Method:    http.MethodPost,
		Body:      m._getRequestBody(slugs...),
		Response:  &apiResp{},
		Headers:   m._getRequestHeaders(ctx),
		Transport: m.Transport,
	}

	transform := func(r interface{}) (imageURLs []string) {
		resp, ok := r.(*apiResp)
		if !ok {
			return imageURLs
		}

		pages := map[string]string{}
		err := json.Unmarshal([]byte(resp.Data.Chapter.Pages), &pages)
		if err != nil {
			return []string{}
		}

		for i := 0; i < len(pages); i++ {
			index := strconv.Itoa(i + 1)
			imageURLs = append(imageURLs, "https://thumb.mghubcdn.com/"+pages[index])
		}

		return imageURLs
	}

	return scrapper.GetPagesListAPI(ctx, q, transform)
}

func getMangaHubConnector() models.IConnector {
	return &mangahub{
		Name:      "Manga Hub",
		Domain:    "mangahub.io",
		IconURL:   "https://mangahub.io/logo-small.png",
		Transport: cloudflarebp.AddCloudFlareByPass((&http.Client{}).Transport),
	}
}
