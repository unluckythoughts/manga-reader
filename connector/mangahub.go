package connector

import (
	"encoding/json"
	"errors"
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

type mangahub models.MangaConnector

func GetMangaHubConnector() models.IMangaConnector {
	return &mangahub{
		Source: models.Source{
			Name:    "Manga Hub",
			Domain:  "mangahub.io",
			IconURL: "https://mangahub.io/logo-small.png",
		},
		Transport: cloudflarebp.AddCloudFlareByPass((&http.Client{}).Transport),
		BaseURL:   "https://api.mghubcdn.com/graphql",
	}
}

func (m *mangahub) GetSource() models.Source {
	return m.Source
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

func (m *mangahub) _getResponse(ctx web.Context, r interface{}) (*mangahubAPIResponseBody, bool) {
	resp, ok := r.(*mangahubAPIResponseBody)
	if !ok {
		return resp, false
	}

	if len(resp.Errors) > 0 {
		ctx.Logger().With("error", resp.Errors[0].Message).Debug("resetting accessToken")
		accessToken = ""
		return resp, false
	}

	return resp, true
}

func (m *mangahub) _getRequestHeaders(ctx web.Context) http.Header {
	if accessToken == "" {
		accessToken = m._fetchAccessKey(ctx)
		ctx.Logger().With("token", accessToken).Debug("setting access token")
	}

	headers := http.Header{}
	headers.Add("origin", "https://mangahub.io")
	headers.Add("referer", "https://mangahub.io/")
	headers.Add("content-type", "application/json")
	headers.Add("x-mhub-access", accessToken)

	return headers
}

type mangahubAPIResponseBody struct {
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors" manga-reader:"error-list"`
	Data struct {
		Search struct {
			Rows []struct {
				ID          int    `json:"id" manga-reader:"manga.id"`
				Title       string `json:"title" manga-reader:"manga.title"`
				Slug        string `json:"slug" manga-reader:"manga.slug"`
				Image       string `json:"image" manga-reader:"manga.image-url"`
				Description string `json:"description" manga-reader:"manga.synopsis"`
			} `json:"rows" manga-reader:"manga-list"`
		} `json:"search"`
		Manga struct {
			ID          int    `json:"id" manga-reader:"manga.id"`
			Title       string `json:"title" manga-reader:"manga.title"`
			Slug        string `json:"slug" manga-reader:"manga.slug"`
			Image       string `json:"image" manga-reader:"manga.image-url"`
			Description string `json:"description" manga-reader:"manga.synopsis"`
			Chapters    []struct {
				ID     int     `json:"id" manga-reader:"manga.chapter.id"`
				Title  string  `json:"title" manga-reader:"manga.chapter.title"`
				Number float64 `json:"number" manga-reader:"manga.chapter.number"`
				Date   string  `json:"date" manga-reader:"manga.chapter.upload-date|2006-01-02T15:04:05.000Z"`
			} `json:"chapters"  manga-reader:"manga.chapter-list"`
		} `json:"manga"`
		Chapter struct {
			Pages string `json:"pages" manga-reader:"manga.chapter.page"`
		} `json:"chapter"`
	} `json:"data"`
}

func (m *mangahub) GetMangaList(ctx web.Context) ([]models.Manga, error) {
	apiResp := mangahubAPIResponseBody{}
	q := scrapper.APIQueryData{
		URL:       m.BaseURL,
		Method:    http.MethodPost,
		Body:      m._getRequestBody(),
		Response:  &apiResp,
		Headers:   m._getRequestHeaders(ctx),
		Transport: m.Transport,
	}

	err := scrapper.GetAPIResponse(ctx, q)
	if err != nil {
		return []models.Manga{}, err
	}

	errorMessages := scrapper.GetErrorsFromTags(apiResp)
	if len(errorMessages) > 0 {
		ctx.Logger().With("error", errorMessages[0]).Debug("resetting accessToken")
		accessToken = ""
		return []models.Manga{}, errors.New(errorMessages[0])
	}

	mangas, err := scrapper.GetMangaListFromTags(apiResp)
	if err != nil {
		return mangas, err
	}

	for i, m := range mangas {
		mangas[i].ImageURL = "https://thumb.mghubcdn.com/" + m.ImageURL
		mangas[i].URL = "https://mangahub.io/manga/" + m.URL
	}

	return mangas, nil
}

func (m *mangahub) GetMangaInfo(ctx web.Context, mangaURL string) (models.Manga, error) {
	slug := strings.Replace(mangaURL, "https://mangahub.io/manga/", "", -1)

	apiResp := mangahubAPIResponseBody{}
	q := scrapper.APIQueryData{
		URL:       m.BaseURL,
		Method:    http.MethodPost,
		Body:      m._getRequestBody(slug),
		Response:  &apiResp,
		Headers:   m._getRequestHeaders(ctx),
		Transport: m.Transport,
	}

	err := scrapper.GetAPIResponse(ctx, q)
	if err != nil {
		return models.Manga{}, err
	}

	manga, err := scrapper.GetMangaInfoFromTags(apiResp)
	if err != nil {
		return manga, err
	}

	manga.URL = "https://mangahub.io/manga/" + manga.Slug
	manga.ImageURL = "https://thumb.mghubcdn.com/" + manga.ImageURL

	for i, c := range manga.Chapters {
		manga.Chapters[i].URL = "https://mangahub.io/chapter/" + manga.Slug + "/chapter-" + c.Number
	}

	return manga, nil
}

func (m *mangahub) GetChapterPages(ctx web.Context, chapterURL string) (models.Pages, error) {
	chapterURL = strings.Replace(chapterURL, "https://mangahub.io/chapter/", "", -1)
	slugs := strings.Split(strings.Replace(chapterURL, "/chapter-", ":", -1), ":")

	apiResp := mangahubAPIResponseBody{}
	q := scrapper.APIQueryData{
		URL:       m.BaseURL,
		Method:    http.MethodPost,
		Body:      m._getRequestBody(slugs...),
		Response:  &apiResp,
		Headers:   m._getRequestHeaders(ctx),
		Transport: m.Transport,
	}

	err := scrapper.GetAPIResponse(ctx, q)
	if err != nil {
		return models.Pages{}, err
	}

	// pages, err := scrapper.GetChapterPagesFromTags(apiResp)
	// if err != nil {
	// 	return pages, err
	// }

	strPages := map[string]string{}
	err = json.Unmarshal([]byte(apiResp.Data.Chapter.Pages), &strPages)
	if err != nil {
		return models.Pages{}, err
	}

	imageURLs := []string{}
	// ordering the image urls map based on index
	for i := 0; i < len(strPages); i++ {
		index := strconv.Itoa(i + 1)
		if len(strPages[index]) > 0 {
			imageURLs = append(imageURLs, "https://img.mghubcdn.com/file/imghub/"+strPages[index])
		}
	}

	if len(imageURLs) > 0 {
		return models.Pages{URLs: scrapper.GetImagesAsDataUrls(ctx, imageURLs)}, nil
	}

	return models.Pages{URLs: imageURLs}, nil
}
