package models

import (
	"net/http"

	"github.com/unluckythoughts/go-microservice/tools/web"
)

type Pages struct {
	Config map[string]interface{} `json:"config"`
	URLs   []string               `json:"urls"`
}

type IConnector interface {
	GetSource() Source
	GetMangaList(ctx web.Context) ([]Manga, error)
	GetMangaInfo(ctx web.Context, mangaURL string) (Manga, error)
	GetChapterPages(ctx web.Context, pageListURL string) (Pages, error)
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
	ChapterUploadDateFormat   string
}

type ChapterInfoSelectors struct {
	URL          string
	PageSelector string
}

type APIQueryData struct {
	URL         string
	Method      string
	Body        []byte
	Headers     http.Header
	PageParam   string
	QueryParams map[string]string
	Response    interface{}
	HasNextPage HasNextPage
	Transport   http.RoundTripper
}

type HasNextPage func(resp interface{}) bool
type MangaListTransform func(interface{}) []Manga
type MangaInfoTransform func(interface{}) Manga
type ChapterListTransform func(interface{}) []Chapter
type PagesListTransform func(interface{}) []string

type Source struct {
	ID        int               `json:"id" gorm:"column:id;primarykey"`
	Name      string            `json:"name" gorm:"column:name"`
	Domain    string            `json:"domain" gorm:"column:domain"`
	IconURL   string            `json:"iconUrl" gorm:"column:icon_url"`
	UpdatedAt string            `json:"updatedAt" gorm:"column:updated_at"`
	Transport http.RoundTripper `json:"-"  gorm:"-"`
}

func (m Source) TableName() string {
	return "source"
}
