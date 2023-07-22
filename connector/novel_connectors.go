package connector

import (
	"github.com/pkg/errors"
	"github.com/unluckythoughts/go-microservice/tools/web"
	"github.com/unluckythoughts/manga-reader/models"
)

var NovelConnectorMap = map[string]models.INovelConnector{}

func GetAllNovelConnectors() map[string]models.INovelConnector {
	allconnectors := append(
		[]models.INovelConnector{},
		GetLightNovelPubConnector(),
	)

	for _, conn := range allconnectors {
		NovelConnectorMap[conn.GetSource().Domain] = conn
	}

	return NovelConnectorMap
}

// func findNovelConnector(domain string) (models.INovelConnector, bool) {
// 	domain = getDomain(domain)
// 	conns := GetAllNovelConnectors()
// 	for i := 0; i < len(conns); i++ {
// 		if strings.Contains(domain, conns[i].GetSource().Domain) {
// 			return conns[i], true
// 		}
// 	}

// 	return nil, false
// }

func NewNovelConnector(ctx web.Context, domain string) (models.INovelConnector, error) {
	conn, ok := NovelConnectorMap[domain]
	if !ok {
		return nil, errors.Errorf("could not find config for %s", domain)
	}

	return conn, nil
}
