package api

import (
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/connector"
	"github.com/unluckythoughts/manga-reader/models"
)

func (h *Handlers) SourceListHandler(r web.Request) (interface{}, error) {
	return connector.GetAllConnectors(), nil
}

func (h *Handlers) SourceMangaListHandler(r web.Request) (interface{}, error) {
	body := models.SourceMangaRequest{}
	err := r.GetValidatedBody(&body)
	if err != nil {
		return nil, err
	}
	return h.s.GetSourceMangaList(r.GetContext(), body.MangaListURL)
}

func (h *Handlers) SourceMangaHandler(r web.Request) (interface{}, error) {
	body := models.SourceMangaRequest{}
	err := r.GetValidatedBody(&body)
	if err != nil {
		return nil, err
	}

	return h.s.GetSourceManga(r.GetContext(), body.MangaURL)
}

func (h *Handlers) SourceMangaChapterHandler(r web.Request) (interface{}, error) {
	body := models.SourceMangaRequest{}
	err := r.GetValidatedBody(&body)
	if err != nil {
		return nil, err
	}

	return h.s.GetSourceMangaChapter(r.GetContext(), body.ChapterURL)
}
