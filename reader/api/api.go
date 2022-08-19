package api

import "github.com/unluckythoughts/manga-reader/reader/service"

type Handlers struct {
	s *service.Service
}

func New(s *service.Service) *Handlers {
	return &Handlers{
		s: s,
	}
}
