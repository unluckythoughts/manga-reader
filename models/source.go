package models

import "net/http"

type MangaListSelectors struct {
	URL                   string
	MangaTitleSelector    string
	MangaImageURLSelector string
	MangaURLSelector      string
	NextPageSelector      string
}

type MangaInfoSelectors struct {
	TitleSelector             string
	ImageURLSelector          string
	SynopsisSelector          string
	ChapterNumberSelector     string
	ChapterTitleSelector      string
	ChapterURLSelector        string
	ChapterUploadDateSelector string
}

type ChapterInfoSelectors struct {
	ImageURLsSelector string
}

type Source struct {
	Name         string
	IconURL      string
	RoundTripper http.RoundTripper
	MangaList    MangaListSelectors
	MangaInfo    MangaInfoSelectors
	ChapterInfo  ChapterInfoSelectors
}
