package api

import (
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
)

func (h *Handlers) SourceListHandler(r web.Request) (interface{}, error) {
	return h.s.GetSourceList(r.GetContext())
}

func (h *Handlers) SourceMangaListHandler(r web.Request) (interface{}, error) {
	body := models.SourceManagaListRequest{}
	err := r.GetValidatedBody(&body)
	if err != nil {
		return nil, err
	}
	return h.s.GetSourceMangaList(r.GetContext(), body.Domain, body.Force)
}

func (h *Handlers) SourceMangaHandler(r web.Request) (interface{}, error) {
	body := models.SourceMangaRequest{}
	err := r.GetValidatedBody(&body)
	if err != nil {
		return nil, err
	}

	return h.s.GetSourceManga(r.GetContext(), body.MangaURL, body.Force)
}

func (h *Handlers) SourceMangaChapterHandler(r web.Request) (interface{}, error) {
	body := models.SourceChapterRequest{}
	err := r.GetValidatedBody(&body)
	if err != nil {
		return nil, err
	}

	return h.s.GetSourceMangaChapter(r.GetContext(), body.ChapterURL)
}
