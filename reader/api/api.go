package api

import (
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/reader/service"
)

type Handlers struct {
	s *service.Service
}

func New(s *service.Service) *Handlers {
	return &Handlers{
		s: s,
	}
}

func isMangaRequest(r web.Request) bool {
	if r.GetRouteParam("type") == "manga" {
		return true
	}

	return false
}

func isNovelRequest(r web.Request) bool {
	if r.GetRouteParam("type") == "novel" {
		return true
	}

	return false
}
