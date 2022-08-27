package connector

import (
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
)

func GetAllConnectors() []models.IConnector {
	return append(
		[]models.IConnector{},
		GetAsuraScansConnector(),
		GetLeviatanScansConnector(),
		GetZeroScansConnector(),
		GetRealmScansConnector(),
		GetReaperScansConnector(),
		// GetMangaHubConnector(),
		GetMangaHasuConnector(),
		GetFlameScansConnector(),
	)
}

func getDomain(link string) string {
	linkURL, err := url.Parse(link)
	if err != nil || linkURL.Hostname() == "" {
		linkURL, err = url.Parse("http://" + link)
		if err != nil || linkURL.Hostname() == "" {
			return link
		}
	}

	return linkURL.Hostname()
}

func findConnector(domain string) (models.IConnector, bool) {
	domain = getDomain(domain)
	conns := GetAllConnectors()
	for i := 0; i < len(conns); i++ {
		if strings.Contains(domain, conns[i].GetSource().Domain) {
			return conns[i], true
		}
	}

	return nil, false
}

func New(ctx web.Context, domain string) (models.IConnector, error) {
	conn, ok := findConnector(domain)
	if !ok {
		return nil, errors.Errorf("could not find config for %s", domain)
	}

	return conn, nil
}
