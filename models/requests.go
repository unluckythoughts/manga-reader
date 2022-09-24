package models

type SourceRequest struct {
	URL   string `json:"url"`
	Force bool   `json:"force"`
}

type SourceChapterRequest struct {
	ChapterURL string `json:"chapterUrl"`
}

type SourceListRequest struct {
	Domain string `json:"domain"`
	Force  bool   `json:"force"`
}

type SearchSourceRequest struct {
	Query string `json:"query"`
}
