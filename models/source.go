package models

import (
	"net/http"

	"github.com/unluckythoughts/go-microservice/tools/web"
)

type IConnector interface {
	Get() Source
	GetDomain() string
	GetName() string
	GetIconURL() string
	GetMangaList(ctx web.Context) ([]Manga, error)
	GetMangaInfo(ctx web.Context, mangaURL string) (Manga, error)
	GetChapterPages(ctx web.Context, pageListURL string) ([]string, error)
}

type MangaListSelectors struct {
	URL                   string
	MangaTitleSelector    string
	MangaImageURLSelector string
	MangaURLSelector      string
	NextPageSelector      string
}

type MangaInfoSelectors struct {
	URL                       string
	TitleSelector             string
	ImageURLSelector          string
	SynopsisSelector          string
	ChapterNumberSelector     string
	ChapterTitleSelector      string
	ChapterURLSelector        string
	ChapterUploadDateSelector string
}

type ChapterInfoSelectors struct {
	URL          string
	PageSelector string
}

type APIQueryData struct {
	URL         string
	PageParam   string
	QueryParams map[string]string
	Response    interface{}
	HasNextPage HasNextPage
}

type HasNextPage func(resp interface{}) bool
type MangaListTransform func(interface{}) []Manga
type MangaInfoTransform func(interface{}) Manga
type ChapterListTransform func(interface{}) []Chapter
type PagesListTransform func(interface{}) []string

type Source struct {
	Name      string            `json:"name"`
	Domain    string            `json:"domain"`
	IconURL   string            `json:"iconUrl"`
	Transport http.RoundTripper `json:"-"`
}
