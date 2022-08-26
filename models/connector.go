package models

import "net/http"

type MangaList struct {
	MangaContainer string
	MangaTitle     string
	MangaURL       string
	MangaImageURL  string
	MangaSlug      string
	MangaOtherID   string
	NextPage       string
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

type Selectors struct {
	List    MangaList
	Info    MangaInfo
	Chapter PageSelectors
}

type Connector struct {
	BaseURL       string
	Transport     http.RoundTripper
	MangaListPath string
	Source
	Selectors
}
