package connector

import (
	"fmt"
	"sync"

	"github.com/unluckythoughts/book-reader/server/models"
)

var (
	lock         = sync.RWMutex{}
	connectorMap = map[string]models.IConnector{}
)

func init() {
	lock.Lock()
	defer lock.Unlock()
	for _, connector := range []models.IConnector{
		// Add connectors here
		NewFreeWebNovelConnector(),
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

func GetConnector(domain string) (models.IConnector, error) {
	lock.RLock()
	defer lock.RUnlock()
	conn, ok := connectorMap[domain]
	if !ok {
		return nil, fmt.Errorf("could not find config for %s", domain)
	}

	return conn, nil
}
