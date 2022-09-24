package models

import "net/http"

type NovelList struct {
	NovelContainer string
	NovelTitle     string
	NovelURL       string
	NovelImageURL  string
	NovelSlug      string
	NovelOtherID   string
	NextPage       string
	LastPage       string
	PageParam      string
}

type NovelInfo struct {
	Title                   string
	ImageURL                string
	Synopsis                string
	Slug                    string
	OtherID                 string
	ChapterListURL          string
	ChapterListNextPage     string
	ChapterContainer        string
	ChapterNumber           string
	ChapterTitle            string
	ChapterURL              string
	ChapterUploadDate       string
	ChapterUploadDateFormat string
}

type NovelChapterTextSelectors struct {
	Paragraph string
}

type NovelSelectors struct {
	List    NovelList
	Info    NovelInfo
	Chapter NovelChapterTextSelectors
}

type NovelConnector struct {
	BaseURL       string
	Transport     http.RoundTripper
	NovelListPath string
	Source
	NovelSelectors
}
