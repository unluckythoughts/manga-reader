package models

import (
	"net/http"

	"github.com/unluckythoughts/go-microservice/tools/web"
)

type IMangaConnector interface {
	GetSource() Source
	GetMangaList(ctx web.Context) ([]Manga, error)
	GetMangaInfo(ctx web.Context, mangaURL string) (Manga, error)
	GetChapterPages(ctx web.Context, pageListURL string) (Pages, error)
}

type Pages struct {
	Config map[string]interface{} `json:"config"`
	URLs   []string               `json:"urls"`
}

type MangaList struct {
	MangaContainer string
	MangaTitle     string
	MangaURL       string
	MangaImageURL  string
	MangaSlug      string
	MangaOtherID   string
	NextPage       string
	LastPage       string
	PageParam      string
}

type MangaInfo struct {
	Title                   string
	ImageURL                string
	Synopsis                string
	Slug                    string
	OtherID                 string
	ChapterContainer        string
	ChapterNumber           string
	ChapterTitle            string
	ChapterURL              string
	ChapterUploadDate       string
	ChapterUploadDateFormat string
}

type PageSelectors struct {
	ImageUrl string
}

type MangaSelectors struct {
	List    MangaList
	Info    MangaInfo
	Chapter PageSelectors
}

type MangaConnector struct {
	BaseURL       string
	Transport     http.RoundTripper
	MangaListPath string
	Source
	MangaSelectors
}
