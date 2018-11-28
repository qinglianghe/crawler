package persist

import (
    "crawler/engine"
    "log"

    "github.com/olivere/elastic"
)

// ItemSaverService
type ItemSaverService struct {
    // Client elastic client
    Client *elastic.Client

    // Index elastic index
    Index string
}

// Save 用于存储Item到elastic
func (s *ItemSaverService) Save(item engine.Item, result *bool) error {
    err := Save(s.Client, item, s.Index)
    if err == nil {
        log.Printf("Item %v saved.", item)
        *result = true
    }
    return err
}
