package connector

import (
	"fmt"
	"sync"

	"github.com/unluckythoughts/manga-reader/server/models"
)

var (
	lock         = sync.RWMutex{}
	ConnectorMap = map[string]models.IConnector{}
)

func init() {
	lock.Lock()
	defer lock.Unlock()
	for _, connector := range []models.IConnector{
		// Add connectors here
	} {
		ConnectorMap[connector.GetDomain()] = connector
	}
}

func GetAllConnectors() map[string]models.IConnector {
	lock.RLock()
	defer lock.RUnlock()
	cMap := make(map[string]models.IConnector)
	for k, v := range ConnectorMap {
		cMap[k] = v
	}
	return cMap
}

func GetConnector(domain string) (models.IConnector, error) {
	lock.RLock()
	defer lock.RUnlock()
	conn, ok := ConnectorMap[domain]
	if !ok {
		return nil, fmt.Errorf("could not find config for %s", domain)
	}

	return conn, nil
}
