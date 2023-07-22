package connector

import (
	"github.com/pkg/errors"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
)

var MangaConnectorMap = map[string]models.IMangaConnector{}

func GetAllMangaConnectors() map[string]models.IMangaConnector {
	allConnectors := append(
		[]models.IMangaConnector{},
		GetAsuraScansConnector(),
		GetLeviatanScansConnector(),
		GetZeroScansConnector(),
		GetRealmScansConnector(),
		GetReaperScansConnector(),
		GetElarcPageConnector(),
		GetInfernalVoidScansConnector(),
		GetLuminousScansConnector(),
		GetMangaReadConnector(),
	)

	for _, conn := range allConnectors {
		MangaConnectorMap[conn.GetSource().Domain] = conn
	}

	return MangaConnectorMap
}

func GetMangaConnector(ctx web.Context, domain string) (models.IMangaConnector, error) {
	conn, ok := MangaConnectorMap[domain]
	if !ok {
		return nil, errors.Errorf("could not find config for %s", domain)
	}

	return conn, nil
}
