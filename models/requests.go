package models

type SourceMangaRequest struct {
	MangaURL string `json:"mangaUrl"`
	Force    bool   `json:"force"`
}

type SourceChapterRequest struct {
	ChapterURL string `json:"chapterUrl"`
}

type SourceManagaListRequest struct {
	Domain string `json:"domain"`
	Force  bool   `json:"force"`
}

type SearchSourceManagaRequest struct {
	Query string `json:"query"`
}
