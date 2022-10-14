package connector

import (
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
)

func GetAllMangaConnectors() []models.IMangaConnector {
	return append(
		[]models.IMangaConnector{},
		GetAsuraScansConnector(),
		GetLeviatanScansConnector(),
		GetZeroScansConnector(),
		GetRealmScansConnector(),
		GetReaperScansConnector(),
		GetElarcPageConnector(),
		GetAquaMangaConnector(),
		GetMangaHasuConnector(),
		GetInfernalVoidScansConnector(),
		GetLuminousScansConnector(),
		GetFlameScansConnector(),
		GetAstraScansConnector(),
		GetMangaClashConnector(),
		GetNitroScansConnector(),
		GetAlphaScansConnector(),
		GetToonilyNetConnector(),
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

func findMangaConnector(domain string) (models.IMangaConnector, bool) {
	domain = getDomain(domain)
	conns := GetAllMangaConnectors()
	for i := 0; i < len(conns); i++ {
		if strings.Contains(domain, conns[i].GetSource().Domain) {
			return conns[i], true
		}
	}

	return nil, false
}

func NewMangaConnector(ctx web.Context, domain string) (models.IMangaConnector, error) {
	conn, ok := findMangaConnector(domain)
	if !ok {
		return nil, errors.Errorf("could not find config for %s", domain)
	}

	return conn, nil
}
