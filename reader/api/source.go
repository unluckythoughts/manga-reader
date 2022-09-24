package api

import (
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
	"go.uber.org/zap"
)

func (h *Handlers) SourceListHandler(r web.Request) (interface{}, error) {
	if isMangaRequest(r) {
		return h.s.GetSourceList(r.GetContext())
	} else if isNovelRequest(r) {
		return h.s.GetNovelSourceList(r.GetContext())
	}

	return nil, web.BadRequest()
}

func (h *Handlers) SourceItemListHandler(r web.Request) (interface{}, error) {
	body := models.SourceListRequest{}
	err := r.GetValidatedBody(&body)
	if err != nil {
		return nil, web.BadRequest(err)
	}

	if isMangaRequest(r) {
		return h.s.GetSourceMangaList(r.GetContext(), body.Domain, body.Force)
	} else if isNovelRequest(r) {
		return h.s.GetSourceNovelList(r.GetContext(), body.Domain, body.Force)
	}

	return nil, web.BadRequest()
}

func (h *Handlers) SourceItemSearchHandler(r web.Request) (interface{}, error) {
	body := models.SearchSourceRequest{}
	err := r.GetValidatedBody(&body)
	if err != nil {
		return nil, web.BadRequest(err)
	}

	if isMangaRequest(r) {
		return h.s.SearchSourceManga(r.GetContext(), body.Query)
	} else if isNovelRequest(r) {
		return h.s.SearchSourceNovel(r.GetContext(), body.Query)
	}

	return nil, web.BadRequest()
}

func (h *Handlers) SourceItemHandler(r web.Request) (interface{}, error) {
	body := models.SourceRequest{}
	err := r.GetValidatedBody(&body)
	if err != nil {
		return nil, web.BadRequest(err)
	}

	if isMangaRequest(r) {
		return h.s.GetSourceManga(r.GetContext(), body.URL, body.Force)
	} else if isNovelRequest(r) {
		return h.s.GetSourceNovel(r.GetContext(), body.URL, body.Force)
	}

	return nil, web.BadRequest()
}

func (h *Handlers) SourceItemChapterHandler(r web.Request) (interface{}, error) {
	body := models.SourceChapterRequest{}
	err := r.GetValidatedBody(&body)
	if err != nil {
		return nil, web.BadRequest(err)
	}

	if isMangaRequest(r) {
		pages, err := h.s.GetSourceMangaChapter(r.GetContext(), body.ChapterURL)
		if err != nil {
			r.GetContext().Logger().With(zap.Error(err)).Error("error getting manga chapter")
			return pages, web.BadRequest(err)
		}
		return pages, nil
	} else if isNovelRequest(r) {
		text, err := h.s.GetSourceNovelChapter(r.GetContext(), body.ChapterURL)
		if err != nil {
			r.GetContext().Logger().With(zap.Error(err)).Error("error getting manga chapter")
			return text, web.BadRequest(err)
		}
		return text, nil
	}

	return nil, web.BadRequest()
}
