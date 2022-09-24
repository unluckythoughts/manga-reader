package connector

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
)

func GetAllNovelConnectors() []models.INovelConnector {
	return append(
		[]models.INovelConnector{},
		GetLightNovelPubConnector(),
	)
}

func findNovelConnector(domain string) (models.INovelConnector, bool) {
	domain = getDomain(domain)
	conns := GetAllNovelConnectors()
	for i := 0; i < len(conns); i++ {
		if strings.Contains(domain, conns[i].GetSource().Domain) {
			return conns[i], true
		}
	}

	return nil, false
}

func NewNovelConnector(ctx web.Context, domain string) (models.INovelConnector, error) {
	conn, ok := findNovelConnector(domain)
	if !ok {
		return nil, errors.Errorf("could not find config for %s", domain)
	}

	return conn, nil
}
