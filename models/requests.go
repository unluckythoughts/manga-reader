package models

type SourceMangaRequest struct {
	MangaListURL string `json:"mangaListUrl"`
	MangaURL     string `json:"mangaUrl"`
	ChapterURL   string `json:"chapterUrl"`
}

type SourceManagaListRequest struct {
	Domain string `json:"domain"`
}
