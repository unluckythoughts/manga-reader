package api

import (
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
	"go.uber.org/zap"
)

func (h *Handlers) SourceListHandler(r web.Request) (interface{}, error) {
	if isMangaRequest(r) {
		return h.s.GetMangaSourceList(r.GetContext())
	} else if isNovelRequest(r) {
		return h.s.GetNovelSourceList(r.GetContext())
	}

	return nil, web.BadRequest()
}

func (h *Handlers) UpdateAllSourceMangasHandler(r web.Request) (interface{}, error) {
	if isMangaRequest(r) {
		return nil, h.s.UpdateAllSourceMangas(r.GetContext())
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
		return h.s.GetSourceMangaList(r.GetContext(), body.ID, body.Force, body.Fix)
	} else if isNovelRequest(r) {
		return h.s.GetSourceNovelList(r.GetContext(), body.ID, body.Force)
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
		return h.s.GetSourceManga(r.GetContext(), body.ID, body.Force)
	} else if isNovelRequest(r) {
		return h.s.GetSourceNovel(r.GetContext(), body.ID, body.Force)
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
		chapter, err := h.s.GetSourceMangaChapter(r.GetContext(), body.ID, body.Force)
		if err != nil {
			r.GetContext().Logger().With(zap.Error(err)).Error("error getting manga chapter")
			return chapter, web.BadRequest(err)
		}
		return chapter, nil
	} else if isNovelRequest(r) {
		text, err := h.s.GetSourceNovelChapter(r.GetContext(), body.ID)
		if err != nil {
			r.GetContext().Logger().With(zap.Error(err)).Error("error getting manga chapter")
			return text, web.BadRequest(err)
		}
		return text, nil
	}

	return nil, web.BadRequest()
}
