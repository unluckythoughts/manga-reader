package connector

import (
	"fmt"
	"sync"

	"github.com/unluckythoughts/book-reader/server/models"
	"github.com/unluckythoughts/go-microservice/v2/tools/logger"
	"github.com/unluckythoughts/go-microservice/v2/utils"
)

var (
	lock         = sync.RWMutex{}
	connectorMap = map[string]models.IConnector{}
)

func init() {
	opts := logger.Options{}
	utils.ParseEnvironmentVars(&opts)
	l := logger.New(opts)
	lock.Lock()
	defer lock.Unlock()
	for _, connector := range []models.IConnector{
		// Add connectors here
		NewFreeWebNovelConnector(l.Named("FreeWebNovel")),
		NewAsuraConnector(l.Named("AsuraScans")),
	} {
		connectorMap[connector.GetName()] = connector
	}
}

func GetAllConnectors() map[string]models.IConnector {
	lock.RLock()
	defer lock.RUnlock()
	cMap := make(map[string]models.IConnector)
	for k, v := range connectorMap {
		cMap[k] = v
	}
	return cMap
}

func GetConnector(name string) (models.IConnector, error) {
	lock.RLock()
	defer lock.RUnlock()
	conn, ok := connectorMap[name]
	if !ok {
		return nil, fmt.Errorf("could not find config for %s", name)
	}

	return conn, nil
}
