package models

type SourceRequest struct {
	ID    uint `json:"id"`
	Force bool `json:"force"`
}

type SourceChapterRequest struct {
	ID    uint `json:"id"`
	Force bool `json:"force"`
}

type SourceListRequest struct {
	ID    uint `json:"id"`
	Force bool `json:"force"`
	Fix   bool `json:"fix"`
}

type SearchSourceRequest struct {
	Query string `json:"query"`
}
